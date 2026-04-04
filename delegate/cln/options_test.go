package cln

import (
	"testing"
	"time"
)

func TestOptions(t *testing.T) {
	var (
		o                             = Options{}
		wantServerInfoReceiveDuration = time.Second
	)
	Apply(&o, WithServerInfoReceiveDuration(wantServerInfoReceiveDuration))

	if o.ServerInfoReceiveDuration != wantServerInfoReceiveDuration {
		t.Errorf("unexpected ServerInfoReceiveDuration, want %v actual %v",
			wantServerInfoReceiveDuration, o.ServerInfoReceiveDuration)
	}
}

func TestKeepAliveOptions(t *testing.T) {
	var (
		o                  = KeepaliveOptions{}
		wantKeepaliveTime  = 2 * time.Second
		wantKeepaliveIntvl = 3 * time.Second
	)
	ApplyKeepalive(&o, WithKeepaliveTime(wantKeepaliveTime),
		WithKeepaliveIntvl(wantKeepaliveIntvl))

	if o.KeepaliveTime != wantKeepaliveTime {
		t.Errorf("unexpected KeepaliveTime, want %v actual %v", wantKeepaliveTime,
			o.KeepaliveTime)
	}

	if o.KeepaliveIntvl != wantKeepaliveIntvl {
		t.Errorf("unexpected KeepaliveIntvl, want %v actual %v", wantKeepaliveIntvl,
			o.KeepaliveIntvl)
	}
}
