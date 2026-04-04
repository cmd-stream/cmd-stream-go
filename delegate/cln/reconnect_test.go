package cln_test

import (
	"testing"

	test "github.com/cmd-stream/cmd-stream-go/test/delegate"
)

func TestReconnectDelegate_Init(t *testing.T) {
	for _, tc := range []test.ReconnectTestCase{
		test.ReconnectInitSuccessTestCase(),
		test.ReconnectInitWrongInfoTestCase(),
		test.ReconnectInitFactoryErrorTestCase(),
		test.ReconnectInitCheckErrorTestCase(),
	} {
		test.RunReconnectTestCase(t, tc)
	}
}

func TestReconnectDelegate_Reconnect(t *testing.T) {
	for _, tc := range []test.ReconnectTestCase{
		test.ReconnectCycleSuccessTestCase(),
		test.ReconnectCycleCloseTestCase(),
		test.ReconnectCycleMismatchTestCase(),
		test.ReconnectCycleCheckErrorTestCase(),
	} {
		test.RunReconnectTestCase(t, tc)
	}
}

func TestReconnectDelegate_Transport(t *testing.T) {
	for _, tc := range []test.ReconnectTestCase{
		test.ReconnectCloseErrorTestCase(),
		test.ReconnectFlushTestCase(),
		test.ReconnectLocalAddrTestCase(),
		test.ReconnectRemoteAddrTestCase(),
		test.ReconnectSetSendDeadlineTestCase(),
		test.ReconnectSetReceiveDeadlineTestCase(),
		test.ReconnectSendTestCase(),
		test.ReconnectSendErrorTestCase(),
		test.ReconnectReceiveTestCase(),
		test.ReconnectReceiveErrorTestCase(),
		test.ReconnectSetReceiveDeadlineErrorTestCase(),
	} {
		test.RunReconnectTestCase(t, tc)
	}
}
