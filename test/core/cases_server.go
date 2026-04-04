package core

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core/srv"
	cmock "github.com/cmd-stream/cmd-stream-go/test/mock/core"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

func ServeNoWorkersTestCase(t *testing.T) ServerTestCase {
	name := "Should return ErrNoWorkers if WorkersCount is 0"

	var (
		listener = cmock.NewListener()
		delegate = cmock.NewServerDelegate()
	)
	return ServerTestCase{
		Name: name,
		Setup: ServerSetup{
			Delegate: delegate,
			Opts: []srv.SetOption{
				srv.WithWorkersCount(0),
			},
			WantErr: srv.NewServerError(srv.ErrNoWorkers),
		},
		Mocks: []*mok.Mock{delegate.Mock, listener.Mock},
	}
}

func ServeSeveralConnectionsTestCase(t *testing.T) ServerTestCase {
	name := "Should be able to serve several connections"

	var (
		wantConn1  = cmock.NewConn()
		wantConn2  = cmock.NewConn()
		wantConn3  = cmock.NewConn()
		listener   = cmock.NewListener()
		delegate   = cmock.NewServerDelegate()
		stopAccept = make(chan struct{})
	)
	listener.RegisterAccept(
		func() (net.Conn, error) {
			return wantConn1, nil
		},
	).RegisterAccept(
		func() (net.Conn, error) {
			return wantConn2, nil
		},
	).RegisterAccept(
		func() (net.Conn, error) {
			return wantConn3, nil
		},
	).RegisterAccept(
		func() (net.Conn, error) {
			<-stopAccept
			return nil, errors.New("listener closed")
		},
	).RegisterClose(func() error {
		close(stopAccept)
		return nil
	})
	delegate.RegisterHandle(
		func(ctx context.Context, conn net.Conn) error {
			asserterror.Equal(t, conn, net.Conn(wantConn1))
			return nil
		},
	).RegisterHandle(
		func(ctx context.Context, conn net.Conn) error {
			asserterror.Equal(t, conn, net.Conn(wantConn2))
			return nil
		},
	).RegisterHandle(
		func(ctx context.Context, conn net.Conn) error {
			asserterror.Equal(t, conn, net.Conn(wantConn3))
			return nil
		},
	)
	return ServerTestCase{
		Name: name,
		Setup: ServerSetup{
			Delegate: delegate,
			Opts: []srv.SetOption{
				srv.WithWorkersCount(1),
			},
			WantErr: srv.NewServerError(srv.ErrShutdown),
		},
		Action: func(t *testing.T, s *srv.Server) {
			errs := make(chan error, 1)
			go func() {
				errs <- s.Serve(listener)
				close(errs)
			}()

			time.Sleep(10 * time.Millisecond)
			err := s.Shutdown()
			asserterror.EqualError(t, err, nil)

			select {
			case err := <-errs:
				asserterror.EqualError(t, err, srv.NewServerError(srv.ErrShutdown))
			case <-time.After(time.Second):
				t.Fatal("timeout waiting for Server.Serve to finish")
			}
		},
		Mocks: []*mok.Mock{wantConn1.Mock, wantConn2.Mock, wantConn3.Mock,
			delegate.Mock, listener.Mock},
	}
}

func ServerShutdownBeforeAcceptTestCase(t *testing.T) ServerTestCase {
	name := "Should be able to shutdown the server before any connection"

	var (
		listener   = cmock.NewListener()
		delegate   = cmock.NewServerDelegate()
		stopAccept = make(chan struct{})
	)
	listener.RegisterAccept(
		func() (net.Conn, error) {
			<-stopAccept
			return nil, errors.New("listener closed")
		},
	).RegisterClose(func() error {
		close(stopAccept)
		return nil
	})

	return ServerTestCase{
		Name: name,
		Setup: ServerSetup{
			Delegate: delegate,
			WantErr:  srv.NewServerError(srv.ErrShutdown),
		},
		Action: func(t *testing.T, s *srv.Server) {
			errs := make(chan error, 1)
			go func() {
				errs <- s.Serve(listener)
				close(errs)
			}()

			time.Sleep(10 * time.Millisecond)
			err := s.Shutdown()
			asserterror.EqualError(t, err, nil)

			select {
			case err := <-errs:
				asserterror.EqualError(t, err, srv.NewServerError(srv.ErrShutdown))
			case <-time.After(time.Second):
				t.Fatal("timeout waiting for Server.Serve to finish")
			}
		},
		Mocks: []*mok.Mock{delegate.Mock, listener.Mock},
	}
}

func ServerCloseBeforeAcceptTestCase(t *testing.T) ServerTestCase {
	name := "Should be able to close the server before any connection"

	var (
		listener   = cmock.NewListener()
		delegate   = cmock.NewServerDelegate()
		stopAccept = make(chan struct{})
	)
	listener.RegisterAccept(
		func() (net.Conn, error) {
			<-stopAccept
			return nil, errors.New("listener closed")
		},
	).RegisterClose(func() error {
		close(stopAccept)
		return nil
	})

	return ServerTestCase{
		Name: name,
		Setup: ServerSetup{
			Delegate: delegate,
			WantErr:  srv.NewServerError(srv.ErrClosed),
		},
		Action: func(t *testing.T, s *srv.Server) {
			errs := make(chan error, 1)
			go func() {
				errs <- s.Serve(listener)
				close(errs)
			}()

			time.Sleep(10 * time.Millisecond)
			err := s.Close()
			asserterror.EqualError(t, err, nil)

			select {
			case err := <-errs:
				asserterror.EqualError(t, err, srv.NewServerError(srv.ErrClosed))
			case <-time.After(time.Second):
				t.Fatal("timeout waiting for Server.Serve to finish")
			}
		},
		Mocks: []*mok.Mock{delegate.Mock, listener.Mock},
	}
}

