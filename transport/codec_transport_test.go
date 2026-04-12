package transport_test

import (
	"testing"

	tspt "github.com/cmd-stream/cmd-stream-go/test/transport"
)

func TestTransport(t *testing.T) {
	for _, tc := range []tspt.TransportTestCase[any, any]{
		tspt.LocalAddrTestCase(),
		tspt.RemoteAddrTestCase(),
		tspt.SetSendDeadlineTestCase(t),
		tspt.SetSendDeadlineErrorTestCase(),
		tspt.SetReceiveDeadlineTestCase(t),
		tspt.SetReceiveDeadlineErrorTestCase(),
		tspt.FlushTestCase(),
		tspt.FlushErrorTestCase(),
		tspt.CloseTestCase(),
		tspt.CloseErrorTestCase(),
	} {
		tspt.RunTransportTestCase(t, tc)
	}

	tspt.RunTransportTestCase(t, tspt.SendTestCase(t))
	tspt.RunTransportTestCase(t, tspt.SendErrorTestCase())
	tspt.RunTransportTestCase(t, tspt.ReceiveTestCase())
	tspt.RunTransportTestCase(t, tspt.ReceiveErrorTestCase())
}
