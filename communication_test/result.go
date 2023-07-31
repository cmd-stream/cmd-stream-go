package ct

import (
	muss "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/ord"
)

type Result struct {
	lastOne bool
}

func (r Result) LastOne() bool {
	return r.lastOne
}

func MarshalResultMUS(result Result, w muss.Writer) (n int, err error) {
	return ord.MarshalBool(result.lastOne, w)
}

func UnmarshalResultMUS(r muss.Reader) (result Result, n int, err error) {
	result.lastOne, n, err = ord.UnmarshalBool(r)
	return
}

func SizeResultMUS(result Result) (size int) {
	return ord.SizeBool(result.lastOne)
}
