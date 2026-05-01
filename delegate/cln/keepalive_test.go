package cln_test

import (
	"testing"

	"github.com/cmd-stream/cmd-stream-go/test"
)

func TestKeepaliveDelegate_Receive(t *testing.T) {
	d := test.KeepaliveDelegate[any]{}
	for _, tc := range []test.KeepaliveTestCase[any]{
		d.SkipPong(t),
	} {
		test.RunKeepaliveTestCase(t, tc)
	}
}

func TestKeepaliveDelegate_Keepalive(t *testing.T) {
	d := test.KeepaliveDelegate[any]{}
	for _, tc := range []test.KeepaliveTestCase[any]{
		d.ShouldSendPing(t),
		d.FlushDelay(t),
		d.CloseCancel(t),
		d.CloseError(t),
		d.SendError(t),
		d.FlushError(t),
	} {
		test.RunKeepaliveTestCase(t, tc)
	}
}
