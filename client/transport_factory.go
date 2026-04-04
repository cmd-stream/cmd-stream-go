package client

import (
	"github.com/cmd-stream/cmd-stream-go/core"
	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
	tspt "github.com/cmd-stream/cmd-stream-go/transport"
	tcln "github.com/cmd-stream/cmd-stream-go/transport/cln"
)

// TransportFactory creates new Transport instances.
//
// It encapsulates the logic for establishing a new connection and applying
// optional transport-level configuration.
type TransportFactory[T any] struct {
	codec   tspt.Codec[core.Cmd[T], core.Result]
	factory ConnFactory
	ops     []tspt.SetOption
}

// NewTransportFactory creates a new TransportFactory.
func NewTransportFactory[T any](codec tspt.Codec[core.Cmd[T], core.Result],
	factory ConnFactory, ops ...tspt.SetOption,
) *TransportFactory[T] {
	return &TransportFactory[T]{
		codec:   codec,
		factory: factory,
		ops:     ops,
	}
}

// New creates a Transport by establishing a new connection.
//
// Returns an error if connection creation fails.
func (f TransportFactory[T]) New() (transport dlgt.ClientTransport[T],
	err error,
) {
	conn, err := f.factory.New()
	if err != nil {
		return
	}
	transport = tcln.New(conn, f.codec, f.ops...)
	return
}
