package srv

import (
	"crypto/tls"
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
		wantTLSConfig        = &tls.Config{}
	)
	Apply(&o,
		WithWorkersCount(wantWorkersCount),
		WithLostConnCallback(wantLostConnCallback),
		WithConnReceiver(wantConnReceiver...),
		WithTLSConfig(wantTLSConfig),
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
		options Options
		wantErr error
	}{
		{
			name:    "valid",
			options: Options{WorkersCount: 1},
			wantErr: nil,
		},
		{
			name:    "no workers",
			options: Options{WorkersCount: 0},
			wantErr: ErrNoWorkers,
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
	o := DefaultOptions()
	asserterror.Equal(t, o.WorkersCount, WorkersCount)
}

func TestConnReceiverOptions(t *testing.T) {
	var (
		o                    = ConnReceiverOptions{}
		wantFirstConnTimeout = time.Second
	)
	ApplyConnReceiver(&o, WithFirstConnTimeout(wantFirstConnTimeout))
	asserterror.Equal(t, o.FirstConnTimeout, wantFirstConnTimeout)
}
