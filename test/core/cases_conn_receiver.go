package core

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core/srv"
	"github.com/cmd-stream/cmd-stream-go/test"
	cmock "github.com/cmd-stream/cmd-stream-go/test/mock/core"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

func FirstSetDeadlineErrorTestCase(t *testing.T) ConnReceiverTestCase {
	name := "Should return an error if the firts connection failed to set a deadline"

	var (
		wantErr              = errors.New("SetDeadline error")
		wantFirstConnTimeout = time.Second
		startTime            = time.Now()
		listener             = cmock.NewListener()
	)
	listener.RegisterSetDeadline(func(deadline time.Time) error {
		asserterror.SameTime(t, deadline, startTime.Add(wantFirstConnTimeout), test.TimeDelta)
		return wantErr
	})
	return ConnReceiverTestCase{
		Name: name,
		Setup: ConnReceiverSetup{
			Listener: listener,
			Conns:    make(chan net.Conn),
			Opts: []srv.SetConnReceiverOption{
				srv.WithFirstConnTimeout(wantFirstConnTimeout),
			},
			WantErr: wantErr,
		},
		Mocks: []*mok.Mock{listener.Mock},
	}
}

// -----------------------------------------------------------------------------

func FirstAcceptErrorTestCase(t *testing.T) ConnReceiverTestCase {
	name := "Should return an error if accepting of the first conn failed"

	var (
		wantErr  = errors.New("accept error")
		listener = cmock.NewListener()
	)
	listener.RegisterAccept(func() (net.Conn, error) {
		return nil, wantErr
	})
	return ConnReceiverTestCase{
		Name: name,
		Setup: ConnReceiverSetup{
			Listener: listener,
			Conns:    make(chan net.Conn),
			Opts:     []srv.SetConnReceiverOption{},
			WantErr:  wantErr,
		},
		Mocks: []*mok.Mock{listener.Mock},
	}
}

// -----------------------------------------------------------------------------

func FirstResetDeadlineErrorTestCase(t *testing.T) ConnReceiverTestCase {
	name := "Should return an error if cancelation of the first connection deadline failed"

	var (
		wantErr  = errors.New("set deadline error")
		wantConn = cmock.NewConn().RegisterClose(
			func() (err error) { return nil },
		)
		wantFirstConnTimeout = time.Second
		startTime            = time.Now()
		listener             = cmock.NewListener()
	)

	listener.RegisterSetDeadline(
		func(deadline time.Time) error {
			asserterror.SameTime(t, deadline, startTime.Add(wantFirstConnTimeout), test.TimeDelta)
			return nil
		},
	).RegisterAccept(
		func() (net.Conn, error) {
			return wantConn, nil
		},
	).RegisterSetDeadline(
		func(deadline time.Time) error {
			asserterror.Equal(t, deadline.IsZero(), true, "deadline should be zero")
			return wantErr
		},
	)
	return ConnReceiverTestCase{
		Name: name,
		Setup: ConnReceiverSetup{
			Listener: listener,
			Conns:    make(chan net.Conn, 1),
			Opts: []srv.SetConnReceiverOption{
				srv.WithFirstConnTimeout(wantFirstConnTimeout),
			},
			WantErr: wantErr,
		},
		Mocks: []*mok.Mock{wantConn.Mock, listener.Mock},
	}
}

// -----------------------------------------------------------------------------

func SecondAcceptErrorTestCase(t *testing.T) ConnReceiverTestCase {
	name := "Should return an error if accepting of the second conn failed"

	var (
		wantErr  = errors.New("accept error")
		wantConn = cmock.NewConn().RegisterClose(
			func() (err error) { return nil },
		)
		wantFirstConnTimeout = time.Second
		startTime            = time.Now()
		listener             = cmock.NewListener()
	)
	listener.RegisterSetDeadline(
		func(deadline time.Time) error {
			asserterror.SameTime(t, deadline, startTime.Add(wantFirstConnTimeout), test.TimeDelta)
			return nil
		},
	).RegisterAccept(
		func() (net.Conn, error) {
			return wantConn, nil
		},
	).RegisterSetDeadline(
		func(deadline time.Time) error {
			asserterror.Equal(t, deadline.IsZero(), true, "deadline should be zero")
			return nil
		},
	).RegisterAccept(
		func() (net.Conn, error) {
			return nil, wantErr
		},
	)
	return ConnReceiverTestCase{
		Name: name,
		Setup: ConnReceiverSetup{
			Listener: listener,
			Conns:    make(chan net.Conn, 1),
			Opts: []srv.SetConnReceiverOption{
				srv.WithFirstConnTimeout(wantFirstConnTimeout),
			},
			WantErr: wantErr,
		},
		During: func(t *testing.T, receiver *srv.ConnReceiver, conns chan net.Conn) {
			select {
			case conn := <-conns:
				conn.Close()
			case <-time.After(time.Second):
				t.Error("timeout waiting for the first connection")
			}
		},
		Mocks: []*mok.Mock{wantConn.Mock, listener.Mock},
	}
}

// -----------------------------------------------------------------------------

