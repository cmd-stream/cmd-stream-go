package mock

import (
	"context"
	"net"

	"github.com/ymz-ncnk/mok"
)

type ServerHandleFn func(ctx context.Context, conn net.Conn) (err error)

type ServerDelegate struct {
	*mok.Mock
}

func NewServerDelegate() ServerDelegate {
	return ServerDelegate{Mock: mok.New("ServerDelegate")}
}

func (m ServerDelegate) RegisterHandle(fn ServerHandleFn) ServerDelegate {
	m.Register("Handle", fn)
	return m
}

func (m ServerDelegate) RegisterHandleN(n int, fn ServerHandleFn) ServerDelegate {
	m.RegisterN("Handle", n, fn)
	return m
}

func (m ServerDelegate) Handle(ctx context.Context, conn net.Conn) (err error) {
	vals, err := m.Call("Handle", ctx, conn)
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}
