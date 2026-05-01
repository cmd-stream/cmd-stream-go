package cln_test

import (
	"testing"
	"time"

	dcln "github.com/cmd-stream/cmd-stream-go/delegate/cln"
)

func TestOptions(t *testing.T) {
	var (
		o                             = dcln.Options{}
		wantServerInfoReceiveDuration = time.Second
	)
	dcln.Apply(&o, dcln.WithServerInfoReceiveDuration(wantServerInfoReceiveDuration))

	if o.ServerInfoReceiveDuration != wantServerInfoReceiveDuration {
		t.Errorf("unexpected ServerInfoReceiveDuration, want %v actual %v",
			wantServerInfoReceiveDuration, o.ServerInfoReceiveDuration)
	}
}

func TestKeepAliveOptions(t *testing.T) {
	var (
		o                  = dcln.KeepaliveOptions{}
		wantKeepaliveTime  = 2 * time.Second
		wantKeepaliveIntvl = 3 * time.Second
	)
	dcln.ApplyKeepalive(&o, dcln.WithKeepaliveTime(wantKeepaliveTime),
		dcln.WithKeepaliveIntvl(wantKeepaliveIntvl))

	if o.KeepaliveTime != wantKeepaliveTime {
		t.Errorf("unexpected KeepaliveTime, want %v actual %v", wantKeepaliveTime,
			o.KeepaliveTime)
	}

	if o.KeepaliveIntvl != wantKeepaliveIntvl {
		t.Errorf("unexpected KeepaliveIntvl, want %v actual %v", wantKeepaliveIntvl,
			o.KeepaliveIntvl)
	}
}
