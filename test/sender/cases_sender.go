package sender

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	grp "github.com/cmd-stream/cmd-stream-go/group"
	sndr "github.com/cmd-stream/cmd-stream-go/sender"
	hks "github.com/cmd-stream/cmd-stream-go/sender/hooks"
	cmock "github.com/cmd-stream/cmd-stream-go/test/mock/core"
	smock "github.com/cmd-stream/cmd-stream-go/test/mock/sender"
	hmock "github.com/cmd-stream/cmd-stream-go/test/mock/sender/hooks"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

func SendSuccessTestCase(t *testing.T) SenderTestCase[any] {
	name := "Send should return no error if successful"

	var (
		wantCmd    = core.Cmd[any](nil)
		seq        = core.Seq(1)
		wantResult = core.Result(nil)
		group      = smock.NewSenderGroup[any]()
	)
	group.RegisterSend(
		func(c core.Cmd[any], r chan<- core.AsyncResult) (
			core.Seq, grp.ClientID, int, error,
		) {
			asserterror.EqualDeep(t, c, wantCmd)
			r <- core.AsyncResult{Seq: seq, Result: wantResult}
			return seq, grp.ClientID(1), 10, nil
		},
	)
	return SenderTestCase[any]{
		Name: name,
		Setup: SenderSetup[any]{
			Group: group,
		},
		Action: func(t *testing.T, s sndr.Sender[any]) {
			result, err := s.Send(context.Background(), wantCmd)
			asserterror.EqualDeep(t, result, wantResult)
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{group.Mock},
	}
}

func SendBeforeSendErrorTestCase() SenderTestCase[any] {
	name := "Send should return an error if Hooks.BeforeSend fails"

	var (
		wantErr = errors.New("BeforeSend error")
		hooks   = hmock.NewHooks[any]().RegisterBeforeSend(
			func(c context.Context, cm core.Cmd[any]) (context.Context, error) {
				return c, wantErr
			},
		)
		factory = hmock.NewHooksFactory[any]()
	)
	factory.RegisterNew(
		func() hks.Hooks[any] { return hooks },
	)
	return SenderTestCase[any]{
		Name: name,
		Setup: SenderSetup[any]{
			Group:   smock.NewSenderGroup[any](),
			Options: []sndr.SetOption[any]{sndr.WithHooksFactory(factory)},
		},
		Action: func(t *testing.T, s sndr.Sender[any]) {
			_, err := s.Send(context.Background(), core.Cmd[any](nil))
			asserterror.EqualDeep(t, err, wantErr)
		},
		Mocks: []*mok.Mock{factory.Mock, hooks.Mock},
	}
}

func SendGroupErrorTestCase() SenderTestCase[any] {
	name := "Send should return an error if SenderGroup.Send fails"

	var (
		wantErr = errors.New("Send error")
		group   = smock.NewSenderGroup[any]()
	)
	group.RegisterSend(
		func(c core.Cmd[any], r chan<- core.AsyncResult) (
			core.Seq, grp.ClientID, int, error,
		) {
			return 0, 0, 0, wantErr
		},
	)
	return SenderTestCase[any]{
		Name: name,
		Setup: SenderSetup[any]{
			Group: group,
		},
		Action: func(t *testing.T, s sndr.Sender[any]) {
			_, err := s.Send(context.Background(), core.Cmd[any](nil))
			asserterror.EqualDeep(t, err, wantErr)
		},
		Mocks: []*mok.Mock{group.Mock},
	}
}

func SendTimeoutTestCase() SenderTestCase[any] {
	name := "Send should return ErrTimeout if no result was received within the timeout"

	var (
		wantCmd = core.Cmd[any](nil)
		group   = smock.NewSenderGroup[any]()
	)
	group.RegisterSend(
		func(c core.Cmd[any], r chan<- core.AsyncResult) (
			core.Seq, grp.ClientID, int, error,
		) {
			// Don't send result to trigger timeout (after context cancel or similar)
			// Here we use Sender's receive internal behavior which uses context.
			return core.Seq(1), grp.ClientID(1), 10, nil
		},
	)
	return SenderTestCase[any]{
		Name: name,
		Setup: SenderSetup[any]{
			Group: group,
		},
		Action: func(t *testing.T, s sndr.Sender[any]) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
			defer cancel()
			_, err := s.Send(ctx, wantCmd)
			asserterror.EqualDeep(t, err, sndr.ErrTimeout)
		},
		Mocks: []*mok.Mock{group.Mock},
	}
}

