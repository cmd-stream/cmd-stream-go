package mock

import (
	"net"

	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
	"github.com/ymz-ncnk/mok"
)

type ServerTransportFactory[T any] struct {
	*mok.Mock
}

func NewServerTransportFactory[T any]() ServerTransportFactory[T] {
	return ServerTransportFactory[T]{Mock: mok.New("TransportFactory")}
}

func (m ServerTransportFactory[T]) New(conn net.Conn) dlgt.ServerTransport[T] {
	vals, err := m.Call("New", conn)
	if err != nil {
		panic(err)
	}
	transport, _ := vals[0].(dlgt.ServerTransport[T])
	return transport
}

func (m ServerTransportFactory[T]) RegisterNew(fn func(conn net.Conn) dlgt.ServerTransport[T]) ServerTransportFactory[T] {
	m.Register("New", fn)
	return m
}
