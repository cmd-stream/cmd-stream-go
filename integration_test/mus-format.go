package intest

import (
	com "github.com/mus-format/common-go"
	muss "github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/ord"
)

// -----------------------------------------------------------------------------
// DTM
// -----------------------------------------------------------------------------

const (
	Cmd1DTM com.DTM = iota
	Cmd2DTM
	Cmd3DTM
)

// -----------------------------------------------------------------------------
// Marshal/Unmarshal/Size functions
// -----------------------------------------------------------------------------

// Result

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
