package test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	grp "github.com/cmd-stream/cmd-stream-go/group"
	"github.com/cmd-stream/cmd-stream-go/sender"
	hks "github.com/cmd-stream/cmd-stream-go/sender/hooks"
	cmock "github.com/cmd-stream/cmd-stream-go/test/mock/core"
	smock "github.com/cmd-stream/cmd-stream-go/test/mock/sender"
	hmock "github.com/cmd-stream/cmd-stream-go/test/mock/sender/hooks"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

type SenderSetup[T any] struct {
	Group   smock.Group[T]
	Options []sender.SetOption[T]
}

type SenderTestCase[T any] struct {
	Name   string
	Setup  SenderSetup[T]
	Action func(t *testing.T, s sender.Sender[T])
	Mocks  []*mok.Mock
}

func RunSenderTestCase[T any](t *testing.T, tc SenderTestCase[T]) {
	t.Run(tc.Name, func(t *testing.T) {
		s := sender.New[T](tc.Setup.Group, tc.Setup.Options...)
		tc.Action(t, s)
		asserterror.EqualDeep(t, mok.CheckCalls(tc.Mocks), mok.EmptyInfomap)
	})
}

type SenderSuite[T any] struct{}

// -----------------------------------------------------------------------------
// Test Cases
// -----------------------------------------------------------------------------

func (SenderSuite[T]) SendSuccess(t *testing.T) SenderTestCase[T] {
	name := "Send should return no error if successful"

	var (
		wantCmd    = core.Cmd[T](nil)
		seq        = core.Seq(1)
		wantResult = core.Result(nil)
		group      = smock.NewGroup[T]()
	)
	group.RegisterSend(
		func(c core.Cmd[T], r chan<- core.AsyncResult) (
			core.Seq, grp.ClientID, int, error,
		) {
			asserterror.EqualDeep(t, c, wantCmd)
			r <- core.AsyncResult{Seq: seq, Result: wantResult}
			return seq, grp.ClientID(1), 10, nil
		},
	)
	return SenderTestCase[T]{
		Name: name,
		Setup: SenderSetup[T]{
			Group: group,
		},
		Action: func(t *testing.T, s sender.Sender[T]) {
			result, err := s.Send(context.Background(), wantCmd)
			asserterror.EqualDeep(t, result, wantResult)
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{group.Mock},
	}
}

func (SenderSuite[T]) SendBeforeSendError(t *testing.T) SenderTestCase[T] {
	name := "Send should return an error if Hooks.BeforeSend fails"

	var (
		wantErr = errors.New("BeforeSend error")
		hooks   = hmock.NewHooks[T]().RegisterBeforeSend(
			func(c context.Context, cm core.Cmd[T]) (context.Context, error) {
				return c, wantErr
			},
		)
		factory = hmock.NewFactory[T]()
	)
	factory.RegisterNew(
		func() hks.Hooks[T] { return hooks },
	)
	return SenderTestCase[T]{
		Name: name,
		Setup: SenderSetup[T]{
			Group:   smock.NewGroup[T](),
			Options: []sender.SetOption[T]{sender.WithHooksFactory(factory)},
		},
		Action: func(t *testing.T, s sender.Sender[T]) {
			_, err := s.Send(context.Background(), core.Cmd[T](nil))
			asserterror.EqualDeep(t, err, wantErr)
		},
		Mocks: []*mok.Mock{factory.Mock, hooks.Mock},
	}
}

func (SenderSuite[T]) SendGroupError(t *testing.T) SenderTestCase[T] {
	name := "Send should return an error if Group.Send fails"

	var (
		wantErr = errors.New("send error")
		group   = smock.NewGroup[T]()
	)
	group.RegisterSend(
		func(c core.Cmd[T], r chan<- core.AsyncResult) (
			core.Seq, grp.ClientID, int, error,
		) {
			return 0, 0, 0, wantErr
		},
	)
	return SenderTestCase[T]{
		Name: name,
		Setup: SenderSetup[T]{
			Group: group,
		},
		Action: func(t *testing.T, s sender.Sender[T]) {
			_, err := s.Send(context.Background(), core.Cmd[T](nil))
			asserterror.EqualDeep(t, err, wantErr)
		},
		Mocks: []*mok.Mock{group.Mock},
	}
}

func (SenderSuite[T]) SendTimeout(t *testing.T) SenderTestCase[T] {
	name := "Send should return ErrTimeout if no result was received within the timeout"

	var (
		wantCmd = core.Cmd[T](nil)
		group   = smock.NewGroup[T]()
	)
	group.RegisterSend(
		func(c core.Cmd[T], r chan<- core.AsyncResult) (
			core.Seq, grp.ClientID, int, error,
		) {
			// Don't send result to trigger timeout (after context cancel or similar)
			// Here we use Sender's receive internal behavior which uses context.
			return core.Seq(1), grp.ClientID(1), 10, nil
		},
	)
	return SenderTestCase[T]{
		Name: name,
		Setup: SenderSetup[T]{
			Group: group,
		},
		Action: func(t *testing.T, s sender.Sender[T]) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
			defer cancel()
			_, err := s.Send(ctx, wantCmd)
			asserterror.EqualDeep(t, err, sender.ErrTimeout)
		},
		Mocks: []*mok.Mock{group.Mock},
	}
}

