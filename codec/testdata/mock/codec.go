package mock

import (
	"github.com/cmd-stream/transport-go"
	"github.com/ymz-ncnk/mok"
)

type Encode[T any] = func(t T, w transport.Writer) (n int, err error)
type Decode[V any] = func(r transport.Reader) (v V, n int, err error)

type Codec[T, V any] struct {
	*mok.Mock
}

func NewCodec[T, V any]() *Codec[T, V] {
	return &Codec[T, V]{mok.New("Codec")}
}

func (c *Codec[T, V]) RegisterEncode(fn Encode[T]) *Codec[T, V] {
	c.Register("Encode", fn)
	return c
}

func (c *Codec[T, V]) RegisterDecode(fn Decode[V]) *Codec[T, V] {
	c.Register("Decode", fn)
	return c
}

func (c *Codec[T, V]) Encode(t T, w transport.Writer) (n int, err error) {
	result, err := c.Call("Encode", t, w)
	if err != nil {
		panic(err)
	}
	n = result[0].(int)
	err, _ = result[1].(error)
	return
}

func (c *Codec[T, V]) Decode(r transport.Reader) (v V, n int, err error) {
	result, err := c.Call("Decode", r)
	if err != nil {
		panic(err)
	}
	v = result[0].(V)
	n = result[1].(int)
	err, _ = result[2].(error)
	return
}
