package transport_test

import (
	"testing"

	"github.com/cmd-stream/cmd-stream-go/test/transport"
)

func TestTransport(t *testing.T) {
	for _, tc := range []transport.TransportTestCase[any, any]{
		transport.LocalAddrTestCase(),
		transport.RemoteAddrTestCase(),
		transport.SetSendDeadlineTestCase(t),
		transport.SetSendDeadlineErrorTestCase(),
		transport.SetReceiveDeadlineTestCase(t),
		transport.SetReceiveDeadlineErrorTestCase(),
		transport.FlushTestCase(),
		transport.FlushErrorTestCase(),
		transport.CloseTestCase(),
		transport.CloseErrorTestCase(),
	} {
		transport.RunTransportTestCase(t, tc)
	}

	transport.RunTransportTestCase(t, transport.SendTestCase(t))
	transport.RunTransportTestCase(t, transport.SendErrorTestCase())
	transport.RunTransportTestCase(t, transport.ReceiveTestCase())
	transport.RunTransportTestCase(t, transport.ReceiveErrorTestCase())
}
