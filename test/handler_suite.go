package test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/cmd-stream/cmd-stream-go/handler"
	"github.com/cmd-stream/cmd-stream-go/test/mock"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

type HandlerTestCase[T any] struct {
	Name   string
	Setup  HandlerSetup[T]
	Action func(t *testing.T, h *handler.Handler[T])
	Mocks  []*mok.Mock
}

type HandlerSetup[T any] struct {
	Invoker handler.Invoker[T]
	Opts    []handler.SetOption
}

func RunHandlerTestCase[T any](t *testing.T, tc HandlerTestCase[T]) {
	t.Run(tc.Name, func(t *testing.T) {
		h := handler.New[T](tc.Setup.Invoker, tc.Setup.Opts...)
		tc.Action(t, h)
		asserterror.EqualDeep(t, mok.CheckCalls(tc.Mocks), mok.EmptyInfomap)
	})
}

type HandlerSuite[T any] struct{}

// -----------------------------------------------------------------------------
// Test Cases
// -----------------------------------------------------------------------------

func (HandlerSuite[T]) HandleSuccess(t *testing.T) HandlerTestCase[T] {
	name := "Handler should be able to handle several cmds and close when ctx done"

	var (
		ctx, cancel = context.WithCancel(context.Background())
		transport   = mock.NewServerTransport[T]()
		invoker     = mock.NewInvoker[T]()

		seq1 = core.Seq(1)
		cmd1 = mock.NewCmd[T]()
		n1   = 10

		seq2 = core.Seq(2)
		cmd2 = mock.NewCmd[T]()
		n2   = 20

		done = make(chan struct{})
	)
	transport.RegisterReceive(
		func() (core.Seq, core.Cmd[T], int, error) {
			return seq1, cmd1, n1, nil
		},
	).RegisterReceive(
		func() (core.Seq, core.Cmd[T], int, error) {
			return seq2, cmd2, n2, nil
		},
	).RegisterReceive(
		func() (core.Seq, core.Cmd[T], int, error) {
			<-done
			return 0, nil, 0, context.Canceled
		},
	)
	invoker.RegisterInvokeN(2,
		func(c context.Context, seq core.Seq, at time.Time, bytesRead int,
			cmd core.Cmd[T], proxy core.Proxy) error {
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
	return HandlerTestCase[T]{
		Name: name,
		Setup: HandlerSetup[T]{
			Invoker: invoker,
			Opts:    []handler.SetOption{},
		},
		Action: func(t *testing.T, h *handler.Handler[T]) {
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

func (HandlerSuite[T]) SetReceiveDeadlineError(t *testing.T) HandlerTestCase[T] {
	name := "If Transport.SetReceiveDeadline fails with an error, Handle should return it"

	var (
		wantErr   = errors.New("Transport.SetReceiveDeadline error")
		transport = mock.NewServerTransport[T]()
		invoker   = mock.NewInvoker[T]()
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return wantErr },
	).RegisterClose(
		func() error { return nil },
	)
	return HandlerTestCase[T]{
		Name: name,
		Setup: HandlerSetup[T]{
			Invoker: invoker,
			Opts:    []handler.SetOption{handler.WithCmdReceiveDuration(time.Second)},
		},
		Action: func(t *testing.T, h *handler.Handler[T]) {
			err := h.Handle(context.Background(), transport)
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{transport.Mock, invoker.Mock},
	}
}

func (HandlerSuite[T]) ReceiveError(t *testing.T) HandlerTestCase[T] {
	name := "If Transport.Receive fails with an error, Handle should return it"

	var (
		wantErr   = errors.New("Transport.Receive error")
		transport = mock.NewServerTransport[T]()
		invoker   = mock.NewInvoker[T]()
	)
	transport.RegisterReceive(
		func() (core.Seq, core.Cmd[T], int, error) { return 0, nil, 2, wantErr },
	).RegisterClose(
		func() error { return nil },
	)
	return HandlerTestCase[T]{
		Name: name,
		Setup: HandlerSetup[T]{
			Invoker: invoker,
			Opts:    []handler.SetOption{},
		},
		Action: func(t *testing.T, h *handler.Handler[T]) {
			err := h.Handle(context.Background(), transport)
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{transport.Mock, invoker.Mock},
	}
}

func (HandlerSuite[T]) InvokeError(t *testing.T) HandlerTestCase[T] {
	name := "If Invoker.Invoke fails with an error, Handle should return it"

	var (
		wantErr   = errors.New("Invoker.Invoke error")
		transport = mock.NewServerTransport[T]()
		invoker   = mock.NewInvoker[T]()
		seq       = core.Seq(1)
		cmd       = mock.NewCmd[T]()
		n         = 10
		done      = make(chan struct{})
	)
	transport.RegisterReceive(
		func() (core.Seq, core.Cmd[T], int, error) { return seq, cmd, n, nil },
	).RegisterReceive(
		func() (core.Seq, core.Cmd[T], int, error) {
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
			cmd core.Cmd[T], proxy core.Proxy) error {
			return wantErr
		},
	)
	return HandlerTestCase[T]{
		Name: name,
		Setup: HandlerSetup[T]{
			Invoker: invoker,
			Opts:    []handler.SetOption{},
		},
		Action: func(t *testing.T, h *handler.Handler[T]) {
			err := h.Handle(context.Background(), transport)
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{transport.Mock, invoker.Mock, cmd.Mock},
	}
}

func (HandlerSuite[T]) OptionAt(t *testing.T) HandlerTestCase[T] {
	name := "If Conf.At == true, Invoker.Invoke should receive not empty 'at' param"

	var (
		ctx, cancel = context.WithCancel(context.Background())
		transport   = mock.NewServerTransport[T]()
		invoker     = mock.NewInvoker[T]()
		done        = make(chan struct{})
		delta       = 100 * time.Millisecond
	)
	transport.RegisterReceive(
		func() (core.Seq, core.Cmd[T], int, error) { return core.Seq(1), mock.NewCmd[T](), 10, nil },
	).RegisterReceive(
		func() (core.Seq, core.Cmd[T], int, error) {
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
			cmd core.Cmd[T], proxy core.Proxy) error {
			asserterror.SameTime(t, at, time.Now(), delta)
			return nil
		},
	)

	return HandlerTestCase[T]{
		Name: name,
		Setup: HandlerSetup[T]{
			Invoker: invoker,
			Opts:    []handler.SetOption{handler.WithAt()},
		},
		Action: func(t *testing.T, h *handler.Handler[T]) {
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

func (HandlerSuite[T]) OptionCmdReceiveDuration(t *testing.T) HandlerTestCase[T] {
	name := "If Conf.CmdReceiveDuration is set, Transport.SetReceiveDeadline should receive it"

	var (
		ctx, cancel            = context.WithCancel(context.Background())
		wantCmdReceiveDuration = time.Second
		transport              = mock.NewServerTransport[T]()
		invoker                = mock.NewInvoker[T]()
		done                   = make(chan struct{})
		startTime              = time.Now()
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) (err error) {
			wantDeadline := startTime.Add(wantCmdReceiveDuration)
			asserterror.SameTime(t, deadline, wantDeadline, TimeDelta)
			return
		},
	).RegisterReceive(
		func() (core.Seq, core.Cmd[T], int, error) {
			<-done
			return 0, nil, 0, context.Canceled
		},
	).RegisterClose(
		func() error {
			close(done)
			return nil
		},
	)
	return HandlerTestCase[T]{
		Name: name,
		Setup: HandlerSetup[T]{
			Invoker: invoker,
			Opts:    []handler.SetOption{handler.WithCmdReceiveDuration(wantCmdReceiveDuration)},
		},
		Action: func(t *testing.T, h *handler.Handler[T]) {
			startTime = time.Now()
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

func (HandlerSuite[T]) CloseWhileInvokingCmds(t *testing.T) HandlerTestCase[T] {
	name := "We should be able to interrupt Handler, while it invokes several Commands"

	var (
		ctx, cancel = context.WithCancel(context.Background())
		transport   = mock.NewServerTransport[T]()
		invoker     = mock.NewInvoker[T]()
		done        = make(chan struct{})
	)
	transport.RegisterReceive(
		func() (core.Seq, core.Cmd[T], int, error) {
			return core.Seq(1), mock.NewCmd[T](), 10, nil
		},
	).RegisterReceive(
		func() (core.Seq, core.Cmd[T], int, error) {
			return core.Seq(2), mock.NewCmd[T](), 10, nil
		},
	).RegisterReceive(
		func() (core.Seq, core.Cmd[T], int, error) {
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
			cmd core.Cmd[T], proxy core.Proxy) error {
			<-ctx.Done()
			return nil
		},
	)
	return HandlerTestCase[T]{
		Name: name,
		Setup: HandlerSetup[T]{
			Invoker: invoker,
			Opts:    []handler.SetOption{},
		},
		Action: func(t *testing.T, h *handler.Handler[T]) {
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
