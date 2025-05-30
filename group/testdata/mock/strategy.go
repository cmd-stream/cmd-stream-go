package mock

import (
	cgrp "github.com/cmd-stream/cmd-stream-go/group"
	"github.com/ymz-ncnk/mok"
)

type NextFn[T any] = func() (cgrp.Client[T], int64)
type SliceFn[T any] = func() []cgrp.Client[T]

func NewDispatchStrategy[T any]() DispatchStrategy[T] {
	return DispatchStrategy[T]{mok.New("DispatchStrategy")}
}

type DispatchStrategy[T any] struct {
	*mok.Mock
}

func (f DispatchStrategy[T]) RegisterSlice(fn SliceFn[T]) DispatchStrategy[T] {
	f.Register("Slice", fn)
	return f
}

func (f DispatchStrategy[T]) RegisterNSlice(n int, fn SliceFn[T]) DispatchStrategy[T] {
	f.RegisterN("Slice", n, fn)
	return f
}

func (f DispatchStrategy[T]) RegisterNext(fn NextFn[T]) DispatchStrategy[T] {
	f.Register("Next", fn)
	return f
}

func (f DispatchStrategy[T]) Next() (client cgrp.Client[T], index int64) {
	result, err := f.Call("Next")
	if err != nil {
		panic(err)
	}
	client = result[0].(cgrp.Client[T])
	index = result[1].(int64)
	return
}

func (f DispatchStrategy[T]) Slice() []cgrp.Client[T] {
	result, err := f.Call("Slice")
	if err != nil {
		panic(err)
	}
	return result[0].([]cgrp.Client[T])
}
