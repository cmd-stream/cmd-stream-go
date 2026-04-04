package delegate

import (
	"testing"

	dlgt "github.com/cmd-stream/cmd-stream-go/delegate"
	"github.com/cmd-stream/cmd-stream-go/delegate/cln"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

type ReconnectTestCase struct {
	Name   string
	Setup  ReconnectSetup
	Action func(t *testing.T, d *cln.ReconnectDelegate[any], initErr error)
	Mocks  []*mok.Mock
}

type ReconnectSetup struct {
	Info    dlgt.ServerInfo
	Factory dlgt.ClientTransportFactory[any]
	Opts    []cln.SetOption
}

func RunReconnectTestCase(t *testing.T, tc ReconnectTestCase) {
	t.Run(tc.Name, func(t *testing.T) {
		d, err := cln.NewReconnect(tc.Setup.Info, tc.Setup.Factory,
			tc.Setup.Opts...)

		tc.Action(t, d, err)
		asserterror.EqualDeep(t, mok.CheckCalls(tc.Mocks), mok.EmptyInfomap)
	})
}
