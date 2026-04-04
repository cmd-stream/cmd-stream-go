package srv

import (
	"crypto/tls"
)

// WorkersCount defines the default number of workers.
const WorkersCount = 8

type Options struct {
	WorkersCount     int
	LostConnCallback LostConnCallback
	ConnReceiver     []SetConnReceiverOption
	TLSConfig        *tls.Config
}

func DefaultOptions() Options {
	return Options{
		WorkersCount: WorkersCount,
	}
}

func (o Options) Validate() error {
	if o.WorkersCount <= 0 {
		return ErrNoWorkers
	}
	return nil
}

type SetOption func(o *Options)

// WithWorkersCount sets the number of workers. Must be greater than 0.
func WithWorkersCount(count int) SetOption {
	return func(o *Options) { o.WorkersCount = count }
}

// WithLostConnCallback sets the callback function to be invoked
// when a connection is lost.
func WithLostConnCallback(callback LostConnCallback) SetOption {
	return func(o *Options) { o.LostConnCallback = callback }
}

// WithConnReceiver configures the ConnReceiver with the specified options.
func WithConnReceiver(ops ...SetConnReceiverOption) SetOption {
	return func(o *Options) { o.ConnReceiver = ops }
}

// WithTLSConfig sets the TLS configuration.
func WithTLSConfig(conf *tls.Config) SetOption {
	return func(o *Options) { o.TLSConfig = conf }
}

func Apply(o *Options, opts ...SetOption) {
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
}
