package group

import (
	"github.com/ymz-ncnk/mok"
)

type (
	Next  func() (t any, index int64)
	Slice func() any
)

func NewDispatchStrategy[T any]() DispatchStrategy[T] {
	return DispatchStrategy[T]{Mock: mok.New("DispatchStrategy")}
}

type DispatchStrategy[T any] struct {
	*mok.Mock
}

func (m DispatchStrategy[T]) RegisterNextN(n int, fn Next) DispatchStrategy[T] {
	m.RegisterN("Next", n, fn)
	return m
}

func (m DispatchStrategy[T]) RegisterNext(fn Next) DispatchStrategy[T] {
	m.Register("Next", fn)
	return m
}

func (m DispatchStrategy[T]) RegisterSliceN(n int, fn Slice) DispatchStrategy[T] {
	m.RegisterN("Slice", n, fn)
	return m
}

func (m DispatchStrategy[T]) RegisterSlice(fn Slice) DispatchStrategy[T] {
	m.Register("Slice", fn)
	return m
}

func (m DispatchStrategy[T]) Next() (t T, index int64) {
	vals, err := m.Call("Next")
	if err != nil {
		panic(err)
	}
	t = vals[0].(T)
	index = vals[1].(int64)
	return
}

func (m DispatchStrategy[T]) Slice() []T {
	vals, err := m.Call("Slice")
	if err != nil {
		panic(err)
	}
	return vals[0].([]T)
}
