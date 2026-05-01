package group

import (
	"github.com/ymz-ncnk/mok"
)

type (
	NextFn[T any]  func() (t T, index int64)
	SliceFn[T any] func() (s []T)
)

func NewDispatchStrategy[T any]() DispatchStrategy[T] {
	return DispatchStrategy[T]{Mock: mok.New("DispatchStrategy")}
}

type DispatchStrategy[T any] struct {
	*mok.Mock
}

func (m DispatchStrategy[T]) RegisterNextN(n int, fn NextFn[T]) DispatchStrategy[T] {
	m.RegisterN("Next", n, fn)
	return m
}

func (m DispatchStrategy[T]) RegisterNext(fn NextFn[T]) DispatchStrategy[T] {
	m.Register("Next", fn)
	return m
}

func (m DispatchStrategy[T]) RegisterSliceN(n int, fn SliceFn[T]) DispatchStrategy[T] {
	m.RegisterN("Slice", n, fn)
	return m
}

func (m DispatchStrategy[T]) RegisterSlice(fn SliceFn[T]) DispatchStrategy[T] {
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
