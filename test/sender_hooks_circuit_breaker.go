package test

import (
	"context"
	"errors"
	"testing"

	"github.com/cmd-stream/cmd-stream-go/core"
	hks "github.com/cmd-stream/cmd-stream-go/sender/hooks"
	cmock "github.com/cmd-stream/cmd-stream-go/test/mock/core"
	hksmock "github.com/cmd-stream/cmd-stream-go/test/mock/sender/hooks"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

type HooksCircuitBreakerSetup[T any] struct {
	CB    hksmock.CircuitBreaker
	Hooks hksmock.Hooks[T]
}

type HooksCircuitBreakerTestCase[T any] struct {
	Name   string
	Setup  HooksCircuitBreakerSetup[T]
	Action func(t *testing.T, h hks.Hooks[T])
	Want   HooksCircuitBreakerWant
}

type HooksCircuitBreakerWant struct {
	Mocks []*mok.Mock
}

func RunHooksCircuitBreakerTestCase[T any](t *testing.T, tc HooksCircuitBreakerTestCase[T]) {
	t.Run(tc.Name, func(t *testing.T) {
		h := hks.NewCircuitBreakerHooks[T](tc.Setup.CB, tc.Setup.Hooks)
		tc.Action(t, h)
		asserterror.EqualDeep(t, mok.CheckCalls(tc.Want.Mocks), mok.EmptyInfomap)
	})
}

type SenderHooksCircuitBreaker[T any] struct{}

// -----------------------------------------------------------------------------
// Test Cases
// -----------------------------------------------------------------------------

func (SenderHooksCircuitBreaker[T]) BeforeSend(t *testing.T) HooksCircuitBreakerTestCase[T] {
	name := "BeforeSend should call CB.Allow and return no error if allowed"

	var (
		wantCtx = context.Background()
		wantCmd = cmock.NewCmd[T]()
		cb      = hksmock.NewCircuitBreaker().RegisterAllow(
			func() bool { return true },
		)
		innerHooks = hksmock.NewHooks[T]().RegisterBeforeSend(
			func(ctx context.Context, cmd core.Cmd[T]) (context.Context, error) {
				asserterror.Equal(t, ctx, wantCtx)
				asserterror.EqualDeep(t, cmd, wantCmd)
				return ctx, nil
			},
		)
	)
	return HooksCircuitBreakerTestCase[T]{
		Name: name,
		Setup: HooksCircuitBreakerSetup[T]{
			CB:    cb,
			Hooks: innerHooks,
		},
		Action: func(t *testing.T, h hks.Hooks[T]) {
			ctx, err := h.BeforeSend(wantCtx, wantCmd)
			asserterror.Equal(t, ctx, wantCtx)
			asserterror.EqualError(t, err, nil)
		},
		Want: HooksCircuitBreakerWant{
			Mocks: []*mok.Mock{cb.Mock, innerHooks.Mock},
		},
	}
}

func (SenderHooksCircuitBreaker[T]) BeforeSendError(t *testing.T) HooksCircuitBreakerTestCase[T] {
	name := "BeforeSend should call CB.Allow and return ErrNotAllowed if not allowed"

	var (
		wantCtx = context.Background()
		wantCmd = cmock.NewCmd[T]()
		cb      = hksmock.NewCircuitBreaker().RegisterAllow(
			func() bool { return false },
		)
	)
	return HooksCircuitBreakerTestCase[T]{
		Name: name,
		Setup: HooksCircuitBreakerSetup[T]{
			CB: cb,
		},
		Action: func(t *testing.T, h hks.Hooks[T]) {
			ctx, err := h.BeforeSend(wantCtx, wantCmd)
			asserterror.Equal(t, ctx, wantCtx)
			asserterror.EqualError(t, err, hks.ErrNotAllowed)
		},
		Want: HooksCircuitBreakerWant{
			Mocks: []*mok.Mock{cb.Mock},
		},
	}
}

func (SenderHooksCircuitBreaker[T]) OnError(t *testing.T) HooksCircuitBreakerTestCase[T] {
	name := "OnError should call CB.Fail"

	var (
		wantCtx     = context.Background()
		wantSentCmd = hks.SentCmd[T]{}
		wantErr     = errors.New("error")
		cb          = hksmock.NewCircuitBreaker().RegisterFail(func() {})
		innerHooks  = hksmock.NewHooks[T]().RegisterOnError(
			func(ctx context.Context, sentCmd hks.SentCmd[T], err error) {
				asserterror.Equal(t, ctx, wantCtx)
				asserterror.EqualDeep(t, sentCmd, wantSentCmd)
				asserterror.EqualError(t, err, wantErr)
			},
		)
	)
	return HooksCircuitBreakerTestCase[T]{
		Name: name,
		Setup: HooksCircuitBreakerSetup[T]{
			CB:    cb,
			Hooks: innerHooks,
		},
		Action: func(t *testing.T, h hks.Hooks[T]) {
			h.OnError(wantCtx, wantSentCmd, wantErr)
		},
		Want: HooksCircuitBreakerWant{
			Mocks: []*mok.Mock{cb.Mock, innerHooks.Mock},
		},
	}
}

func (SenderHooksCircuitBreaker[T]) OnResult(t *testing.T) HooksCircuitBreakerTestCase[T] {
	name := "OnResult should call CB.Success"

	var (
		wantCtx        = context.Background()
		wantSentCmd    = hks.SentCmd[T]{}
		wantRecvResult = hks.ReceivedResult{}
		wantErr        = errors.New("error")
		cb             = hksmock.NewCircuitBreaker().RegisterSuccess(func() {})
		innerHooks     = hksmock.NewHooks[T]().RegisterOnResult(
			func(ctx context.Context, sentCmd hks.SentCmd[T],
				recvResult hks.ReceivedResult, err error,
			) {
				asserterror.Equal(t, ctx, wantCtx)
				asserterror.EqualDeep(t, sentCmd, wantSentCmd)
				asserterror.EqualDeep(t, recvResult, wantRecvResult)
				asserterror.EqualError(t, err, wantErr)
			},
		)
	)
	return HooksCircuitBreakerTestCase[T]{
		Name: name,
		Setup: HooksCircuitBreakerSetup[T]{
			CB:    cb,
			Hooks: innerHooks,
		},
		Action: func(t *testing.T, h hks.Hooks[T]) {
			h.OnResult(wantCtx, wantSentCmd, wantRecvResult, wantErr)
		},
		Want: HooksCircuitBreakerWant{
			Mocks: []*mok.Mock{cb.Mock, innerHooks.Mock},
		},
	}
}

func (SenderHooksCircuitBreaker[T]) OnTimeout(t *testing.T) HooksCircuitBreakerTestCase[T] {
	name := "OnTimeout should call CB.Fail"

	var (
		wantCtx     = context.Background()
		wantSentCmd = hks.SentCmd[T]{}
		wantErr     = errors.New("error")
		cb          = hksmock.NewCircuitBreaker().RegisterFail(func() {})
		innerHooks  = hksmock.NewHooks[T]().RegisterOnTimeout(
			func(ctx context.Context, sentCmd hks.SentCmd[T], err error) {
				asserterror.Equal(t, ctx, wantCtx)
				asserterror.EqualDeep(t, sentCmd, wantSentCmd)
				asserterror.EqualError(t, err, wantErr)
			},
		)
	)
	return HooksCircuitBreakerTestCase[T]{
		Name: name,
		Setup: HooksCircuitBreakerSetup[T]{
			CB:    cb,
			Hooks: innerHooks,
		},
		Action: func(t *testing.T, h hks.Hooks[T]) {
			h.OnTimeout(wantCtx, wantSentCmd, wantErr)
		},
		Want: HooksCircuitBreakerWant{
			Mocks: []*mok.Mock{cb.Mock, innerHooks.Mock},
		},
	}
}
