package cs

import (
	"github.com/cmd-stream/transport-go"
)

// Codec represents a general codec.
type Codec[T, V any] interface {
	Encode(t T, w transport.Writer) (err error)
	Decode(r transport.Reader) (v V, err error)
}