func (SenderSuite[T]) SendWithDeadlineSuccess(t *testing.T) SenderTestCase[T] {
	name := "SendWithDeadline should return no error if successful"

	var (
		wantDeadline = time.Now().Add(time.Hour)
		wantCmd      = core.Cmd[T](nil)
		seq          = core.Seq(1)
		wantResult   = core.Result(nil)
		group        = smock.NewGroup[T]()
	)
	group.RegisterSendWithDeadline(
		func(d time.Time, c core.Cmd[T], r chan<- core.AsyncResult) (
			core.Seq, grp.ClientID, int, error,
		) {
			asserterror.Equal(t, d, wantDeadline)
			asserterror.EqualDeep(t, c, wantCmd)
			r <- core.AsyncResult{Seq: seq, Result: wantResult}
			return core.Seq(1), grp.ClientID(1), 10, nil
		},
	)
	return SenderTestCase[T]{
		Name: name,
		Setup: SenderSetup[T]{
			Group: group,
		},
		Action: func(t *testing.T, s sender.Sender[T]) {
			result, err := s.SendWithDeadline(context.Background(), wantDeadline, wantCmd)
			asserterror.EqualDeep(t, result, wantResult)
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{group.Mock},
	}
}

func (SenderSuite[T]) SendWithDeadlineBeforeSendError(t *testing.T) SenderTestCase[T] {
	name := "SendWithDeadline should return an error if Hooks.BeforeSend fails"

	var (
		wantDeadline = time.Now().Add(time.Hour)
		wantErr      = errors.New("BeforeSend error")
		hooks        = hmock.NewHooks[T]().RegisterBeforeSend(
			func(c context.Context, cm core.Cmd[T]) (context.Context, error) {
				return c, wantErr
			},
		)
		factory = hmock.NewFactory[T]().RegisterNew(
			func() hks.Hooks[T] { return hooks },
		)
		group = smock.NewGroup[T]()
	)
	return SenderTestCase[T]{
		Name: name,
		Setup: SenderSetup[T]{
			Group:   group,
			Options: []sender.SetOption[T]{sender.WithHooksFactory(factory)},
		},
		Action: func(t *testing.T, s sender.Sender[T]) {
			_, err := s.SendWithDeadline(context.Background(), wantDeadline, core.Cmd[T](nil))
			asserterror.EqualDeep(t, err, wantErr)
		},
		Mocks: []*mok.Mock{factory.Mock, hooks.Mock},
	}
}

func (SenderSuite[T]) SendWithDeadlineGroupError(t *testing.T) SenderTestCase[T] {
	name := "SendWithDeadline should return an error if Group.SendWithDeadline fails"

	var (
		wantDeadline = time.Now().Add(time.Hour)
		wantErr      = errors.New("SendWithDeadline error")
		group        = smock.NewGroup[T]().RegisterSendWithDeadline(
			func(d time.Time, c core.Cmd[T], r chan<- core.AsyncResult) (
				core.Seq, grp.ClientID, int, error,
			) {
				return 0, 0, 0, wantErr
			},
		)
	)
	return SenderTestCase[T]{
		Name: name,
		Setup: SenderSetup[T]{
			Group: group,
		},
		Action: func(t *testing.T, s sender.Sender[T]) {
			_, err := s.SendWithDeadline(context.Background(), wantDeadline, core.Cmd[T](nil))
			asserterror.EqualDeep(t, err, wantErr)
		},
		Mocks: []*mok.Mock{group.Mock},
	}
}

func (SenderSuite[T]) SendWithDeadlineTimeout(t *testing.T) SenderTestCase[T] {
	name := "SendWithDeadline should return ErrTimeout if no result was received within the timeout"

	var (
		wantDeadline = time.Now().Add(time.Hour)
		seq          = core.Seq(1)
		group        = smock.NewGroup[T]().RegisterSendWithDeadline(
			func(d time.Time, c core.Cmd[T], r chan<- core.AsyncResult) (
				core.Seq, grp.ClientID, int, error,
			) {
				return seq, grp.ClientID(1), 10, nil
			},
		)
	)
	return SenderTestCase[T]{
		Name: name,
		Setup: SenderSetup[T]{
			Group: group,
		},
		Action: func(t *testing.T, s sender.Sender[T]) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
			defer cancel()
			_, err := s.SendWithDeadline(ctx, wantDeadline, core.Cmd[T](nil))
			asserterror.EqualDeep(t, err, sender.ErrTimeout)
		},
		Mocks: []*mok.Mock{group.Mock},
	}
}

