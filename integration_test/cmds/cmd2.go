package cmds

import (
	"context"
	"time"

	"github.com/cmd-stream/cmd-stream-go/integration_test/results"
	"github.com/cmd-stream/core-go"
	"github.com/mus-format/dts-stream-go"
	"github.com/mus-format/mus-stream-go"
)

type Cmd2 struct{}

func (c Cmd2) Exec(ctx context.Context, seq core.Seq, at time.Time,
	receiver struct{},
	proxy core.Proxy,
) (err error) {
	time.Sleep(500 * time.Millisecond)
	_, err = proxy.Send(seq, results.NewResult(true))
	return
}

func (c Cmd2) MarshalTypedMUS(w mus.Writer) (n int, err error) {
	return dts.DTMSer.Marshal(Cmd2DTM, w)
}

func (c Cmd2) SizeTypedMUS() (size int) {
	return dts.DTMSer.Size(Cmd2DTM)
}
