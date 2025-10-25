package intest

import (
	"github.com/cmd-stream/cmd-stream-go/integration_test/results"
	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/transport-go"
	exts "github.com/mus-format/ext-mus-stream-go"
)

type ClientCodec struct{}

func (c ClientCodec) Encode(cmd core.Cmd[struct{}],
	w transport.Writer,
) (n int, err error) {
	if m, ok := cmd.(exts.MarshallerTypedMUS); ok {
		n, err = m.MarshalTypedMUS(w)
		return
	}
	panic("cmd doesn't implement the ext.MarshallerTypedMUS interface")
}

func (c ClientCodec) Decode(r transport.Reader) (result core.Result, n int,
	err error,
) {
	result, n, err = results.ResultMUS.Unmarshal(r)
	return
}
