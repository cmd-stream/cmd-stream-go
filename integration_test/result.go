package intest

import muss "github.com/mus-format/mus-stream-go"

type Result struct {
	lastOne bool
}

func (r Result) LastOne() bool {
	return r.lastOne
}

func (r Result) MarshalTypedMUS(w muss.Writer) (n int, err error) {
	return ResultMUS.Marshal(r, w)
}

func (r Result) SizeTypedMUS() (size int) {
	return ResultMUS.Size(r)
}
