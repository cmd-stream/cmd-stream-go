package cser

import (
	"github.com/cmd-stream/base-go"
	cs "github.com/cmd-stream/cmd-stream-go"
	"github.com/cmd-stream/delegate-go"
	"github.com/cmd-stream/transport-go"
)

// Codec defines a generic server-side Codec interface responsible for encoding
// Results and decoding Commands.
//
//   - Encode is used by the server to send Results. If encoding fails, the
//     server closes the corresponding client connection.
//   - Decode is used by the server to receive Commands. If decoding fails, the
//     server closes the corresponding client connection.
type Codec[T any] cs.Codec[base.Result, base.Cmd[T]]

type codecAdapter[T any] struct {
	c Codec[T]
}

func (c codecAdapter[T]) Encode(seq base.Seq, result base.Result,
	w transport.Writer) (err error) {
	if _, err = cs.SeqMUS.Marshal(seq, w); err != nil {
		return
	}
	if seq == 0 { // It is a delegate.PongResult.
		return
	}
	return c.c.Encode(result, w)
}

func (c codecAdapter[T]) Decode(r transport.Reader) (seq base.Seq,
	cmd base.Cmd[T], err error) {
	if seq, _, err = cs.SeqMUS.Unmarshal(r); err != nil {
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
