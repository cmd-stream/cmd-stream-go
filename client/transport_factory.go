package cln

import (
	"github.com/cmd-stream/core-go"
	dcln "github.com/cmd-stream/delegate-go/client"
	"github.com/cmd-stream/transport-go"
	tcln "github.com/cmd-stream/transport-go/client"
)

// NewTransportFactory creates a new TransportFactory.
func NewTransportFactory[T any](codec transport.Codec[core.Cmd[T], core.Result],
	factory ConnFactory, ops ...transport.SetOption) *TransportFactory[T] {
	return &TransportFactory[T]{
		codec:   codec,
		factory: factory,
		ops:     ops,
	}
}

// TransportFactory creates new Transport instances.
//
// It encapsulates the logic for establishing a new connection and applying
// optional transport-level configuration.
type TransportFactory[T any] struct {
	codec   transport.Codec[core.Cmd[T], core.Result]
	factory ConnFactory
	ops     []transport.SetOption
}

// New creates a Transport by establishing a new connection.
//
// Returns an error if connection creation fails.
func (f TransportFactory[T]) New() (transport dcln.Transport[T],
	err error) {
	conn, err := f.factory.New()
	if err != nil {
		return
	}
	transport = tcln.New(conn, f.codec, f.ops...)
	return
}
