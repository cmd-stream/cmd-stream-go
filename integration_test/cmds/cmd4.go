package cmds

import (
	"context"
	"time"

	"github.com/cmd-stream/base-go"
	"github.com/cmd-stream/cmd-stream-go/integration_test/results"
	dts "github.com/mus-format/dts-stream-go"
	muss "github.com/mus-format/mus-stream-go"
)

type Cmd4 struct{}

func (c Cmd4) Exec(ctx context.Context, seq base.Seq, at time.Time,
	receiver struct{},
	proxy base.Proxy,
) (err error) {
	_, err = proxy.Send(seq, results.NewResult(true))
	return
}

func (c Cmd4) MarshalTypedMUS(w muss.Writer) (n int, err error) {
	return dts.DTMSer.Marshal(Cmd4DTM, w)
}

func (c Cmd4) SizeTypedMUS() (size int) {
	return dts.DTMSer.Size(Cmd4DTM)
}
