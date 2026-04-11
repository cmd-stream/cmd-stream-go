package delegate

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	ccln "github.com/cmd-stream/cmd-stream-go/core/cln"
	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
	"github.com/cmd-stream/cmd-stream-go/delegate/cln"
	mockcln "github.com/cmd-stream/cmd-stream-go/test/mock/delegate"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

type dummyAddr string

func (a dummyAddr) Network() string { return "tcp" }
func (a dummyAddr) String() string  { return string(a) }

func ReconnectInitSuccessTestCase() ReconnectTestCase {
	name := "Should check ServerInfo"
	var (
		wantInfo  = dlgt.ServerInfo("server-1")
		transport = mockcln.NewClientTransport[any]()
		factory   = mockcln.NewClientTransportFactory[any]()
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (dlgt.ServerInfo, error) { return wantInfo, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	)
	factory.RegisterNew(
		func() (dlgt.ClientTransport[any], error) { return transport, nil },
	)
	return ReconnectTestCase{
		Name: name,
		Setup: ReconnectSetup{
			Info:    wantInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[any], initErr error) {
			asserterror.EqualError(t, initErr, nil)
		},
		Mocks: []*mok.Mock{factory.Mock, transport.Mock},
	}
}

func ReconnectInitWrongInfoTestCase() ReconnectTestCase {
	name := "Should return error if wrong ServerInfo was received"
	var (
		wantInfo        = dlgt.ServerInfo("server-1")
		wrongServerInfo = dlgt.ServerInfo("server-2")
		transport       = mockcln.NewClientTransport[any]()
		factory         = mockcln.NewClientTransportFactory[any]()
	)

	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (dlgt.ServerInfo, error) { return wrongServerInfo, nil },
	).RegisterClose(func() error { return nil })
	factory.RegisterNew(
		func() (dlgt.ClientTransport[any], error) { return transport, nil },
	)
	return ReconnectTestCase{
		Name: name,
		Setup: ReconnectSetup{
			Info:    wantInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[any], initErr error) {
			asserterror.EqualError(t, initErr, cln.ErrServerInfoMismatch)
		},
		Mocks: []*mok.Mock{factory.Mock, transport.Mock},
	}
}

func ReconnectInitFactoryErrorTestCase() ReconnectTestCase {
	name := "If TransportFactory.New fails with an error, NewReconnect should return it"

	var (
		wantErr = cln.ErrServerInfoMismatch // Any error will do
		factory = mockcln.NewClientTransportFactory[any]()
	)
	factory.RegisterNew(
		func() (dlgt.ClientTransport[any], error) { return nil, wantErr },
	)
	return ReconnectTestCase{
		Name: name,
		Setup: ReconnectSetup{
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[any], initErr error) {
			asserterror.EqualError(t, initErr, wantErr)
		},
		Mocks: []*mok.Mock{factory.Mock},
	}
}

func ReconnectInitCheckErrorTestCase() ReconnectTestCase {
	name := "If ServerInfo check fails with an error, NewReconnect should return it"

	var (
		wantInfo  = dlgt.ServerInfo("server-1")
		wantErr   = errors.New("checkServerInfo error")
		transport = mockcln.NewClientTransport[any]()
		factory   = mockcln.NewClientTransportFactory[any]()
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (dlgt.ServerInfo, error) { return nil, wantErr },
	).RegisterClose(func() error { return nil })
	factory.RegisterNew(
		func() (dlgt.ClientTransport[any], error) { return transport, nil },
	)
	return ReconnectTestCase{
		Name: name,
		Setup: ReconnectSetup{
			Info:    wantInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[any], initErr error) {
			asserterror.EqualError(t, initErr, wantErr)
		},
		Mocks: []*mok.Mock{factory.Mock, transport.Mock},
	}
}

