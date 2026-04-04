package srv

import "time"

type ConnReceiverOptions struct {
	FirstConnTimeout time.Duration
}

type SetConnReceiverOption func(o *ConnReceiverOptions)

// WithFirstConnTimeout sets the timeout for the first connection attempt.
func WithFirstConnTimeout(d time.Duration) SetConnReceiverOption {
	return func(o *ConnReceiverOptions) { o.FirstConnTimeout = d }
}

func ApplyConnReceiver(o *ConnReceiverOptions, opts ...SetConnReceiverOption) {
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
}
