package group_test

import (
	"testing"

	"github.com/cmd-stream/cmd-stream-go/test"
)

func TestGroup(t *testing.T) {
	g := test.GroupSuite[any]{}
	for _, tc := range []test.GroupTestCase[any]{
		g.Send(t),
		g.SendWithDeadline(t),
		g.Has(t),
		g.Forget(t),
		g.Error(t),
		g.Close(t),
	} {
		test.RunGroupTestCase(t, tc)
	}
}
