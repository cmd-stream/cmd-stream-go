package delegate

import (
	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
	"github.com/ymz-ncnk/mok"
)

type ClientTransportFactory[T any] struct {
	*mok.Mock
}

func NewClientTransportFactory[T any]() ClientTransportFactory[T] {
	return ClientTransportFactory[T]{mok.New("TransportFactory")}
}

func (m ClientTransportFactory[T]) RegisterNew(fn func() (dlgt.ClientTransport[T], error)) ClientTransportFactory[T] {
	m.Register("New", fn)
	return m
}

func (m ClientTransportFactory[T]) RegisterNewN(n int, fn func() (dlgt.ClientTransport[T], error)) ClientTransportFactory[T] {
	m.RegisterN("New", n, fn)
	return m
}

func (m ClientTransportFactory[T]) New() (dlgt.ClientTransport[T], error) {
	vals, err := m.Call("New")
	if err != nil {
		return nil, err
	}
	transport, _ := vals[0].(dlgt.ClientTransport[T])
	err, _ = vals[1].(error)
	return transport, err
}
