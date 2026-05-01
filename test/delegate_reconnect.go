package test

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	ccln "github.com/cmd-stream/cmd-stream-go/core/cln"
	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
	"github.com/cmd-stream/cmd-stream-go/delegate/cln"
	"github.com/cmd-stream/cmd-stream-go/test/mock"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

type ReconnectTestCase[T any] struct {
	Name   string
	Setup  ReconnectSetup[T]
	Action func(t *testing.T, d *cln.ReconnectDelegate[T], initErr error)
	Mocks  []*mok.Mock
}

type ReconnectSetup[T any] struct {
	Info    dlgt.ServerInfo
	Factory dlgt.ClientTransportFactory[T]
	Opts    []cln.SetOption
}

func RunReconnectTestCase[T any](t *testing.T, tc ReconnectTestCase[T]) {
	t.Run(tc.Name, func(t *testing.T) {
		d, err := cln.NewReconnect[T](tc.Setup.Info, tc.Setup.Factory,
			tc.Setup.Opts...)

		tc.Action(t, d, err)
		asserterror.EqualDeep(t, mok.CheckCalls(tc.Mocks), mok.EmptyInfomap)
	})
}

type ReconnectDelegate[T any] struct{}

// -----------------------------------------------------------------------------
// Test Cases
// -----------------------------------------------------------------------------

type dummyAddr string

func (a dummyAddr) Network() string { return "tcp" }
func (a dummyAddr) String() string  { return string(a) }

func (ReconnectDelegate[T]) InitSuccess(t *testing.T) ReconnectTestCase[T] {
	name := "Should check ServerInfo"
	var (
		wantInfo  = dlgt.ServerInfo("server-1")
		transport = mock.NewClientTransport[T]()
		factory   = mock.NewClientTransportFactory[T]()
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (dlgt.ServerInfo, error) { return wantInfo, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	)
	factory.RegisterNew(
		func() (dlgt.ClientTransport[T], error) { return transport, nil },
	)
	return ReconnectTestCase[T]{
		Name: name,
		Setup: ReconnectSetup[T]{
			Info:    wantInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, nil)
		},
		Mocks: []*mok.Mock{factory.Mock, transport.Mock},
	}
}

func (ReconnectDelegate[T]) InitWrongInfo(t *testing.T) ReconnectTestCase[T] {
	name := "Should return error if wrong ServerInfo was received"
	var (
		wantInfo        = dlgt.ServerInfo("server-1")
		wrongServerInfo = dlgt.ServerInfo("server-2")
		transport       = mock.NewClientTransport[T]()
		factory         = mock.NewClientTransportFactory[T]()
	)

	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (dlgt.ServerInfo, error) { return wrongServerInfo, nil },
	).RegisterClose(func() error { return nil })
	factory.RegisterNew(
		func() (dlgt.ClientTransport[T], error) { return transport, nil },
	)
	return ReconnectTestCase[T]{
		Name: name,
		Setup: ReconnectSetup[T]{
			Info:    wantInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, cln.ErrServerInfoMismatch)
		},
		Mocks: []*mok.Mock{factory.Mock, transport.Mock},
	}
}

func (ReconnectDelegate[T]) InitFactoryError(t *testing.T) ReconnectTestCase[T] {
	name := "If TransportFactory.New fails with an error, NewReconnect should return it"

	var (
		wantErr = cln.ErrServerInfoMismatch // Any error will do
		factory = mock.NewClientTransportFactory[T]()
	)
	factory.RegisterNew(
		func() (dlgt.ClientTransport[T], error) { return nil, wantErr },
	)
	return ReconnectTestCase[T]{
		Name: name,
		Setup: ReconnectSetup[T]{
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, wantErr)
		},
		Mocks: []*mok.Mock{factory.Mock},
	}
}

func (ReconnectDelegate[T]) InitCheckError(t *testing.T) ReconnectTestCase[T] {
	name := "If ServerInfo check fails with an error, NewReconnect should return it"

	var (
		wantInfo  = dlgt.ServerInfo("server-1")
		wantErr   = errors.New("checkServerInfo error")
		transport = mock.NewClientTransport[T]()
		factory   = mock.NewClientTransportFactory[T]()
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (dlgt.ServerInfo, error) { return nil, wantErr },
	).RegisterClose(func() error { return nil })
	factory.RegisterNew(
		func() (dlgt.ClientTransport[T], error) { return transport, nil },
	)
	return ReconnectTestCase[T]{
		Name: name,
		Setup: ReconnectSetup[T]{
			Info:    wantInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, wantErr)
		},
		Mocks: []*mok.Mock{factory.Mock, transport.Mock},
	}
}

