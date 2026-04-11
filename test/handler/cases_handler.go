package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/cmd-stream/cmd-stream-go/handler"
	"github.com/cmd-stream/cmd-stream-go/test"
	cmock "github.com/cmd-stream/cmd-stream-go/test/mock/core"
	dmock "github.com/cmd-stream/cmd-stream-go/test/mock/delegate"
	hmock "github.com/cmd-stream/cmd-stream-go/test/mock/handler"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

func HandleSuccessTestCase(t *testing.T) HandlerTestCase[any] {
	name := "Handler should be able to handle several cmds and close when ctx done"

	var (
		ctx, cancel = context.WithCancel(context.Background())
		transport   = dmock.NewServerTransport[any]()
		invoker     = hmock.NewInvoker[any]()

		seq1 = core.Seq(1)
		cmd1 = cmock.NewCmd[any]()
		n1   = 10

		seq2 = core.Seq(2)
		cmd2 = cmock.NewCmd[any]()
		n2   = 20

		done = make(chan struct{})
	)
	transport.RegisterReceive(
		func() (core.Seq, core.Cmd[any], int, error) {
			return seq1, cmd1, n1, nil
		},
	).RegisterReceive(
		func() (core.Seq, core.Cmd[any], int, error) {
			return seq2, cmd2, n2, nil
		},
	).RegisterReceive(
		func() (core.Seq, core.Cmd[any], int, error) {
			<-done
			return 0, nil, 0, context.Canceled
		},
	)
	invoker.RegisterInvokeN(2,
		func(c context.Context, seq core.Seq, at time.Time, bytesRead int,
			cmd core.Cmd[any], proxy core.Proxy) error {
			switch seq {
			case 1:
				asserterror.EqualDeep(t, cmd, cmd1)
				asserterror.Equal(t, bytesRead, n1)
			case 2:
				asserterror.EqualDeep(t, cmd, cmd2)
				asserterror.Equal(t, bytesRead, n2)
			default:
				t.Errorf("unexpected seq: %d", seq)
			}
			return nil
		},
	)
	transport.RegisterClose(func() error {
		close(done)
		return nil
	})
	return HandlerTestCase[any]{
		Name: name,
		Setup: HandlerSetup[any]{
			Invoker: invoker,
			Opts:    []handler.SetOption{},
		},
		Action: func(t *testing.T, h *handler.Handler[any]) {
			go func() {
				time.Sleep(100 * time.Millisecond)
				cancel()
			}()
			err := h.Handle(ctx, transport)
			asserterror.EqualError(t, err, context.Canceled)
		},
		Mocks: []*mok.Mock{transport.Mock, invoker.Mock, cmd1.Mock, cmd2.Mock},
	}
}

func SetReceiveDeadlineErrorTestCase(t *testing.T) HandlerTestCase[any] {
	name := "If Transport.SetReceiveDeadline fails with an error, Handle should return it"

	var (
		wantErr   = errors.New("Transport.SetReceiveDeadline error")
		transport = dmock.NewServerTransport[any]()
		invoker   = hmock.NewInvoker[any]()
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return wantErr },
	).RegisterClose(
		func() error { return nil },
	)
	return HandlerTestCase[any]{
		Name: name,
		Setup: HandlerSetup[any]{
			Invoker: invoker,
			Opts:    []handler.SetOption{handler.WithCmdReceiveDuration(time.Second)},
		},
		Action: func(t *testing.T, h *handler.Handler[any]) {
			err := h.Handle(context.Background(), transport)
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{transport.Mock, invoker.Mock},
	}
}

