package sender

import (
	hks "github.com/cmd-stream/cmd-stream-go/sender/hooks"
)

// Options represents the configurations for a Sender.
type Options[T any] struct {
	HooksFactory hks.HooksFactory[T]
}

// DefaultOptions returns the default Sender configuration.
func DefaultOptions[T any]() Options[T] {
	return Options[T]{
		HooksFactory: hks.NoopHooksFactory[T]{},
	}
}

// SetOption defines a function for configuring Options.
type SetOption[T any] func(o *Options[T])

// WithHooksFactory sets a factory that creates new hooks for each send
// operation. Hooks can customize behavior during the sending process, such as
// logging or instrumentation.
func WithHooksFactory[T any](factory hks.HooksFactory[T]) SetOption[T] {
	return func(o *Options[T]) {
		o.HooksFactory = factory
	}
}

// Apply applies the given options to the Options struct.
func Apply[T any](o *Options[T], opts ...SetOption[T]) {
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
}
