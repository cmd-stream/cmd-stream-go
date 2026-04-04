package srv

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
)

const (
	inProgress int = iota
	shutdown
	closed
)

// NewConnReceiver creates a new ConnReceiver.
func NewConnReceiver(listener core.Listener, conns chan net.Conn,
	opts ...SetConnReceiverOption,
) *ConnReceiver {
	o := ConnReceiverOptions{}
	ApplyConnReceiver(&o, opts...)

	return &ConnReceiver{
		listener: listener,
		conns:    conns,
		options:  o,
		stopped:  make(chan struct{}),
	}
}

// ConnReceiver listens for incoming connections and adds them to the conns
// channel.
//
// It can wait for the first connection for a specified duration, after which it
// stops. ConnReceiver also implements the jointwork.Task interface, allowing
// it to work in conjunction with Workers.
type ConnReceiver struct {
	listener core.Listener
	conns    chan net.Conn
	options  ConnReceiverOptions
	stopped  chan struct{}
	state    int
	mu       sync.Mutex
}

func (r *ConnReceiver) Run() (err error) {
	defer func() {
		err = r.resolveErr(err)
		r.cleanup()
	}()
	err = r.acceptFirstConn()
	if err != nil {
		return
	}
	return r.acceptConns()
}

// Shutdown stops ConnReceiver - the Run() method returns nil, which allows
// Workers to finish their work.
func (r *ConnReceiver) Shutdown() (err error) {
	return r.terminate(shutdown)
}

// Stop stops ConnReceiver - the Run() method returns ErrClosed.
func (r *ConnReceiver) Stop() (err error) {
	return r.terminate(closed)
}

func (r *ConnReceiver) acceptFirstConn() (err error) {
	if r.options.FirstConnTimeout != 0 {
		defer func() {
			if err == nil {
				err = r.listener.SetDeadline(time.Time{})
			}
		}()
		deadline := time.Now().Add(r.options.FirstConnTimeout)
		err = r.listener.SetDeadline(deadline)
		if err != nil {
			return err
		}
	}
	conn, err := r.listener.Accept()
	if err != nil {
		return err
	}
	return r.queueConn(conn)
}

func (r *ConnReceiver) acceptConns() (err error) {
	var conn net.Conn
	for {
		conn, err = r.listener.Accept()
		if err != nil {
			return
		}
		err = r.queueConn(conn)
		if err != nil {
			return
		}
	}
}

func (r *ConnReceiver) queueConn(conn net.Conn) error {
	select {
	case <-r.stopped:
		_ = conn.Close()
		return ErrClosed
	case r.conns <- conn:
		return nil
	}
}

func (r *ConnReceiver) terminate(state int) (err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state == inProgress {
		if err = r.listener.Close(); err != nil {
			return
		}
		r.state = state
		close(r.stopped)
	}
	return
}

func (r *ConnReceiver) resolveErr(err error) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	switch r.state {
	case inProgress:
		return err
	case shutdown:
		return nil
	case closed:
		return ErrClosed
	default:
		panic(fmt.Sprintf("CMD-STREAM INTERNAL BUG: unexpected state %d", r.state))
	}
}

func (r *ConnReceiver) cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()

	close(r.conns)
	if r.state == shutdown {
		return
	}
	for conn := range r.conns {
		_ = conn.Close()
	}
}
