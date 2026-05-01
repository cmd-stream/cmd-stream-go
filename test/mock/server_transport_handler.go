package mock

import (
	"context"

	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
	"github.com/ymz-ncnk/mok"
)

type ServerTransportHandler[T any] struct {
	*mok.Mock
}

func NewServerTransportHandler[T any]() ServerTransportHandler[T] {
	return ServerTransportHandler[T]{Mock: mok.New("TransportHandler")}
}

func (m ServerTransportHandler[T]) Handle(ctx context.Context,
	transport dlgt.ServerTransport[T],
) error {
	vals, err := m.Call("Handle", ctx, transport)
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return err
}

func (m ServerTransportHandler[T]) RegisterHandle(fn func(ctx context.Context, transport dlgt.ServerTransport[T]) error) ServerTransportHandler[T] {
	m.Register("Handle", fn)
	return m
}
