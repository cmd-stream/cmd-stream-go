package ccln

import (
	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/cmd-stream-go/codec"
	"github.com/cmd-stream/delegate-go"
	"github.com/cmd-stream/transport-go"
)

type keepaliveCodecAdapter[T any] struct {
	c Codec[T]
}

func (c keepaliveCodecAdapter[T]) Encode(seq base.Seq, cmd base.Cmd[T],
	w transport.Writer) (n int, err error) {
	if n, err = codec.SeqMUS.Marshal(seq, w); err != nil {
		return
	}
	if seq == 0 { // It is a delegate.PingCmd.
		return
	}
	var n1 int
	n1, err = c.c.Encode(cmd, w)
	n += n1
	return
}

func (c keepaliveCodecAdapter[T]) Decode(r transport.Reader) (seq base.Seq,
	result base.Result, n int, err error) {
	if seq, n, err = codec.SeqMUS.Unmarshal(r); err != nil {
		return
	}
	if seq == 0 {
		result = delegate.PongResult{}
		return
	}
	var n1 int
	result, n1, err = c.c.Decode(r)
	n += n1
	return
}
