// Package codec defines a generic interface for encoding and decoding data
// exchanged over transport connection in cmd-stream-go.
//
// It provides the Codec interface, which abstracts how Commands and Results
// are serialized and deserialized when sent between client and server.
package codec

import (
	"github.com/cmd-stream/transport-go"
)

// Codec defines a generic interface for encoding and decoding data transmitted
// over a transport connection.
type Codec[T, V any] interface {
	Encode(t T, w transport.Writer) (n int, err error)
	Decode(r transport.Reader) (v V, n int, err error)
}
