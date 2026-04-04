package handler

import "time"

type Options struct {
	CmdReceiveDuration time.Duration
	At                 bool
}

// WithCmdReceiveDuration sets the maximum time the Handler will wait for a
// command. If no command arrives within this duration, the connection is
// closed. A duration of 0 means the Handler will wait indefinitely.
func WithCmdReceiveDuration(d time.Duration) SetOption {
	return func(o *Options) { o.CmdReceiveDuration = d }
}

// WithAt enables the "at" flag. When enabled, the Handler passes the command's
// received timestamp to Invoker.Invoke().
func WithAt() SetOption {
	return func(o *Options) { o.At = true }
}

type SetOption func(o *Options)

func Apply(o *Options, opts ...SetOption) {
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
}
