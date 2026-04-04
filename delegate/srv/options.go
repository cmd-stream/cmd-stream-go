package srv

import "time"

type Options struct {
	ServerInfoSendDuration time.Duration
}

type SetOption func(o *Options)

// WithServerInfoSendDuration specifies how long the server will try to send
// ServerInfo to the client. If == 0, it will try forever.
func WithServerInfoSendDuration(d time.Duration) SetOption {
	return func(o *Options) { o.ServerInfoSendDuration = d }
}

func Apply(o *Options, opts ...SetOption) {
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
}
