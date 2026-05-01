package client

import (
	"github.com/cmd-stream/cmd-stream-go/core"
	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
	tspt "github.com/cmd-stream/cmd-stream-go/transport"
)

type keepaliveCodecAdapter[T any] struct {
	c Codec[T]
}

func (c keepaliveCodecAdapter[T]) Encode(seq core.Seq, cmd core.Cmd[T],
	w tspt.Writer,
) (n int, err error) {
	if n, err = core.SeqMUS.Marshal(seq, w); err != nil {
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

func (c keepaliveCodecAdapter[T]) Decode(r tspt.Reader) (seq core.Seq,
	result core.Result, n int, err error,
) {
	if seq, n, err = core.SeqMUS.Unmarshal(r); err != nil {
		return
	}
	if seq == 0 {
		result = dlgt.PongResult{}
		return
	}
	var n1 int
	result, n1, err = c.c.Decode(r)
	n += n1
	return
}
