package group

import (
	"testing"
	"time"

	grp "github.com/cmd-stream/cmd-stream-go/group"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

type GroupTestCase struct {
	Name     string
	Strategy grp.DispatchStrategy[grp.GroupClient[any]]
	Action   func(t *testing.T, g grp.Group[any])
	Mocks    []*mok.Mock
}

func RunGroupTestCase(t *testing.T, tc GroupTestCase) {
	t.Run(tc.Name, func(t *testing.T) {
		group := grp.New(tc.Strategy)
		tc.Action(t, group)
		_ = group.Close()
		select {
		case <-group.Done():
		case <-time.After(time.Second):
			t.Fatal("timeout waiting for group to be done")
		}
		asserterror.EqualDeep(t, mok.CheckCalls(tc.Mocks), mok.EmptyInfomap)
	})
}
