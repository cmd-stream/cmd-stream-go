package srv_test

import (
	"testing"
	"time"

	dsrv "github.com/cmd-stream/cmd-stream-go/delegate/srv"
	asserterror "github.com/ymz-ncnk/assert/error"
)

func TestOptions(t *testing.T) {
	var (
		o            = dsrv.Options{}
		wantDuration = time.Second
	)
	dsrv.Apply(&o, dsrv.WithServerInfoSendDuration(wantDuration))
	asserterror.Equal(t, o.ServerInfoSendDuration, wantDuration)
}
