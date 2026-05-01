package mock

import (
	"context"

	"github.com/cmd-stream/cmd-stream-go/core"
	hks "github.com/cmd-stream/cmd-stream-go/sender/hooks"
	"github.com/ymz-ncnk/mok"
)

type (
	BeforeSendFn[T any] func(ctx context.Context, cmd core.Cmd[T]) (context.Context, error)
	OnErrorFn[T any]    func(ctx context.Context, sentCmd hks.SentCmd[T], err error)
	OnResultFn[T any]   func(ctx context.Context, sentCmd hks.SentCmd[T],
		recvResult hks.ReceivedResult, err error)
	OnTimeoutFn[T any] func(ctx context.Context, sentCmd hks.SentCmd[T], err error)
)

func NewHooks[T any]() Hooks[T] {
	return Hooks[T]{Mock: mok.New("Hooks")}
}

type Hooks[T any] struct {
	*mok.Mock
}

func (m Hooks[T]) RegisterBeforeSendN(n int, fn BeforeSendFn[T]) Hooks[T] {
	m.RegisterN("BeforeSend", n, fn)
	return m
}

func (m Hooks[T]) RegisterBeforeSend(fn BeforeSendFn[T]) Hooks[T] {
	m.Register("BeforeSend", fn)
	return m
}

func (m Hooks[T]) RegisterOnErrorN(n int, fn OnErrorFn[T]) Hooks[T] {
	m.RegisterN("OnError", n, fn)
	return m
}

func (m Hooks[T]) RegisterOnError(fn OnErrorFn[T]) Hooks[T] {
	m.Register("OnError", fn)
	return m
}

func (m Hooks[T]) RegisterOnResultN(n int, fn OnResultFn[T]) Hooks[T] {
	m.RegisterN("OnResult", n, fn)
	return m
}

func (m Hooks[T]) RegisterOnResult(fn OnResultFn[T]) Hooks[T] {
	m.Register("OnResult", fn)
	return m
}

func (m Hooks[T]) RegisterOnTimeoutN(n int, fn OnTimeoutFn[T]) Hooks[T] {
	m.RegisterN("OnTimeout", n, fn)
	return m
}

func (m Hooks[T]) RegisterOnTimeout(fn OnTimeoutFn[T]) Hooks[T] {
	m.Register("OnTimeout", fn)
	return m
}

func (m Hooks[T]) BeforeSend(ctx context.Context, cmd core.Cmd[T]) (context.Context, error) {
	vals, err := m.Call("BeforeSend", ctx, cmd)
	if err != nil {
		panic(err)
	}
	ctx, _ = vals[0].(context.Context)
	err, _ = vals[1].(error)
	return ctx, err
}

func (m Hooks[T]) OnError(ctx context.Context, sentCmd hks.SentCmd[T], err error) {
	_, e := m.Call("OnError", ctx, sentCmd, err)
	if e != nil {
		panic(e)
	}
}

func (m Hooks[T]) OnResult(ctx context.Context, sentCmd hks.SentCmd[T],
	recvResult hks.ReceivedResult, err error,
) {
	_, e := m.Call("OnResult", ctx, sentCmd, recvResult, err)
	if e != nil {
		panic(e)
	}
}

func (m Hooks[T]) OnTimeout(ctx context.Context, sentCmd hks.SentCmd[T], err error) {
	_, e := m.Call("OnTimeout", ctx, sentCmd, err)
	if e != nil {
		panic(e)
	}
}
