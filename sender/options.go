package sender

import (
	"github.com/cmd-stream/cmd-stream-go/sender/hooks"
	hks "github.com/cmd-stream/cmd-stream-go/sender/hooks"
)

type Options[T any] struct {
	HooksFactory hooks.HooksFactory[T]
}

func DefaultOptions[T any]() Options[T] {
	return Options[T]{
		HooksFactory: hks.NoopHooksFactory[T]{},
	}
}

type SetOption[T any] func(o *Options[T])

// WithHooksFactory sets a factory that creates new hooks for each send
// operation. Hooks can customize behavior during the sending process, such as
// logging or instrumentation.
func WithHooksFactory[T any](factory hooks.HooksFactory[T]) SetOption[T] {
	return func(o *Options[T]) {
		o.HooksFactory = factory
	}
}

func Apply[T any](o *Options[T], opts ...SetOption[T]) {
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
}
