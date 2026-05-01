package srv_test

import (
	"testing"

	"github.com/cmd-stream/cmd-stream-go/test"
)

func TestServerInfoDelegate(t *testing.T) {
	d := test.DelegateServer[any]{}
	for _, tc := range []test.ServerDelegateTestCase[any]{
		d.ZeroLenServerInfo(t),
		d.SendServerInfoError(t),
		d.TransportHandleError(t),
		d.SendServerInfoTransportDeadlineError(t),
		d.HandleConnSuccess(t),
	} {
		test.RunServerDelegateTestCase(t, tc)
	}
}
