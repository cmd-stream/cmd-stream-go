package core

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/cmd-stream/cmd-stream-go/core/srv"
	cmock "github.com/cmd-stream/cmd-stream-go/test/mock/core"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

func WorkerRunSeveralConnectionsTestCase(t *testing.T) WorkerTestCase {
	name := "Should be able to handle several connections without LostConnCallback"

	var (
		wantConn1 = cmock.NewConn()
		wantConn2 = cmock.NewConn()
		delegate  = cmock.NewServerDelegate()
	)
	delegate.RegisterHandle(func(ctx context.Context, conn net.Conn) error {
		asserterror.Equal(t, conn, net.Conn(wantConn1))
		return nil
	}).RegisterHandle(func(ctx context.Context, conn net.Conn) error {
		asserterror.Equal(t, conn, net.Conn(wantConn2))
		return nil
	})
	return WorkerTestCase{
		Name: name,
		Setup: WorkerSetup{
			Conns:    make(chan net.Conn, 2),
			Delegate: delegate,
			WantErr:  nil,
		},
		During: func(t *testing.T, w *srv.Worker, conns chan net.Conn) {
			conns <- wantConn1
			conns <- wantConn2
			close(conns)
		},
		Mocks: []*mok.Mock{wantConn1.Mock, wantConn2.Mock, delegate.Mock},
	}
}

func WorkerLostConnCallbackTestCase(t *testing.T) WorkerTestCase {
	name := "Should call LostConnCallback if connection handling failed"

	var (
		wantErr1   = errors.New("handle conn 1 failed")
		wantErr2   = errors.New("handle conn 2 failed")
		wantAddr1  = &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9001}
		wantAddr2  = &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9002}
		wantConn1  = cmock.NewConn()
		wantConn2  = cmock.NewConn()
		delegate   = cmock.NewServerDelegate()
		callbackCh = make(chan struct{}, 2)
		callback   = func(addr net.Addr, err error) {
			if len(callbackCh) == 0 {
				asserterror.Equal(t, addr, net.Addr(wantAddr1))
				asserterror.EqualError(t, err, wantErr1)
			} else {
				asserterror.Equal(t, addr, net.Addr(wantAddr2))
				asserterror.EqualError(t, err, wantErr2)
			}
			callbackCh <- struct{}{}
		}
	)
	wantConn1.RegisterRemoteAddr(
		func() net.Addr { return wantAddr1 },
	)
	wantConn2.RegisterRemoteAddr(
		func() net.Addr { return wantAddr2 },
	)
	delegate.RegisterHandle(
		func(ctx context.Context, conn net.Conn) error {
			asserterror.Equal(t, conn, net.Conn(wantConn1))
			return wantErr1
		},
	).RegisterHandle(
		func(ctx context.Context, conn net.Conn) error {
			asserterror.Equal(t, conn, net.Conn(wantConn2))
			return wantErr2
		},
	)
	return WorkerTestCase{
		Name: name,
		Setup: WorkerSetup{
			Conns:    make(chan net.Conn, 2),
			Delegate: delegate,
			Callback: callback,
			WantErr:  nil,
		},
		During: func(t *testing.T, w *srv.Worker, conns chan net.Conn) {
			conns <- wantConn1
			conns <- wantConn2
			close(conns)
		},
		Mocks: []*mok.Mock{wantConn1.Mock, wantConn2.Mock, delegate.Mock},
	}
}

func WorkerStopTestCase(t *testing.T) WorkerTestCase {
	name := "Should be able to close the worker"

	var delegate = cmock.NewServerDelegate()

	return WorkerTestCase{
		Name: name,
		Setup: WorkerSetup{
			Conns:    make(chan net.Conn),
			Delegate: delegate,
			WantErr:  srv.ErrClosed,
		},
		During: func(t *testing.T, w *srv.Worker, conns chan net.Conn) {
			err := w.Stop()
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{delegate.Mock},
	}
}

func WorkerShutdownTestCase(t *testing.T) WorkerTestCase {
	name := "Should be able to shutdown the worker"

	var delegate = cmock.NewServerDelegate()

	return WorkerTestCase{
		Name: name,
		Setup: WorkerSetup{
			Conns:    make(chan net.Conn),
			Delegate: delegate,
			WantErr:  nil,
		},
		During: func(t *testing.T, w *srv.Worker, conns chan net.Conn) {
			close(conns)
		},
		Mocks: []*mok.Mock{delegate.Mock},
	}
}
