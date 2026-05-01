package cln_test

import (
	"testing"

	"github.com/cmd-stream/cmd-stream-go/test"
)

func TestReconnectDelegate_Init(t *testing.T) {
	d := test.ReconnectDelegate[any]{}
	for _, tc := range []test.ReconnectTestCase[any]{
		d.InitSuccess(t),
		d.InitWrongInfo(t),
		d.InitFactoryError(t),
		d.InitCheckError(t),
	} {
		test.RunReconnectTestCase(t, tc)
	}
}

func TestReconnectDelegate_Reconnect(t *testing.T) {
	d := test.ReconnectDelegate[any]{}
	for _, tc := range []test.ReconnectTestCase[any]{
		d.CycleSuccess(t),
		d.CycleClose(t),
		d.CycleMismatch(t),
		d.CycleCheckError(t),
	} {
		test.RunReconnectTestCase(t, tc)
	}
}

func TestReconnectDelegate_Transport(t *testing.T) {
	d := test.ReconnectDelegate[any]{}
	for _, tc := range []test.ReconnectTestCase[any]{
		d.CloseError(t),
		d.Flush(t),
		d.LocalAddr(t),
		d.RemoteAddr(t),
		d.SetSendDeadline(t),
		d.SetReceiveDeadline(t),
		d.Send(t),
		d.SendError(t),
		d.Receive(t),
		d.ReceiveError(t),
		d.SetReceiveDeadlineError(t),
	} {
		test.RunReconnectTestCase(t, tc)
	}
}
