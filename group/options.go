package group

import (
	"errors"

	cln "github.com/cmd-stream/cmd-stream-go/client"
)

// Options defines the configuration settings for creating a ClientGroup.
type Options[T any] struct {
	Factory    DispatchStrategyFactory[T]
	Reconnect  bool
	ClientOpts []cln.SetOption
}

// DefaultOptions returns the default group configuration.
func DefaultOptions[T any]() Options[T] {
	return Options[T]{
		Factory: RoundRobinStrategyFactory[T]{},
	}
}

// Validate ensures the group configuration is valid.
func (o Options[T]) Validate() error {
	if o.Factory == nil {
		return errors.New("factory is nil")
	}
	return nil
}

// SetOption defines a function for configuring Options.
type SetOption[T any] func(o *Options[T])

// WithFactory sets the dispatch strategy factory for the client group.
//
// The dispatch strategy determines how Commands are distributed among clients.
// For example, a round-robin strategy will rotate client usage evenly.
func WithFactory[T any](factory DispatchStrategyFactory[T]) SetOption[T] {
	return func(o *Options[T]) { o.Factory = factory }
}

// WithReconnect enables automatic reconnection for all clients in the group.
//
// When this option is set, reconnect-capable clients are created, which attempt
// to re-establish the connection if it's lost during communication.
func WithReconnect[T any]() SetOption[T] {
	return func(o *Options[T]) { o.Reconnect = true }
}

// WithClient sets client-specific options to be applied when initializing
// each client in the group.
func WithClient[T any](opts ...cln.SetOption) SetOption[T] {
	return func(o *Options[T]) { o.ClientOpts = append(o.ClientOpts, opts...) }
}

// Apply applies the given options to the Options struct.
func Apply[T any](o *Options[T], opts ...SetOption[T]) error {
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
	return o.Validate()
}
