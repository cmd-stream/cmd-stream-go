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

func NewGroupClient[T any]() GroupClient[T] {
	return GroupClient[T]{Mock: mok.New("GroupClient")}
}

type GroupClient[T any] struct {
	*mok.Mock
}

func (m GroupClient[T]) RegisterSendN(n int, fn Send) GroupClient[T] {
	m.RegisterN("Send", n, fn)
	return m
}

func (m GroupClient[T]) RegisterSend(fn Send) GroupClient[T] {
	m.Register("Send", fn)
	return m
}

func (m GroupClient[T]) RegisterSendWithDeadlineN(n int, fn SendWithDeadline) GroupClient[T] {
	m.RegisterN("SendWithDeadline", n, fn)
	return m
}

func (m GroupClient[T]) RegisterSendWithDeadline(fn SendWithDeadline) GroupClient[T] {
	m.Register("SendWithDeadline", fn)
	return m
}

func (m GroupClient[T]) RegisterHasN(n int, fn Has) GroupClient[T] {
	m.RegisterN("Has", n, fn)
	return m
}

func (m GroupClient[T]) RegisterHas(fn Has) GroupClient[T] {
	m.Register("Has", fn)
	return m
}

func (m GroupClient[T]) RegisterForgetN(n int, fn Forget) GroupClient[T] {
	m.RegisterN("Forget", n, fn)
	return m
}

func (m GroupClient[T]) RegisterForget(fn Forget) GroupClient[T] {
	m.Register("Forget", fn)
	return m
}

func (m GroupClient[T]) RegisterErrorN(n int, fn Error) GroupClient[T] {
	m.RegisterN("Error", n, fn)
	return m
}

func (m GroupClient[T]) RegisterError(fn Error) GroupClient[T] {
	m.Register("Error", fn)
	return m
}

func (m GroupClient[T]) RegisterCloseN(n int, fn Close) GroupClient[T] {
	m.RegisterN("Close", n, fn)
	return m
}

func (m GroupClient[T]) RegisterClose(fn Close) GroupClient[T] {
	m.Register("Close", fn)
	return m
}

func (m GroupClient[T]) RegisterDoneN(n int, fn Done) GroupClient[T] {
	m.RegisterN("Done", n, fn)
	return m
}

func (m GroupClient[T]) RegisterDone(fn Done) GroupClient[T] {
	m.Register("Done", fn)
	return m
}

func (m GroupClient[T]) Send(cmd core.Cmd[T], results chan<- core.AsyncResult) (seq core.Seq, n int, err error) {
	vals, err := m.Call("Send", cmd, results)
	if err != nil {
		panic(err)
	}
	seq = vals[0].(core.Seq)
	n = vals[1].(int)
	err, _ = vals[2].(error)
	return
}

func (m GroupClient[T]) SendWithDeadline(deadline time.Time, cmd core.Cmd[T], results chan<- core.AsyncResult) (seq core.Seq, n int, err error) {
	vals, err := m.Call("SendWithDeadline", deadline, cmd, results)
	if err != nil {
		panic(err)
	}
	seq = vals[0].(core.Seq)
	n = vals[1].(int)
	err, _ = vals[2].(error)
	return
}

func (m GroupClient[T]) Has(seq core.Seq) bool {
	vals, err := m.Call("Has", seq)
	if err != nil {
		panic(err)
	}
	return vals[0].(bool)
}

func (m GroupClient[T]) Forget(seq core.Seq) {
	_, err := m.Call("Forget", seq)
	if err != nil {
		panic(err)
	}
}

func (m GroupClient[T]) Error() (err error) {
	vals, err := m.Call("Error")
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (m GroupClient[T]) Close() (err error) {
	vals, err := m.Call("Close")
	if err != nil {
		panic(err)
	}
	err, _ = vals[0].(error)
	return
}

func (m GroupClient[T]) Done() <-chan struct{} {
	vals, err := m.Call("Done")
	if err != nil {
		panic(err)
	}
	return vals[0].(<-chan struct{})
}