func SendWithDeadlineSuccessTestCase(t *testing.T) SenderTestCase[any] {
	name := "SendWithDeadline should return no error if successful"

	var (
		wantDeadline = time.Now().Add(time.Hour)
		wantCmd      = core.Cmd[any](nil)
		seq          = core.Seq(1)
		wantResult   = core.Result(nil)
		group        = smock.NewSenderGroup[any]()
	)
	group.RegisterSendWithDeadline(
		func(d time.Time, c core.Cmd[any], r chan<- core.AsyncResult) (
			core.Seq, grp.ClientID, int, error,
		) {
			asserterror.Equal(t, d, wantDeadline)
			asserterror.EqualDeep(t, c, wantCmd)
			r <- core.AsyncResult{Seq: seq, Result: wantResult}
			return core.Seq(1), grp.ClientID(1), 10, nil
		},
	)
	return SenderTestCase[any]{
		Name: name,
		Setup: SenderSetup[any]{
			Group: group,
		},
		Action: func(t *testing.T, s sndr.Sender[any]) {
			result, err := s.SendWithDeadline(context.Background(), wantDeadline, wantCmd)
			asserterror.EqualDeep(t, result, wantResult)
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{group.Mock},
	}
}

func SendWithDeadlineBeforeSendErrorTestCase() SenderTestCase[any] {
	name := "SendWithDeadline should return an error if Hooks.BeforeSend fails"

	var (
		wantDeadline = time.Now().Add(time.Hour)
		wantErr      = errors.New("BeforeSend error")
		hooks        = hmock.NewHooks[any]().RegisterBeforeSend(
			func(c context.Context, cm core.Cmd[any]) (context.Context, error) {
				return c, wantErr
			},
		)
		factory = hmock.NewHooksFactory[any]().RegisterNew(
			func() hks.Hooks[any] { return hooks },
		)
		group = smock.NewSenderGroup[any]()
	)
	return SenderTestCase[any]{
		Name: name,
		Setup: SenderSetup[any]{
			Group:   group,
			Options: []sndr.SetOption[any]{sndr.WithHooksFactory(factory)},
		},
		Action: func(t *testing.T, s sndr.Sender[any]) {
			_, err := s.SendWithDeadline(context.Background(), wantDeadline, core.Cmd[any](nil))
			asserterror.EqualDeep(t, err, wantErr)
		},
		Mocks: []*mok.Mock{factory.Mock, hooks.Mock},
	}
}

func SendWithDeadlineGroupErrorTestCase() SenderTestCase[any] {
	name := "SendWithDeadline should return an error if SenderGroup.SendWithDeadline fails"

	var (
		wantDeadline = time.Now().Add(time.Hour)
		wantErr      = errors.New("SendWithDeadline error")
		group        = smock.NewSenderGroup[any]().RegisterSendWithDeadline(
			func(d time.Time, c core.Cmd[any], r chan<- core.AsyncResult) (
				core.Seq, grp.ClientID, int, error,
			) {
				return 0, 0, 0, wantErr
			},
		)
	)
	return SenderTestCase[any]{
		Name: name,
		Setup: SenderSetup[any]{
			Group: group,
		},
		Action: func(t *testing.T, s sndr.Sender[any]) {
			_, err := s.SendWithDeadline(context.Background(), wantDeadline, core.Cmd[any](nil))
			asserterror.EqualDeep(t, err, wantErr)
		},
		Mocks: []*mok.Mock{group.Mock},
	}
}

func SendWithDeadlineTimeoutTestCase() SenderTestCase[any] {
	name := "SendWithDeadline should return ErrTimeout if no result was received within the timeout"

	var (
		wantDeadline = time.Now().Add(time.Hour)
		seq          = core.Seq(1)
		group        = smock.NewSenderGroup[any]().RegisterSendWithDeadline(
			func(d time.Time, c core.Cmd[any], r chan<- core.AsyncResult) (
				core.Seq, grp.ClientID, int, error,
			) {
				return seq, grp.ClientID(1), 10, nil
			},
		)
	)
	return SenderTestCase[any]{
		Name: name,
		Setup: SenderSetup[any]{
			Group: group,
		},
		Action: func(t *testing.T, s sndr.Sender[any]) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
			defer cancel()
			_, err := s.SendWithDeadline(ctx, wantDeadline, core.Cmd[any](nil))
			asserterror.EqualDeep(t, err, sndr.ErrTimeout)
		},
		Mocks: []*mok.Mock{group.Mock},
	}
}

