package sender

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

func BeforeSendTestCase(t *testing.T) HooksCircuitBreakerTestCase {
	name := "BeforeSend should call CB.Allow and return no error if allowed"

	var (
		wantCtx = context.Background()
		wantCmd = cmock.NewCmd[any]()
		cb      = hksmock.NewCircuitBreaker().RegisterAllow(
			func() bool { return true },
		)
		innerHooks = hksmock.NewHooks[any]().RegisterBeforeSend(
			func(ctx context.Context, cmd core.Cmd[any]) (context.Context, error) {
				asserterror.Equal(t, ctx, wantCtx)
				asserterror.EqualDeep(t, cmd, wantCmd)
				return ctx, nil
			},
		)
	)
	return HooksCircuitBreakerTestCase{
		Name: name,
		Setup: HooksCircuitBreakerSetup{
			CB:    cb,
			Hooks: innerHooks,
		},
		Action: func(t *testing.T, h hks.Hooks[any]) {
			ctx, err := h.BeforeSend(wantCtx, wantCmd)
			asserterror.Equal(t, ctx, wantCtx)
			asserterror.EqualError(t, err, nil)
		},
		Want: HooksCircuitBreakerWant{
			Mocks: []*mok.Mock{cb.Mock, innerHooks.Mock},
		},
	}
}

func BeforeSendErrorTestCase() HooksCircuitBreakerTestCase {
	name := "BeforeSend should call CB.Allow and return ErrNotAllowed if not allowed"

	var (
		wantCtx = context.Background()
		wantCmd = cmock.NewCmd[any]()
		cb      = hksmock.NewCircuitBreaker().RegisterAllow(
			func() bool { return false },
		)
	)
	return HooksCircuitBreakerTestCase{
		Name: name,
		Setup: HooksCircuitBreakerSetup{
			CB: cb,
		},
		Action: func(t *testing.T, h hks.Hooks[any]) {
			ctx, err := h.BeforeSend(wantCtx, wantCmd)
			asserterror.Equal(t, ctx, wantCtx)
			asserterror.EqualError(t, err, hks.ErrNotAllowed)
		},
		Want: HooksCircuitBreakerWant{
			Mocks: []*mok.Mock{cb.Mock},
		},
	}
}

func OnErrorTestCase(t *testing.T) HooksCircuitBreakerTestCase {
	name := "OnError should call CB.Fail"

	var (
		wantCtx     = context.Background()
		wantSentCmd = hks.SentCmd[any]{}
		wantErr     = errors.New("error")
		cb          = hksmock.NewCircuitBreaker().RegisterFail(func() {})
		innerHooks  = hksmock.NewHooks[any]().RegisterOnError(
			func(ctx context.Context, sentCmd hks.SentCmd[any], err error) {
				asserterror.Equal(t, ctx, wantCtx)
				asserterror.EqualDeep(t, sentCmd, wantSentCmd)
				asserterror.EqualError(t, err, wantErr)
			},
		)
	)
	return HooksCircuitBreakerTestCase{
		Name: name,
		Setup: HooksCircuitBreakerSetup{
			CB:    cb,
			Hooks: innerHooks,
		},
		Action: func(t *testing.T, h hks.Hooks[any]) {
			h.OnError(wantCtx, wantSentCmd, wantErr)
		},
		Want: HooksCircuitBreakerWant{
			Mocks: []*mok.Mock{cb.Mock, innerHooks.Mock},
		},
	}
}

func OnResultTestCase(t *testing.T) HooksCircuitBreakerTestCase {
	name := "OnResult should call CB.Success"

	var (
		wantCtx        = context.Background()
		wantSentCmd    = hks.SentCmd[any]{}
		wantRecvResult = hks.ReceivedResult{}
		wantErr        = errors.New("error")
		cb             = hksmock.NewCircuitBreaker().RegisterSuccess(func() {})
		innerHooks     = hksmock.NewHooks[any]().RegisterOnResult(
			func(ctx context.Context, sentCmd hks.SentCmd[any],
				recvResult hks.ReceivedResult, err error,
			) {
				asserterror.Equal(t, ctx, wantCtx)
				asserterror.EqualDeep(t, sentCmd, wantSentCmd)
				asserterror.EqualDeep(t, recvResult, wantRecvResult)
				asserterror.EqualError(t, err, wantErr)
			},
		)
	)
	return HooksCircuitBreakerTestCase{
		Name: name,
		Setup: HooksCircuitBreakerSetup{
			CB:    cb,
			Hooks: innerHooks,
		},
		Action: func(t *testing.T, h hks.Hooks[any]) {
			h.OnResult(wantCtx, wantSentCmd, wantRecvResult, wantErr)
		},
		Want: HooksCircuitBreakerWant{
			Mocks: []*mok.Mock{cb.Mock, innerHooks.Mock},
		},
	}
}

func OnTimeoutTestCase(t *testing.T) HooksCircuitBreakerTestCase {
	name := "OnTimeout should call CB.Fail"

	var (
		wantCtx     = context.Background()
		wantSentCmd = hks.SentCmd[any]{}
		wantErr     = errors.New("error")
		cb          = hksmock.NewCircuitBreaker().RegisterFail(func() {})
		innerHooks  = hksmock.NewHooks[any]().RegisterOnTimeout(
			func(ctx context.Context, sentCmd hks.SentCmd[any], err error) {
				asserterror.Equal(t, ctx, wantCtx)
				asserterror.EqualDeep(t, sentCmd, wantSentCmd)
				asserterror.EqualError(t, err, wantErr)
			},
		)
	)
	return HooksCircuitBreakerTestCase{
		Name: name,
		Setup: HooksCircuitBreakerSetup{
			CB:    cb,
			Hooks: innerHooks,
		},
		Action: func(t *testing.T, h hks.Hooks[any]) {
			h.OnTimeout(wantCtx, wantSentCmd, wantErr)
		},
		Want: HooksCircuitBreakerWant{
			Mocks: []*mok.Mock{cb.Mock, innerHooks.Mock},
		},
	}
}
