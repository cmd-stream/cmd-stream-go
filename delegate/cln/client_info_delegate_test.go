package cln_test

import (
	"testing"

	"github.com/cmd-stream/cmd-stream-go/delegate/cln"
	"github.com/cmd-stream/cmd-stream-go/test"
	asserterror "github.com/ymz-ncnk/assert/error"
)

func TestClientInfoDelegate(t *testing.T) {
	d := test.DelegateClient[any]{}
	for _, tc := range []test.ClientDelegateTestCase[any]{
		d.NewCheckServerInfo(t),
		d.NewSetReceiveDeadlineError(t),
		d.NewReceiveServerInfoError(t),
		d.NewServerInfoMismatch(t),
		d.Send(t),
		d.SendError(t),
		d.Receive(t),
		d.LocalAddr(t),
		d.RemoteAddr(t),
		d.Close(t),
	} {
		test.RunClientDelegateTestCase(t, tc)
	}
}

func TestClientInfoDelegate_Options(t *testing.T) {
	var (
		wantO    = cln.Options{}
		delegate = cln.NewWithoutInfo[any](nil)
	)
	o := delegate.Options()
	asserterror.Equal(t, o, wantO)
}
