// Package client provides a thread-safe, asynchronous cmd-stream client.
//
// The Client uses a core.ClientDelegate for sending Commands, receiving Results, and
// managing connection state.
package cln

import (
	"context"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
)

// UnexpectedResultCallback processes unexpected Results received from the server.
//
// It is invoked when the sequence number of a Result does not match the sequence
// number of any Command sent by the client that is awaiting a Result.
type UnexpectedResultCallback func(seq core.Seq, result core.Result)

// New creates a new client.
func New[T any](delegate core.ClientDelegate[T], opts ...SetOption) *Client[T] {
	o := Options{}
	Apply(&o, opts...)
	var (
		ctx, cancel        = context.WithCancel(context.Background())
		flagFl      uint32 = 0
		client             = &Client[T]{
			cancel:   cancel,
			delegate: delegate,
			options:  o,
			pending: pending[T]{
				m: make(map[core.Seq]chan<- core.AsyncResult),
			},
			done:   make(chan struct{}),
			flagFl: &flagFl,
			chFl:   make(chan error, 1),
		}
	)
	if keepaliveDelegate, ok := delegate.(core.KeepaliveDelegate[T]); ok {
		keepaliveDelegate.Keepalive(&client.muSn)
	}
	go receive(ctx, client)
	return client
}

// It utilizes core.ClientDelegate for communication tasks such as sending Commands,
// receiving Results, and managing deadlines. If the connection is lost, the
// client will close, and Client.Error() will return the corresponding connection
// error.
//
// To close the client, use Client.Close(). You can track the completion of this
// process by checking Client.Done():
//
//	err = client.Close()
//	...
//	select {
//	case <-time.NewTimer(time.Second).C:
//		err = errors.New("timeout exceeded")
//		...
//	case <-client.Done():
//	}
type Client[T any] struct {
	cancel   context.CancelFunc
	delegate core.ClientDelegate[T]
	options  Options
	seq      core.Seq
	done     chan struct{}

	flagFl *uint32
	chFl   chan error

	pending pending[T]
	err     errStatus
	state   state
	muSn    sync.Mutex
}

// Send transmits a Command to the server.
//
// Received Results from the server are added to the results channel. If the
// channel lacks sufficient capacity, retrieving results for all Commands may
// hang.
//
// Each Command is assigned a unique sequence number, starting from 1:
//   - The first Command is sent with `seq == 1`, the second with `seq == 2`, etc.
//   - `seq == 0` is reserved for the Ping-Pong mechanism, which maintains
//     connection liveness.
//
// Returns the sequence number of the Command and any error encountered
// (non-nil if the Command was not sent successfully).
func (c *Client[T]) Send(cmd core.Cmd[T], results chan<- core.AsyncResult) (
	seq core.Seq, n int, err error,
) {
	if c.state.Closed() {
		err = ErrClosed
		return
	}
	var chFl chan error
	c.muSn.Lock()
	chFl = c.chFl
	seq = c.nextSeq()
	c.pending.add(seq, results)

	n, err = c.delegate.Send(seq, cmd)
	if err != nil {
		c.muSn.Unlock()
		c.Forget(seq)
		err = NewClientError(err)
		return
	}
	c.muSn.Unlock()

	if err = c.flush(seq, chFl); err != nil {
		err = NewClientError(err)
	}
	return
}

// nextSeq should be c.muSn protected.
func (c *Client[T]) nextSeq() core.Seq {
	c.seq++
	return c.seq
}

// SendWithDeadline sends a Command with a specified deadline.
//
// This method behaves like Send but allows setting a deadline for Command
// execution. Use it when you need to enforce a time limit on the Command's
// processing.
func (c *Client[T]) SendWithDeadline(deadline time.Time,
	cmd core.Cmd[T],
	results chan<- core.AsyncResult,
) (seq core.Seq, n int, err error) {
	if c.state.Closed() {
		err = ErrClosed
		return
	}

	var chFl chan error
	c.muSn.Lock()
	chFl = c.chFl
	seq = c.nextSeq()
	c.pending.add(seq, results)

	err = c.delegate.SetSendDeadline(deadline)
	if err != nil {
		c.muSn.Unlock()
		c.Forget(seq)
		err = NewClientError(err)
		return
	}

	n, err = c.delegate.Send(seq, cmd)
	if err != nil {
		c.muSn.Unlock()
		c.Forget(seq)
		err = NewClientError(err)
		return
	}
	c.muSn.Unlock()

	if err = c.flush(seq, chFl); err != nil {
		err = NewClientError(err)
	}
	return
}