func ReceiveErrorTestCase(t *testing.T) HandlerTestCase[any] {
	name := "If Transport.Receive fails with an error, Handle should return it"

	var (
		wantErr   = errors.New("Transport.Receive error")
		transport = dmock.NewServerTransport[any]()
		invoker   = hmock.NewInvoker[any]()
	)
	transport.RegisterReceive(
		func() (core.Seq, core.Cmd[any], int, error) { return 0, nil, 2, wantErr },
	).RegisterClose(
		func() error { return nil },
	)
	return HandlerTestCase[any]{
		Name: name,
		Setup: HandlerSetup[any]{
			Invoker: invoker,
			Opts:    []handler.SetOption{},
		},
		Action: func(t *testing.T, h *handler.Handler[any]) {
			err := h.Handle(context.Background(), transport)
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{transport.Mock, invoker.Mock},
	}
}

func InvokeErrorTestCase(t *testing.T) HandlerTestCase[any] {
	name := "If Invoker.Invoke fails with an error, Handle should return it"

	var (
		wantErr   = errors.New("Invoker.Invoke error")
		transport = dmock.NewServerTransport[any]()
		invoker   = hmock.NewInvoker[any]()
		seq       = core.Seq(1)
		cmd       = cmock.NewCmd[any]()
		n         = 10
		done      = make(chan struct{})
	)
	transport.RegisterReceive(
		func() (core.Seq, core.Cmd[any], int, error) { return seq, cmd, n, nil },
	).RegisterReceive(
		func() (core.Seq, core.Cmd[any], int, error) {
			<-done
			return 0, nil, 0, context.Canceled
		},
	).RegisterClose(
		func() error {
			close(done)
			return nil
		},
	)
	invoker.RegisterInvoke(
		func(ctx context.Context, seq core.Seq, at time.Time, bytesRead int,
			cmd core.Cmd[any], proxy core.Proxy) error {
			return wantErr
		},
	)
	return HandlerTestCase[any]{
		Name: name,
		Setup: HandlerSetup[any]{
			Invoker: invoker,
			Opts:    []handler.SetOption{},
		},
		Action: func(t *testing.T, h *handler.Handler[any]) {
			err := h.Handle(context.Background(), transport)
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{transport.Mock, invoker.Mock, cmd.Mock},
	}
}

func OptionAtTestCase(t *testing.T) HandlerTestCase[any] {
	name := "If Conf.At == true, Invoker.Invoke should receive not empty 'at' param"

	var (
		ctx, cancel = context.WithCancel(context.Background())
		transport   = dmock.NewServerTransport[any]()
		invoker     = hmock.NewInvoker[any]()
		done        = make(chan struct{})
		delta       = 100 * time.Millisecond
	)
	transport.RegisterReceive(
		func() (core.Seq, core.Cmd[any], int, error) { return core.Seq(1), cmock.NewCmd[any](), 10, nil },
	).RegisterReceive(
		func() (core.Seq, core.Cmd[any], int, error) {
			<-done
			return 0, nil, 0, context.Canceled
		},
	).RegisterClose(
		func() error {
			close(done)
			return nil
		},
	)
	invoker.RegisterInvoke(
		func(ctx context.Context, seq core.Seq, at time.Time, bytesRead int,
			cmd core.Cmd[any], proxy core.Proxy) error {
			asserterror.SameTime(t, at, time.Now(), delta)
			return nil
		},
	)

	return HandlerTestCase[any]{
		Name: name,
		Setup: HandlerSetup[any]{
			Invoker: invoker,
			Opts:    []handler.SetOption{handler.WithAt()},
		},
		Action: func(t *testing.T, h *handler.Handler[any]) {
			go func() {
				time.Sleep(100 * time.Millisecond)
				cancel()
			}()
			err := h.Handle(ctx, transport)
			asserterror.EqualError(t, err, context.Canceled)
		},
		Mocks: []*mok.Mock{transport.Mock, invoker.Mock},
	}
}

func OptionCmdReceiveDurationTestCase(t *testing.T) HandlerTestCase[any] {
	name := "If Conf.CmdReceiveDuration is set, Transport.SetReceiveDeadline should receive it"

	var (
		ctx, cancel            = context.WithCancel(context.Background())
		wantCmdReceiveDuration = time.Second
		transport              = dmock.NewServerTransport[any]()
		invoker                = hmock.NewInvoker[any]()
		done                   = make(chan struct{})
		startTime              = time.Now()
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) (err error) {
			wantDeadline := startTime.Add(wantCmdReceiveDuration)
			asserterror.SameTime(t, deadline, wantDeadline, test.TimeDelta)
			return
		},
	).RegisterReceive(
		func() (core.Seq, core.Cmd[any], int, error) {
			<-done
			return 0, nil, 0, context.Canceled
		},
	).RegisterClose(
		func() error {
			close(done)
			return nil
		},
	)
	return HandlerTestCase[any]{
		Name: name,
		Setup: HandlerSetup[any]{
			Invoker: invoker,
			Opts:    []handler.SetOption{handler.WithCmdReceiveDuration(wantCmdReceiveDuration)},
		},
		Action: func(t *testing.T, h *handler.Handler[any]) {
			go func() {
				time.Sleep(100 * time.Millisecond)
				cancel()
			}()
			err := h.Handle(ctx, transport)
			asserterror.EqualError(t, err, context.Canceled)
		},
		Mocks: []*mok.Mock{transport.Mock, invoker.Mock},
	}
}

func CloseWhileInvokingCmdsTestCase(t *testing.T) HandlerTestCase[any] {
	name := "We should be able to interrupt Handler, while it invokes several Commands"

	var (
		ctx, cancel = context.WithCancel(context.Background())
		transport   = dmock.NewServerTransport[any]()
		invoker     = hmock.NewInvoker[any]()
		done        = make(chan struct{})
	)
	transport.RegisterReceive(
		func() (core.Seq, core.Cmd[any], int, error) {
			return core.Seq(1), cmock.NewCmd[any](), 10, nil
		},
	).RegisterReceive(
		func() (core.Seq, core.Cmd[any], int, error) {
			return core.Seq(2), cmock.NewCmd[any](), 10, nil
		},
	).RegisterReceive(
		func() (core.Seq, core.Cmd[any], int, error) {
			<-done
			return 0, nil, 0, context.Canceled
		},
	).RegisterClose(
		func() error {
			close(done)
			return nil
		},
	)
	invoker.RegisterInvokeN(2,
		func(ctx context.Context, seq core.Seq, at time.Time, bytesRead int,
			cmd core.Cmd[any], proxy core.Proxy) error {
			<-ctx.Done()
			return nil
		},
	)
	return HandlerTestCase[any]{
		Name: name,
		Setup: HandlerSetup[any]{
			Invoker: invoker,
			Opts:    []handler.SetOption{},
		},
		Action: func(t *testing.T, h *handler.Handler[any]) {
			go func() {
				time.Sleep(100 * time.Millisecond)
				cancel()
			}()
			err := h.Handle(ctx, transport)
			asserterror.EqualError(t, err, context.Canceled)
		},
		Mocks: []*mok.Mock{transport.Mock, invoker.Mock},
	}
}
