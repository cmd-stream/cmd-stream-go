package client

import (
	"github.com/cmd-stream/base-go"
	cs "github.com/cmd-stream/cmd-stream-go"
	"github.com/cmd-stream/delegate-go"
	"github.com/cmd-stream/transport-go"
)

type keepaliveCodecAdapter[T any] struct {
	c Codec[T]
}

func (c keepaliveCodecAdapter[T]) Encode(seq base.Seq, cmd base.Cmd[T],
	w transport.Writer) (err error) {
	if _, err = cs.MarshalSeqMUS(seq, w); err != nil {
		return
	}
	if seq == 0 { // It is a delegate.PingCmd.
		return
	}
	return c.c.Encode(cmd, w)
}

func (c keepaliveCodecAdapter[T]) Decode(r transport.Reader) (seq base.Seq,
	result base.Result, err error) {
	if seq, _, err = cs.UnmarshalSeqMUS(r); err != nil {
		return
	}
	if seq == 0 {
		result = delegate.PongResult{}
		return
	}
	result, err = c.c.Decode(r)
	return
}

func (c keepaliveCodecAdapter[T]) Size(cmd base.Cmd[T]) int {
	return c.c.Size(cmd)
}
