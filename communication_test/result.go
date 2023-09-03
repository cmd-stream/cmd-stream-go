package ct

type Result struct {
	lastOne bool
}

func (r Result) LastOne() bool {
	return r.lastOne
}
