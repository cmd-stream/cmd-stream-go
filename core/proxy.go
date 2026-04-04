package core

import (
	"net"
	"time"

	"github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/varint"
)

// Proxy represents a server transport proxy, enabling Commands to send Results
// back.
//
// Implementation of this interface must be thread-safe.
type Proxy interface {
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Send(seq Seq, result Result) (n int, err error)
	SendWithDeadline(deadline time.Time, seq Seq, result Result) (n int, err error)
}

// -----------------------------------------------------------------------------

// Seq represents the sequence number of a Command.
//
// The sequence number ensures that each Command can be uniquely identified and
// mapped to its corresponding Results.
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
