package server

import (
	"github.com/cmd-stream/base-go"
	cs "github.com/cmd-stream/cmd-stream-go"
	"github.com/cmd-stream/delegate-go"
	"github.com/cmd-stream/transport-go"
)

// Codec represents a generic server Codec inteface. It encodes Results and
// decodes Commands.
//
// Encode method is used by the server to send results. If Encode fails with an
// error, the server closes the coresponding client connection.
//
// Decode method is used by the server to receive commands. If it fails with an
// error, the server closes the corresponding client connection.
type Codec[T any] cs.Codec[base.Result, base.Cmd[T]]

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
