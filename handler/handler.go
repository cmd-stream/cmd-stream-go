// Package handler provides a server-side connection handler for cmd-stream
// library.
package handler

import (
	"context"
	"sync"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
)

// Handler implements the dlgt.ServerTransportHandler interface.
//
// It receives Commands sequentially and executes each in a separate goroutine
// using the Invoker.
//
// If an error occurs, the Handler closes the transport connection.
type Handler[T any] struct {
	invoker Invoker[T]
	options Options
}

// New creates a new Handler.
func New[T any](invoker Invoker[T], opts ...SetOption) *Handler[T] {
	o := Options{}
	Apply(&o, opts...)
	return &Handler[T]{
		invoker: invoker,
		options: o,
	}
}

// Handle processes incoming commands on the given transport.
func (h *Handler[T]) Handle(ctx context.Context, transport dlgt.ServerTransport[T]) error {
	var (
		wg             = &sync.WaitGroup{}
		ownCtx, cancel = context.WithCancelCause(ctx)
	)
	wg.Add(1)
	go h.receiveLoop(ownCtx, cancel, transport, wg)
	<-ownCtx.Done()

	_ = transport.Close()
	wg.Wait()
	return context.Cause(ownCtx)
}

func (h *Handler[T]) receiveLoop(ctx context.Context, cancel context.CancelCauseFunc,
	transport dlgt.ServerTransport[T],
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	proxy := NewProxy(transport)
	for {
		if h.options.CmdReceiveDuration != 0 {
			deadline := time.Now().Add(h.options.CmdReceiveDuration)
			if err := transport.SetReceiveDeadline(deadline); err != nil {
				cancel(err)
				return
			}
		}
		seq, cmd, n, err := transport.Receive()
		if err != nil {
			cancel(err)
			return
		}
		var at time.Time
		if h.options.At {
			at = time.Now()
		}
		wg.Add(1)
		go func(seq core.Seq, cmd core.Cmd[T], n int, at time.Time) {
			defer wg.Done()
			if err := h.invoker.Invoke(ctx, seq, at, n, cmd, proxy); err != nil {
				cancel(err)
			}
		}(seq, cmd, n, at)
	}
}
