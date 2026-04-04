package group_test

import (
	"testing"

	test "github.com/cmd-stream/cmd-stream-go/test/group"
)

func TestGroup(t *testing.T) {
	for _, tc := range []test.GroupTestCase{
		test.GroupSendTestCase(t),
		test.GroupSendWithDeadlineTestCase(t),
		test.GroupHasTestCase(t),
		test.GroupForgetTestCase(t),
		test.GroupErrorTestCase(t),
		test.GroupCloseTestCase(t),
	} {
		test.RunGroupTestCase(t, tc)
	}
}
