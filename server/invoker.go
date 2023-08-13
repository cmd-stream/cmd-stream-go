package server

import (
	"context"
	"time"

	base "github.com/cmd-stream/base-go"
)

// DefInvoker is a default implementation of the handler.Invoker interface.
type DefInvoker[T any] struct {
	receiver T
}

func (i DefInvoker[T]) Invoke(ctx context.Context, at time.Time, seq base.Seq,
	cmd base.Cmd[T],
	proxy base.Proxy,
) error {
	return cmd.Exec(ctx, at, seq, i.receiver, proxy)
}
