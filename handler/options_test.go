package handler_test

import (
	"testing"
	"time"

	hdlr "github.com/cmd-stream/cmd-stream-go/handler"
	asserterror "github.com/ymz-ncnk/assert/error"
)

func TestOptions(t *testing.T) {
	var (
		o            = hdlr.Options{}
		wantDuration = time.Second
		wantAt       = true
	)
	hdlr.Apply(&o, hdlr.WithCmdReceiveDuration(wantDuration), hdlr.WithAt())
	asserterror.Equal(t, o.CmdReceiveDuration, wantDuration)
	asserterror.Equal(t, o.At, wantAt)
}