func ServerCloseAfterAcceptTestCase(t *testing.T) ServerTestCase {
	name := "Should be able to close the server after accepting a connection"

	var (
		wantConn1  = cmock.NewConn()
		listener   = cmock.NewListener()
		delegate   = cmock.NewServerDelegate()
		stopAccept = make(chan struct{})
	)

	listener.RegisterAccept(
		func() (net.Conn, error) {
			return wantConn1, nil
		},
	).RegisterAccept(
		func() (net.Conn, error) {
			<-stopAccept
			return nil, errors.New("listener closed")
		},
	).RegisterClose(func() error {
		close(stopAccept)
		return nil
	})

	delegate.RegisterHandle(func(ctx context.Context, conn net.Conn) error {
		asserterror.Equal(t, conn, net.Conn(wantConn1))
		return nil
	})

	return ServerTestCase{
		Name: name,
		Setup: ServerSetup{
			Delegate: delegate,
			Opts: []srv.SetOption{
				srv.WithWorkersCount(1),
			},
			WantErr: srv.NewServerError(srv.ErrClosed),
		},
		Action: func(t *testing.T, s *srv.Server) {
			errs := make(chan error, 1)
			go func() {
				errs <- s.Serve(listener)
				close(errs)
			}()

			time.Sleep(10 * time.Millisecond)
			err := s.Close()
			asserterror.EqualError(t, err, nil)

			select {
			case err := <-errs:
				asserterror.EqualError(t, err, srv.NewServerError(srv.ErrClosed))
			case <-time.After(time.Second):
				t.Fatal("timeout waiting for Server.Serve to finish")
			}
		},
		Mocks: []*mok.Mock{wantConn1.Mock, delegate.Mock, listener.Mock},
	}
}

func ServerShutdownAfterAcceptTestCase(t *testing.T) ServerTestCase {
	name := "Should be able to shutdown the server after accepting a connection"

	var (
		wantConn1  = cmock.NewConn()
		listener   = cmock.NewListener()
		delegate   = cmock.NewServerDelegate()
		stopAccept = make(chan struct{})
	)
	listener.RegisterAccept(
		func() (net.Conn, error) {
			return wantConn1, nil
		},
	).RegisterAccept(
		func() (net.Conn, error) {
			<-stopAccept
			return nil, errors.New("listener closed")
		},
	).RegisterClose(func() error {
		close(stopAccept)
		return nil
	})
	delegate.RegisterHandle(func(ctx context.Context, conn net.Conn) error {
		asserterror.Equal(t, conn, net.Conn(wantConn1))
		return nil
	})
	return ServerTestCase{
		Name: name,
		Setup: ServerSetup{
			Delegate: delegate,
			Opts: []srv.SetOption{
				srv.WithWorkersCount(1),
			},
			WantErr: srv.NewServerError(srv.ErrShutdown),
		},
		Action: func(t *testing.T, s *srv.Server) {
			errs := make(chan error, 1)
			go func() {
				errs <- s.Serve(listener)
				close(errs)
			}()

			time.Sleep(10 * time.Millisecond)
			err := s.Shutdown()
			asserterror.EqualError(t, err, nil)

			select {
			case err := <-errs:
				asserterror.EqualError(t, err, srv.NewServerError(srv.ErrShutdown))
			case <-time.After(time.Second):
				t.Fatal("timeout waiting for Server.Serve to finish")
			}
		},
		Mocks: []*mok.Mock{wantConn1.Mock, delegate.Mock, listener.Mock},
	}
}

func ListenAndServeFailOnInvalidAddrTestCase(t *testing.T) ServerTestCase {
	name := "ListenAndServe should fail on invalid address"

	var delegate = cmock.NewServerDelegate()
	return ServerTestCase{
		Name: name,
		Setup: ServerSetup{
			Delegate: delegate,
			WantErr:  srv.NewServerError(errors.New("listen tcp: lookup tcp/addr: unknown port")),
		},
		Action: func(t *testing.T, s *srv.Server) {
			_ = s.ListenAndServe("invalid:addr")
		},
		Mocks: []*mok.Mock{delegate.Mock},
	}
}

func ServerShutdownFailIfNotServingTestCase(t *testing.T) ServerTestCase {
	name := "Shutdown should fail if server is not serving"

	return ServerTestCase{
		Name: name,
		Setup: ServerSetup{
			Opts: []srv.SetOption{srv.WithWorkersCount(1)},
		},
		Action: func(t *testing.T, s *srv.Server) {
			err := s.Shutdown()
			asserterror.EqualError(t, err, srv.NewServerError(srv.ErrNotServing))
		},
	}
}

func ServerCloseFailIfNotServingTestCase(t *testing.T) ServerTestCase {
	name := "Close should fail if server is not serving"

	return ServerTestCase{
		Name: name,
		Setup: ServerSetup{
			Opts: []srv.SetOption{srv.WithWorkersCount(1)},
		},
		Action: func(t *testing.T, s *srv.Server) {
			err := s.Close()
			asserterror.EqualError(t, err, srv.NewServerError(srv.ErrNotServing))
		},
	}
}

func ServerNegativeWorkersCountTestCase(t *testing.T) ServerTestCase {
	name := "Should fail if WorkersCount is negative"

	return ServerTestCase{
		Name: name,
		Setup: ServerSetup{
			Opts:    []srv.SetOption{srv.WithWorkersCount(-1)},
			WantErr: srv.NewServerError(srv.ErrNoWorkers),
		},
	}
}
