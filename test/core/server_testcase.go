package core

import (
	"testing"

	"github.com/cmd-stream/cmd-stream-go/core/srv"
	cmock "github.com/cmd-stream/cmd-stream-go/test/mock/core"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

type ServerTestCase struct {
	Name   string
	Setup  ServerSetup
	Action func(t *testing.T, s *srv.Server)
	Mocks  []*mok.Mock
}

type ServerSetup struct {
	Delegate cmock.ServerDelegate
	Opts     []srv.SetOption
	WantErr  error
}

func RunServerTestCase(t *testing.T, tc ServerTestCase) {
	t.Run(tc.Name, func(t *testing.T) {
		s, nErr := srv.New(tc.Setup.Delegate, tc.Setup.Opts...)
		if nErr != nil {
			asserterror.EqualError(t, nErr, tc.Setup.WantErr)
			return
		}

		if tc.Action != nil {
			tc.Action(t, s)
		}

		asserterror.EqualDeep(t, mok.CheckCalls(tc.Mocks), mok.EmptyInfomap)
	})
}
