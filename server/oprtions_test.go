package cser

import (
	"reflect"
	"testing"

	bser "github.com/cmd-stream/base-go/server"
	dser "github.com/cmd-stream/delegate-go/server"
	"github.com/cmd-stream/handler-go"
	tcom "github.com/cmd-stream/transport-go/common"
)

func TestOptions(t *testing.T) {
	var (
		o             = Options{}
		wantBase      = []bser.SetOption{}
		wantDelegate  = []dser.SetOption{}
		wantHandler   = []handler.SetOption{}
		wantTransport = []tcom.SetOption{}
	)
	Apply([]SetOption{
		WithBase(wantBase...),
		WithDelegate(wantDelegate...),
		WithHandler(wantHandler...),
		WithTransport(wantTransport...),
	}, &o)

	if !reflect.DeepEqual(o.Base, wantBase) {
		t.Errorf("unexpected Base, want %v actual %v", wantBase, o.Base)
	}

	if !reflect.DeepEqual(o.Delegate, wantDelegate) {
		t.Errorf("unexpected Delegate, want %v actual %v", wantDelegate,
			o.Delegate)
	}

	if !reflect.DeepEqual(o.Handler, wantHandler) {
		t.Errorf("unexpected Handler, want %v actual %v", wantHandler,
			o.Handler)
	}

	if !reflect.DeepEqual(o.Transport, wantTransport) {
		t.Errorf("unexpected Transport, want %v actual %v", wantTransport,
			o.Transport)
	}

}
