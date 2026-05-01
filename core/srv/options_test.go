package srv_test

import (
	"crypto/tls"
	"net"
	"testing"
	"time"

	csrv "github.com/cmd-stream/cmd-stream-go/core/srv"
	asserterror "github.com/ymz-ncnk/assert/error"
)

func TestOptions(t *testing.T) {
	var (
		o                    = csrv.Options{}
		wantWorkersCount     = 1
		wantLostConnCallback = func(addr net.Addr, err error) {}
		wantConnReceiver     = []csrv.SetConnReceiverOption{}
		wantTLSConfig        = &tls.Config{}
	)
	csrv.Apply(&o,
		csrv.WithWorkersCount(wantWorkersCount),
		csrv.WithLostConnCallback(wantLostConnCallback),
		csrv.WithConnReceiver(wantConnReceiver...),
		csrv.WithTLSConfig(wantTLSConfig),
		nil,
	)

	asserterror.Equal(t, o.WorkersCount, wantWorkersCount)
	asserterror.EqualDeep(t, o.ConnReceiver, wantConnReceiver)
	asserterror.Equal(t, o.TLSConfig, wantTLSConfig)

	if o.LostConnCallback == nil {
		t.Errorf("LostConnCallback == nil")
	}
}

func TestOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		options csrv.Options
		wantErr error
	}{
		{
			name:    "valid",
			options: csrv.Options{WorkersCount: 1},
			wantErr: nil,
		},
		{
			name:    "no workers",
			options: csrv.Options{WorkersCount: 0},
			wantErr: csrv.ErrNoWorkers,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.options.Validate(); err != tt.wantErr {
				t.Errorf("Options.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultOptions(t *testing.T) {
	o := csrv.DefaultOptions()
	asserterror.Equal(t, o.WorkersCount, csrv.WorkersCount)
}

func TestConnReceiverOptions(t *testing.T) {
	var (
		o                    = csrv.ConnReceiverOptions{}
		wantFirstConnTimeout = time.Second
	)
	csrv.ApplyConnReceiver(&o, csrv.WithFirstConnTimeout(wantFirstConnTimeout))
	asserterror.Equal(t, o.FirstConnTimeout, wantFirstConnTimeout)
}
