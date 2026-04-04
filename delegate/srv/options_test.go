package srv

import (
	"testing"
	"time"

	asserterror "github.com/ymz-ncnk/assert/error"
)

func TestOptions(t *testing.T) {
	var (
		o            = Options{}
		wantDuration = time.Second
	)
	Apply(&o, WithServerInfoSendDuration(wantDuration))
	asserterror.Equal(t, o.ServerInfoSendDuration, wantDuration)
}
