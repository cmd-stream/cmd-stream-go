package mock

import (
	"net"

	ccln "github.com/cmd-stream/cmd-stream-go/client"
	"github.com/ymz-ncnk/mok"
)

type ClientFactoryNew[T any] = func(codec ccln.Codec[T], conn net.Conn,
	o ccln.Options) (ccln.Client[T], error)

func NewClientFactory[T any]() ClientFactory[T] {
	return ClientFactory[T]{mok.New("Client")}
}

type ClientFactory[T any] struct {
	*mok.Mock
}

func (f ClientFactory[T]) RegisterNew(fn ClientFactoryNew[T]) ClientFactory[T] {
	f.Register("New", fn)
	return f
}

func (f ClientFactory[T]) New(codec ccln.Codec[T], conn net.Conn,
	o ccln.Options) (client ccln.Client[T], err error) {
	result, err := f.Call("New", mok.SafeVal[ccln.Codec[T]](codec),
		mok.SafeVal[net.Conn](conn), o)
	if err != nil {
		return nil, err
	}
	client, _ = result[0].(ccln.Client[T])
	err, _ = result[1].(error)
	return
}
