package intest

import (
	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
	dts "github.com/mus-format/dts-stream-go"
	exts "github.com/mus-format/ext-mus-stream-go"
)

type ServerCodec struct{}

func (c ServerCodec) Encode(result base.Result, w transport.Writer) (
	err error) {
	if m, ok := result.(exts.MarshallerTypedMUS); ok {
		_, err = m.MarshalTypedMUS(w)
		return
	}
	panic("result doesn't implement the MarshallerMUS interface")
}

func (c ServerCodec) Decode(r transport.Reader) (cmd base.Cmd[Receiver],
	err error) {
	dtm, _, err := dts.DTMSer.Unmarshal(r)
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
