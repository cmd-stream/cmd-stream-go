package ccln

import (
	"github.com/cmd-stream/base-go"
	cs "github.com/cmd-stream/cmd-stream-go"
	"github.com/cmd-stream/transport-go"
)

// Codec represents a generic client Codec interface responsible for encoding
// Commands and decoding Results.
//
//   - Encode is used by the client to send Commands to the server. If encoding
//     fails, Client.Send() will return the corresponding error.
//   - Decode is used by the client to receive Results from the server. If
//     decoding fails, the client will be closed.
type Codec[T any] cs.Codec[base.Cmd[T], base.Result]

type codecAdapter[T any] struct {
	c Codec[T]
}

func (c codecAdapter[T]) Encode(seq base.Seq, cmd base.Cmd[T],
	w transport.Writer) (err error) {
	if _, err = cs.SeqMUS.Marshal(seq, w); err != nil {
		return
	}
	return c.c.Encode(cmd, w)
}

func (c codecAdapter[T]) Decode(r transport.Reader) (seq base.Seq,
	result base.Result, err error) {
	if seq, _, err = cs.SeqMUS.Unmarshal(r); err != nil {
		return
	}
	result, err = c.c.Decode(r)
	return
}
