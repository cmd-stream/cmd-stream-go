package core

import (
	"net"
	"time"

	"github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/varint"
)

// Proxy provides an interface for Commands to interact with the server
// transport and send results back to the client.
//
// A Proxy instance is typically short-lived, created specifically for a single
// Command execution. It encapsulates the network context and sequence
// information required for Result routing.
//
// Implementations must be thread-safe, as a Command may execute asynchronously
// and send Results from multiple goroutines.
type Proxy interface {
	// LocalAddr returns the local network address.
	LocalAddr() net.Addr
	// RemoteAddr returns the remote network address.
	RemoteAddr() net.Addr
	// ReceivedAt returns the time when the Command was received by the server.
	// Returns a zero time if the information is not available.
	ReceivedAt() time.Time
	// Send sends a Result back to the client.
	Send(result Result) (n int, err error)
	// SendWithDeadline sends a Result back to the client with a specified deadline.
	SendWithDeadline(deadline time.Time, result Result) (n int, err error)
	// Seq returns the sequence number of the Command.
	Seq() Seq
}

// -----------------------------------------------------------------------------

// Seq represents the sequence number of a Command.
//
// Sequence numbers provide a unique identifier to map Results back to their
// original Commands.
type Seq int64

// -----------------------------------------------------------------------------

// SeqMUS is a Seq MUS serializer.
var SeqMUS = seqMUS{}

type seqMUS struct{}

func (s seqMUS) Marshal(seq Seq, w mus.Writer) (n int, err error) {
	return varint.PositiveInt64.Marshal(int64(seq), w)
}

func (s seqMUS) Unmarshal(r mus.Reader) (seq Seq, n int, err error) {
	num, n, err := varint.PositiveInt64.Unmarshal(r)
	seq = Seq(num)
	return
}

func (s seqMUS) Size(seq Seq) (size int) {
	return varint.PositiveInt64.Size(int64(seq))
}

func (s seqMUS) Skip(r mus.Reader) (n int, err error) {
	return varint.PositiveInt64.Skip(r)
}
