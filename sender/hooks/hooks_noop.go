package hooks

import (
	"context"

	"github.com/cmd-stream/cmd-stream-go/core"
)

// NoopHooksFactory implements Factory with NoopHooks.
type NoopHooksFactory[T any] struct {
	hooks NoopHooks[T]
}

// New returns a new NoopHooks instance.
func (f NoopHooksFactory[T]) New() Hooks[T] {
	return f.hooks
}

// NoopHooks implements Hooks with empty methods.
type NoopHooks[T any] struct{}

// BeforeSend returns the context as-is.
func (h NoopHooks[T]) BeforeSend(ctx context.Context, _ core.Cmd[T]) (
	context.Context, error,
) {
	return ctx, nil
}

// OnError does nothing.
func (h NoopHooks[T]) OnError(_ context.Context, _ SentCmd[T],
	_ error) {
}

// OnResult does nothing.
func (h NoopHooks[T]) OnResult(_ context.Context, _ SentCmd[T],
	_ ReceivedResult, _ error) {
}

// OnTimeout does nothing.
func (h NoopHooks[T]) OnTimeout(_ context.Context, _ SentCmd[T],
	_ error) {
}
