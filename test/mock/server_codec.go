package mock

import (
	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/cmd-stream/cmd-stream-go/transport"
	"github.com/ymz-ncnk/mok"
)

type (
	ServerDecodeFn[T any] func(r transport.Reader) (seq core.Seq, cmd core.Cmd[T], n int, err error)
	ServerEncodeFn        func(seq core.Seq, result core.Result, w transport.Writer) (n int, err error)
)

func NewServerCodec[T any]() ServerCodec[T] {
	return ServerCodec[T]{
		Mock: mok.New("ServerCodec"),
	}
}

type ServerCodec[T any] struct {
	*mok.Mock
}

func (c ServerCodec[T]) RegisterDecode(fn ServerDecodeFn[T]) ServerCodec[T] {
	c.Register("Decode", fn)
	return c
}

func (c ServerCodec[T]) RegisterEncode(fn ServerEncodeFn) ServerCodec[T] {
	c.Register("Encode", fn)
	return c
}

func (c ServerCodec[T]) RegisterSize(
	fn func(result core.Result) (size int),
) ServerCodec[T] {
	c.Register("Size", fn)
	return c
}

func (c ServerCodec[T]) Decode(r transport.Reader) (seq core.Seq, cmd core.Cmd[T],
	n int, err error,
) {
	vals, err := c.Call("Decode", r)
	if err != nil {
		panic(err)
	}
	seq = vals[0].(core.Seq)
	cmd, _ = vals[1].(core.Cmd[T])
	n = vals[2].(int)
	err, _ = vals[3].(error)
	return
}

func (c ServerCodec[T]) Encode(seq core.Seq, result core.Result, w transport.Writer) (
	n int, err error,
) {
	vals, err := c.Call("Encode", seq, result, w)
	if err != nil {
		panic(err)
	}
	n = vals[0].(int)
	err, _ = vals[1].(error)
	return
}

func (c ServerCodec[T]) Size(result core.Result) (size int) {
	vals, err := c.Call("Size", result)
	if err != nil {
		panic(err)
	}
	size = vals[0].(int)
	return
}