func (ReconnectDelegate[T]) CycleSuccess(t *testing.T) ReconnectTestCase[T] {
	name := "Should reconnect"

	var (
		initTransport = mock.NewClientTransport[T]()
		wantTransport = mock.NewClientTransport[T]()
		factory       = mock.NewClientTransportFactory[T]()
		serverInfo    = dlgt.ServerInfo("server-1")
	)
	initTransport.RegisterSetReceiveDeadline(
		func(deadline time.Time) (err error) { return nil },
	).RegisterReceiveServerInfo(
		func() (i dlgt.ServerInfo, err error) { return serverInfo, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) (err error) { return nil },
	)
	wantTransport.RegisterSetReceiveDeadline(
		func(deadline time.Time) (err error) { return nil },
	).RegisterReceiveServerInfo(
		func() (i dlgt.ServerInfo, err error) { return serverInfo, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) (err error) { return nil },
	)
	factory.RegisterNew(
		func() (dlgt.ClientTransport[T], error) { return initTransport, nil },
	).RegisterNew(
		func() (dlgt.ClientTransport[T], error) { return nil, errors.New("transport creation error") },
	).RegisterNew(
		func() (dlgt.ClientTransport[T], error) { return nil, errors.New("transport creation error") },
	).RegisterNew(
		func() (dlgt.ClientTransport[T], error) { return wantTransport, nil },
	)
	return ReconnectTestCase[T]{
		Name: name,
		Setup: ReconnectSetup[T]{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			err := d.Reconnect()
			asserterror.EqualError(t, err, nil)
			asserterror.Equal(t, d.Transport(), (dlgt.ClientTransport[T])(wantTransport))
		},
		Mocks: []*mok.Mock{initTransport.Mock, wantTransport.Mock, factory.Mock},
	}
}

func (ReconnectDelegate[T]) CycleClose(t *testing.T) ReconnectTestCase[T] {
	name := "Reconnect should return ErrClosed, if the delegate is closed"

	var (
		initTransport = mock.NewClientTransport[T]()
		factory       = mock.NewClientTransportFactory[T]()
		serverInfo    = dlgt.ServerInfo("server-1")
	)
	initTransport.RegisterSetReceiveDeadline(
		func(deadline time.Time) (err error) { return nil },
	).RegisterReceiveServerInfo(
		func() (i dlgt.ServerInfo, err error) { return serverInfo, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) (err error) { return nil },
	).RegisterClose(
		func() error { return nil },
	)
	factory.RegisterNew(
		func() (dlgt.ClientTransport[T], error) {
			return initTransport, nil
		},
	).RegisterNew(
		func() (dlgt.ClientTransport[T], error) {
			time.Sleep(100 * time.Millisecond)
			return nil, errors.New("transport creation error")
		},
	)
	return ReconnectTestCase[T]{
		Name: name,
		Setup: ReconnectSetup[T]{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			go func() {
				time.Sleep(50 * time.Millisecond)
				err := d.Close()
				asserterror.EqualError(t, err, nil)
			}()
			err := d.Reconnect()
			asserterror.EqualError(t, err, ccln.ErrClosed)
			asserterror.Equal(t, d.Transport(), (dlgt.ClientTransport[T])(initTransport))
		},
		Mocks: []*mok.Mock{initTransport.Mock, factory.Mock},
	}
}

func (ReconnectDelegate[T]) CycleMismatch(t *testing.T) ReconnectTestCase[T] {
	name := "If ServerInfo check fails with the ErrServerInfoMismatch, Reconnect should return it"

	var (
		initTransport = mock.NewClientTransport[T]()
		newTransport  = mock.NewClientTransport[T]()
		factory       = mock.NewClientTransportFactory[T]()
		serverInfo    = dlgt.ServerInfo("server-1")
		wrongInfo     = dlgt.ServerInfo("wrong-server")
	)
	initTransport.RegisterSetReceiveDeadline(
		func(deadline time.Time) (err error) { return nil },
	).RegisterReceiveServerInfo(
		func() (i dlgt.ServerInfo, err error) { return serverInfo, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) (err error) { return nil },
	)
	newTransport.RegisterSetReceiveDeadline(
		func(deadline time.Time) (err error) { return nil },
	).RegisterReceiveServerInfo(
		func() (i dlgt.ServerInfo, err error) { return wrongInfo, nil },
	).RegisterClose(
		func() error { return nil },
	)
	factory.RegisterNew(
		func() (dlgt.ClientTransport[T], error) {
			return initTransport, nil
		},
	).RegisterNew(
		func() (dlgt.ClientTransport[T], error) {
			return newTransport, nil
		},
	)
	return ReconnectTestCase[T]{
		Name: name,
		Setup: ReconnectSetup[T]{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			err := d.Reconnect()
			asserterror.EqualError(t, err, cln.ErrServerInfoMismatch)
			asserterror.Equal(t, d.Transport(), (dlgt.ClientTransport[T])(initTransport))
		},
		Mocks: []*mok.Mock{initTransport.Mock, newTransport.Mock, factory.Mock},
	}
}

