package delegate

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
	"github.com/cmd-stream/cmd-stream-go/delegate/srv"
	"github.com/cmd-stream/cmd-stream-go/test"
	dmock "github.com/cmd-stream/cmd-stream-go/test/mock/delegate"
	asserterror "github.com/ymz-ncnk/assert/error"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
	"github.com/ymz-ncnk/mok"
)

func HandleConnSuccessTestCase(t *testing.T) ServerDelegateTestCase[any] {
	name := "If Handle succeeds, it should call all components correctly"

	var (
		wantCtx                = context.Background()
		wantConn               = &net.TCPConn{}
		wantInfo               = dlgt.ServerInfo([]byte{1, 2, 3})
		serverInfoSendDuration = time.Second
		startTime              time.Time
		factory                = dmock.NewServerTransportFactory[any]()
		handler                = dmock.NewServerTransportHandler[any]()
		wantTransport          = dmock.NewServerTransport[any]()
	)

	factory.RegisterNew(
		func(conn net.Conn) dlgt.ServerTransport[any] {
			startTime = time.Now()
			asserterror.EqualDeep(t, conn, wantConn)
			return wantTransport
		},
	)
	wantTransport.RegisterSetSendDeadline(
		func(deadline time.Time) error {
			wantDeadline := startTime.Add(serverInfoSendDuration)
			asserterror.SameTime(t, deadline, wantDeadline, test.TimeDelta)
			return nil
		},
	).RegisterSendServerInfo(
		func(info dlgt.ServerInfo) error {
			asserterror.EqualDeep(t, info, wantInfo)
			return nil
		},
	)
	handler.RegisterHandle(
		func(ctx context.Context, transport dlgt.ServerTransport[any]) error {
			asserterror.Equal(t, ctx, wantCtx)
			asserterror.EqualDeep(t, transport, wantTransport)
			return nil
		},
	)
	return ServerDelegateTestCase[any]{
		Name: name,
		Setup: ServerDelegateSetup[any]{
			Info:    wantInfo,
			Factory: factory,
			Handler: handler,
			Opts: []srv.SetOption{
				srv.WithServerInfoSendDuration(serverInfoSendDuration),
			},
		},
		Action: func(t *testing.T, d srv.ServerInfoDelegate[any], initErr error) {
			assertfatal.Equal(t, initErr, nil)
			err := d.Handle(wantCtx, wantConn)
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{factory.Mock, wantTransport.Mock, handler.Mock},
	}
}

func SendServerInfoErrorTestCase() ServerDelegateTestCase[any] {
	name := "If send ServerInfo fails with an error, Handle should return it"

	var (
		info      = dlgt.ServerInfo([]byte{1, 2, 3})
		wantErr   = errors.New("send server info error")
		factory   = dmock.NewServerTransportFactory[any]()
		handler   = dmock.NewServerTransportHandler[any]()
		transport = dmock.NewServerTransport[any]()
	)
	factory.RegisterNew(
		func(conn net.Conn) dlgt.ServerTransport[any] { return transport },
	)
	transport.RegisterSendServerInfo(
		func(gotInfo dlgt.ServerInfo) error { return wantErr },
	).RegisterClose(
		func() error { return errors.New("should be ignored") },
	)
	return ServerDelegateTestCase[any]{
		Name: name,
		Setup: ServerDelegateSetup[any]{
			Info:    info,
			Factory: factory,
			Handler: handler,
		},
		Action: func(t *testing.T, d srv.ServerInfoDelegate[any], initErr error) {
			assertfatal.EqualError(t, initErr, nil)
			err := d.Handle(context.Background(), &net.TCPConn{})
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{factory.Mock, transport.Mock, handler.Mock},
	}
}

func ZeroLenServerInfoTestCase() ServerDelegateTestCase[any] {
	name := "If ServerInfo len is zero, New should return an error"

	return ServerDelegateTestCase[any]{
		Name: name,
		Setup: ServerDelegateSetup[any]{
			Info:    dlgt.ServerInfo([]byte{}),
			Factory: dmock.NewServerTransportFactory[any](),
			Handler: dmock.NewServerTransportHandler[any](),
		},
		Action: func(t *testing.T, d srv.ServerInfoDelegate[any], initErr error) {
			asserterror.EqualError(t, initErr, srv.ErrEmptyInfo)
		},
		Mocks: nil,
	}
}

func TransportHandleErrorTestCase() ServerDelegateTestCase[any] {
	name := "If Transport.Handle fails with an error, Handle should return it"

	var (
		wantCtx       = context.Background()
		conn          = &net.TCPConn{}
		wantInfo      = dlgt.ServerInfo([]byte{1, 2, 3})
		wantErr       = errors.New("transport handle error")
		factory       = dmock.NewServerTransportFactory[any]()
		handler       = dmock.NewServerTransportHandler[any]()
		wantTransport = dmock.NewServerTransport[any]()
	)
	factory.RegisterNew(
		func(gotConn net.Conn) dlgt.ServerTransport[any] { return wantTransport },
	)
	wantTransport.RegisterSendServerInfo(
		func(info dlgt.ServerInfo) error { return nil },
	)
	handler.RegisterHandle(
		func(ctx context.Context, transport dlgt.ServerTransport[any]) error {
			return wantErr
		},
	)
	return ServerDelegateTestCase[any]{
		Name: name,
		Setup: ServerDelegateSetup[any]{
			Info:    wantInfo,
			Factory: factory,
			Handler: handler,
		},
		Action: func(t *testing.T, d srv.ServerInfoDelegate[any], initErr error) {
			assertfatal.EqualError(t, initErr, nil)
			err := d.Handle(wantCtx, conn)
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{factory.Mock, wantTransport.Mock, handler.Mock},
	}
}

func SendServerInfoTransportDeadlineErrorTestCase() ServerDelegateTestCase[any] {
	name := "If Transport.SetSendDeadline fails with an error on ServerInfo send, Handle should return it"

	var (
		ctx       = context.Background()
		conn      = &net.TCPConn{}
		info      = dlgt.ServerInfo([]byte{1, 2, 3})
		wantErr   = errors.New("set send deadline error")
		factory   = dmock.NewServerTransportFactory[any]()
		handler   = dmock.NewServerTransportHandler[any]()
		transport = dmock.NewServerTransport[any]()
	)
	factory.RegisterNew(
		func(gotConn net.Conn) dlgt.ServerTransport[any] { return transport },
	)
	transport.RegisterSetSendDeadline(
		func(deadline time.Time) error { return wantErr },
	).RegisterClose(
		func() error { return nil },
	)
	return ServerDelegateTestCase[any]{
		Name: name,
		Setup: ServerDelegateSetup[any]{
			Info:    info,
			Factory: factory,
			Handler: handler,
			Opts: []srv.SetOption{
				srv.WithServerInfoSendDuration(time.Second),
			},
		},
		Action: func(t *testing.T, d srv.ServerInfoDelegate[any], initErr error) {
			assertfatal.EqualError(t, initErr, nil)
			err := d.Handle(ctx, conn)
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{factory.Mock, transport.Mock, handler.Mock},
	}
}
