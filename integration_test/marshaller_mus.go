package intest

import muss "github.com/mus-format/mus-stream-go"

type MarshallerMUS interface {
	MarshalMUS(w muss.Writer) (n int, err error)
}
