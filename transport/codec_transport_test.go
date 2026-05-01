package transport_test

import (
	"testing"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/cmd-stream/cmd-stream-go/test"
)

func TestTransport(t *testing.T) {
	s := test.TransportSuite[any, any]{}
	for _, tc := range []test.TransportTestCase[any, any]{
		s.LocalAddr(t),
		s.RemoteAddr(t),
		s.SetSendDeadline(t),
		s.SetSendDeadlineError(t),
		s.SetReceiveDeadline(t),
		s.SetReceiveDeadlineError(t),
		s.Flush(t),
		s.FlushError(t),
		s.Close(t),
		s.CloseError(t),
	} {
		test.RunTransportTestCase(t, tc)
	}

	c := test.TransportSuite[core.Cmd[any], core.Result]{}
	for _, tc := range []test.TransportTestCase[core.Cmd[any], core.Result]{
		c.Send(t),
		c.SendError(t),
		c.Receive(t),
		c.ReceiveError(t),
	} {
		test.RunTransportTestCase(t, tc)
	}
}
