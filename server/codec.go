package server

import (
	"github.com/cmd-stream/base-go"
	cs "github.com/cmd-stream/cmd-stream-go"
	"github.com/cmd-stream/delegate-go"
	"github.com/cmd-stream/transport-go"
)

// Codec helps the Server decode commands and encode the results.
type Codec[T any] interface {
	Encode(result base.Result, w transport.Writer) (err error)
	Decode(r transport.Reader) (cmd base.Cmd[T], err error)
}

type codecAdapter[T any] struct {
	c Codec[T]
}

func (c codecAdapter[T]) Encode(seq base.Seq, result base.Result,
	w transport.Writer) (err error) {
	if _, err = cs.MarshalSeqMUS(seq, w); err != nil {
		return
	}
	if seq == 0 { // It is a delegate.PongResult.
		return
	}
	return c.c.Encode(result, w)
}

func (c codecAdapter[T]) Decode(r transport.Reader) (seq base.Seq,
	cmd base.Cmd[T], err error) {
	if seq, _, err = cs.UnmarshalSeqMUS(r); err != nil {
		return
	}
	if seq == 0 {
		cmd = delegate.PingCmd[T]{}
		return
	}
	cmd, err = c.c.Decode(r)
	return
}

func (c codecAdapter[T]) Size(result base.Result) int {
	panic("unimplemented")
}
