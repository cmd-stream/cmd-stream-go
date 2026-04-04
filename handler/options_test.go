package handler

import (
	"testing"
	"time"

	asserterror "github.com/ymz-ncnk/assert/error"
)

func TestOptions(t *testing.T) {
	var (
		o            = Options{}
		wantDuration = time.Second
		wantAt       = true
	)
	Apply(&o, WithCmdReceiveDuration(wantDuration), WithAt())
	asserterror.Equal(t, o.CmdReceiveDuration, wantDuration)
	asserterror.Equal(t, o.At, wantAt)
}
