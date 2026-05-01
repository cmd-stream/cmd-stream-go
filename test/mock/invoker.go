package mock

import (
	"context"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/ymz-ncnk/mok"
)

type InvokeFn[T any] func(ctx context.Context, seq core.Seq, at time.Time,
	bytesRead int, cmd core.Cmd[T], proxy core.Proxy) (err error)

type Invoker[T any] struct {
	*mok.Mock
}

func NewInvoker[T any]() Invoker[T] {
	return Invoker[T]{mok.New("Invoker")}
}

func (i Invoker[T]) RegisterInvoke(fn InvokeFn[T]) Invoker[T] {
	i.Register("Invoke", fn)
	return i
}

func (i Invoker[T]) RegisterInvokeN(n int, fn InvokeFn[T]) Invoker[T] {
	i.RegisterN("Invoke", n, fn)
	return i
}

func (i Invoker[T]) Invoke(ctx context.Context, seq core.Seq, at time.Time,
	bytesRead int, cmd core.Cmd[T], proxy core.Proxy,
) (err error) {
	vals, err := i.Call("Invoke", ctx, seq, at, bytesRead, cmd, proxy)
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}
