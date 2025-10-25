package intest

import (
	"github.com/cmd-stream/cmd-stream-go/integration_test/cmds"
	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/transport-go"
	dts "github.com/mus-format/dts-stream-go"
	exts "github.com/mus-format/ext-mus-stream-go"
)

type ServerCodec struct{}

func (c ServerCodec) Encode(result core.Result, w transport.Writer) (
	n int, err error,
) {
	if m, ok := result.(exts.MarshallerTypedMUS); ok {
		n, err = m.MarshalTypedMUS(w)
		return
	}
	panic("result doesn't implement the MarshallerMUS interface")
}

func (c ServerCodec) Decode(r transport.Reader) (cmd core.Cmd[struct{}],
	n int, err error,
) {
	dtm, n, err := dts.DTMSer.Unmarshal(r)
	if err != nil {
		return
	}
	switch dtm {
	case cmds.Cmd1DTM:
		cmd = cmds.Cmd1{}
	case cmds.Cmd2DTM:
		cmd = cmds.Cmd2{}
	case cmds.Cmd3DTM:
		cmd = cmds.Cmd3{}
	case cmds.Cmd4DTM:
		cmd = cmds.Cmd4{}
	}
	return
}
