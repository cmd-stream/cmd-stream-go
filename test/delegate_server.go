package test

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
	"github.com/cmd-stream/cmd-stream-go/delegate/srv"
	"github.com/cmd-stream/cmd-stream-go/test/mock"
	asserterror "github.com/ymz-ncnk/assert/error"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
	"github.com/ymz-ncnk/mok"
)

type ServerDelegateTestCase[T any] struct {
	Name   string
	Setup  ServerDelegateSetup[T]
	Action func(t *testing.T, d srv.ServerInfoDelegate[T], initErr error)
	Mocks  []*mok.Mock
}

type ServerDelegateSetup[T any] struct {
	Info    dlgt.ServerInfo
	Factory mock.ServerTransportFactory[T]
	Handler mock.ServerTransportHandler[T]
	Opts    []srv.SetOption
}

func RunServerDelegateTestCase[T any](t *testing.T, tc ServerDelegateTestCase[T]) {
	t.Run(tc.Name, func(t *testing.T) {
		d, err := srv.New[T](tc.Setup.Info, tc.Setup.Factory, tc.Setup.Handler,
			tc.Setup.Opts...)

		tc.Action(t, d, err)
		asserterror.EqualDeep(t, mok.CheckCalls(tc.Mocks), mok.EmptyInfomap)
	})
}

type DelegateServer[T any] struct{}

// -----------------------------------------------------------------------------
// Test Cases
// -----------------------------------------------------------------------------

