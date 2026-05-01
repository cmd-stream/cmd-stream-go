package test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
	cln "github.com/cmd-stream/cmd-stream-go/delegate/cln"
	"github.com/cmd-stream/cmd-stream-go/test/mock"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

type KeepaliveTestCase[T any] struct {
	Name   string
	Setup  KeepaliveSetup[T]
	Action func(t *testing.T, d *cln.KeepaliveDelegate[T])
	Mocks  []*mok.Mock
}

type KeepaliveSetup[T any] struct {
	Delegate mock.ClientDelegate[T]
	Opts     []cln.SetKeepaliveOption
}

func RunKeepaliveTestCase[T any](t *testing.T, tc KeepaliveTestCase[T]) {
	t.Run(tc.Name, func(t *testing.T) {
		d := cln.NewKeepalive[T](tc.Setup.Delegate, tc.Setup.Opts...)

		tc.Action(t, d)
		asserterror.EqualDeep(t, mok.CheckCalls(tc.Mocks), mok.EmptyInfomap)
	})
}

type KeepaliveDelegate[T any] struct{}

// -----------------------------------------------------------------------------
// Test Cases
// -----------------------------------------------------------------------------

func (KeepaliveDelegate[T]) ShouldSendPing(t *testing.T) KeepaliveTestCase[T] {
	name := "Should send Ping Commands if no Commands was sent"

	var (
		done               = make(chan struct{})
		wantCmd            = dlgt.PingCmd[T]{}
		start              time.Time
		wantKeepaliveTime  = 2 * 200 * time.Millisecond
		wantKeepaliveIntvl = 200 * time.Millisecond
		delegateMock       = mock.NewClientDelegate[T]()
	)
	delegateMock.RegisterSetSendDeadlineN(2,
		func(deadline time.Time) error {
			asserterror.SameTime(t, deadline, time.Now().Add(wantKeepaliveIntvl),
				TimeDelta)
			return nil
		},
	).RegisterSend(
		func(seq core.Seq, cmd core.Cmd[T]) (int, error) {
			wantTime := start.Add(wantKeepaliveTime)
			asserterror.SameTime(t, time.Now(), wantTime, TimeDelta)
			asserterror.Equal(t, seq, core.Seq(0))
			asserterror.EqualDeep(t, cmd, wantCmd)
			return 1, nil
		},
	).RegisterFlush(
		func() error { return nil },
	).RegisterSend(
		func(seq core.Seq, cmd core.Cmd[T]) (int, error) {
			wantTime := start.Add(wantKeepaliveTime).Add(wantKeepaliveIntvl)
			asserterror.SameTime(t, time.Now(), wantTime, TimeDelta)
			asserterror.Equal(t, seq, core.Seq(0))
			asserterror.EqualDeep(t, cmd, wantCmd)
			return 1, nil
		},
	).RegisterFlush(
		func() error { defer close(done); return nil },
	).RegisterClose(
		func() error { return nil },
	)
	return KeepaliveTestCase[T]{
		Name: name,
		Setup: KeepaliveSetup[T]{
			Delegate: delegateMock,
			Opts: []cln.SetKeepaliveOption{
				cln.WithKeepaliveTime(wantKeepaliveTime),
				cln.WithKeepaliveIntvl(wantKeepaliveIntvl),
			},
		},
		Action: func(t *testing.T, d *cln.KeepaliveDelegate[T]) {
			start = time.Now()
			d.Keepalive(&sync.Mutex{})
			select {
			case <-done:
			case <-time.After(wantKeepaliveTime + wantKeepaliveIntvl + time.Second):
				t.Fatal("test lasts too long")
			}
			err := d.Close()
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{delegateMock.Mock},
	}
}

func (KeepaliveDelegate[T]) FlushDelay(t *testing.T) KeepaliveTestCase[T] {
	name := "Command flushing should delay a ping"

	var (
		done               = make(chan struct{})
		start              time.Time
		flushDelay         = 200 * time.Millisecond
		wantKeepaliveTime  = 2 * 200 * time.Millisecond
		wantKeepaliveIntvl = 200 * time.Millisecond
		delegateMock       = mock.NewClientDelegate[T]()
	)
	delegateMock.RegisterFlush(
		func() error { return nil },
	).RegisterSetSendDeadline(
		func(deadline time.Time) error {
			asserterror.SameTime(t, deadline, time.Now().Add(wantKeepaliveIntvl),
				TimeDelta)
			return nil
		},
	).RegisterSend(
		func(seq core.Seq, cmd core.Cmd[T]) (int, error) {
			// flushDelay + wantKeepaliveTime
			wantTime := start.Add(flushDelay).Add(wantKeepaliveTime)
			asserterror.SameTime(t, time.Now(), wantTime, TimeDelta)
			return 0, nil
		},
	).RegisterFlush(
		func() error { defer close(done); return nil },
	).RegisterClose(
		func() error { return nil },
	)
	return KeepaliveTestCase[T]{
		Name: name,
		Setup: KeepaliveSetup[T]{
			Delegate: delegateMock,
			Opts: []cln.SetKeepaliveOption{
				cln.WithKeepaliveTime(wantKeepaliveTime),
				cln.WithKeepaliveIntvl(wantKeepaliveIntvl),
			},
		},
		Action: func(t *testing.T, d *cln.KeepaliveDelegate[T]) {
			start = time.Now()
			d.Keepalive(&sync.Mutex{})

			time.Sleep(flushDelay)
			err := d.Flush()
			asserterror.EqualError(t, err, nil)

			select {
			case <-done:
			case <-time.After(wantKeepaliveTime + flushDelay + time.Second):
				t.Fatal("test lasts too long")
			}
			err = d.Close()
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{delegateMock.Mock},
	}
}

func (KeepaliveDelegate[T]) CloseCancel(t *testing.T) KeepaliveTestCase[T] {
	name := "Close should cancel ping sending"

	var delegateMock = mock.NewClientDelegate[T]()
	delegateMock.RegisterClose(
		func() error { return nil },
	)
	return KeepaliveTestCase[T]{
		Name: name,
		Setup: KeepaliveSetup[T]{
			Delegate: delegateMock,
			Opts: []cln.SetKeepaliveOption{
				cln.WithKeepaliveTime(200 * time.Millisecond),
				cln.WithKeepaliveIntvl(200 * time.Millisecond),
			},
		},
		Action: func(t *testing.T, d *cln.KeepaliveDelegate[T]) {
			d.Keepalive(&sync.Mutex{})

			err := d.Close()
			asserterror.EqualError(t, err, nil)
			time.Sleep(300 * time.Millisecond) // wait more than KeepaliveTime (200ms)
		},
		Mocks: []*mok.Mock{delegateMock.Mock},
	}
}

func (KeepaliveDelegate[T]) CloseError(t *testing.T) KeepaliveTestCase[T] {
	name := "If ClientDelegate.Close fails with an error, Close should return it and ping should not be canceled"

	var (
		done         = make(chan struct{})
		wantErr      = errors.New("close error")
		delegateMock = mock.NewClientDelegate[T]()
	)
	delegateMock.RegisterClose(
		func() error { return wantErr },
	).RegisterSetSendDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterSend(
		func(seq core.Seq, cmd core.Cmd[T]) (int, error) { return 1, nil },
	).RegisterFlush(
		func() error { defer close(done); return nil },
	).RegisterClose(
		func() error { return nil },
	)
	return KeepaliveTestCase[T]{
		Name: name,
		Setup: KeepaliveSetup[T]{
			Delegate: delegateMock,
			Opts: []cln.SetKeepaliveOption{
				cln.WithKeepaliveTime(200 * time.Millisecond),
				cln.WithKeepaliveIntvl(200 * time.Millisecond),
			},
		},
		Action: func(t *testing.T, d *cln.KeepaliveDelegate[T]) {
			d.Keepalive(&sync.Mutex{})

			err := d.Close()
			asserterror.EqualError(t, err, wantErr)

			// Second Close returns nil, cleanly shutting down the loop
			select {
			case <-done:
			case <-time.After(time.Second):
				t.Fatal("test lasts too long")
			}
			err = d.Close()
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{delegateMock.Mock},
	}
}

func (KeepaliveDelegate[T]) SendError(t *testing.T) KeepaliveTestCase[T] {
	name := "If ping sending fails with an error, connection should NOT be closed and ping should be retried"

	var (
		done     = make(chan struct{})
		delegate = mock.NewClientDelegate[T]()
	)
	// First attempt fails.
	delegate.RegisterSetSendDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterSend(
		func(seq core.Seq, cmd core.Cmd[T]) (int, error) {
			return 1, errors.New("send error")
		},
	)
	// Second attempt succeeds - proving the loop is still alive.
	delegate.RegisterSetSendDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterSend(
		func(seq core.Seq, cmd core.Cmd[T]) (int, error) {
			return 1, nil
		},
	).RegisterFlush(
		func() error {
			defer close(done)
			return nil
		},
	).RegisterClose(
		func() error { return nil },
	)
	return KeepaliveTestCase[T]{
		Name: name,
		Setup: KeepaliveSetup[T]{
			Delegate: delegate,
			Opts: []cln.SetKeepaliveOption{
				cln.WithKeepaliveTime(100 * time.Millisecond),
				cln.WithKeepaliveIntvl(100 * time.Millisecond),
			},
		},
		Action: func(t *testing.T, d *cln.KeepaliveDelegate[T]) {
			d.Keepalive(&sync.Mutex{})

			select {
			case <-done:
			case <-time.After(time.Second):
				t.Fatal("test lasts too long")
			}
			err := d.Close()
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{delegate.Mock},
	}
}

func (KeepaliveDelegate[T]) FlushError(t *testing.T) KeepaliveTestCase[T] {
	name := "If ClientDelegate.Flush fails with an error, Flush should return it and ping sending should not be delayed"

	var (
		done         = make(chan struct{})
		wantErr      = errors.New("flush error")
		delegateMock = mock.NewClientDelegate[T]()
	)

	delegateMock.RegisterFlush(
		func() error { return wantErr },
	).RegisterSetSendDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterSend(
		func(seq core.Seq, cmd core.Cmd[T]) (int, error) {
			return 1, nil
		},
	).RegisterFlush(
		func() error {
			defer close(done)
			return nil
		},
	).RegisterClose(
		func() error { return nil },
	)

	return KeepaliveTestCase[T]{
		Name: name,
		Setup: KeepaliveSetup[T]{
			Delegate: delegateMock,
			Opts: []cln.SetKeepaliveOption{
				cln.WithKeepaliveTime(200 * time.Millisecond),
				cln.WithKeepaliveIntvl(200 * time.Millisecond),
			},
		},
		Action: func(t *testing.T, d *cln.KeepaliveDelegate[T]) {
			d.Keepalive(&sync.Mutex{})

			err := d.Flush()
			asserterror.EqualError(t, err, wantErr)

			select {
			case <-done:
			case <-time.After(time.Second):
				t.Fatal("test lasts too long")
			}
			err = d.Close()
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{delegateMock.Mock},
	}
}

func (KeepaliveDelegate[T]) SkipPong(t *testing.T) KeepaliveTestCase[T] {
	name := "Should skip Pong Result"

	var (
		wantSeq      = core.Seq(1)
		wantResult   = mock.NewResult()
		wantN        = 1
		delegateMock = mock.NewClientDelegate[T]()
	)
	delegateMock.RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			return 0, dlgt.PongResult{}, 1, nil
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			return wantSeq, wantResult, wantN, nil
		},
	)
	return KeepaliveTestCase[T]{
		Name: name,
		Setup: KeepaliveSetup[T]{
			Delegate: delegateMock,
		},
		Action: func(t *testing.T, d *cln.KeepaliveDelegate[T]) {
			seq, result, n, err := d.Receive()

			asserterror.Equal(t, seq, wantSeq)
			asserterror.EqualDeep(t, result, wantResult)
			asserterror.Equal(t, n, wantN)
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{delegateMock.Mock, wantResult.Mock},
	}
}
