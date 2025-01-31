package client

import (
	"github.com/cmd-stream/base-go"
	cs "github.com/cmd-stream/cmd-stream-go"
	"github.com/cmd-stream/transport-go"
)

// Сodec represents a generic client Codec interface. It encodes Commands and
// decodes Results.
//
// Encode method is used by the client to send Commands to the server. If
// Encode fails with an error, the Client.Send() method will return it.
//
// Decode method is used by the client to receive Resulsts from the server. If
// Decode fails with an error, the client will be closed.
type Codec[T any] cs.Codec[base.Cmd[T], base.Result]

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
