package intest

import muss "github.com/mus-format/mus-stream-go"

type Result struct {
	lastOne bool
}

func (r Result) LastOne() bool {
	return r.lastOne
}

func (r Result) MarshalMUS(w muss.Writer) (n int, err error) {
	return ResultMUS.Marshal(r, w)
}
