package test

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core/srv"
	"github.com/cmd-stream/cmd-stream-go/test/mock"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

type ConnReceiverTestCase struct {
	Name   string
	Setup  ConnReceiverSetup
	During func(t *testing.T, receiver *srv.ConnReceiver, conns chan net.Conn)
	Mocks  []*mok.Mock
}

type ConnReceiverSetup struct {
	Listener mock.Listener
	Conns    chan net.Conn
	Opts     []srv.SetConnReceiverOption
	WantErr  error
}

func RunConnReceiverTestCase(t *testing.T, tc ConnReceiverTestCase) {
	t.Run(tc.Name, func(t *testing.T) {
		receiver := srv.NewConnReceiver(tc.Setup.Listener, tc.Setup.Conns,
			tc.Setup.Opts...)

		errs := make(chan error, 1)
		go func() {
			if err := receiver.Run(); err != nil {
				errs <- err
			}
			close(errs)
		}()

		if tc.During != nil {
			tc.During(t, receiver, tc.Setup.Conns)
		}

		select {
		case err := <-errs:
			asserterror.EqualError(t, err, tc.Setup.WantErr)
		case <-time.After(time.Second):
			t.Fatal("timeout waiting for ConnReceiver to finish")
		}

		asserterror.EqualDeep(t, mok.CheckCalls(tc.Mocks), mok.EmptyInfomap)
	})
}

type connReceiver struct{}

var ConnReceiver connReceiver

// -----------------------------------------------------------------------------
// Test Cases
// -----------------------------------------------------------------------------

func (connReceiver) FirstSetDeadlineError(t *testing.T) ConnReceiverTestCase {
	name := "Should return an error if the first connection failed to set a deadline"

	var (
		wantErr              = errors.New("SetDeadline error")
		wantFirstConnTimeout = time.Second
		startTime            = time.Now()
		listener             = mock.NewListener()
	)
	listener.RegisterSetDeadline(func(deadline time.Time) error {
		asserterror.SameTime(t, deadline, startTime.Add(wantFirstConnTimeout), TimeDelta)
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

func (connReceiver) FirstAcceptError(t *testing.T) ConnReceiverTestCase {
	name := "Should return an error if accepting of the first conn failed"

	var (
		wantErr  = errors.New("accept error")
		listener = mock.NewListener()
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

func (connReceiver) FirstResetDeadlineError(t *testing.T) ConnReceiverTestCase {
	name := "Should return an error if cancelation of the first connection deadline failed"

	var (
		wantErr  = errors.New("set deadline error")
		wantConn = mock.NewConn().RegisterClose(
			func() (err error) { return nil },
		)
		wantFirstConnTimeout = time.Second
		startTime            = time.Now()
		listener             = mock.NewListener()
	)

	listener.RegisterSetDeadline(
		func(deadline time.Time) error {
			asserterror.SameTime(t, deadline, startTime.Add(wantFirstConnTimeout), TimeDelta)
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

func (connReceiver) SecondAcceptError(t *testing.T) ConnReceiverTestCase {
	name := "Should return an error if accepting of the second conn failed"

	var (
		wantErr  = errors.New("accept error")
		wantConn = mock.NewConn().RegisterClose(
			func() (err error) { return nil },
		)
		wantFirstConnTimeout = time.Second
		startTime            = time.Now()
		listener             = mock.NewListener()
	)
	listener.RegisterSetDeadline(
		func(deadline time.Time) error {
			asserterror.SameTime(t, deadline, startTime.Add(wantFirstConnTimeout), TimeDelta)
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
				err := conn.Close()
				asserterror.EqualError(t, err, nil)
			case <-time.After(time.Second):
				t.Error("timeout waiting for the first connection")
			}
		},
		Mocks: []*mok.Mock{wantConn.Mock, listener.Mock},
	}
}

// -----------------------------------------------------------------------------

func (connReceiver) RunSeveralConnections(t *testing.T) ConnReceiverTestCase {
	name := "Should be able to accept several connections"

	var (
		wantConn1 = mock.NewConn().RegisterClose(
			func() (err error) { return nil },
		)
		wantConn2 = mock.NewConn().RegisterClose(
			func() (err error) { return nil },
		)
		wantConn3 = mock.NewConn().RegisterClose(
			func() (err error) { return nil },
		)
		listener   = mock.NewListener()
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
					err := conn.Close()
					asserterror.EqualError(t, err, nil)
				case <-time.After(time.Second):
					t.Errorf("timeout waiting for connection %v", i+1)
					return
				}
			}
			err := receiver.Stop()
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{wantConn1.Mock, wantConn2.Mock, wantConn3.Mock, listener.Mock},
	}
}

// -----------------------------------------------------------------------------

func (connReceiver) StopWhileAccepting(t *testing.T) ConnReceiverTestCase {
	name := "Should be able to close while Listener.Accept"

	var (
		listener   = mock.NewListener()
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
			time.Sleep(TimeDelta)
			err := receiver.Stop()
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{listener.Mock},
	}
}

// -----------------------------------------------------------------------------

func (connReceiver) ShutdownWhileAccepting(t *testing.T) ConnReceiverTestCase {
	name := "Should be able to shutdown while Listener.Accept"

	var (
		listener   = mock.NewListener()
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
			time.Sleep(TimeDelta)
			err := receiver.Shutdown()
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{listener.Mock},
	}
}

func (connReceiver) StopWhileQueuing(t *testing.T) ConnReceiverTestCase {
	name := "Should be able to close while queuing a connection"

	var (
		wantConn = mock.NewConn().RegisterClose(
			func() (err error) { return nil },
		)
		listener = mock.NewListener()
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
			time.Sleep(TimeDelta)
			err := receiver.Stop()
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{wantConn.Mock, listener.Mock},
	}
}

func (connReceiver) ShutdownWhileQueuing(t *testing.T) ConnReceiverTestCase {
	name := "Should be able to shutdown the ConnReceiver while queuing a connection"

	var (
		wantConn = mock.NewConn().RegisterClose(
			func() (err error) { return nil },
		)
		listener = mock.NewListener()
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
			time.Sleep(TimeDelta)
			err := receiver.Shutdown()
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{wantConn.Mock, listener.Mock},
	}
}
