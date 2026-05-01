package mock

import "github.com/ymz-ncnk/mok"

type (
	Allow   func() bool
	Fail    func()
	Success func()
)

func NewCircuitBreaker() CircuitBreaker {
	return CircuitBreaker{Mock: mok.New("CircuitBreaker")}
}

type CircuitBreaker struct {
	*mok.Mock
}

func (m CircuitBreaker) RegisterAllowN(n int, fn Allow) CircuitBreaker {
	m.RegisterN("Allow", n, fn)
	return m
}

func (m CircuitBreaker) RegisterAllow(fn Allow) CircuitBreaker {
	m.Register("Allow", fn)
	return m
}

func (m CircuitBreaker) RegisterFailN(n int, fn Fail) CircuitBreaker {
	m.RegisterN("Fail", n, fn)
	return m
}

func (m CircuitBreaker) RegisterFail(fn Fail) CircuitBreaker {
	m.Register("Fail", fn)
	return m
}

func (m CircuitBreaker) RegisterSuccessN(n int, fn Success) CircuitBreaker {
	m.RegisterN("Success", n, fn)
	return m
}

func (m CircuitBreaker) RegisterSuccess(fn Success) CircuitBreaker {
	m.Register("Success", fn)
	return m
}

func (m CircuitBreaker) Allow() bool {
	vals, err := m.Call("Allow")
	if err != nil {
		panic(err)
	}
	return vals[0].(bool)
}

func (m CircuitBreaker) Fail() {
	_, err := m.Call("Fail")
	if err != nil {
		panic(err)
	}
}

func (m CircuitBreaker) Success() {
	_, err := m.Call("Success")
	if err != nil {
		panic(err)
	}
}
