package srv

import (
	"net"

	"github.com/cmd-stream/core-go"
	dsrv "github.com/cmd-stream/delegate-go/server"
	transport "github.com/cmd-stream/transport-go"
	tsrv "github.com/cmd-stream/transport-go/server"
)

// NewTransportFactory creates a new TransportFactory using the provided codec
// and optional transport-level configuration options.
func NewTransportFactory[T any](codec transport.Codec[core.Result, core.Cmd[T]],
	ops ...transport.SetOption,
) *TransportFactory[T] {
	return &TransportFactory[T]{codec, ops}
}

// TransportFactory implements the delegate.ServerTransportFactory interface.
//
// It creates Transports that handle encoding Results / decoding Commands over
// a network connection.
type TransportFactory[T any] struct {
	codec transport.Codec[core.Result, core.Cmd[T]]
	ops   []transport.SetOption
}

func (f TransportFactory[T]) New(conn net.Conn) dsrv.Transport[T] {
	return tsrv.New(conn, f.codec, f.ops...)
}
