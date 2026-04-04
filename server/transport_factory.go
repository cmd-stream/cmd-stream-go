package server

import (
	"net"

	"github.com/cmd-stream/cmd-stream-go/core"
	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
	tspt "github.com/cmd-stream/cmd-stream-go/transport"
	tsrv "github.com/cmd-stream/cmd-stream-go/transport/srv"
)

// TransportFactory implements the delegate.ServerTransportFactory interface.
//
// It creates Transports that handle encoding Results / decoding Commands over
// a network connection.
type TransportFactory[T any] struct {
	codec tspt.Codec[core.Result, core.Cmd[T]]
	opts  []tspt.SetOption
}

// NewTransportFactory creates a new TransportFactory using the provided codec
// and optional transport-level configuration options.
func NewTransportFactory[T any](codec tspt.Codec[core.Result, core.Cmd[T]],
	opts ...tspt.SetOption,
) *TransportFactory[T] {
	return &TransportFactory[T]{codec, opts}
}

func (f TransportFactory[T]) New(conn net.Conn) dlgt.ServerTransport[T] {
	return tsrv.New(conn, f.codec, f.opts...)
}
