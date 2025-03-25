package ccln

import (
	"reflect"
	"testing"

	dcln "github.com/cmd-stream/delegate-go/client"
	tcom "github.com/cmd-stream/transport-go/common"
)

func TestOptions(t *testing.T) {
	var (
		o             = Options{}
		wantDelegate  = []dcln.SetOption{}
		wantKeepalive = []dcln.SetKeepaliveOption{}
		wantTransport = []tcom.SetOption{}
	)
	Apply([]SetOption{
		WithDelegate(wantDelegate...),
		WithKeepalive(wantKeepalive...),
		WithTransport(wantTransport...),
	}, &o)

	if !reflect.DeepEqual(o.Delegate, wantDelegate) {
		t.Errorf("unexpected Delegate, want %v actual %v", wantDelegate,
			o.Delegate)
	}

	if !reflect.DeepEqual(o.Keepalive, wantKeepalive) {
		t.Errorf("unexpected Keepalive, want %v actual %v", wantKeepalive,
			o.Keepalive)
	}

	if !reflect.DeepEqual(o.Transport, wantTransport) {
		t.Errorf("unexpected Transport, want %v actual %v", wantTransport,
			o.Transport)
	}

}
