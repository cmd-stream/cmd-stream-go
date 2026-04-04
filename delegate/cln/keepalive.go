package cln

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/cmd-stream/cmd-stream-go/delegate"
)

// KeepaliveDelegate implements the core.ClientDelegate interface.
//
// When there are no Commands to send, it initiates a Ping-Pong exchange with
// the server. It sends a Ping Command and expects a Pong Result, both
// represented as a single zero byte (like a ball being passed).
type KeepaliveDelegate[T any] struct {
	core.ClientDelegate[T]
	options    KeepaliveOptions
	alive      chan struct{}
	done       chan struct{}
	mu         sync.Mutex
	closedFlag uint32
}

// NewKeepalive creates a new KeepaliveDelegate.
func NewKeepalive[T any](delegate core.ClientDelegate[T],
	opts ...SetKeepaliveOption,
) *KeepaliveDelegate[T] {
	o := DefaultKeepaliveOptions()
	ApplyKeepalive(&o, opts...)
	return &KeepaliveDelegate[T]{
		ClientDelegate: delegate,
		options:        o,
		alive:          make(chan struct{}),
		done:           make(chan struct{}),
	}
}

func (d *KeepaliveDelegate[T]) Receive() (seq core.Seq, result core.Result,
	n int, err error,
) {
	for {
		seq, result, n, err = d.ClientDelegate.Receive()
		if err != nil {
			return
		}
		if _, ok := result.(delegate.PongResult); !ok {
			return
		}
	}
}

func (d *KeepaliveDelegate[T]) Flush() (err error) {
	if err = d.ClientDelegate.Flush(); err != nil {
		return
	}
	select {
	case d.alive <- struct{}{}:
	default:
	}
	return
}

func (d *KeepaliveDelegate[T]) Keepalive(muSn *sync.Mutex) {
	go keepalive(d, muSn)
}

func (d *KeepaliveDelegate[T]) Close() (err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if atomic.LoadUint32(&d.closedFlag) == 1 {
		return
	}
	if err = d.ClientDelegate.Close(); err != nil {
		return
	}
	atomic.StoreUint32(&d.closedFlag, 1)
	close(d.done)
	return
}

func keepalive[T any](d *KeepaliveDelegate[T], muSn *sync.Mutex) {
	timer := time.NewTimer(d.options.KeepaliveTime)
	for {
		select {
		case <-d.done:
			return
		case <-timer.C:
			if _, err := ping(muSn, 0, d); err != nil {
				d.Close()
				return
			}
			timer.Reset(d.options.KeepaliveIntvl)
		case <-d.alive:
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(d.options.KeepaliveTime)
		}
	}
}

func ping[T any](muSn *sync.Mutex, seq core.Seq, d *KeepaliveDelegate[T]) (
	n int, err error,
) {
	muSn.Lock()
	if err = d.SetSendDeadline(time.Now().Add(d.options.KeepaliveIntvl)); err != nil {
		muSn.Unlock()
		return
	}
	if n, err = d.Send(seq, delegate.PingCmd[T]{}); err != nil {
		muSn.Unlock()
		return
	}
	muSn.Unlock()
	return n, d.ClientDelegate.Flush()
}
