package testkit

const ResultSize = 18

type Result struct {
	LastOneFlag bool `json:"last_one"`
}

func (r Result) LastOne() bool {
	return r.LastOneFlag
}
