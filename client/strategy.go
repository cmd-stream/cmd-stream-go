package ccln

import "sync/atomic"

// DispatchStrategy is a dispatch strategy.
type DispatchStrategy[T any] interface {
	Next() (t T, index int64)
	Slice() []T
}

// NewRoundRobinStrategy creates a new RoundRobinStrategy.
func NewRoundRobinStrategy[T any](sl []T) RoundRobinStrategy[T] {
	var i int64
	return RoundRobinStrategy[T]{sl: sl, length: int64(len(sl)), i: &i}
}

// RoundRobinStrategy implements a round-robin dispatch strategy.
type RoundRobinStrategy[T any] struct {
	sl     []T
	length int64
	i      *int64
}

// Next returns the next element and its index in the slice, following a
// round-robin strategy. The index is incremented atomically to ensure
// thread-safety in concurrent environments.
func (s RoundRobinStrategy[T]) Next() (t T, index int64) {
	index = (atomic.AddInt64(s.i, 1) - 1) % s.length
	return s.sl[index], index
}

// Slice returns the slice of elements underlying this RoundRobinStrategy.
func (s RoundRobinStrategy[T]) Slice() []T {
	return s.sl
}
