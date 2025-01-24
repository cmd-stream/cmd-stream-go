package cs

import (
	"github.com/cmd-stream/base-go"
	muss "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/varint"
)

// MarshalSeqMUS marshals a sequence number to the MUS format.
func MarshalSeqMUS(seq base.Seq, w muss.Writer) (n int, err error) {
	return varint.MarshalInt64(int64(seq), w)
}

// UnmarshalSeqMUS unmarshals a sequence number from the MUS format.
func UnmarshalSeqMUS(r muss.Reader) (seq base.Seq, n int, err error) {
	num, n, err := varint.UnmarshalInt64(r)
	seq = base.Seq(num)
	return
}

// SizeSeqMUS returns the size of a sequence number in the MUS format.
func SizeSeqMUS(seq base.Seq) (size int) {
	return varint.SizeInt64(int64(seq))
}
