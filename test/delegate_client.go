package test

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
	"github.com/cmd-stream/cmd-stream-go/delegate/cln"
	"github.com/cmd-stream/cmd-stream-go/test/mock"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

type ClientDelegateTestCase[T any] struct {
	Name   string
	Setup  ClientDelegateSetup[T]
	Action func(t *testing.T, d cln.ClientInfoDelegate[T], initErr error)
	Mocks  []*mok.Mock
}

type ClientDelegateSetup[T any] struct {
	Info      dlgt.ServerInfo
	Transport mock.ClientTransport[T]
	Opts      []cln.SetOption
}

func RunClientDelegateTestCase[T any](t *testing.T, tc ClientDelegateTestCase[T]) {
	t.Run(tc.Name, func(t *testing.T) {
		d, err := cln.New(tc.Setup.Info, tc.Setup.Transport, tc.Setup.Opts...)

		tc.Action(t, d, err)
		asserterror.EqualDeep(t, mok.CheckCalls(tc.Mocks), mok.EmptyInfomap)
	})
}

type DelegateClient[T any] struct{}

// -----------------------------------------------------------------------------
// Test Cases
// -----------------------------------------------------------------------------

func (DelegateClient[T]) NewCheckServerInfo(t *testing.T) ClientDelegateTestCase[T] {
	name := "New should check ServerInfo"

	var (
		wantInfo  = dlgt.ServerInfo{1, 2, 3}
		transport = mock.NewClientTransport[T]()
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (info dlgt.ServerInfo, err error) { return wantInfo, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	)
	return ClientDelegateTestCase[T]{
		Name: name,
		Setup: ClientDelegateSetup[T]{
			Info:      wantInfo,
			Transport: transport,
		},
		Action: func(t *testing.T, d cln.ClientInfoDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, nil)
		},
		Mocks: []*mok.Mock{transport.Mock},
	}
}

func (DelegateClient[T]) NewSetReceiveDeadlineError(t *testing.T) ClientDelegateTestCase[T] {
	name := "If Transport.SetReceiveDeadline fails with an error before receive ServerInfo, New should return it"

	var (
		wantInfo  = dlgt.ServerInfo{1, 2, 3}
		wantErr   = errors.New("Transport.SetReceiveDeadline")
		transport = mock.NewClientTransport[T]()
	)

	transport.RegisterSetReceiveDeadline(func(deadline time.Time) error {
		return wantErr
	})

	return ClientDelegateTestCase[T]{
		Name: name,
		Setup: ClientDelegateSetup[T]{
			Info:      wantInfo,
			Transport: transport,
		},
		Action: func(t *testing.T, d cln.ClientInfoDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, wantErr)
		},
		Mocks: []*mok.Mock{transport.Mock},
	}
}

func (DelegateClient[T]) NewReceiveServerInfoError(t *testing.T) ClientDelegateTestCase[T] {
	name := "If Transport.ReceiveServerInfo fails with an error, New should return it"

	var (
		wantInfo  = dlgt.ServerInfo{1, 2, 3}
		wantErr   = errors.New("Transport.ReceiveServerInfo error")
		transport = mock.NewClientTransport[T]()
	)

	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (info dlgt.ServerInfo, err error) { return nil, wantErr },
	)

	return ClientDelegateTestCase[T]{
		Name: name,
		Setup: ClientDelegateSetup[T]{
			Info:      wantInfo,
			Transport: transport,
		},
		Action: func(t *testing.T, d cln.ClientInfoDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, wantErr)
		},
		Mocks: []*mok.Mock{transport.Mock},
	}
}

func (DelegateClient[T]) NewServerInfoMismatch(t *testing.T) ClientDelegateTestCase[T] {
	name := "If wrong ServerInfo was received, New should return error"

	var (
		wantInfo        = dlgt.ServerInfo{1, 2, 3}
		wrongServerInfo = dlgt.ServerInfo{1}
		transport       = mock.NewClientTransport[T]()
	)

	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (info dlgt.ServerInfo, err error) { return wrongServerInfo, nil },
	)

	return ClientDelegateTestCase[T]{
		Name: name,
		Setup: ClientDelegateSetup[T]{
			Info:      wantInfo,
			Transport: transport,
		},
		Action: func(t *testing.T, d cln.ClientInfoDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, cln.ErrServerInfoMismatch)
		},
		Mocks: []*mok.Mock{transport.Mock},
	}
}

