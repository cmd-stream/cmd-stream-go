package sender

import (
	"crypto/tls"
	"errors"

	grp "github.com/cmd-stream/cmd-stream-go/group"
)

// MakeOptions defines the configuration for creating a new Sender and its
// underlying client group.
type MakeOptions[T any] struct {
	Group        []grp.SetOption[T]
	Sender       []SetOption[T]
	TLSConfig    *tls.Config
	ClientsCount int
}

func DefaultMakeOptions[T any]() MakeOptions[T] {
	return MakeOptions[T]{
		ClientsCount: 1,
	}
}

func (o MakeOptions[T]) Validate() error {
	if o.ClientsCount <= 0 {
		return errors.New("clients count must be positive")
	}
	return nil
}

type SetMakeOption[T any] func(o *MakeOptions[T])

// WithGroup adds options for the underlying client group.
func WithGroup[T any](opts ...grp.SetOption[T]) SetMakeOption[T] {
	return func(o *MakeOptions[T]) { o.Group = append(o.Group, opts...) }
}

// WithSender adds options for the Sender.
func WithSender[T any](opts ...SetOption[T]) SetMakeOption[T] {
	return func(o *MakeOptions[T]) { o.Sender = append(o.Sender, opts...) }
}

// WithTLSConfig sets the TLS configuration for the server connections.
func WithTLSConfig[T any](conf *tls.Config) SetMakeOption[T] {
	return func(o *MakeOptions[T]) { o.TLSConfig = conf }
}

// WithClientsCount sets the number of clients to be created in the group.
func WithClientsCount[T any](count int) SetMakeOption[T] {
	return func(o *MakeOptions[T]) { o.ClientsCount = count }
}

func ApplyMake[T any](o *MakeOptions[T], opts ...SetMakeOption[T]) error {
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
	return o.Validate()
}
