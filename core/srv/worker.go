package srv

import (
	"context"
	"net"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/ymz-ncnk/jointwork-go"
)

// WorkersFactory creates server Workers.
type WorkersFactory interface {
	New(count int, conns <-chan net.Conn, delegate core.ServerDelegate,
		callback LostConnCallback) []jointwork.Task
}

// Worker is a server worker.
//
// It receives connections from the conns channel one at a time and processes
// them using the Delegate. Additionally, it implements jointwork.Task, allowing
// it to work in conjunction with ConnReceiver.
//
// If connection processing fails with an error, the error is passed to
// LostConnCallback.
type Worker struct {
	ctx      context.Context
	cancel   context.CancelFunc
	conns    <-chan net.Conn
	delegate core.ServerDelegate
	callback LostConnCallback
}

// NewWorker creates a new Worker.
func NewWorker(conns <-chan net.Conn, delegate core.ServerDelegate,
	callback LostConnCallback,
) *Worker {
	var (
		success = false
		worker  = &Worker{
			conns:    conns,
			delegate: delegate,
			callback: callback,
		}
	)
	worker.ctx, worker.cancel = context.WithCancel(context.Background())
	defer func() {
		if !success {
			worker.cancel()
		}
	}()
	success = true
	return worker
}

func (w *Worker) Run() (err error) {
	var (
		conn net.Conn
		more bool
	)
	for {
		select {
		case <-w.ctx.Done():
			return ErrClosed
		case conn, more = <-w.conns:
			// If Shutdown was called w.Stop() will not be called, in this case
			// conns channel is closed.
			if !more {
				return nil
			}
			err = w.delegate.Handle(w.ctx, conn)
			if err != nil {
				if err == context.Canceled {
					err = ErrClosed
				}
				if w.callback != nil {
					w.callback(conn.RemoteAddr(), err)
				}
				if err == ErrClosed {
					return
				}
			}
		}
	}
}

func (w *Worker) Stop() (err error) {
	w.cancel()
	return
}

type workersFactory struct{}

func (f workersFactory) New(numWorkers int, conns <-chan net.Conn,
	delegate core.ServerDelegate, callback LostConnCallback,
) (workers []jointwork.Task) {
	workers = make([]jointwork.Task, numWorkers)
	for i := range numWorkers {
		workers[i] = NewWorker(conns, delegate, callback)
	}
	return
}