func (DelegateClient[T]) Send(t *testing.T) ClientDelegateTestCase[T] {
	name := "Send should call Transport.Send"

	var (
		wantSeq   core.Seq = 1
		wantCmd            = mock.NewCmd[T]()
		wantN              = 2
		transport          = mock.NewClientTransport[T]()
	)

	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (info dlgt.ServerInfo, err error) { return dlgt.ServerInfo{}, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterSend(
		func(seq core.Seq, cmd core.Cmd[T]) (n int, err error) {
			asserterror.Equal(t, seq, wantSeq)
			asserterror.EqualDeep(t, cmd, wantCmd)
			return wantN, nil
		},
	)

	return ClientDelegateTestCase[T]{
		Name: name,
		Setup: ClientDelegateSetup[T]{
			Transport: transport,
		},
		Action: func(t *testing.T, d cln.ClientInfoDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			n, err := d.Send(wantSeq, wantCmd)
			asserterror.Equal(t, n, wantN)
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{transport.Mock},
	}
}

func (DelegateClient[T]) SendError(t *testing.T) ClientDelegateTestCase[T] {
	name := "If Transport.Send fails with an error, Send should return it"

	var (
		wantErr   = errors.New("send error")
		transport = mock.NewClientTransport[T]()
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (info dlgt.ServerInfo, err error) { return dlgt.ServerInfo{}, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterSend(
		func(seq core.Seq, cmd core.Cmd[T]) (n int, err error) {
			return 0, wantErr
		},
	)
	return ClientDelegateTestCase[T]{
		Name: name,
		Setup: ClientDelegateSetup[T]{
			Transport: transport,
		},
		Action: func(t *testing.T, d cln.ClientInfoDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			_, err := d.Send(0, nil)
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{transport.Mock},
	}
}

func (DelegateClient[T]) Receive(t *testing.T) ClientDelegateTestCase[T] {
	name := "Receive should return values from Transport.Receive"

	var (
		wantSeq    core.Seq = 1
		wantResult          = mock.NewResult()
		wantN               = 3
		wantErr             = errors.New("receive failed")
		transport           = mock.NewClientTransport[T]()
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (info dlgt.ServerInfo, err error) { return dlgt.ServerInfo{}, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceive(
		func() (seq core.Seq, r core.Result, n int, err error) {
			return wantSeq, wantResult, wantN, wantErr
		},
	)
	return ClientDelegateTestCase[T]{
		Name: name,
		Setup: ClientDelegateSetup[T]{
			Transport: transport,
		},
		Action: func(t *testing.T, d cln.ClientInfoDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			seq, result, n, err := d.Receive()
			asserterror.Equal(t, seq, wantSeq)
			asserterror.EqualDeep(t, result, wantResult)
			asserterror.Equal(t, n, wantN)
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{transport.Mock, wantResult.Mock},
	}
}

func (DelegateClient[T]) LocalAddr(t *testing.T) ClientDelegateTestCase[T] {
	name := "LocalAddr should return Transport.LocalAddr"

	var (
		wantAddr  = &net.IPAddr{IP: net.ParseIP("127.0.0.1")}
		transport = mock.NewClientTransport[T]()
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (info dlgt.ServerInfo, err error) { return dlgt.ServerInfo{}, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterLocalAddr(
		func() net.Addr { return wantAddr },
	)
	return ClientDelegateTestCase[T]{
		Name: name,
		Setup: ClientDelegateSetup[T]{
			Transport: transport,
		},
		Action: func(t *testing.T, d cln.ClientInfoDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			addr := d.LocalAddr()
			asserterror.Equal(t, addr, (net.Addr)(wantAddr))
		},
		Mocks: []*mok.Mock{transport.Mock},
	}
}

func (DelegateClient[T]) RemoteAddr(t *testing.T) ClientDelegateTestCase[T] {
	name := "RemoteAddr should return Transport.RemoteAddr"

	var (
		wantAddr  = &net.IPAddr{IP: net.ParseIP("127.0.0.1")}
		transport = mock.NewClientTransport[T]()
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (info dlgt.ServerInfo, err error) { return dlgt.ServerInfo{}, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterRemoteAddr(
		func() net.Addr { return wantAddr },
	)
	return ClientDelegateTestCase[T]{
		Name: name,
		Setup: ClientDelegateSetup[T]{
			Transport: transport,
		},
		Action: func(t *testing.T, d cln.ClientInfoDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			addr := d.RemoteAddr()
			asserterror.Equal(t, addr, (net.Addr)(wantAddr))
		},
		Mocks: []*mok.Mock{transport.Mock},
	}
}

func (DelegateClient[T]) Close(t *testing.T) ClientDelegateTestCase[T] {
	name := "Close should call Transport.Close"

	var (
		wantErr   = errors.New("close error")
		transport = mock.NewClientTransport[T]()
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (info dlgt.ServerInfo, err error) { return dlgt.ServerInfo{}, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterClose(
		func() error { return wantErr },
	)
	return ClientDelegateTestCase[T]{
		Name: name,
		Setup: ClientDelegateSetup[T]{
			Transport: transport,
		},
		Action: func(t *testing.T, d cln.ClientInfoDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			err := d.Close()
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{transport.Mock},
	}
}
