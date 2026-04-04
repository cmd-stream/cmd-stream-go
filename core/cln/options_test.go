package cln

import (
	"testing"

	"github.com/cmd-stream/cmd-stream-go/core"
)

func TestOptions(t *testing.T) {
	var (
		o                                     = Options{}
		wantCallback UnexpectedResultCallback = func(seq core.Seq, result core.Result) {}
	)
	Apply(&o, WithUnexpectedResultCallback(wantCallback))

	if o.UnexpectedResultCallback == nil {
		t.Errorf("UnexpectedResultCallback == nil")
	}
}
