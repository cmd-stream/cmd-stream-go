package srv

import "time"

// ConnReceiverOptions represents the configuration for a ConnReceiver.
type ConnReceiverOptions struct {
	FirstConnTimeout time.Duration
}

// SetConnReceiverOption defines a function for configuring ConnReceiverOptions.
type SetConnReceiverOption func(o *ConnReceiverOptions)

// WithFirstConnTimeout sets the timeout for the first connection attempt.
func WithFirstConnTimeout(d time.Duration) SetConnReceiverOption {
	return func(o *ConnReceiverOptions) { o.FirstConnTimeout = d }
}

// ApplyConnReceiver applies the given options to the ConnReceiverOptions struct.
func ApplyConnReceiver(o *ConnReceiverOptions, opts ...SetConnReceiverOption) {
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
}
