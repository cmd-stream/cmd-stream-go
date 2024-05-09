package ct

import (
	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/transport-go"
	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/mus-stream-dts-go"
)

type ClientCodec struct{}

func (c ClientCodec) Encode(cmd base.Cmd[Receiver],
	w transport.Writer) (err error) {
	var dtm com.DTM
	switch cmd.(type) {
	case Cmd1:
		dtm = Cmd1DTM
	case Cmd2:
		dtm = Cmd2DTM
	case Cmd3:
		dtm = Cmd3DTM
	default:
		panic("unexpected cmd type")
	}
	_, err = dts.MarshalDTM(dtm, w)
	return
}

func (c ClientCodec) Decode(r transport.Reader) (result base.Result,
	err error) {
	result, _, err = UnmarshalResultMUS(r)
	return
}

func (c ClientCodec) Size(cmd base.Cmd[Receiver]) (size int) {
	return 0
}