func (ReconnectDelegate[T]) CycleCheckError(t *testing.T) ReconnectTestCase[T] {
	name := "If ServerInfo check fails with an error, Reconnect should try again"

	var (
		initTransport = mock.NewClientTransport[T]()
		failTransport = mock.NewClientTransport[T]()
		succTransport = mock.NewClientTransport[T]()
		factory       = mock.NewClientTransportFactory[T]()
		serverInfo    = dlgt.ServerInfo("server-1")
	)
	initTransport.RegisterSetReceiveDeadline(
		func(deadline time.Time) (err error) { return nil },
	).RegisterReceiveServerInfo(
		func() (i dlgt.ServerInfo, err error) { return serverInfo, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) (err error) { return nil },
	)
	failTransport.RegisterSetReceiveDeadline(
		func(deadline time.Time) (err error) { return nil },
	).RegisterReceiveServerInfo(
		func() (i dlgt.ServerInfo, err error) { return dlgt.ServerInfo(""), errors.New("check error") },
	).RegisterClose(
		func() error { return nil },
	)
	succTransport.RegisterSetReceiveDeadline(
		func(deadline time.Time) (err error) { return nil },
	).RegisterReceiveServerInfo(
		func() (i dlgt.ServerInfo, err error) { return serverInfo, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) (err error) { return nil },
	)
	factory.RegisterNew(
		func() (dlgt.ClientTransport[T], error) {
			return initTransport, nil
		},
	).RegisterNew(
		func() (dlgt.ClientTransport[T], error) {
			return failTransport, nil
		},
	).RegisterNew(
		func() (dlgt.ClientTransport[T], error) {
			return succTransport, nil
		},
	)
	return ReconnectTestCase[T]{
		Name: name,
		Setup: ReconnectSetup[T]{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			err := d.Reconnect()
			asserterror.EqualError(t, err, nil)
			asserterror.Equal(t, d.Transport(), (dlgt.ClientTransport[T])(succTransport))
		},
		Mocks: []*mok.Mock{initTransport.Mock, failTransport.Mock, succTransport.Mock, factory.Mock},
	}
}

