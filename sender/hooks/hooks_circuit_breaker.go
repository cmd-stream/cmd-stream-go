package hooks

import (
	"context"

	"github.com/cmd-stream/cmd-stream-go/core"
)

// CircuitBreakerHooks checks whether the circuit breaker allows the operation
// before sending. If not, it returns ErrNotAllowed, otherwise the
// corresponding method of the inner Hooks is called.
type CircuitBreakerHooks[T any] struct {
	cb    CircuitBreaker
	hooks Hooks[T]
}

// NewCircuitBreakerHooks creates a new CircuitBreakerHooks.
func NewCircuitBreakerHooks[T any](cb CircuitBreaker,
	hooks Hooks[T],
) CircuitBreakerHooks[T] {
	return CircuitBreakerHooks[T]{cb, hooks}
}

// BeforeSend checks with the circuit breaker before sending.
func (h CircuitBreakerHooks[T]) BeforeSend(ctx context.Context, cmd core.Cmd[T]) (
	context.Context, error,
) {
	if !h.cb.Allow() {
		return ctx, ErrNotAllowed
	}
	return h.hooks.BeforeSend(ctx, cmd)
}

// OnError records a failure and calls the inner hook.
func (h CircuitBreakerHooks[T]) OnError(ctx context.Context, sentCmd SentCmd[T],
	err error,
) {
	h.cb.Fail()
	h.hooks.OnError(ctx, sentCmd, err)
}

// OnResult records a success and calls the inner hook.
func (h CircuitBreakerHooks[T]) OnResult(ctx context.Context, sentCmd SentCmd[T],
	recvResult ReceivedResult, err error,
) {
	h.cb.Success()
	h.hooks.OnResult(ctx, sentCmd, recvResult, err)
}

// OnTimeout records a failure and calls the inner hook.
func (h CircuitBreakerHooks[T]) OnTimeout(ctx context.Context, sentCmd SentCmd[T],
	err error,
) {
	h.cb.Fail()
	h.hooks.OnTimeout(ctx, sentCmd, err)
}

// -----------------------------------------------------------------------------

// CircuitBreakerHooksFactory can be used to create hooks that incorporate
// circuit breaker logic during the command sending process.
type CircuitBreakerHooksFactory[T any] struct {
	cb      CircuitBreaker
	factory HooksFactory[T]
}

// NewCircuitBreakerHooksFactory creates a new CircuitBreakerHooksFactory.
func NewCircuitBreakerHooksFactory[T any](cb CircuitBreaker,
	factory HooksFactory[T],
) CircuitBreakerHooksFactory[T] {
	return CircuitBreakerHooksFactory[T]{cb, factory}
}

// New returns a new CircuitBreakerHooks instance.
func (f CircuitBreakerHooksFactory[T]) New() Hooks[T] {
	return NewCircuitBreakerHooks(f.cb, f.factory.New())
}
