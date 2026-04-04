package cln_test

import (
	"testing"

	"github.com/cmd-stream/cmd-stream-go/delegate/cln"
	test "github.com/cmd-stream/cmd-stream-go/test/delegate"
	asserterror "github.com/ymz-ncnk/assert/error"
)

func TestClientInfoDelegate(t *testing.T) {
	for _, tc := range []test.ClientDelegateTestCase[any]{
		test.NewTestCase[any](t),
		test.NewSetReceiveDeadlineErrorTestCase[any](t),
		test.NewReceiveServerInfoErrorTestCase[any](t),
		test.NewServerInfoMismatchTestCase[any](t),
		test.SendTestCase[any](t),
		test.SendErrorTestCase[any](t),
		test.ReceiveTestCase[any](t),
		test.LocalAddrTestCase[any](t),
		test.RemoteAddrTestCase[any](t),
		test.CloseTestCase[any](t),
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
