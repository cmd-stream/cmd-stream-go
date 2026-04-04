package srv

import (
	"net"
	"testing"
	"time"

	asserterror "github.com/ymz-ncnk/assert/error"
)

func TestOptions(t *testing.T) {
	var (
		o                    = Options{}
		wantWorkersCount     = 1
		wantLostConnCallback = func(addr net.Addr, err error) {}
		wantConnReceiver     = []SetConnReceiverOption{}
	)
	Apply(&o,
		WithWorkersCount(wantWorkersCount),
		WithLostConnCallback(wantLostConnCallback),
		WithConnReceiver(wantConnReceiver...),
	)

	asserterror.Equal(t, o.WorkersCount, wantWorkersCount)
	asserterror.EqualDeep(t, o.ConnReceiver, wantConnReceiver)

	if o.LostConnCallback == nil {
		t.Errorf("LostConnCallback == nil")
	}
}

func TestConnReceiverOptions(t *testing.T) {
	var (
		o                    = ConnReceiverOptions{}
		wantFirstConnTimeout = time.Second
	)
	ApplyConnReceiver(&o, WithFirstConnTimeout(wantFirstConnTimeout))
	asserterror.Equal(t, o.FirstConnTimeout, wantFirstConnTimeout)
}
