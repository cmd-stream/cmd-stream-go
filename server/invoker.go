package cser

import (
	"context"
	"time"

	"github.com/cmd-stream/base-go"
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

func (i Invoker[T]) Invoke(ctx context.Context, at time.Time, seq base.Seq,
	cmd base.Cmd[T], proxy base.Proxy) error {
	return cmd.Exec(ctx, at, seq, i.receiver, proxy)
}
