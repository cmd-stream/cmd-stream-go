package handler

import (
	"context"

	"github.com/cmd-stream/cmd-stream-go/core"
)

// Invoker executes Commands on the server.
//
// A single Invoker instance handles Commands concurrently in separate goroutines,
// so it must be thread-safe. The Invoke method can act as a central handler for
// common operations across multiple Commands (e.g. logging).
type Invoker[T any] interface {
	Invoke(ctx context.Context, bytesRead int, cmd core.Cmd[T], proxy core.Proxy) error
}

// InvokerFn is a functional implementation of the Invoker interface.
type InvokerFn[T any] func(ctx context.Context, bytesRead int, cmd core.Cmd[T],
	proxy core.Proxy) error

// Invoke executes the given command.
func (i InvokerFn[T]) Invoke(ctx context.Context, bytesRead int, cmd core.Cmd[T],
	proxy core.Proxy) error {
	return i(ctx, bytesRead, cmd, proxy)
}
