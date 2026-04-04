package delegate

import (
	"context"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
)

type PingCmd[T any] struct{}

func (c PingCmd[T]) Exec(ctx context.Context, seq core.Seq, at time.Time,
	receiver T, proxy core.Proxy,
) (err error) {
	_, err = proxy.SendWithDeadline(time.Time{}, seq, PongResult{})
	return
}

type PongResult struct{}

func (r PongResult) LastOne() bool {
	return true
}
