package client

import (
	"testing"

	dcln "github.com/cmd-stream/cmd-stream-go/delegate/cln"
	tspt "github.com/cmd-stream/cmd-stream-go/transport"
	asserterror "github.com/ymz-ncnk/assert/error"
)

func TestOptions(t *testing.T) {
	var (
		o             = Options{}
		wantDelegate  = []dcln.SetOption{func(o *dcln.Options) {}}
		wantKeepalive = []dcln.SetKeepaliveOption{func(o *dcln.KeepaliveOptions) {}}
		wantTransport = []tspt.SetOption{func(o *tspt.Options) {}}
	)
	Apply(&o, WithDelegate(wantDelegate...),
		WithKeepalive(wantKeepalive...),
		WithTransport(wantTransport...),
	)
	asserterror.Equal(t, len(o.Delegate), len(wantDelegate))
	asserterror.Equal(t, len(o.Keepalive), len(wantKeepalive))
	asserterror.Equal(t, len(o.Transport), len(wantTransport))
}
