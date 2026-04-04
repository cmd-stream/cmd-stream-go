package delegate

import (
	"testing"

	cln "github.com/cmd-stream/cmd-stream-go/delegate/cln"
	cmock "github.com/cmd-stream/cmd-stream-go/test/mock/core"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

type KeepaliveTestCase struct {
	Name   string
	Setup  KeepaliveSetup
	Action func(t *testing.T, d *cln.KeepaliveDelegate[any])
	Mocks  []*mok.Mock
}

type KeepaliveSetup struct {
	Delegate cmock.ClientDelegate
	Opts     []cln.SetKeepaliveOption
}

func RunKeepaliveTestCase(t *testing.T, tc KeepaliveTestCase) {
	t.Run(tc.Name, func(t *testing.T) {
		d := cln.NewKeepalive(tc.Setup.Delegate, tc.Setup.Opts...)

		tc.Action(t, d)
		asserterror.EqualDeep(t, mok.CheckCalls(tc.Mocks), mok.EmptyInfomap)
	})
}
