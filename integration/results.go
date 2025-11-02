package integration_test

import (
	"github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/ord"
)

func NewResult(lastOne bool) Result {
	return Result{lastOne: lastOne}
}

type Result struct {
	lastOne bool
}

func (r Result) LastOne() bool {
	return r.lastOne
}

func (r Result) MarshalTypedMUS(w mus.Writer) (n int, err error) {
	return ResultMUS.Marshal(r, w)
}

func (r Result) SizeTypedMUS() (size int) {
	return ResultMUS.Size(r)
}

var ResultMUS = resultMUS{}

type resultMUS struct{}

func (s resultMUS) Marshal(result Result, w mus.Writer) (n int, err error) {
	return ord.Bool.Marshal(result.lastOne, w)
}

func (s resultMUS) Unmarshal(r mus.Reader) (result Result, n int, err error) {
	result.lastOne, n, err = ord.Bool.Unmarshal(r)
	return
}

func (s resultMUS) Size(result Result) (size int) {
	return ord.SizeBool(result.lastOne)
}