func (DelegateServer[T]) HandleConnSuccess(t *testing.T) ServerDelegateTestCase[T] {
	name := "If Handle succeeds, it should call all components correctly"

	var (
		wantCtx                = context.Background()
		wantConn               = &net.TCPConn{}
		wantInfo               = dlgt.ServerInfo([]byte{1, 2, 3})
		serverInfoSendDuration = time.Second
		startTime              time.Time
		factory                = mock.NewServerTransportFactory[T]()
		handler                = mock.NewServerTransportHandler[T]()
		wantTransport          = mock.NewServerTransport[T]()
	)

	factory.RegisterNew(
		func(conn net.Conn) dlgt.ServerTransport[T] {
			startTime = time.Now()
			asserterror.EqualDeep(t, conn, wantConn)
			return wantTransport
		},
	)
	wantTransport.RegisterSetSendDeadline(
		func(deadline time.Time) error {
			wantDeadline := startTime.Add(serverInfoSendDuration)
			asserterror.SameTime(t, deadline, wantDeadline, TimeDelta)
			return nil
		},
	).RegisterSendServerInfo(
		func(info dlgt.ServerInfo) error {
			asserterror.EqualDeep(t, info, wantInfo)
			return nil
		},
	)
	handler.RegisterHandle(
		func(ctx context.Context, transport dlgt.ServerTransport[T]) error {
			asserterror.Equal(t, ctx, wantCtx)
			asserterror.EqualDeep(t, transport, wantTransport)
			return nil
		},
	)
	return ServerDelegateTestCase[T]{
		Name: name,
		Setup: ServerDelegateSetup[T]{
			Info:    wantInfo,
			Factory: factory,
			Handler: handler,
			Opts: []srv.SetOption{
				srv.WithServerInfoSendDuration(serverInfoSendDuration),
			},
		},
		Action: func(t *testing.T, d srv.ServerInfoDelegate[T], initErr error) {
			assertfatal.Equal(t, initErr, nil)
			err := d.Handle(wantCtx, wantConn)
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{factory.Mock, wantTransport.Mock, handler.Mock},
	}
}

func (DelegateServer[T]) SendServerInfoError(t *testing.T) ServerDelegateTestCase[T] {
	name := "If send ServerInfo fails with an error, Handle should return it"

	var (
		info      = dlgt.ServerInfo([]byte{1, 2, 3})
		wantErr   = errors.New("send server info error")
		factory   = mock.NewServerTransportFactory[T]()
		handler   = mock.NewServerTransportHandler[T]()
		transport = mock.NewServerTransport[T]()
	)
	factory.RegisterNew(
		func(conn net.Conn) dlgt.ServerTransport[T] { return transport },
	)
	transport.RegisterSendServerInfo(
		func(gotInfo dlgt.ServerInfo) error { return wantErr },
	).RegisterClose(
		func() error { return errors.New("should be ignored") },
	)
	return ServerDelegateTestCase[T]{
		Name: name,
		Setup: ServerDelegateSetup[T]{
			Info:    info,
			Factory: factory,
			Handler: handler,
		},
		Action: func(t *testing.T, d srv.ServerInfoDelegate[T], initErr error) {
			assertfatal.EqualError(t, initErr, nil)
			err := d.Handle(context.Background(), &net.TCPConn{})
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{factory.Mock, transport.Mock, handler.Mock},
	}
}

func (DelegateServer[T]) ZeroLenServerInfo(t *testing.T) ServerDelegateTestCase[T] {
	name := "If ServerInfo len is zero, New should return an error"

	return ServerDelegateTestCase[T]{
		Name: name,
		Setup: ServerDelegateSetup[T]{
			Info:    dlgt.ServerInfo([]byte{}),
			Factory: mock.NewServerTransportFactory[T](),
			Handler: mock.NewServerTransportHandler[T](),
		},
		Action: func(t *testing.T, d srv.ServerInfoDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, srv.ErrEmptyInfo)
		},
		Mocks: nil,
	}
}

func (DelegateServer[T]) TransportHandleError(t *testing.T) ServerDelegateTestCase[T] {
	name := "If Transport.Handle fails with an error, Handle should return it"

	var (
		wantCtx       = context.Background()
		conn          = &net.TCPConn{}
		wantInfo      = dlgt.ServerInfo([]byte{1, 2, 3})
		wantErr       = errors.New("transport handle error")
		factory       = mock.NewServerTransportFactory[T]()
		handler       = mock.NewServerTransportHandler[T]()
		wantTransport = mock.NewServerTransport[T]()
	)
	factory.RegisterNew(
		func(gotConn net.Conn) dlgt.ServerTransport[T] { return wantTransport },
	)
	wantTransport.RegisterSendServerInfo(
		func(info dlgt.ServerInfo) error { return nil },
	)
	handler.RegisterHandle(
		func(ctx context.Context, transport dlgt.ServerTransport[T]) error {
			return wantErr
		},
	)
	return ServerDelegateTestCase[T]{
		Name: name,
		Setup: ServerDelegateSetup[T]{
			Info:    wantInfo,
			Factory: factory,
			Handler: handler,
		},
		Action: func(t *testing.T, d srv.ServerInfoDelegate[T], initErr error) {
			assertfatal.EqualError(t, initErr, nil)
			err := d.Handle(wantCtx, conn)
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{factory.Mock, wantTransport.Mock, handler.Mock},
	}
}

func (DelegateServer[T]) SendServerInfoTransportDeadlineError(t *testing.T) ServerDelegateTestCase[T] {
	name := "If Transport.SetSendDeadline fails with an error on ServerInfo send, Handle should return it"

	var (
		ctx       = context.Background()
		conn      = &net.TCPConn{}
		info      = dlgt.ServerInfo([]byte{1, 2, 3})
		wantErr   = errors.New("set send deadline error")
		factory   = mock.NewServerTransportFactory[T]()
		handler   = mock.NewServerTransportHandler[T]()
		transport = mock.NewServerTransport[T]()
	)
	factory.RegisterNew(
		func(gotConn net.Conn) dlgt.ServerTransport[T] { return transport },
	)
	transport.RegisterSetSendDeadline(
		func(deadline time.Time) error { return wantErr },
	).RegisterClose(
		func() error { return nil },
	)
	return ServerDelegateTestCase[T]{
		Name: name,
		Setup: ServerDelegateSetup[T]{
			Info:    info,
			Factory: factory,
			Handler: handler,
			Opts: []srv.SetOption{
				srv.WithServerInfoSendDuration(time.Second),
			},
		},
		Action: func(t *testing.T, d srv.ServerInfoDelegate[T], initErr error) {
			assertfatal.EqualError(t, initErr, nil)
			err := d.Handle(ctx, conn)
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{factory.Mock, transport.Mock, handler.Mock},
	}
}
