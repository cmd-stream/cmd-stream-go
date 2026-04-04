package client

import (
	"github.com/cmd-stream/cmd-stream-go/codec"
	"github.com/cmd-stream/cmd-stream-go/core"
	tdpt "github.com/cmd-stream/cmd-stream-go/transport"
)

// Codec defines a generic client codec interface for encoding Commands and
// decoding Results.
//   - Encode is used by the client to send Commands to the server. If encoding
//     fails, Client.Send returns the corresponding error.
//   - Decode is used by the client to receive Results from the server. If
//     decoding fails, the client is closed automatically.
type Codec[T any] codec.Codec[core.Cmd[T], core.Result]

// AdaptCodec adapts the provided Codec.
func AdaptCodec[T any](codec Codec[T], o Options) tdpt.Codec[core.Cmd[T], core.Result] {
	if o.Keepalive != nil {
		return keepaliveCodecAdapter[T]{codec}
	}
	return codecAdapter[T]{codec}
}

type codecAdapter[T any] struct {
	c Codec[T]
}

func (c codecAdapter[T]) Encode(seq core.Seq, cmd core.Cmd[T],
	w tdpt.Writer,
) (n int, err error) {
	if n, err = core.SeqMUS.Marshal(seq, w); err != nil {
		return
	}
	var n1 int
	n1, err = c.c.Encode(cmd, w)
	n += n1
	return
}

func (c codecAdapter[T]) Decode(r tdpt.Reader) (seq core.Seq,
	result core.Result, n int, err error,
) {
	if seq, n, err = core.SeqMUS.Unmarshal(r); err != nil {
		return
	}
	var n1 int
	result, n1, err = c.c.Decode(r)
	n += n1
	return
}
