package ccln

import (
	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/cmd-stream-go/codec"
	"github.com/cmd-stream/transport-go"
)

// Codec defines a generic client codec interface for encoding Commands and
// decoding Results.
//   - Encode is used by the client to send Commands to the server. If encoding
//     fails, Client.Send returns the corresponding error.
//   - Decode is used by the client to receive Results from the server. If
//     decoding fails, the client is closed automatically.
type Codec[T any] codec.Codec[base.Cmd[T], base.Result]

// AdaptCodec adapts the provided Codec.
func AdaptCodec[T any](codec Codec[T], o Options) transport.Codec[base.Cmd[T], base.Result] {
	if o.Keepalive != nil {
		return keepaliveCodecAdapter[T]{codec}
	}
	return codecAdapter[T]{codec}
}

type codecAdapter[T any] struct {
	c Codec[T]
}

func (c codecAdapter[T]) Encode(seq base.Seq, cmd base.Cmd[T],
	w transport.Writer) (n int, err error) {
	if n, err = codec.SeqMUS.Marshal(seq, w); err != nil {
		return
	}
	var n1 int
	n1, err = c.c.Encode(cmd, w)
	n += n1
	return
}

func (c codecAdapter[T]) Decode(r transport.Reader) (seq base.Seq,
	result base.Result, n int, err error) {
	if seq, n, err = codec.SeqMUS.Unmarshal(r); err != nil {
		return
	}
	var n1 int
	result, n1, err = c.c.Decode(r)
	n += n1
	return
}
