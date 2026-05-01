package mock

import (
	"context"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/ymz-ncnk/mok"
)

type ExecFn[T any] func(ctx context.Context, seq core.Seq, at time.Time, receiver T,
	proxy core.Proxy) (err error)
type TimeoutFn func() (timeout time.Duration)

type Cmd[T any] struct {
	*mok.Mock
}

func NewCmd[T any]() Cmd[T] {
	return Cmd[T]{mok.New("Cmd")}
}

func (c Cmd[T]) RegisterExec(fn ExecFn[T]) Cmd[T] {
	c.Register("Exec", fn)
	return c
}

func (c Cmd[T]) Exec(ctx context.Context, seq core.Seq, at time.Time, receiver T,
	proxy core.Proxy,
) (err error) {
	vals, err := c.Call("Exec", ctx, seq, at, receiver, proxy)
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}
