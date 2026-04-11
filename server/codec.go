package server

import (
	"github.com/cmd-stream/cmd-stream-go/codec"
	"github.com/cmd-stream/cmd-stream-go/core"
	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
	tspt "github.com/cmd-stream/cmd-stream-go/transport"
)

// Codec defines a generic server codec interface for encoding Results and
// decoding Commands.
//   - Encode is used by the server to send Results back to the client. If
//     encoding fails, the server closes the corresponding client connection.
//   - Decode is used by the server to receive Commands. If decoding fails, the
//     server closes the corresponding client connection.
type Codec[T any] codec.Codec[core.Result, core.Cmd[T]]

// AdaptCodec adapts the provided Codec.
// AdaptCodec adapts the provided Codec.
func AdaptCodec[T any](codec Codec[T], _ Options) tspt.Codec[core.Result, core.Cmd[T]] {
	return codecAdapter[T]{codec}
}

type codecAdapter[T any] struct {
	c Codec[T]
}

func (c codecAdapter[T]) Encode(seq core.Seq, result core.Result,
	w tspt.Writer,
) (n int, err error) {
	if _, err = core.SeqMUS.Marshal(seq, w); err != nil {
		return
	}
	if seq == 0 { // It is a delegate.PongResult.
		return
	}
	n, err = c.c.Encode(result, w)
	return
}

func (c codecAdapter[T]) Decode(r tspt.Reader) (seq core.Seq,
	cmd core.Cmd[T], n int, err error,
) {
	if seq, _, err = core.SeqMUS.Unmarshal(r); err != nil {
		return
	}
	if seq == 0 {
		cmd = dlgt.PingCmd[T]{}
		return
	}
	cmd, n, err = c.c.Decode(r)
	return
}
