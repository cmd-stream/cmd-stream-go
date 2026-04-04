package cln_test

import (
	"testing"

	cln "github.com/cmd-stream/cmd-stream-go/test/delegate"
)

func TestKeepaliveDelegate_Receive(t *testing.T) {
	for _, tc := range []cln.KeepaliveTestCase{
		cln.KeepaliveSkipPongTestCase(t),
	} {
		cln.RunKeepaliveTestCase(t, tc)
	}
}

func TestKeepaliveDelegate_Keepalive(t *testing.T) {
	for _, tc := range []cln.KeepaliveTestCase{
		cln.KeepaliveShouldSendPingTestCase(t),
		cln.KeepaliveFlushDelayTestCase(t),
		cln.KeepaliveCloseCancelTestCase(t),
		cln.KeepaliveCloseErrorTestCase(t),
		cln.KeepaliveSendErrorTestCase(t),
		cln.KeepaliveFlushErrorTestCase(t),
	} {
		cln.RunKeepaliveTestCase(t, tc)
	}
}
