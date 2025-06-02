package server

import (
	"context"
	"time"

	"github.com/cmd-stream/core-go"
)

// NewInvoker creates a new Invoker.
func NewInvoker[T any](receiver T) Invoker[T] {
	return Invoker[T]{receiver}
}

// Invoker is the default implementation of the handler.Invoker interface.
//
// It performs no additional actions other than executing the provided Command.
type Invoker[T any] struct {
	receiver T
}

func (i Invoker[T]) Invoke(ctx context.Context, seq core.Seq, at time.Time,
	bytesRead int, cmd core.Cmd[T], proxy core.Proxy) error {
	return cmd.Exec(ctx, seq, at, i.receiver, proxy)
}
