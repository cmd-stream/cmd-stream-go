package mock

import (
	"github.com/cmd-stream/transport-go"
	"github.com/ymz-ncnk/mok"
)

type Encode[T any] = func(t T, w transport.Writer) (err error)
type Decode[V any] = func(r transport.Reader) (v V, err error)

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

func (c *Codec[T, V]) Encode(t T, w transport.Writer) (err error) {
	result, err := c.Call("Encode", t, w)
	if err != nil {
		return err
	}
	err, _ = result[0].(error)
	return
}

func (c *Codec[T, V]) Decode(r transport.Reader) (v V, err error) {
	result, err := c.Call("Decode", r)
	if err != nil {
		return v, err
	}
	v = result[0].(V)
	err, _ = result[1].(error)
	return
}
