package cln_test

import (
	"testing"

	test "github.com/cmd-stream/cmd-stream-go/test/core"
)

func TestInit(t *testing.T) {
	for _, tc := range []test.ClientTestCase[any]{
		test.ReconnectTestCase(),
		test.ReconnectOnEOFTestCase(),
		test.NoReconnectOnCloseTestCase(),
		test.ReconnectFailTestCase(),
		test.KeepaliveTestCase(),
	} {
		test.RunClientTestCase(t, tc)
	}
}

func TestSend(t *testing.T) {
	for _, tc := range []test.ClientTestCase[any]{
		test.SendSuccessTestCase(),
		test.HasTestCase(),
		test.ForgetTestCase(),
		test.ForgetOnFailTestCase(),
		test.ClosedOnReceiveErrorTestCase(),
	} {
		test.RunClientTestCase(t, tc)
	}

	for _, tc := range []test.MultiSendTestCase[any]{
		test.MultiSuccessTestCase(),
		test.IncrementSeqTestCase(),
		test.MultiResultSuccessTestCase(),
		test.PartialResultsTestCase(),
		test.IncrementSeqAfterFailTestCase(),
		test.ErrForAllCmdsOnFlushFailTestCase(),
	} {
		test.RunMultiSendTestCase(t, tc)
	}
}

func TestSendWithDeadline(t *testing.T) {
	for _, tc := range []test.ClientTestCase[any]{
		test.SendWDTestCase(),
		test.SendWDFailSetDeadlineTestCase(),
		test.SendWDFailTestCase(),
		test.ForgetOnSendWDFailSetDeadlineTestCase(),
		test.ForgetOnSendWDFailSendTestCase(),
	} {
		test.RunClientTestCase(t, tc)
	}

	test.RunMultiSendTestCase(t, test.IncrementSeqOnSendWDFailTestCase())
}

func TestClose(t *testing.T) {
	for _, tc := range []test.ClientTestCase[any]{
		test.CloseSuccessTestCase(),
		test.CloseDuringQueueResultTestCase(),
		test.CloseDelegateFailTestCase(),
	} {
		test.RunClientTestCase(t, tc)
	}
}
