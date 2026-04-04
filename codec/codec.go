// Package codec defines base interfaces and utilities for encoding and
// decoding data within cmd-stream.
package codec

import (
	tspt "github.com/cmd-stream/cmd-stream-go/transport"
)

// Codec defines a generic interface for encoding and decoding data.
// In cmd-stream, it serves as a high-level abstraction that doesn't concern
// itself with transport-specific details like sequence numbers.
type Codec[T, V any] interface {
	Encode(t T, w tspt.Writer) (n int, err error)
	Decode(r tspt.Reader) (v V, n int, err error)
}
