package cmds

import (
	"context"
	"time"

	"github.com/cmd-stream/cmd-stream-go/integration_test/results"
	"github.com/cmd-stream/core-go"
	"github.com/mus-format/dts-stream-go"
	"github.com/mus-format/mus-stream-go"
)

type Cmd4 struct{}

func (c Cmd4) Exec(ctx context.Context, seq core.Seq, at time.Time,
	receiver struct{},
	proxy core.Proxy,
) (err error) {
	_, err = proxy.Send(seq, results.NewResult(true))
	return
}

func (c Cmd4) MarshalTypedMUS(w mus.Writer) (n int, err error) {
	return dts.DTMSer.Marshal(Cmd4DTM, w)
}

func (c Cmd4) SizeTypedMUS() (size int) {
	return dts.DTMSer.Size(Cmd4DTM)
}
