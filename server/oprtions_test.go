package srv

import (
	"reflect"
	"testing"

	csrv "github.com/cmd-stream/core-go/server"
	dsrv "github.com/cmd-stream/delegate-go/server"
	"github.com/cmd-stream/handler-go"
	"github.com/cmd-stream/transport-go"
)

func TestOptions(t *testing.T) {
	var (
		o             = Options{}
		wantCore      = []csrv.SetOption{}
		wantDelegate  = []dsrv.SetOption{}
		wantHandler   = []handler.SetOption{}
		wantTransport = []transport.SetOption{}
	)
	Apply([]SetOption{
		WithCore(wantCore...),
		WithDelegate(wantDelegate...),
		WithHandler(wantHandler...),
		WithTransport(wantTransport...),
	}, &o)

	if !reflect.DeepEqual(o.Base, wantCore) {
		t.Errorf("unexpected Base, want %v actual %v", wantCore, o.Base)
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