func ReconnectCycleSuccessTestCase() ReconnectTestCase {
	name := "Should reconnect"

	var (
		initTransport = mockcln.NewClientTransport[any]()
		wantTransport = mockcln.NewClientTransport[any]()
		factory       = mockcln.NewClientTransportFactory[any]()
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
		func() (dlgt.ClientTransport[any], error) { return initTransport, nil },
	).RegisterNew(
		func() (dlgt.ClientTransport[any], error) { return nil, errors.New("transport creation error") },
	).RegisterNew(
		func() (dlgt.ClientTransport[any], error) { return nil, errors.New("transport creation error") },
	).RegisterNew(
		func() (dlgt.ClientTransport[any], error) { return wantTransport, nil },
	)
	return ReconnectTestCase{
		Name: name,
		Setup: ReconnectSetup{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[any], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			err := d.Reconnect()
			asserterror.EqualError(t, err, nil)
			asserterror.Equal(t, d.Transport(), (dlgt.ClientTransport[any])(wantTransport))
		},
		Mocks: []*mok.Mock{initTransport.Mock, wantTransport.Mock, factory.Mock},
	}
}

func ReconnectCycleCloseTestCase() ReconnectTestCase {
	name := "Reconnect should return ErrClosed, if the delegate is closed"

	var (
		initTransport = mockcln.NewClientTransport[any]()
		factory       = mockcln.NewClientTransportFactory[any]()
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
		func() (dlgt.ClientTransport[any], error) {
			return initTransport, nil
		},
	).RegisterNew(
		func() (dlgt.ClientTransport[any], error) {
			time.Sleep(100 * time.Millisecond)
			return nil, errors.New("transport creation error")
		},
	)
	return ReconnectTestCase{
		Name: name,
		Setup: ReconnectSetup{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[any], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			go func() {
				time.Sleep(50 * time.Millisecond)
				err := d.Close()
				asserterror.EqualError(t, err, nil)
			}()
			err := d.Reconnect()
			asserterror.EqualError(t, err, ccln.ErrClosed)
			asserterror.Equal(t, d.Transport(), (dlgt.ClientTransport[any])(initTransport))
		},
		Mocks: []*mok.Mock{initTransport.Mock, factory.Mock},
	}
}

func ReconnectCycleMismatchTestCase() ReconnectTestCase {
	name := "If ServerInfo check fails with the ErrServerInfoMismatch, Reconnect should return it"

	var (
		initTransport = mockcln.NewClientTransport[any]()
		newTransport  = mockcln.NewClientTransport[any]()
		factory       = mockcln.NewClientTransportFactory[any]()
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
		func() (dlgt.ClientTransport[any], error) {
			return initTransport, nil
		},
	).RegisterNew(
		func() (dlgt.ClientTransport[any], error) {
			return newTransport, nil
		},
	)
	return ReconnectTestCase{
		Name: name,
		Setup: ReconnectSetup{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[any], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			err := d.Reconnect()
			asserterror.EqualError(t, err, cln.ErrServerInfoMismatch)
			asserterror.Equal(t, d.Transport(), (dlgt.ClientTransport[any])(initTransport))
		},
		Mocks: []*mok.Mock{initTransport.Mock, newTransport.Mock, factory.Mock},
	}
}

func ReconnectCycleCheckErrorTestCase() ReconnectTestCase {
	name := "If ServerInfo check fails with an error, Reconnect should try again"

	var (
		initTransport = mockcln.NewClientTransport[any]()
		failTransport = mockcln.NewClientTransport[any]()
		succTransport = mockcln.NewClientTransport[any]()
		factory       = mockcln.NewClientTransportFactory[any]()
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
		func() (dlgt.ClientTransport[any], error) {
			return initTransport, nil
		},
	).RegisterNew(
		func() (dlgt.ClientTransport[any], error) {
			return failTransport, nil
		},
	).RegisterNew(
		func() (dlgt.ClientTransport[any], error) {
			return succTransport, nil
		},
	)
	return ReconnectTestCase{
		Name: name,
		Setup: ReconnectSetup{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[any], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			err := d.Reconnect()
			asserterror.EqualError(t, err, nil)
			asserterror.Equal(t, d.Transport(), (dlgt.ClientTransport[any])(succTransport))
		},
		Mocks: []*mok.Mock{initTransport.Mock, failTransport.Mock, succTransport.Mock, factory.Mock},
	}
}

