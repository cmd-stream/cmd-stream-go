package cln

type Options struct {
	UnexpectedResultCallback UnexpectedResultCallback
}

type SetOption func(o *Options)

// WithUnexpectedResultCallback sets the the unexpected result callback.
func WithUnexpectedResultCallback(callback UnexpectedResultCallback) SetOption {
	return func(o *Options) { o.UnexpectedResultCallback = callback }
}

func Apply(o *Options, opts ...SetOption) {
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
}
