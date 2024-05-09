package ct

import (
	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
	dts "github.com/mus-format/mus-stream-dts-go"
)

type ServerCodec struct{}

func (c ServerCodec) Encode(result base.Result, w transport.Writer) (
	err error) {
	switch r := result.(type) {
	case Result:
		_, err = MarshalResultMUS(r, w)
	default:
		panic("unexpected result")
	}
	return
}

func (c ServerCodec) Decode(r transport.Reader) (cmd base.Cmd[Receiver],
	err error) {
	dtm, _, err := dts.UnmarshalDTM(r)
	switch dtm {
	case Cmd1DTM:
		cmd = Cmd1{}
	case Cmd2DTM:
		cmd = Cmd2{}
	case Cmd3DTM:
		cmd = Cmd3{}
	}
	return
}