func (ReconnectDelegate[T]) SetSendDeadline(t *testing.T) ReconnectTestCase[T] {
	name := "SetSendDeadline should call corresponding Transport.SetSendDeadline"

	var (
		wantDeadline = time.Now().Add(time.Second)
		transport    = mock.NewClientTransport[T]()
		factory      = mock.NewClientTransportFactory[T]()
		serverInfo   = dlgt.ServerInfo("server-1")
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (dlgt.ServerInfo, error) { return serverInfo, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterSetSendDeadline(func(deadline time.Time) error {
		if !deadline.Equal(wantDeadline) {
			return errors.New("wrong deadline")
		}
		return nil
	})
	factory.RegisterNew(
		func() (dlgt.ClientTransport[T], error) { return transport, nil },
	)
	return ReconnectTestCase[T]{
		Name: name,
		Setup: ReconnectSetup[T]{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			err := d.SetSendDeadline(wantDeadline)
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{transport.Mock, factory.Mock},
	}
}

func (ReconnectDelegate[T]) CloseError(t *testing.T) ReconnectTestCase[T] {
	name := "If Transport.Close fails with an error, Close should return it"

	var (
		wantErr    = errors.New("close error")
		transport  = mock.NewClientTransport[T]()
		factory    = mock.NewClientTransportFactory[T]()
		serverInfo = dlgt.ServerInfo("server-1")
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (dlgt.ServerInfo, error) { return serverInfo, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterClose(func() error {
		return wantErr
	})
	factory.RegisterNew(
		func() (dlgt.ClientTransport[T], error) { return transport, nil },
	)
	return ReconnectTestCase[T]{
		Name: name,
		Setup: ReconnectSetup[T]{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			err := d.Close()
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{transport.Mock, factory.Mock},
	}
}

func (ReconnectDelegate[T]) Flush(t *testing.T) ReconnectTestCase[T] {
	name := "Flush should call corresponding Transport.Flush"

	var (
		transport  = mock.NewClientTransport[T]()
		factory    = mock.NewClientTransportFactory[T]()
		serverInfo = dlgt.ServerInfo("server-1")
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (dlgt.ServerInfo, error) { return serverInfo, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterFlush(func() error {
		return nil
	})
	factory.RegisterNew(
		func() (dlgt.ClientTransport[T], error) { return transport, nil },
	)
	return ReconnectTestCase[T]{
		Name: name,
		Setup: ReconnectSetup[T]{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			err := d.Flush()
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{transport.Mock, factory.Mock},
	}
}

func (ReconnectDelegate[T]) LocalAddr(t *testing.T) ReconnectTestCase[T] {
	name := "LocalAddr should return Transport.LocalAddr"

	var (
		wantAddr   = dummyAddr("127.0.0.1:1234")
		transport  = mock.NewClientTransport[T]()
		factory    = mock.NewClientTransportFactory[T]()
		serverInfo = dlgt.ServerInfo("server-1")
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (dlgt.ServerInfo, error) { return serverInfo, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterLocalAddr(func() net.Addr {
		return wantAddr
	})
	factory.RegisterNew(
		func() (dlgt.ClientTransport[T], error) { return transport, nil },
	)
	return ReconnectTestCase[T]{
		Name: name,
		Setup: ReconnectSetup[T]{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			addr := d.LocalAddr()
			if addr != wantAddr {
				t.Errorf("wrong addr: want %v, got %v", wantAddr, addr)
			}
		},
		Mocks: []*mok.Mock{transport.Mock, factory.Mock},
	}
}

func (ReconnectDelegate[T]) RemoteAddr(t *testing.T) ReconnectTestCase[T] {
	name := "RemoteAddr should return Transport.RemoteAddr"

	var (
		wantAddr   = dummyAddr("127.0.0.1:5678")
		transport  = mock.NewClientTransport[T]()
		factory    = mock.NewClientTransportFactory[T]()
		serverInfo = dlgt.ServerInfo("server-1")
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (dlgt.ServerInfo, error) { return serverInfo, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterRemoteAddr(
		func() net.Addr { return wantAddr },
	)
	factory.RegisterNew(
		func() (dlgt.ClientTransport[T], error) { return transport, nil },
	)
	return ReconnectTestCase[T]{
		Name: name,
		Setup: ReconnectSetup[T]{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			addr := d.RemoteAddr()
			asserterror.EqualDeep(t, addr, wantAddr)
		},
		Mocks: []*mok.Mock{transport.Mock, factory.Mock},
	}
}

func (ReconnectDelegate[T]) Send(t *testing.T) ReconnectTestCase[T] {
	name := "Transport.Send should send same seq and cmd as Send"

	var (
		wantSeq    = core.Seq(123)
		wantCmd    core.Cmd[T]
		transport  = mock.NewClientTransport[T]()
		factory    = mock.NewClientTransportFactory[T]()
		serverInfo = dlgt.ServerInfo("server-1")
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (dlgt.ServerInfo, error) { return serverInfo, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterSend(func(seq core.Seq, cmd core.Cmd[T]) (int, error) {
		if seq != wantSeq {
			return 0, errors.New("wrong seq")
		}
		if cmd != wantCmd {
			return 0, errors.New("wrong cmd")
		}
		return 10, nil
	})
	factory.RegisterNew(
		func() (dlgt.ClientTransport[T], error) { return transport, nil },
	)
	return ReconnectTestCase[T]{
		Name: name,
		Setup: ReconnectSetup[T]{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			n, err := d.Send(wantSeq, wantCmd)
			asserterror.EqualError(t, err, nil)
			asserterror.Equal(t, n, 10)
		},
		Mocks: []*mok.Mock{transport.Mock, factory.Mock},
	}
}

func (ReconnectDelegate[T]) SendError(t *testing.T) ReconnectTestCase[T] {
	name := "If Transport.Send fails with an error, Send should return it"

	var (
		wantErr    = errors.New("send error")
		transport  = mock.NewClientTransport[T]()
		factory    = mock.NewClientTransportFactory[T]()
		serverInfo = dlgt.ServerInfo("server-1")
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (dlgt.ServerInfo, error) { return serverInfo, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterSend(func(seq core.Seq, cmd core.Cmd[T]) (int, error) {
		return 0, wantErr
	})
	factory.RegisterNew(
		func() (dlgt.ClientTransport[T], error) { return transport, nil },
	)
	return ReconnectTestCase[T]{
		Name: name,
		Setup: ReconnectSetup[T]{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			_, err := d.Send(0, nil)
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{transport.Mock, factory.Mock},
	}
}

func (ReconnectDelegate[T]) Receive(t *testing.T) ReconnectTestCase[T] {
	name := "Receive should return same seq and result as Transport.Receive"

	var (
		wantSeq    = core.Seq(456)
		wantResult = (core.Result)(nil)
		wantN      = 20
		transport  = mock.NewClientTransport[T]()
		factory    = mock.NewClientTransportFactory[T]()
		serverInfo = dlgt.ServerInfo("server-1")
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (dlgt.ServerInfo, error) { return serverInfo, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceive(func() (core.Seq, core.Result, int, error) {
		return wantSeq, wantResult, wantN, nil
	})
	factory.RegisterNew(
		func() (dlgt.ClientTransport[T], error) { return transport, nil },
	)
	return ReconnectTestCase[T]{
		Name: name,
		Setup: ReconnectSetup[T]{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			seq, res, n, err := d.Receive()
			asserterror.EqualError(t, err, nil)
			asserterror.Equal(t, seq, wantSeq)
			asserterror.Equal(t, res, wantResult)
			asserterror.Equal(t, n, wantN)
		},
		Mocks: []*mok.Mock{transport.Mock, factory.Mock},
	}
}

func (ReconnectDelegate[T]) ReceiveError(t *testing.T) ReconnectTestCase[T] {
	name := "If Transport.Receive fails with an error, Receive should return it"

	var (
		wantErr    = errors.New("receive error")
		transport  = mock.NewClientTransport[T]()
		factory    = mock.NewClientTransportFactory[T]()
		serverInfo = dlgt.ServerInfo("server-1")
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (dlgt.ServerInfo, error) { return serverInfo, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceive(func() (core.Seq, core.Result, int, error) {
		return 0, nil, 0, wantErr
	})
	factory.RegisterNew(
		func() (dlgt.ClientTransport[T], error) { return transport, nil },
	)
	return ReconnectTestCase[T]{
		Name: name,
		Setup: ReconnectSetup[T]{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			_, _, _, err := d.Receive()
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{transport.Mock, factory.Mock},
	}
}

func (ReconnectDelegate[T]) SetReceiveDeadline(t *testing.T) ReconnectTestCase[T] {
	name := "SetReceiveDeadline should call corresponding Transport.SetReceiveDeadline"

	var (
		wantDeadline = time.Now().Add(time.Second)
		transport    = mock.NewClientTransport[T]()
		factory      = mock.NewClientTransportFactory[T]()
		serverInfo   = dlgt.ServerInfo("server-1")
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (dlgt.ServerInfo, error) { return serverInfo, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) error {
			if !deadline.Equal(wantDeadline) {
				return errors.New("wrong deadline")
			}
			return nil
		},
	)
	factory.RegisterNew(
		func() (dlgt.ClientTransport[T], error) { return transport, nil },
	)
	return ReconnectTestCase[T]{
		Name: name,
		Setup: ReconnectSetup[T]{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			err := d.SetReceiveDeadline(wantDeadline)
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{transport.Mock, factory.Mock},
	}
}

func (ReconnectDelegate[T]) SetReceiveDeadlineError(t *testing.T) ReconnectTestCase[T] {
	name := "If Transport.SetReceiveDeadline fails with an error, SetReceiveDeadline should return it"

	var (
		wantErr    = errors.New("SetReceiveDeadline error")
		transport  = mock.NewClientTransport[T]()
		factory    = mock.NewClientTransportFactory[T]()
		serverInfo = dlgt.ServerInfo("server-1")
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (dlgt.ServerInfo, error) { return serverInfo, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterSetReceiveDeadline(func(deadline time.Time) error {
		return wantErr
	})
	factory.RegisterNew(
		func() (dlgt.ClientTransport[T], error) { return transport, nil },
	)
	return ReconnectTestCase[T]{
		Name: name,
		Setup: ReconnectSetup[T]{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[T], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			err := d.SetReceiveDeadline(time.Time{})
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{transport.Mock, factory.Mock},
	}
}
