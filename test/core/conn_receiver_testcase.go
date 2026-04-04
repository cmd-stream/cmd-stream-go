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

type ConnReceiverTestCase struct {
	Name   string
	Setup  ConnReceiverSetup
	During func(t *testing.T, receiver *srv.ConnReceiver, conns chan net.Conn)
	Mocks  []*mok.Mock
}

type ConnReceiverSetup struct {
	Listener cmock.Listener
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
