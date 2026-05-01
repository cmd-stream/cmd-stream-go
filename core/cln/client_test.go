package cln_test

import (
	"testing"

	"github.com/cmd-stream/cmd-stream-go/test"
)

func TestInit(t *testing.T) {
	for _, tc := range []test.ClientTestCase[any]{
		test.Client.Reconnect(),
		test.Client.ReconnectOnEOF(),
		test.Client.NoReconnectOnClose(),
		test.Client.ReconnectFail(),
		test.Client.Keepalive(),
	} {
		test.RunClientTestCase(t, tc)
	}
}

func TestSend(t *testing.T) {
	for _, tc := range []test.ClientTestCase[any]{
		test.Client.SendSuccess(),
		test.Client.Has(),
		test.Client.Forget(),
		test.Client.ForgetOnFail(),
		test.Client.ClosedOnReceiveError(),
		test.Client.UnexpectedResult(),
		test.Client.UnexpectedResultCallback(),
	} {
		test.RunClientTestCase(t, tc)
	}

	for _, tc := range []test.MultiSendTestCase[any]{
		test.Client.MultiSendSuccess(),
		test.Client.IncrementSeq(),
		test.Client.MultiResultSuccess(),
		test.Client.PartialResults(),
		test.Client.IncrementSeqAfterFail(),
		test.Client.ErrForAllCmdsOnFlushFail(),
	} {
		test.RunMultiSendTestCase(t, tc)
	}
}

func TestSendWithDeadline(t *testing.T) {
	for _, tc := range []test.ClientTestCase[any]{
		test.Client.SendWithDeadline(),
		test.Client.SendWithDeadlineFailSetDeadline(),
		test.Client.SendWithDeadlineFail(),
		test.Client.ForgetOnSendWithDeadlineFailSetDeadline(),
		test.Client.ForgetOnSendWithDeadlineFail(),
	} {
		test.RunClientTestCase(t, tc)
	}

	test.RunMultiSendTestCase(t, test.Client.IncrementSeqOnSendWithDeadlineFail())
}

func TestClose(t *testing.T) {
	for _, tc := range []test.ClientTestCase[any]{
		test.Client.CloseSuccess(),
		test.Client.CloseDuringQueueResult(),
		test.Client.CloseDelegateFail(),
	} {
		test.RunClientTestCase(t, tc)
	}
}
