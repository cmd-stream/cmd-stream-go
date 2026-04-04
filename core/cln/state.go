package cln

import (
	"sync"

	"github.com/cmd-stream/cmd-stream-go/core"
)

const (
	inProgress int = iota
	closed
)

type state struct {
	sync.Mutex
	v int
}

func (s *state) Closed() (isClosed bool) {
	s.Lock()
	isClosed = s.v == closed
	s.Unlock()
	return
}

func (s *state) SetClosed() (ok bool) {
	s.Lock()
	if s.v == closed {
		s.Unlock()
		return
	}
	s.v = closed
	s.Unlock()
	return true
}

// -----------------------------------------------------------------------------

type pending[T any] struct {
	sync.Mutex
	m map[core.Seq]chan<- core.AsyncResult
}

func (p *pending[T]) add(seq core.Seq, results chan<- core.AsyncResult) {
	p.Lock()
	p.m[seq] = results
	p.Unlock()
}

func (p *pending[T]) remove(seq core.Seq) {
	p.Lock()
	delete(p.m, seq)
	p.Unlock()
}

func (p *pending[T]) get(seq core.Seq) (results chan<- core.AsyncResult, pst bool) {
	p.Lock()
	results, pst = p.m[seq]
	p.Unlock()
	return
}

func (p *pending[T]) pop(seq core.Seq) (results chan<- core.AsyncResult, pst bool) {
	p.Lock()
	results, pst = p.m[seq]
	if pst {
		delete(p.m, seq)
	}
	p.Unlock()
	return
}

func (p *pending[T]) rangeAndRemove(fn func(seq core.Seq, results chan<- core.AsyncResult)) {
	p.Lock()
	for seq, results := range p.m {
		fn(seq, results)
		delete(p.m, seq)
	}
	p.Unlock()
}

func (p *pending[T]) failAll(cause error) {
	p.rangeAndRemove(func(seq core.Seq, results chan<- core.AsyncResult) {
		select {
		case results <- core.AsyncResult{Seq: seq, Error: cause}:
		default:
		}
	})
}

// -----------------------------------------------------------------------------

type errStatus struct {
	sync.Mutex
	v error
}

func (e *errStatus) Set(v error) {
	e.Lock()
	e.v = v
	e.Unlock()
}

func (e *errStatus) Get() (err error) {
	e.Lock()
	err = e.v
	e.Unlock()
	return
}
