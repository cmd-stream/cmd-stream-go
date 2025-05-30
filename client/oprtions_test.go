package cln

import (
	"testing"

	dcln "github.com/cmd-stream/delegate-go/client"
	"github.com/cmd-stream/transport-go"
	asserterror "github.com/ymz-ncnk/assert/error"
)

func TestOptions(t *testing.T) {
	var (
		o             = Options{}
		wantDelegate  = []dcln.SetOption{}
		wantKeepalive = []dcln.SetKeepaliveOption{}
		wantTransport = []transport.SetOption{}
	)
	Apply([]SetOption{
		WithDelegate(wantDelegate...),
		WithKeepalive(wantKeepalive...),
		WithTransport(wantTransport...),
	}, &o)

	asserterror.EqualDeep(o.Delegate, wantDelegate, t)
	asserterror.EqualDeep(o.Keepalive, wantKeepalive, t)
	asserterror.EqualDeep(o.Transport, wantTransport, t)

}
