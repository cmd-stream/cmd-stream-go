package cln

// Options represents the client configuration options.
type Options struct {
	UnexpectedResultCallback UnexpectedResultCallback
}

// SetOption defines a function for configuring Options.
type SetOption func(o *Options)

// WithUnexpectedResultCallback sets the the unexpected result callback.
func WithUnexpectedResultCallback(callback UnexpectedResultCallback) SetOption {
	return func(o *Options) { o.UnexpectedResultCallback = callback }
}

// Apply applies the given options to the Options struct.
func Apply(o *Options, opts ...SetOption) {
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
}
