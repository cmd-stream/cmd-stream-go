package mock

import (
	"github.com/cmd-stream/cmd-stream-go/core"
	tspt "github.com/cmd-stream/cmd-stream-go/transport"
	"github.com/ymz-ncnk/mok"
)

type (
	DecodeFn[V any] func(r tspt.Reader) (seq core.Seq, val V, n int, err error)
	EncodeFn[T any] func(seq core.Seq, val T, w tspt.Writer) (n int, err error)
)

func NewCodec[T, V any]() Codec[T, V] {
	return Codec[T, V]{
		Mock: mok.New("Codec"),
	}
}

type Codec[T, V any] struct {
	*mok.Mock
}

func (c Codec[T, V]) RegisterDecode(fn DecodeFn[V]) Codec[T, V] {
	c.Register("Decode", fn)
	return c
}

func (c Codec[T, V]) RegisterEncode(fn EncodeFn[T]) Codec[T, V] {
	c.Register("Encode", fn)
	return c
}

func (c Codec[T, V]) Decode(r tspt.Reader) (seq core.Seq, val V, n int, err error) {
	vals, err := c.Call("Decode", r)
	if err != nil {
		panic(err)
	}
	seq = vals[0].(core.Seq)
	val, _ = vals[1].(V)
	n = vals[2].(int)
	err, _ = vals[3].(error)
	return
}

func (c Codec[T, V]) Encode(seq core.Seq, val T, w tspt.Writer) (
	n int, err error,
) {
	vals, err := c.Call("Encode", seq, val, w)
	if err != nil {
		panic(err)
	}
	n = vals[0].(int)
	err, _ = vals[1].(error)
	return
}
