package client

import (
	"github.com/cmd-stream/base-go"
	cs "github.com/cmd-stream/cmd-stream-go"
	"github.com/cmd-stream/transport-go"
)

// Сodec helps the Сlient to encode commands and decode the results.
//
// Size method should return the size of a command in bytes.
type Codec[T any] interface {
	Encode(cmd base.Cmd[T], w transport.Writer) (err error)
	Decode(r transport.Reader) (result base.Result, err error)
	Size(cmd base.Cmd[T]) int
}

type codecAdapter[T any] struct {
	c Codec[T]
}

func (c codecAdapter[T]) Encode(seq base.Seq, cmd base.Cmd[T],
	w transport.Writer) (err error) {
	if _, err = cs.MarshalSeqMUS(seq, w); err != nil {
		return
	}
	return c.c.Encode(cmd, w)
}

func (c codecAdapter[T]) Decode(r transport.Reader) (seq base.Seq,
	result base.Result, err error) {
	if seq, _, err = cs.UnmarshalSeqMUS(r); err != nil {
		return
	}
	result, err = c.c.Decode(r)
	return
}

func (c codecAdapter[T]) Size(cmd base.Cmd[T]) int {
	return c.c.Size(cmd)
}
