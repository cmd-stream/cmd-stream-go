package results

import "github.com/mus-format/mus-stream-go"

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
