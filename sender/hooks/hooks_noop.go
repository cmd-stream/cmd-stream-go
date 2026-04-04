package hooks

import (
	"context"

	"github.com/cmd-stream/cmd-stream-go/core"
)

type NoopHooksFactory[T any] struct {
	hooks NoopHooks[T]
}

func (f NoopHooksFactory[T]) New() Hooks[T] {
	return f.hooks
}

type NoopHooks[T any] struct{}

func (h NoopHooks[T]) BeforeSend(ctx context.Context, cmd core.Cmd[T]) (
	context.Context, error,
) {
	return ctx, nil
}

func (h NoopHooks[T]) OnError(ctx context.Context, sentCmd SentCmd[T],
	err error) {
}

func (h NoopHooks[T]) OnResult(ctx context.Context, sentCmd SentCmd[T],
	recvResult ReceivedResult, err error) {
}

func (h NoopHooks[T]) OnTimeout(ctx context.Context, sentCmd SentCmd[T],
	err error) {
}