// Has checks if the Command with the specified sequence number has been sent
// by the Client and still waiting for the Result.
func (c *Client[T]) Has(seq core.Seq) bool {
	_, pst := c.pending.get(seq)
	return pst
}

// Forget makes the Client to forget about the Command which still waiting for
// the result.
//
// After calling Forget, all the results of the corresponding Command will be
// handled with UnexpectedResultCallback.
func (c *Client[T]) Forget(seq core.Seq) {
	c.pending.remove(seq)
}

// Done returns a channel that is closed when the Client terminates.
func (c *Client[T]) Done() <-chan struct{} {
	return c.done
}

// Error returns a connection error.
func (c *Client[T]) Error() (err error) {
	return c.err.Get()
}

// Close terminates the underlying connection and closes the Client.
//
// All Commands waiting for the results will receive an error
// (AsyncResult.Error != nil).
func (c *Client[T]) Close() (err error) {
	if !c.state.SetClosed() {
		return
	}
	c.cancel()
	err = c.delegate.Close()
	if err != nil {
		err = NewClientError(err)
	}
	return
}

func (c *Client[T]) receive(ctx context.Context) (err error) {
	defer func() {
		c.pending.failAll(err)
	}()
	var (
		seq     core.Seq
		result  core.Result
		n       int
		results chan<- core.AsyncResult
		pst     bool
	)
	for {
		seq, result, n, err = c.delegate.Receive()
		if err != nil {
			return
		}
		if result.LastOne() {
			results, pst = c.pending.pop(seq)
		} else {
			results, pst = c.pending.get(seq)
		}
		if !pst && c.options.UnexpectedResultCallback != nil {
			c.options.UnexpectedResultCallback(seq, result)
			continue
		}
		select {
		case <-ctx.Done():
			return context.Canceled
		case results <- core.AsyncResult{Seq: seq, BytesRead: n, Result: result}:
			continue
		}
	}
}

func (c *Client[T]) flush(seq core.Seq, chFl chan error) (err error) {
	if swapped := atomic.CompareAndSwapUint32(c.flagFl, 0, 1); !swapped {
		err = <-chFl
		if err != nil {
			chFl <- err
			c.Forget(seq)
		}
		return
	}
	c.muSn.Lock()
	err = c.delegate.Flush()
	if err != nil {
		c.chFl <- err
	} else {
		close(c.chFl)
	}
	c.resetChFl()
	c.muSn.Unlock()
	if err != nil {
		c.Forget(seq)
	}
	return
}

func (c *Client[T]) resetChFl() {
	c.chFl = make(chan error, 1)
	atomic.CompareAndSwapUint32(c.flagFl, 1, 0)
}

func (c *Client[T]) exit(cause error) {
	if c.state.SetClosed() {
		_ = c.delegate.Close()
	}
	c.cancel()
	c.err.Set(cause)
	close(c.done)
}

func (c *Client[T]) correctErr(err error) error {
	if c.state.Closed() {
		return ErrClosed
	}
	return err
}

func receive[T any](ctx context.Context, client *Client[T]) {
Start:
	err := client.receive(ctx)
	if err != nil {
		err = client.correctErr(err)
		if _, ok := err.(net.Error); ok || err == io.EOF { // TODO Test EOF.
			if reconnectDelegate, ok := client.delegate.(core.ReconnectDelegate[T]); ok {
				if err = reconnectDelegate.Reconnect(); err == nil {
					goto Start
				}
			}
		}
	}
	client.exit(err)
}
