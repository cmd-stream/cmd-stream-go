package cln

import "time"

const (
// KeepaliveTime is the default inactivity period before a ping is sent.
	KeepaliveTime  = 3 * time.Second
// KeepaliveIntvl is the default interval between pings.
	KeepaliveIntvl = time.Second
)

// Options defines the configuration for ClientInfoDelegate.
type Options struct {
	ServerInfoReceiveDuration time.Duration
}

// SetOption is a functional option for configuring ClientInfoDelegate.
type SetOption func(o *Options)

// WithServerInfoReceiveDuration sets the duration the client will wait
// for the ServerInfo. If set to 0, the client waits indefinitely.
func WithServerInfoReceiveDuration(d time.Duration) SetOption {
	return func(o *Options) { o.ServerInfoReceiveDuration = d }
}

// Apply applies the given options to the Options struct.
func Apply(o *Options, opts ...SetOption) {
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
}

// KeepaliveOptions defines the configuration for keepalive behavior.
type KeepaliveOptions struct {
	KeepaliveTime  time.Duration
	KeepaliveIntvl time.Duration
}

// DefaultKeepaliveOptions returns the default keepalive configuration.
func DefaultKeepaliveOptions() KeepaliveOptions {
	return KeepaliveOptions{
		KeepaliveTime:  KeepaliveTime,
		KeepaliveIntvl: KeepaliveIntvl,
	}
}

// SetKeepaliveOption is a functional option for configuring keepalive behavior.
type SetKeepaliveOption func(o *KeepaliveOptions)

// WithKeepaliveTime sets the inactivity period after which the client
// starts sending Ping Commands to the server if no Commands have been sent.
func WithKeepaliveTime(d time.Duration) SetKeepaliveOption {
	return func(o *KeepaliveOptions) { o.KeepaliveTime = d }
}

// WithKeepaliveIntvl sets the interval between consecutive Ping Commands
// sent by the client.
func WithKeepaliveIntvl(d time.Duration) SetKeepaliveOption {
	return func(o *KeepaliveOptions) { o.KeepaliveIntvl = d }
}

// ApplyKeepalive applies the given options to the KeepaliveOptions struct.
func ApplyKeepalive(o *KeepaliveOptions, opts ...SetKeepaliveOption) {
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
}
