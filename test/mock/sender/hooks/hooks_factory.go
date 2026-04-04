package hooks

import (
	hks "github.com/cmd-stream/cmd-stream-go/sender/hooks"
	"github.com/ymz-ncnk/mok"
)

type New[T any] func() hks.Hooks[T]

func NewHooksFactory[T any]() HooksFactory[T] {
	return HooksFactory[T]{mok.New("HooksFactory")}
}

type HooksFactory[T any] struct {
	*mok.Mock
}

func (m HooksFactory[T]) RegisterNew(fn New[T]) HooksFactory[T] {
	m.Register("New", fn)
	return m
}

func (m HooksFactory[T]) New() hks.Hooks[T] {
	vals, err := m.Call("New")
	if err != nil {
		panic(err)
	}
	return vals[0].(hks.Hooks[T])
}
