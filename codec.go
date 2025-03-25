package cs

import (
	"github.com/cmd-stream/transport-go"
)

// Codec defines an interface for encoding and decoding data in communication.
//
//   - Encode serializes a value of type T and writes it using the provided
//     transport.Writer.
//   - Decode reads data from the transport.Reader and deserializes it into a
//     value of type V.
//
// Implementations of this interface allow for customized serialization
// strategies.
type Codec[T, V any] interface {
	Encode(t T, w transport.Writer) (err error)
	Decode(r transport.Reader) (v V, err error)
}
