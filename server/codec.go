package server

import (
	"github.com/cmd-stream/cmd-stream-go/codec"
	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/delegate-go"
	"github.com/cmd-stream/transport-go"
)

// Codec defines a generic server codec interface for encoding Results and
// decoding Commands.
//   - Encode is used by the server to send Results back to the client. If
//     encoding fails, the server closes the corresponding client connection.
//   - Decode is used by the server to receive Commands. If decoding fails, the
//     server closes the corresponding client connection.
type Codec[T any] codec.Codec[core.Result, core.Cmd[T]]

// AdaptCodec adapts the provided Codec.
func AdaptCodec[T any](codec Codec[T], o Options) transport.Codec[core.Result, core.Cmd[T]] {
	return codecAdapter[T]{codec}
}

type codecAdapter[T any] struct {
	c Codec[T]
}

func (c codecAdapter[T]) Encode(seq core.Seq, result core.Result,
	w transport.Writer,
) (n int, err error) {
	if _, err = codec.SeqMUS.Marshal(seq, w); err != nil {
		return
	}
	if seq == 0 { // It is a delegate.PongResult.
		return
	}
	n, err = c.c.Encode(result, w)
	return
}

func (c codecAdapter[T]) Decode(r transport.Reader) (seq core.Seq,
	cmd core.Cmd[T], n int, err error,
) {
	if seq, _, err = codec.SeqMUS.Unmarshal(r); err != nil {
		return
	}
	if seq == 0 {
		cmd = delegate.PingCmd[T]{}
		return
	}
	cmd, n, err = c.c.Decode(r)
	return
}
