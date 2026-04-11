package group

import (
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/ymz-ncnk/mok"
)

type (
	Send             func(cmd core.Cmd[any], results chan<- core.AsyncResult) (seq core.Seq, n int, err error)
	SendWithDeadline func(deadline time.Time, cmd core.Cmd[any], results chan<- core.AsyncResult) (seq core.Seq, n int, err error)
	Has              func(seq core.Seq) bool
	Forget           func(seq core.Seq)
	Error            func() (err error)
	Close            func() (err error)
	Done             func() <-chan struct{}
)

func NewClient[T any]() Client[T] {
	return Client[T]{Mock: mok.New("Client")}
}

type Client[T any] struct {
	*mok.Mock
}

func (m Client[T]) RegisterSendN(n int, fn Send) Client[T] {
	m.RegisterN("Send", n, fn)
	return m
}

func (m Client[T]) RegisterSend(fn Send) Client[T] {
	m.Register("Send", fn)
	return m
}

func (m Client[T]) RegisterSendWithDeadlineN(n int, fn SendWithDeadline) Client[T] {
	m.RegisterN("SendWithDeadline", n, fn)
	return m
}

func (m Client[T]) RegisterSendWithDeadline(fn SendWithDeadline) Client[T] {
	m.Register("SendWithDeadline", fn)
	return m
}

func (m Client[T]) RegisterHasN(n int, fn Has) Client[T] {
	m.RegisterN("Has", n, fn)
	return m
}

func (m Client[T]) RegisterHas(fn Has) Client[T] {
	m.Register("Has", fn)
	return m
}

func (m Client[T]) RegisterForgetN(n int, fn Forget) Client[T] {
	m.RegisterN("Forget", n, fn)
	return m
}

func (m Client[T]) RegisterForget(fn Forget) Client[T] {
	m.Register("Forget", fn)
	return m
}

func (m Client[T]) RegisterErrorN(n int, fn Error) Client[T] {
	m.RegisterN("Error", n, fn)
	return m
}

func (m Client[T]) RegisterError(fn Error) Client[T] {
	m.Register("Error", fn)
	return m
}

func (m Client[T]) RegisterCloseN(n int, fn Close) Client[T] {
	m.RegisterN("Close", n, fn)
	return m
}

func (m Client[T]) RegisterClose(fn Close) Client[T] {
	m.Register("Close", fn)
	return m
}

func (m Client[T]) RegisterDoneN(n int, fn Done) Client[T] {
	m.RegisterN("Done", n, fn)
	return m
}

func (m Client[T]) RegisterDone(fn Done) Client[T] {
	m.Register("Done", fn)
	return m
}

func (m Client[T]) Send(cmd core.Cmd[T], results chan<- core.AsyncResult) (seq core.Seq, n int, err error) {
	vals, err := m.Call("Send", cmd, results)
	if err != nil {
		panic(err)
	}
	seq = vals[0].(core.Seq)
	n = vals[1].(int)
	err, _ = vals[2].(error)
	return
}

func (m Client[T]) SendWithDeadline(deadline time.Time, cmd core.Cmd[T], results chan<- core.AsyncResult) (seq core.Seq, n int, err error) {
	vals, err := m.Call("SendWithDeadline", deadline, cmd, results)
	if err != nil {
		panic(err)
	}
	seq = vals[0].(core.Seq)
	n = vals[1].(int)
	err, _ = vals[2].(error)
	return
}

func (m Client[T]) Has(seq core.Seq) bool {
	vals, err := m.Call("Has", seq)
	if err != nil {
		panic(err)
	}
	return vals[0].(bool)
}

func (m Client[T]) Forget(seq core.Seq) {
	_, err := m.Call("Forget", seq)
	if err != nil {
		panic(err)
	}
}

func (m Client[T]) Error() (err error) {
	vals, err := m.Call("Error")
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (m Client[T]) Close() (err error) {
	vals, err := m.Call("Close")
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (m Client[T]) Done() <-chan struct{} {
	vals, err := m.Call("Done")
	if err != nil {
		panic(err)
	}
	return vals[0].(<-chan struct{})
}
