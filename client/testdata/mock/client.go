package mock

import (
	"time"

	"github.com/cmd-stream/base-go"
	"github.com/ymz-ncnk/mok"
)

type SendFn[T any] = func(cmd base.Cmd[T], results chan<- base.AsyncResult) (seq base.Seq, err error)
type SendWithDeadlineFn[T any] = func(deadline time.Time, cmd base.Cmd[T],
	results chan<- base.AsyncResult) (seq base.Seq, err error)
type HasFn = func(seq base.Seq) bool
type ForgetFn = func(seq base.Seq)
type DoneFn = func() <-chan struct{}
type ErrFn = func() (err error)
type CloseFn = func() (err error)

func NewClient[T any]() Client[T] {
	return Client[T]{mok.New("Client")}
}

type Client[T any] struct {
	*mok.Mock
}

func (c Client[T]) RegisterSend(fn SendFn[T]) Client[T] {
	c.Register("Send", fn)
	return c
}

func (c Client[T]) RegisterSendWithDeadline(fn SendWithDeadlineFn[T]) Client[T] {
	c.Register("SendWithDeadline", fn)
	return c
}

func (c Client[T]) RegisterHas(fn HasFn) Client[T] {
	c.Register("Has", fn)
	return c
}

func (c Client[T]) RegisterForget(fn ForgetFn) Client[T] {
	c.Register("Forget", fn)
	return c
}

func (c Client[T]) RegisterDone(fn DoneFn) Client[T] {
	c.Register("Done", fn)
	return c
}

func (c Client[T]) RegisterErr(fn ErrFn) Client[T] {
	c.Register("Err", fn)
	return c
}

func (c Client[T]) RegisterClose(fn CloseFn) Client[T] {
	c.Register("Close", fn)
	return c
}

func (c Client[T]) Send(cmd base.Cmd[T], results chan<- base.AsyncResult) (
	seq base.Seq, err error) {
	result, err := c.Call("Send", mok.SafeVal[base.Cmd[T]](cmd), results)
	if err != nil {
		panic(err)
	}
	seq = result[0].(base.Seq)
	err, _ = result[1].(error)
	return
}

func (c Client[T]) SendWithDeadline(deadline time.Time, cmd base.Cmd[T],
	results chan<- base.AsyncResult) (seq base.Seq, err error) {
	result, err := c.Call("SendWithDeadline", deadline,
		mok.SafeVal[base.Cmd[T]](cmd), results)
	if err != nil {
		panic(err)
	}
	seq = result[0].(base.Seq)
	err, _ = result[1].(error)
	return
}

func (c Client[T]) Has(seq base.Seq) bool {
	result, err := c.Call("Has", seq)
	if err != nil {
		panic(err)
	}
	return result[0].(bool)
}

func (c Client[T]) Forget(seq base.Seq) {
	_, err := c.Call("Forget", seq)
	if err != nil {
		panic(err)
	}
}

func (c Client[T]) Done() <-chan struct{} {
	result, err := c.Call("Done")
	if err != nil {
		return nil
	}
	return result[0].(<-chan struct{})
}

func (c Client[T]) Err() (err error) {
	result, err := c.Call("Err")
	if err != nil {
		panic(err)
	}
	err, _ = result[0].(error)
	return
}

func (c Client[T]) Close() (err error) {
	result, err := c.Call("Close")
	if err != nil {
		panic(err)
	}
	err, _ = result[0].(error)
	return
}
