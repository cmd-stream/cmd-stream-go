package srv

import "time"

// Options represents configuration for ServerInfoDelegate.
type Options struct {
	ServerInfoSendDuration time.Duration
}

// SetOption defines a function for configuring Options.
type SetOption func(o *Options)

// WithServerInfoSendDuration specifies how long the server will try to send
// ServerInfo to the client. If == 0, it will try forever.
func WithServerInfoSendDuration(d time.Duration) SetOption {
	return func(o *Options) { o.ServerInfoSendDuration = d }
}

// Apply applies the given options to the Options struct.
func Apply(o *Options, opts ...SetOption) {
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
}
