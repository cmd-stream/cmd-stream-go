package hooks

// CircuitBreaker defines the interface for the Circuit Breaker Pattern.
type CircuitBreaker interface {
	Allow() bool
	Fail()
	Success()
}
