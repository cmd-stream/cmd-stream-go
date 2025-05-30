package csrv

import (
	"net"

	base "github.com/cmd-stream/base-go"
	dser "github.com/cmd-stream/delegate-go/server"
	transport "github.com/cmd-stream/transport-go"
	tser "github.com/cmd-stream/transport-go/server"
)

// NewTransportFactory creates a new TransportFactory using the provided codec
// and optional transport-level configuration options.
func NewTransportFactory[T any](codec transport.Codec[base.Result, base.Cmd[T]],
	ops ...transport.SetOption,
) *TransportFactory[T] {
	return &TransportFactory[T]{codec, ops}
}

// TransportFactory implements the delegate.ServerTransportFactory interface.
//
// It creates Transports that handle encoding Results / decoding Commands over
// a network connection.
type TransportFactory[T any] struct {
	codec transport.Codec[base.Result, base.Cmd[T]]
	ops   []transport.SetOption
}

func (f TransportFactory[T]) New(conn net.Conn) dser.Transport[T] {
	return tser.New(conn, f.codec, f.ops...)
}
