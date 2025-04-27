package mock

import (
	"net"

	"github.com/ymz-ncnk/mok"
)

type ConnFactoryNew = func() (net.Conn, error)

func NewConnFactory[T any]() ConnFactory[T] {
	return ConnFactory[T]{mok.New("ConnFactory")}
}

type ConnFactory[T any] struct {
	*mok.Mock
}

func (f ConnFactory[T]) RegisterNew(fn ConnFactoryNew) ConnFactory[T] {
	f.Register("New", fn)
	return f
}

func (f ConnFactory[T]) New() (conn net.Conn, err error) {
	result, err := f.Call("New")
	if err != nil {
		return nil, err
	}
	conn = result[0].(net.Conn)
	err, _ = result[1].(error)
	return
}
