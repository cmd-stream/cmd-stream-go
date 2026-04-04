package delegate

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
	cln "github.com/cmd-stream/cmd-stream-go/delegate/cln"
	"github.com/cmd-stream/cmd-stream-go/test"
	cmock "github.com/cmd-stream/cmd-stream-go/test/mock/core"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

func KeepaliveShouldSendPingTestCase(t *testing.T) KeepaliveTestCase {
	name := "Should send Ping Commands if no Commands was sent"

	var (
		done               = make(chan struct{})
		wantCmd            = dlgt.PingCmd[any]{}
		start              time.Time
		wantKeepaliveTime  = 2 * 200 * time.Millisecond
		wantKeepaliveIntvl = 200 * time.Millisecond
		delegateMock       = cmock.NewClientDelegate()
	)
	delegateMock.RegisterSetSendDeadlineN(2,
		func(deadline time.Time) error {
			asserterror.SameTime(t, deadline, time.Now().Add(wantKeepaliveIntvl),
				test.TimeDelta)
			return nil
		},
	).RegisterSend(
		func(seq core.Seq, cmd core.Cmd[any]) (int, error) {
			wantTime := start.Add(wantKeepaliveTime)
			asserterror.SameTime(t, time.Now(), wantTime, test.TimeDelta)
			asserterror.Equal(t, seq, core.Seq(0))
			asserterror.EqualDeep(t, cmd, wantCmd)
			return 1, nil
		},
	).RegisterFlush(
		func() error { return nil },
	).RegisterSend(
		func(seq core.Seq, cmd core.Cmd[any]) (int, error) {
			wantTime := start.Add(wantKeepaliveTime).Add(wantKeepaliveIntvl)
			asserterror.SameTime(t, time.Now(), wantTime, test.TimeDelta)
			asserterror.Equal(t, seq, core.Seq(0))
			asserterror.EqualDeep(t, cmd, wantCmd)
			return 1, nil
		},
	).RegisterFlush(
		func() error { defer close(done); return nil },
	).RegisterClose(
		func() error { return nil },
	)
	return KeepaliveTestCase{
		Name: name,
		Setup: KeepaliveSetup{
			Delegate: delegateMock,
			Opts: []cln.SetKeepaliveOption{
				cln.WithKeepaliveTime(wantKeepaliveTime),
				cln.WithKeepaliveIntvl(wantKeepaliveIntvl),
			},
		},
		Action: func(t *testing.T, d *cln.KeepaliveDelegate[any]) {
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

func KeepaliveFlushDelayTestCase(t *testing.T) KeepaliveTestCase {
	name := "Command flushing should delay a ping"

	var (
		done               = make(chan struct{})
		start              time.Time
		flushDelay         = 200 * time.Millisecond
		wantKeepaliveTime  = 2 * 200 * time.Millisecond
		wantKeepaliveIntvl = 200 * time.Millisecond
		delegateMock       = cmock.NewClientDelegate()
	)
	delegateMock.RegisterFlush(
		func() error { return nil },
	).RegisterSetSendDeadline(
		func(deadline time.Time) error {
			asserterror.SameTime(t, deadline, time.Now().Add(wantKeepaliveIntvl),
				test.TimeDelta)
			return nil
		},
	).RegisterSend(
		func(seq core.Seq, cmd core.Cmd[any]) (int, error) {
			// flushDelay + wantKeepaliveTime
			wantTime := start.Add(flushDelay).Add(wantKeepaliveTime)
			asserterror.SameTime(t, time.Now(), wantTime, test.TimeDelta)
			return 0, nil
		},
	).RegisterFlush(
		func() error { defer close(done); return nil },
	).RegisterClose(
		func() error { return nil },
	)
	return KeepaliveTestCase{
		Name: name,
		Setup: KeepaliveSetup{
			Delegate: delegateMock,
			Opts: []cln.SetKeepaliveOption{
				cln.WithKeepaliveTime(wantKeepaliveTime),
				cln.WithKeepaliveIntvl(wantKeepaliveIntvl),
			},
		},
		Action: func(t *testing.T, d *cln.KeepaliveDelegate[any]) {
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

func KeepaliveCloseCancelTestCase(t *testing.T) KeepaliveTestCase {
	name := "Close should cancel ping sending"

	var delegateMock = cmock.NewClientDelegate()
	delegateMock.RegisterClose(
		func() error { return nil },
	)
	return KeepaliveTestCase{
		Name: name,
		Setup: KeepaliveSetup{
			Delegate: delegateMock,
			Opts: []cln.SetKeepaliveOption{
				cln.WithKeepaliveTime(200 * time.Millisecond),
				cln.WithKeepaliveIntvl(200 * time.Millisecond),
			},
		},
		Action: func(t *testing.T, d *cln.KeepaliveDelegate[any]) {
			d.Keepalive(&sync.Mutex{})

			err := d.Close()
			asserterror.EqualError(t, err, nil)
			time.Sleep(300 * time.Millisecond) // wait more than KeepaliveTime (200ms)
		},
		Mocks: []*mok.Mock{delegateMock.Mock},
	}
}

func KeepaliveCloseErrorTestCase(t *testing.T) KeepaliveTestCase {
	name := "If ClientDelegate.Close fails with an error, Close should return it and ping should not be canceled"

	var (
		done         = make(chan struct{})
		wantErr      = errors.New("close error")
		delegateMock = cmock.NewClientDelegate()
	)
	delegateMock.RegisterClose(
		func() error { return wantErr },
	).RegisterSetSendDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterSend(
		func(seq core.Seq, cmd core.Cmd[any]) (int, error) { return 1, nil },
	).RegisterFlush(
		func() error { defer close(done); return nil },
	).RegisterClose(
		func() error { return nil },
	)
	return KeepaliveTestCase{
		Name: name,
		Setup: KeepaliveSetup{
			Delegate: delegateMock,
			Opts: []cln.SetKeepaliveOption{
				cln.WithKeepaliveTime(200 * time.Millisecond),
				cln.WithKeepaliveIntvl(200 * time.Millisecond),
			},
		},
		Action: func(t *testing.T, d *cln.KeepaliveDelegate[any]) {
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

func KeepaliveSendErrorTestCase(t *testing.T) KeepaliveTestCase {
	name := "If ping sending fails with an error, connection should be closed"

	var (
		done     = make(chan struct{})
		delegate = cmock.NewClientDelegate()
	)
	delegate.RegisterSetSendDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterSend(
		func(seq core.Seq, cmd core.Cmd[any]) (int, error) {
			return 1, errors.New("send error")
		},
	).RegisterClose(
		func() error {
			defer close(done)
			return nil
		},
	)
	return KeepaliveTestCase{
		Name: name,
		Setup: KeepaliveSetup{
			Delegate: delegate,
			Opts: []cln.SetKeepaliveOption{
				cln.WithKeepaliveTime(200 * time.Millisecond),
				cln.WithKeepaliveIntvl(200 * time.Millisecond),
			},
		},
		Action: func(t *testing.T, d *cln.KeepaliveDelegate[any]) {
			d.Keepalive(&sync.Mutex{})

			select {
			case <-done:
			case <-time.After(time.Second):
				t.Fatal("test lasts too long")
			}
		},
		Mocks: []*mok.Mock{delegate.Mock},
	}
}

func KeepaliveFlushErrorTestCase(t *testing.T) KeepaliveTestCase {
	name := "If ClientDelegate.Flush fails with an error, Flush should return it and ping sending should not be delayed"

	var (
		done         = make(chan struct{})
		wantErr      = errors.New("flush error")
		delegateMock = cmock.NewClientDelegate()
	)

	delegateMock.RegisterFlush(
		func() error { return wantErr },
	).RegisterSetSendDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterSend(
		func(seq core.Seq, cmd core.Cmd[any]) (int, error) {
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

	return KeepaliveTestCase{
		Name: name,
		Setup: KeepaliveSetup{
			Delegate: delegateMock,
			Opts: []cln.SetKeepaliveOption{
				cln.WithKeepaliveTime(200 * time.Millisecond),
				cln.WithKeepaliveIntvl(200 * time.Millisecond),
			},
		},
		Action: func(t *testing.T, d *cln.KeepaliveDelegate[any]) {
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

func KeepaliveSkipPongTestCase(t *testing.T) KeepaliveTestCase {
	name := "Should skip Pong Result"

	var (
		wantSeq      = core.Seq(1)
		wantResult   = cmock.NewResult()
		wantN        = 1
		delegateMock = cmock.NewClientDelegate()
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
	return KeepaliveTestCase{
		Name: name,
		Setup: KeepaliveSetup{
			Delegate: delegateMock,
		},
		Action: func(t *testing.T, d *cln.KeepaliveDelegate[any]) {
			seq, result, n, err := d.Receive()

			asserterror.Equal(t, seq, wantSeq)
			asserterror.EqualDeep(t, result, wantResult)
			asserterror.Equal(t, n, wantN)
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{delegateMock.Mock, wantResult.Mock},
	}
}
