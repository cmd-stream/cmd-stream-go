package cser

import (
	"net"

	base "github.com/cmd-stream/base-go"
	delegate "github.com/cmd-stream/delegate-go"
	transport "github.com/cmd-stream/transport-go"
	tcom "github.com/cmd-stream/transport-go/common"
	tser "github.com/cmd-stream/transport-go/server"
)

// TransportFactory implements the delegate.ServerTransportFactory interface.
type TransportFactory[T any] struct {
	Codec transport.Codec[base.Result, base.Cmd[T]]
	Ops   []tcom.SetOption
}

func (f TransportFactory[T]) New(conn net.Conn) delegate.ServerTransport[T] {
	return tser.New(conn, f.Codec, f.Ops...)
}
