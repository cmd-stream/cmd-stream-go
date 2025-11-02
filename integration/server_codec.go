package integration_test

import (
	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/transport-go"
	"github.com/mus-format/dts-stream-go"
	"github.com/mus-format/ext-stream-go"
)

type ServerCodec struct{}

func (c ServerCodec) Encode(result core.Result, w transport.Writer) (
	n int, err error,
) {
	if m, ok := result.(ext.MarshallerTypedMUS); ok {
		n, err = m.MarshalTypedMUS(w)
		return
	}
	panic("cmdstream server codec: result doesn't implement the MarshallerMUS interface")
}

func (c ServerCodec) Decode(r transport.Reader) (cmd core.Cmd[struct{}],
	n int, err error,
) {
	dtm, n, err := dts.DTMSer.Unmarshal(r)
	if err != nil {
		return
	}
	switch dtm {
	case Cmd1DTM:
		cmd = Cmd1{}
	case Cmd2DTM:
		cmd = Cmd2{}
	case Cmd3DTM:
		cmd = Cmd3{}
	case Cmd4DTM:
		cmd = Cmd4{}
	}
	return
}
