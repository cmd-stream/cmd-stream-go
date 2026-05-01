package mock

import (
	hks "github.com/cmd-stream/cmd-stream-go/sender/hooks"
	"github.com/ymz-ncnk/mok"
)

type New[T any] func() hks.Hooks[T]

func NewFactory[T any]() Factory[T] {
	return Factory[T]{mok.New("Factory")}
}

type Factory[T any] struct {
	*mok.Mock
}

func (m Factory[T]) RegisterNew(fn New[T]) Factory[T] {
	m.Register("New", fn)
	return m
}

func (m Factory[T]) New() hks.Hooks[T] {
	vals, err := m.Call("New")
	if err != nil {
		panic(err)
	}
	return vals[0].(hks.Hooks[T])
}
