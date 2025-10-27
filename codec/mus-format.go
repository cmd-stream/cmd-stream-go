package codec

import (
	"github.com/cmd-stream/core-go"
	"github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/varint"
)

// SeqMUS is a core.Seq MUS serializer.
var SeqMUS = seqMUS{}

type seqMUS struct{}

func (s seqMUS) Marshal(seq core.Seq, w mus.Writer) (n int, err error) {
	return varint.PositiveInt64.Marshal(int64(seq), w)
}

func (s seqMUS) Unmarshal(r mus.Reader) (seq core.Seq, n int, err error) {
	num, n, err := varint.PositiveInt64.Unmarshal(r)
	seq = core.Seq(num)
	return
}

func (s seqMUS) Size(seq core.Seq) (size int) {
	return varint.PositiveInt64.Size(int64(seq))
}

func (s seqMUS) Skip(r mus.Reader) (n int, err error) {
	return varint.PositiveInt64.Skip(r)
}
