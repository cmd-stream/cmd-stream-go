package server

import (
	"net"

	base "github.com/cmd-stream/base-go"
	delegate "github.com/cmd-stream/delegate-go"
	transport "github.com/cmd-stream/transport-go"
	transport_common "github.com/cmd-stream/transport-go/common"
	transport_server "github.com/cmd-stream/transport-go/server"
)

// TransportFactory is an implementation of the delegate.ServerTransportFactory.
type TransportFactory[T any] struct {
	Conf  transport_common.Conf
	Codec transport.Codec[base.Result, base.Cmd[T]]
}

func (f TransportFactory[T]) New(conn net.Conn) delegate.ServerTransport[T] {
	return transport_server.New(f.Conf, conn, f.Codec)
}
