package results

import (
	muss "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/ord"
)

var ResultMUS = resultMUS{}

type resultMUS struct{}

func (s resultMUS) Marshal(result Result, w muss.Writer) (n int, err error) {
	return ord.Bool.Marshal(result.lastOne, w)
}

func (s resultMUS) Unmarshal(r muss.Reader) (result Result, n int, err error) {
	result.lastOne, n, err = ord.Bool.Unmarshal(r)
	return
}

func (s resultMUS) Size(result Result) (size int) {
	return ord.SizeBool(result.lastOne)
}
