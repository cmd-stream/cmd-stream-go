package cln_test

import (
	"testing"

	"github.com/cmd-stream/cmd-stream-go/core"
	ccln "github.com/cmd-stream/cmd-stream-go/core/cln"
)

func TestOptions(t *testing.T) {
	var (
		o                                          = ccln.Options{}
		wantCallback ccln.UnexpectedResultCallback = func(seq core.Seq, result core.Result) {}
	)
	ccln.Apply(&o, ccln.WithUnexpectedResultCallback(wantCallback))

	if o.UnexpectedResultCallback == nil {
		t.Errorf("UnexpectedResultCallback == nil")
	}
}
