package cmds

import (
	"context"
	"time"

	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/cmd-stream-go/integration_test/results"
	dts "github.com/mus-format/dts-stream-go"
	muss "github.com/mus-format/mus-stream-go"
)

type Cmd3 struct{}

func (c Cmd3) Exec(ctx context.Context, seq base.Seq, at time.Time,
	receiver struct{},
	proxy base.Proxy,
) (err error) {
	time.Sleep(500 * time.Millisecond)
	_, err = proxy.Send(seq, results.NewResult(true))
	return
}

func (c Cmd3) MarshalTypedMUS(w muss.Writer) (n int, err error) {
	return dts.DTMSer.Marshal(Cmd3DTM, w)
}

func (c Cmd3) SizeTypedMUS() (size int) {
	return dts.DTMSer.Size(Cmd3DTM)
}
