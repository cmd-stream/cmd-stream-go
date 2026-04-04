package transport

type Options struct {
	WriterBufSize int
	ReaderBufSize int
}

type SetOption func(o *Options)

// WithWriterBufSize sets the buffer size for the Writer. If set to <= 0, the
// default bufio.Writer size is used.
func WithWriterBufSize(size int) SetOption {
	return func(o *Options) { o.WriterBufSize = size }
}

// WithReaderBufSize sets the buffer size for the Reader. If set to <= 0, the
// default bufio.Reader size is used.
func WithReaderBufSize(size int) SetOption {
	return func(o *Options) { o.ReaderBufSize = size }
}

func Apply(o *Options, opts ...SetOption) {
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
}
