package core

import (
	"net"
	"testing"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core/srv"
	cmock "github.com/cmd-stream/cmd-stream-go/test/mock/core"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

type WorkerTestCase struct {
	Name   string
	Setup  WorkerSetup
	During func(t *testing.T, w *srv.Worker, conns chan net.Conn)
	Mocks  []*mok.Mock
}

type WorkerSetup struct {
	Conns    chan net.Conn
	Delegate cmock.ServerDelegate
	Callback srv.LostConnCallback
	WantErr  error
}

func RunWorkerTestCase(t *testing.T, tc WorkerTestCase) {
	t.Run(tc.Name, func(t *testing.T) {
		var (
			w    = srv.NewWorker(tc.Setup.Conns, tc.Setup.Delegate, tc.Setup.Callback)
			errs = make(chan error, 1)
		)
		go func() {
			errs <- w.Run()
			close(errs)
		}()
		if tc.During != nil {
			tc.During(t, w, tc.Setup.Conns)
		}
		select {
		case err := <-errs:
			asserterror.EqualError(t, err, tc.Setup.WantErr)
		case <-time.After(time.Second):
			t.Fatal("timeout waiting for Worker to finish")
		}
		asserterror.EqualDeep(t, mok.CheckCalls(tc.Mocks), mok.EmptyInfomap)
	})
}
