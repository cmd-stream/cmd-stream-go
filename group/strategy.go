package group

// DispatchStrategy is a dispatch strategy.
type DispatchStrategy[T any] interface {
	Next() (t T, index int64)
	Slice() []T
}

// DispatchStrategyFactory is a factory for a dispatch strategy.
type DispatchStrategyFactory[T any] interface {
	New(clients []Client[T]) DispatchStrategy[Client[T]]
}