func ReconnectSetSendDeadlineTestCase() ReconnectTestCase {
	name := "SetSendDeadline should call corresponding Transport.SetSendDeadline"

	var (
		wantDeadline = time.Now().Add(time.Second)
		transport    = mockcln.NewClientTransport[any]()
		factory      = mockcln.NewClientTransportFactory[any]()
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
		func() (dlgt.ClientTransport[any], error) { return transport, nil },
	)
	return ReconnectTestCase{
		Name: name,
		Setup: ReconnectSetup{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[any], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			err := d.SetSendDeadline(wantDeadline)
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{transport.Mock, factory.Mock},
	}
}

func ReconnectCloseErrorTestCase() ReconnectTestCase {
	name := "If Transport.Close fails with an error, Close should return it"

	var (
		wantErr    = errors.New("close error")
		transport  = mockcln.NewClientTransport[any]()
		factory    = mockcln.NewClientTransportFactory[any]()
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
		func() (dlgt.ClientTransport[any], error) { return transport, nil },
	)
	return ReconnectTestCase{
		Name: name,
		Setup: ReconnectSetup{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[any], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			err := d.Close()
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{transport.Mock, factory.Mock},
	}
}

func ReconnectFlushTestCase() ReconnectTestCase {
	name := "Flush should call corresponding Transport.Flush"

	var (
		transport  = mockcln.NewClientTransport[any]()
		factory    = mockcln.NewClientTransportFactory[any]()
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
		func() (dlgt.ClientTransport[any], error) { return transport, nil },
	)
	return ReconnectTestCase{
		Name: name,
		Setup: ReconnectSetup{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[any], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			err := d.Flush()
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{transport.Mock, factory.Mock},
	}
}

func ReconnectLocalAddrTestCase() ReconnectTestCase {
	name := "LocalAddr should return Transport.LocalAddr"

	var (
		wantAddr   = dummyAddr("127.0.0.1:1234")
		transport  = mockcln.NewClientTransport[any]()
		factory    = mockcln.NewClientTransportFactory[any]()
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
		func() (dlgt.ClientTransport[any], error) { return transport, nil },
	)
	return ReconnectTestCase{
		Name: name,
		Setup: ReconnectSetup{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[any], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			addr := d.LocalAddr()
			if addr != wantAddr {
				t.Errorf("wrong addr: want %v, got %v", wantAddr, addr)
			}
		},
		Mocks: []*mok.Mock{transport.Mock, factory.Mock},
	}
}

func ReconnectRemoteAddrTestCase() ReconnectTestCase {
	name := "RemoteAddr should return Transport.RemoteAddr"

	var (
		wantAddr   = dummyAddr("127.0.0.1:5678")
		transport  = mockcln.NewClientTransport[any]()
		factory    = mockcln.NewClientTransportFactory[any]()
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
		func() (dlgt.ClientTransport[any], error) { return transport, nil },
	)
	return ReconnectTestCase{
		Name: name,
		Setup: ReconnectSetup{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[any], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			addr := d.RemoteAddr()
			asserterror.EqualDeep(t, addr, wantAddr)
		},
		Mocks: []*mok.Mock{transport.Mock, factory.Mock},
	}
}

func ReconnectSendTestCase() ReconnectTestCase {
	name := "Transport.Send should send same seq and cmd as Send"

	var (
		wantSeq    = core.Seq(123)
		wantCmd    = (core.Cmd[any])(nil)
		transport  = mockcln.NewClientTransport[any]()
		factory    = mockcln.NewClientTransportFactory[any]()
		serverInfo = dlgt.ServerInfo("server-1")
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (dlgt.ServerInfo, error) { return serverInfo, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterSend(func(seq core.Seq, cmd core.Cmd[any]) (int, error) {
		if seq != wantSeq {
			return 0, errors.New("wrong seq")
		}
		if cmd != wantCmd {
			return 0, errors.New("wrong cmd")
		}
		return 10, nil
	})
	factory.RegisterNew(
		func() (dlgt.ClientTransport[any], error) { return transport, nil },
	)
	return ReconnectTestCase{
		Name: name,
		Setup: ReconnectSetup{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[any], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			n, err := d.Send(wantSeq, wantCmd)
			asserterror.EqualError(t, err, nil)
			asserterror.Equal(t, n, 10)
		},
		Mocks: []*mok.Mock{transport.Mock, factory.Mock},
	}
}

func ReconnectSendErrorTestCase() ReconnectTestCase {
	name := "If Transport.Send fails with an error, Send should return it"

	var (
		wantErr    = errors.New("send error")
		transport  = mockcln.NewClientTransport[any]()
		factory    = mockcln.NewClientTransportFactory[any]()
		serverInfo = dlgt.ServerInfo("server-1")
	)
	transport.RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterReceiveServerInfo(
		func() (dlgt.ServerInfo, error) { return serverInfo, nil },
	).RegisterSetReceiveDeadline(
		func(deadline time.Time) error { return nil },
	).RegisterSend(func(seq core.Seq, cmd core.Cmd[any]) (int, error) {
		return 0, wantErr
	})
	factory.RegisterNew(
		func() (dlgt.ClientTransport[any], error) { return transport, nil },
	)
	return ReconnectTestCase{
		Name: name,
		Setup: ReconnectSetup{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[any], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			_, err := d.Send(0, nil)
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{transport.Mock, factory.Mock},
	}
}

func ReconnectReceiveTestCase() ReconnectTestCase {
	name := "Receive should return same seq and result as Transport.Receive"

	var (
		wantSeq    = core.Seq(456)
		wantResult = (core.Result)(nil)
		wantN      = 20
		transport  = mockcln.NewClientTransport[any]()
		factory    = mockcln.NewClientTransportFactory[any]()
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
		func() (dlgt.ClientTransport[any], error) { return transport, nil },
	)
	return ReconnectTestCase{
		Name: name,
		Setup: ReconnectSetup{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[any], initErr error) {
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

func ReconnectReceiveErrorTestCase() ReconnectTestCase {
	name := "If Transport.Receive fails with an error, Receive should return it"

	var (
		wantErr    = errors.New("receive error")
		transport  = mockcln.NewClientTransport[any]()
		factory    = mockcln.NewClientTransportFactory[any]()
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
		func() (dlgt.ClientTransport[any], error) { return transport, nil },
	)
	return ReconnectTestCase{
		Name: name,
		Setup: ReconnectSetup{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[any], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			_, _, _, err := d.Receive()
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{transport.Mock, factory.Mock},
	}
}

func ReconnectSetReceiveDeadlineTestCase() ReconnectTestCase {
	name := "SetReceiveDeadline should call corresponding Transport.SetReceiveDeadline"

	var (
		wantDeadline = time.Now().Add(time.Second)
		transport    = mockcln.NewClientTransport[any]()
		factory      = mockcln.NewClientTransportFactory[any]()
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
		func() (dlgt.ClientTransport[any], error) { return transport, nil },
	)
	return ReconnectTestCase{
		Name: name,
		Setup: ReconnectSetup{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[any], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			err := d.SetReceiveDeadline(wantDeadline)
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{transport.Mock, factory.Mock},
	}
}

func ReconnectSetReceiveDeadlineErrorTestCase() ReconnectTestCase {
	name := "If Transport.SetReceiveDeadline fails with an error, SetReceiveDeadline should return it"

	var (
		wantErr    = errors.New("SetReceiveDeadline error")
		transport  = mockcln.NewClientTransport[any]()
		factory    = mockcln.NewClientTransportFactory[any]()
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
		func() (dlgt.ClientTransport[any], error) { return transport, nil },
	)
	return ReconnectTestCase{
		Name: name,
		Setup: ReconnectSetup{
			Info:    serverInfo,
			Factory: factory,
		},
		Action: func(t *testing.T, d *cln.ReconnectDelegate[any], initErr error) {
			asserterror.EqualError(t, initErr, nil)
			err := d.SetReceiveDeadline(time.Time{})
			asserterror.EqualError(t, err, wantErr)
		},
		Mocks: []*mok.Mock{transport.Mock, factory.Mock},
	}
}