func RunSeveralConnectionsTestCase(t *testing.T) ConnReceiverTestCase {
	name := "Should be able to accept several connections"

	var (
		wantConn1 = cmock.NewConn().RegisterClose(
			func() (err error) { return nil },
		)
		wantConn2 = cmock.NewConn().RegisterClose(
			func() (err error) { return nil },
		)
		wantConn3 = cmock.NewConn().RegisterClose(
			func() (err error) { return nil },
		)
		listener   = cmock.NewListener()
		stopAccept = make(chan struct{})
	)
	listener.RegisterAccept(
		func() (net.Conn, error) { return wantConn1, nil },
	).RegisterAccept(
		func() (net.Conn, error) { return wantConn2, nil },
	).RegisterAccept(
		func() (net.Conn, error) { return wantConn3, nil },
	).RegisterAccept(
		func() (net.Conn, error) {
			<-stopAccept
			return nil, errors.New("listener closed")
		},
	).RegisterClose(
		func() error {
			close(stopAccept)
			return nil
		},
	)
	return ConnReceiverTestCase{
		Name: name,
		Setup: ConnReceiverSetup{
			Listener: listener,
			Conns:    make(chan net.Conn),
			Opts:     []srv.SetConnReceiverOption{},
			WantErr:  srv.ErrClosed,
		},
		During: func(t *testing.T, receiver *srv.ConnReceiver, conns chan net.Conn) {
			for i := range 3 {
				select {
				case conn, ok := <-conns:
					if !ok {
						t.Errorf("channel closed unexpectedly at connection %v", i+1)
						return
					}
					conn.Close()
				case <-time.After(time.Second):
					t.Errorf("timeout waiting for connection %v", i+1)
					return
				}
			}
			receiver.Stop()
		},
		Mocks: []*mok.Mock{wantConn1.Mock, wantConn2.Mock, wantConn3.Mock, listener.Mock},
	}
}

// -----------------------------------------------------------------------------

func StopWhileAcceptingTestCase(t *testing.T) ConnReceiverTestCase {
	name := "Should be able to close while Listener.Accept"

	var (
		listener   = cmock.NewListener()
		stopAccept = make(chan struct{})
	)
	listener.RegisterAccept(
		func() (net.Conn, error) {
			<-stopAccept
			return nil, errors.New("listener closed")
		},
	).RegisterClose(
		func() error {
			close(stopAccept)
			return nil
		},
	)
	return ConnReceiverTestCase{
		Name: name,
		Setup: ConnReceiverSetup{
			Listener: listener,
			Conns:    make(chan net.Conn),
			Opts:     []srv.SetConnReceiverOption{},
			WantErr:  srv.ErrClosed,
		},
		During: func(t *testing.T, receiver *srv.ConnReceiver, conns chan net.Conn) {
			time.Sleep(test.TimeDelta)
			receiver.Stop()
		},
		Mocks: []*mok.Mock{listener.Mock},
	}
}

// -----------------------------------------------------------------------------

func ShutdownWhileAcceptingTestCase(t *testing.T) ConnReceiverTestCase {
	name := "Should be able to shutdown while Listener.Accept"

	var (
		listener   = cmock.NewListener()
		stopAccept = make(chan struct{})
	)
	listener.RegisterAccept(
		func() (net.Conn, error) {
			<-stopAccept
			return nil, errors.New("listener closed")
		},
	).RegisterClose(
		func() error {
			close(stopAccept)
			return nil
		},
	)
	return ConnReceiverTestCase{
		Name: name,
		Setup: ConnReceiverSetup{
			Listener: listener,
			Conns:    make(chan net.Conn),
			Opts:     []srv.SetConnReceiverOption{},
			WantErr:  nil,
		},
		During: func(t *testing.T, receiver *srv.ConnReceiver, conns chan net.Conn) {
			time.Sleep(test.TimeDelta)
			receiver.Shutdown()
		},
		Mocks: []*mok.Mock{listener.Mock},
	}
}

func StopWhileQueuingTestCase(t *testing.T) ConnReceiverTestCase {
	name := "Should be able to close while queuing a connection"

	var (
		wantConn = cmock.NewConn().RegisterClose(
			func() (err error) { return nil },
		)
		listener = cmock.NewListener()
	)
	listener.RegisterAccept(
		func() (net.Conn, error) { return wantConn, nil },
	).RegisterClose(
		func() error { return nil },
	)
	return ConnReceiverTestCase{
		Name: name,
		Setup: ConnReceiverSetup{
			Listener: listener,
			Conns:    make(chan net.Conn), // unbuffered
			Opts:     []srv.SetConnReceiverOption{},
			WantErr:  srv.ErrClosed,
		},
		During: func(t *testing.T, receiver *srv.ConnReceiver, conns chan net.Conn) {
			time.Sleep(test.TimeDelta)
			receiver.Stop()
		},
		Mocks: []*mok.Mock{wantConn.Mock, listener.Mock},
	}
}

func ShutdownWhileQueuingTestCase(t *testing.T) ConnReceiverTestCase {
	name := "Should be able to shutdown the ConnReceiver while queuing a connection"

	var (
		wantConn = cmock.NewConn().RegisterClose(
			func() (err error) { return nil },
		)
		listener = cmock.NewListener()
	)
	listener.RegisterAccept(
		func() (net.Conn, error) { return wantConn, nil },
	).RegisterClose(
		func() error { return nil },
	)
	return ConnReceiverTestCase{
		Name: name,
		Setup: ConnReceiverSetup{
			Listener: listener,
			Conns:    make(chan net.Conn), // unbuffered
			Opts:     []srv.SetConnReceiverOption{},
			WantErr:  nil,
		},
		During: func(t *testing.T, receiver *srv.ConnReceiver, conns chan net.Conn) {
			time.Sleep(test.TimeDelta)
			receiver.Shutdown()
		},
		Mocks: []*mok.Mock{wantConn.Mock, listener.Mock},
	}
}