func (SenderSuite[T]) SendMultiSuccess(t *testing.T) SenderTestCase[T] {
	name := "SendMulti should return no error if successful"

	var (
		wantCmd     = core.Cmd[T](nil)
		seq1        = core.Seq(1)
		seq2        = core.Seq(2)
		wantResult1 = cmock.NewResult().RegisterLastOne(func() bool { return false })
		wantResult2 = cmock.NewResult().RegisterLastOne(func() bool { return true })
		group       = smock.NewGroup[T]()
		wantResults = []core.Result{wantResult1, wantResult2}
		count       = 0
		handler     = sender.ResultHandlerFn(
			func(result core.Result, err error) error {
				asserterror.EqualDeep(t, result, wantResults[count])
				count++
				return nil
			},
		)
	)
	group.RegisterSend(
		func(c core.Cmd[T], r chan<- core.AsyncResult) (
			core.Seq, grp.ClientID, int, error,
		) {
			asserterror.EqualDeep(t, c, wantCmd)
			r <- core.AsyncResult{Seq: seq1, Result: wantResult1}
			r <- core.AsyncResult{Seq: seq2, Result: wantResult2}
			return seq1, grp.ClientID(1), 10, nil
		},
	)
	return SenderTestCase[T]{
		Name: name,
		Setup: SenderSetup[T]{
			Group: group,
		},
		Action: func(t *testing.T, s sender.Sender[T]) {
			err := s.SendMulti(context.Background(), wantCmd, 2, handler)
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{group.Mock, wantResult1.Mock, wantResult2.Mock},
	}
}

func (SenderSuite[T]) SendMultiBeforeSendError(t *testing.T) SenderTestCase[T] {
	name := "SendMulti should return an error if Hooks.BeforeSend fails"

	var (
		wantErr = errors.New("BeforeSend error")
		hooks   = hmock.NewHooks[T]()
		factory = hmock.NewFactory[T]()
	)
	hooks.RegisterBeforeSend(
		func(c context.Context, cm core.Cmd[T]) (context.Context, error) {
			return c, wantErr
		},
	)
	factory.RegisterNew(
		func() hks.Hooks[T] { return hooks },
	)
	return SenderTestCase[T]{
		Name: name,
		Setup: SenderSetup[T]{
			Group:   smock.NewGroup[T](),
			Options: []sender.SetOption[T]{sender.WithHooksFactory(factory)},
		},
		Action: func(t *testing.T, s sender.Sender[T]) {
			err := s.SendMulti(context.Background(), core.Cmd[T](nil), 1, nil)
			asserterror.EqualDeep(t, err, wantErr)
		},
		Mocks: []*mok.Mock{factory.Mock, hooks.Mock},
	}
}

func (SenderSuite[T]) SendMultiGroupError(t *testing.T) SenderTestCase[T] {
	name := "SendMulti should return an error if Group.Send fails"

	var (
		wantErr = errors.New("send error")
		group   = smock.NewGroup[T]()
	)
	group.RegisterSend(
		func(c core.Cmd[T], r chan<- core.AsyncResult) (
			core.Seq, grp.ClientID, int, error,
		) {
			return 0, 0, 0, wantErr
		},
	)
	return SenderTestCase[T]{
		Name: name,
		Setup: SenderSetup[T]{
			Group: group,
		},
		Action: func(t *testing.T, s sender.Sender[T]) {
			err := s.SendMulti(context.Background(), core.Cmd[T](nil), 1, nil)
			asserterror.EqualDeep(t, err, wantErr)
		},
		Mocks: []*mok.Mock{group.Mock},
	}
}

func (SenderSuite[T]) SendMultiTimeout(t *testing.T) SenderTestCase[T] {
	name := "SendMulti should report ErrTimeout to the ResultHandler if no result was received within the timeout"

	var (
		wantCmd = core.Cmd[T](nil)
		seq     = core.Seq(1)
		group   = smock.NewGroup[T]()
		handler = sender.ResultHandlerFn(
			func(result core.Result, err error) error {
				asserterror.EqualDeep(t, err, sender.ErrTimeout)
				return nil
			},
		)
	)
	group.RegisterSend(
		func(c core.Cmd[T], r chan<- core.AsyncResult) (
			core.Seq, grp.ClientID, int, error,
		) {
			return seq, grp.ClientID(1), 10, nil
		},
	)
	return SenderTestCase[T]{
		Name: name,
		Setup: SenderSetup[T]{
			Group: group,
		},
		Action: func(t *testing.T, s sender.Sender[T]) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
			defer cancel()
			err := s.SendMulti(ctx, wantCmd, 1, handler)
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{group.Mock},
	}
}

func (SenderSuite[T]) SendMultiWithDeadlineSuccess(t *testing.T) SenderTestCase[T] {
	name := "SendMultiWithDeadline should return no error if successful"

	var (
		deadline    = time.Now().Add(time.Hour)
		wantCmd     = core.Cmd[T](nil)
		seq1        = core.Seq(1)
		seq2        = core.Seq(2)
		wantResult1 = cmock.NewResult().RegisterLastOne(func() bool { return false })
		wantResult2 = cmock.NewResult().RegisterLastOne(func() bool { return true })
		group       = smock.NewGroup[T]()
		wantResults = []core.Result{wantResult1, wantResult2}
		count       = 0
		handler     = sender.ResultHandlerFn(
			func(result core.Result, err error) error {
				asserterror.EqualDeep(t, result, wantResults[count])
				count++
				return nil
			},
		)
	)
	group.RegisterSendWithDeadline(
		func(d time.Time, c core.Cmd[T], r chan<- core.AsyncResult) (
			core.Seq, grp.ClientID, int, error,
		) {
			asserterror.Equal(t, d, deadline)
			asserterror.EqualDeep(t, c, wantCmd)
			r <- core.AsyncResult{Seq: seq1, Result: wantResult1}
			r <- core.AsyncResult{Seq: seq2, Result: wantResult2}
			return seq1, grp.ClientID(1), 10, nil
		},
	)
	return SenderTestCase[T]{
		Name: name,
		Setup: SenderSetup[T]{
			Group: group,
		},
		Action: func(t *testing.T, s sender.Sender[T]) {
			err := s.SendMultiWithDeadline(context.Background(), deadline, wantCmd, 2,
				handler)
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{group.Mock, wantResult1.Mock, wantResult2.Mock},
	}
}

func (SenderSuite[T]) SendMultiWithDeadlineBeforeSendError(t *testing.T) SenderTestCase[T] {
	name := "SendMultiWithDeadline should return an error if Hooks.BeforeSend fails"

	var (
		deadline = time.Now().Add(time.Hour)
		wantErr  = errors.New("BeforeSend error")
		hooks    = hmock.NewHooks[T]().RegisterBeforeSend(
			func(c context.Context, cm core.Cmd[T]) (context.Context, error) {
				return c, wantErr
			},
		)
		factory = hmock.NewFactory[T]()
	)
	factory.RegisterNew(
		func() hks.Hooks[T] { return hooks },
	)
	return SenderTestCase[T]{
		Name: name,
		Setup: SenderSetup[T]{
			Group:   smock.NewGroup[T](),
			Options: []sender.SetOption[T]{sender.WithHooksFactory(factory)},
		},
		Action: func(t *testing.T, s sender.Sender[T]) {
			err := s.SendMultiWithDeadline(context.Background(), deadline, core.Cmd[T](nil), 1, nil)
			asserterror.EqualDeep(t, err, wantErr)
		},
		Mocks: []*mok.Mock{factory.Mock, hooks.Mock},
	}
}

func (SenderSuite[T]) SendMultiWithDeadlineGroupError(t *testing.T) SenderTestCase[T] {
	name := "SendMultiWithDeadline should return an error if Group.SendWithDeadline fails"

	var (
		wantErr = errors.New("SendWithDeadline error")
		group   = smock.NewGroup[T]()
	)
	group.RegisterSendWithDeadline(
		func(d time.Time, c core.Cmd[T], r chan<- core.AsyncResult) (
			core.Seq, grp.ClientID, int, error,
		) {
			return 0, 0, 0, wantErr
		},
	)
	return SenderTestCase[T]{
		Name: name,
		Setup: SenderSetup[T]{
			Group: group,
		},
		Action: func(t *testing.T, s sender.Sender[T]) {
			err := s.SendMultiWithDeadline(context.Background(), time.Now().Add(time.Hour),
				core.Cmd[T](nil), 1, nil)
			asserterror.EqualDeep(t, err, wantErr)
		},
		Mocks: []*mok.Mock{group.Mock},
	}
}

func (SenderSuite[T]) SendMultiWithDeadlineTimeout(t *testing.T) SenderTestCase[T] {
	name := "SendMultiWithDeadline should report ErrTimeout to the ResultHandler if no result was received within the timeout"

	var (
		group   = smock.NewGroup[T]()
		handler = sender.ResultHandlerFn(
			func(result core.Result, err error) error {
				asserterror.EqualDeep(t, err, sender.ErrTimeout)
				return nil
			},
		)
	)
	group.RegisterSendWithDeadline(
		func(d time.Time, c core.Cmd[T], r chan<- core.AsyncResult) (
			core.Seq, grp.ClientID, int, error,
		) {
			return core.Seq(1), grp.ClientID(1), 10, nil
		},
	)
	return SenderTestCase[T]{
		Name: name,
		Setup: SenderSetup[T]{
			Group: group,
		},
		Action: func(t *testing.T, s sender.Sender[T]) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
			defer cancel()
			err := s.SendMultiWithDeadline(ctx, time.Now().Add(time.Hour),
				core.Cmd[T](nil), 1, handler)
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{group.Mock},
	}
}