func SendMultiSuccessTestCase(t *testing.T) SenderTestCase[any] {
	name := "SendMulti should return no error if successful"

	var (
		wantCmd     = core.Cmd[any](nil)
		seq1        = core.Seq(1)
		seq2        = core.Seq(2)
		wantResult1 = cmock.NewResult().RegisterLastOne(func() bool { return false })
		wantResult2 = cmock.NewResult().RegisterLastOne(func() bool { return true })
		group       = smock.NewSenderGroup[any]()
		wantResults = []core.Result{wantResult1, wantResult2}
		count       = 0
		handler     = sndr.ResultHandlerFn(
			func(result core.Result, err error) error {
				asserterror.EqualDeep(t, result, wantResults[count])
				count++
				return nil
			},
		)
	)
	group.RegisterSend(
		func(c core.Cmd[any], r chan<- core.AsyncResult) (
			core.Seq, grp.ClientID, int, error,
		) {
			asserterror.EqualDeep(t, c, wantCmd)
			r <- core.AsyncResult{Seq: seq1, Result: wantResult1}
			r <- core.AsyncResult{Seq: seq2, Result: wantResult2}
			return seq1, grp.ClientID(1), 10, nil
		},
	)
	return SenderTestCase[any]{
		Name: name,
		Setup: SenderSetup[any]{
			Group: group,
		},
		Action: func(t *testing.T, s sndr.Sender[any]) {
			err := s.SendMulti(context.Background(), wantCmd, 2, handler)
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{group.Mock, wantResult1.Mock, wantResult2.Mock},
	}
}

func SendMultiBeforeSendErrorTestCase() SenderTestCase[any] {
	name := "SendMulti should return an error if Hooks.BeforeSend fails"

	var (
		wantErr = errors.New("BeforeSend error")
		hooks   = hmock.NewHooks[any]()
		factory = hmock.NewHooksFactory[any]()
	)
	hooks.RegisterBeforeSend(
		func(c context.Context, cm core.Cmd[any]) (context.Context, error) {
			return c, wantErr
		},
	)
	factory.RegisterNew(
		func() hks.Hooks[any] { return hooks },
	)
	return SenderTestCase[any]{
		Name: name,
		Setup: SenderSetup[any]{
			Group:   smock.NewSenderGroup[any](),
			Options: []sndr.SetOption[any]{sndr.WithHooksFactory(factory)},
		},
		Action: func(t *testing.T, s sndr.Sender[any]) {
			err := s.SendMulti(context.Background(), core.Cmd[any](nil), 1, nil)
			asserterror.EqualDeep(t, err, wantErr)
		},
		Mocks: []*mok.Mock{factory.Mock, hooks.Mock},
	}
}

func SendMultiGroupErrorTestCase() SenderTestCase[any] {
	name := "SendMulti should return an error if SenderGroup.Send fails"

	var (
		wantErr = errors.New("Send error")
		group   = smock.NewSenderGroup[any]()
	)
	group.RegisterSend(
		func(c core.Cmd[any], r chan<- core.AsyncResult) (
			core.Seq, grp.ClientID, int, error,
		) {
			return 0, 0, 0, wantErr
		},
	)
	return SenderTestCase[any]{
		Name: name,
		Setup: SenderSetup[any]{
			Group: group,
		},
		Action: func(t *testing.T, s sndr.Sender[any]) {
			err := s.SendMulti(context.Background(), core.Cmd[any](nil), 1, nil)
			asserterror.EqualDeep(t, err, wantErr)
		},
		Mocks: []*mok.Mock{group.Mock},
	}
}

func SendMultiTimeoutTestCase(t *testing.T) SenderTestCase[any] {
	name := "SendMulti should report ErrTimeout to the ResultHandler if no result was received within the timeout"

	var (
		wantCmd = core.Cmd[any](nil)
		seq     = core.Seq(1)
		group   = smock.NewSenderGroup[any]()
		handler = sndr.ResultHandlerFn(
			func(result core.Result, err error) error {
				asserterror.EqualDeep(t, err, sndr.ErrTimeout)
				return nil
			},
		)
	)
	group.RegisterSend(
		func(c core.Cmd[any], r chan<- core.AsyncResult) (
			core.Seq, grp.ClientID, int, error,
		) {
			return seq, grp.ClientID(1), 10, nil
		},
	)
	return SenderTestCase[any]{
		Name: name,
		Setup: SenderSetup[any]{
			Group: group,
		},
		Action: func(t *testing.T, s sndr.Sender[any]) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
			defer cancel()
			err := s.SendMulti(ctx, wantCmd, 1, handler)
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{group.Mock},
	}
}

func SendMultiWithDeadlineSuccessTestCase(t *testing.T) SenderTestCase[any] {
	name := "SendMultiWithDeadline should return no error if successful"

	var (
		deadline    = time.Now().Add(time.Hour)
		wantCmd     = core.Cmd[any](nil)
		seq1        = core.Seq(1)
		seq2        = core.Seq(2)
		wantResult1 = cmock.NewResult().RegisterLastOne(func() bool { return false })
		wantResult2 = cmock.NewResult().RegisterLastOne(func() bool { return true })
		group       = smock.NewSenderGroup[any]()
		wantResults = []core.Result{wantResult1, wantResult2}
		count       = 0
		handler     = sndr.ResultHandlerFn(
			func(result core.Result, err error) error {
				asserterror.EqualDeep(t, result, wantResults[count])
				count++
				return nil
			},
		)
	)
	group.RegisterSendWithDeadline(
		func(d time.Time, c core.Cmd[any], r chan<- core.AsyncResult) (
			core.Seq, grp.ClientID, int, error,
		) {
			asserterror.Equal(t, d, deadline)
			asserterror.EqualDeep(t, c, wantCmd)
			r <- core.AsyncResult{Seq: seq1, Result: wantResult1}
			r <- core.AsyncResult{Seq: seq2, Result: wantResult2}
			return seq1, grp.ClientID(1), 10, nil
		},
	)
	return SenderTestCase[any]{
		Name: name,
		Setup: SenderSetup[any]{
			Group: group,
		},
		Action: func(t *testing.T, s sndr.Sender[any]) {
			err := s.SendMultiWithDeadline(context.Background(), deadline, wantCmd, 2,
				handler)
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{group.Mock, wantResult1.Mock, wantResult2.Mock},
	}
}

func SendMultiWithDeadlineBeforeSendErrorTestCase() SenderTestCase[any] {
	name := "SendMultiWithDeadline should return an error if Hooks.BeforeSend fails"

	var (
		deadline = time.Now().Add(time.Hour)
		wantErr  = errors.New("BeforeSend error")
		hooks    = hmock.NewHooks[any]().RegisterBeforeSend(
			func(c context.Context, cm core.Cmd[any]) (context.Context, error) {
				return c, wantErr
			},
		)
		factory = hmock.NewHooksFactory[any]()
	)
	factory.RegisterNew(
		func() hks.Hooks[any] { return hooks },
	)
	return SenderTestCase[any]{
		Name: name,
		Setup: SenderSetup[any]{
			Group:   smock.NewSenderGroup[any](),
			Options: []sndr.SetOption[any]{sndr.WithHooksFactory(factory)},
		},
		Action: func(t *testing.T, s sndr.Sender[any]) {
			err := s.SendMultiWithDeadline(context.Background(), deadline, core.Cmd[any](nil), 1, nil)
			asserterror.EqualDeep(t, err, wantErr)
		},
		Mocks: []*mok.Mock{factory.Mock, hooks.Mock},
	}
}

func SendMultiWithDeadlineGroupErrorTestCase() SenderTestCase[any] {
	name := "SendMultiWithDeadline should return an error if SenderGroup.SendWithDeadline fails"

	var (
		wantErr = errors.New("SendWithDeadline error")
		group   = smock.NewSenderGroup[any]()
	)
	group.RegisterSendWithDeadline(
		func(d time.Time, c core.Cmd[any], r chan<- core.AsyncResult) (
			core.Seq, grp.ClientID, int, error,
		) {
			return 0, 0, 0, wantErr
		},
	)
	return SenderTestCase[any]{
		Name: name,
		Setup: SenderSetup[any]{
			Group: group,
		},
		Action: func(t *testing.T, s sndr.Sender[any]) {
			err := s.SendMultiWithDeadline(context.Background(), time.Now().Add(time.Hour),
				core.Cmd[any](nil), 1, nil)
			asserterror.EqualDeep(t, err, wantErr)
		},
		Mocks: []*mok.Mock{group.Mock},
	}
}

func SendMultiWithDeadlineTimeoutTestCase(t *testing.T) SenderTestCase[any] {
	name := "SendMultiWithDeadline should report ErrTimeout to the ResultHandler if no result was received within the timeout"

	var (
		group   = smock.NewSenderGroup[any]()
		handler = sndr.ResultHandlerFn(
			func(result core.Result, err error) error {
				asserterror.EqualDeep(t, err, sndr.ErrTimeout)
				return nil
			},
		)
	)
	group.RegisterSendWithDeadline(
		func(d time.Time, c core.Cmd[any], r chan<- core.AsyncResult) (
			core.Seq, grp.ClientID, int, error,
		) {
			return core.Seq(1), grp.ClientID(1), 10, nil
		},
	)
	return SenderTestCase[any]{
		Name: name,
		Setup: SenderSetup[any]{
			Group: group,
		},
		Action: func(t *testing.T, s sndr.Sender[any]) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
			defer cancel()
			err := s.SendMultiWithDeadline(ctx, time.Now().Add(time.Hour),
				core.Cmd[any](nil), 1, handler)
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{group.Mock},
	}
}
