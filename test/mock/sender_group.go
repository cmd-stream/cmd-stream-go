package mock

import (
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	grp "github.com/cmd-stream/cmd-stream-go/group"
	"github.com/ymz-ncnk/mok"
)

type Group[T any] struct {
	*mok.Mock
}

func NewGroup[T any]() Group[T] {
	return Group[T]{mok.New("Group")}
}

func (m Group[T]) Send(cmd core.Cmd[T], results chan<- core.AsyncResult) (
	seq core.Seq, clientID grp.ClientID, n int, err error,
) {
	vals, err := m.Call("Send", cmd, results)
	if err != nil {
		return
	}
	seq = vals[0].(core.Seq)
	clientID = vals[1].(grp.ClientID)
	n = vals[2].(int)
	err, _ = vals[3].(error)
	return
}

func (m Group[T]) SendWithDeadline(deadline time.Time, cmd core.Cmd[T],
	results chan<- core.AsyncResult,
) (seq core.Seq, clientID grp.ClientID, n int, err error) {
	vals, err := m.Call("SendWithDeadline", deadline, cmd, results)
	if err != nil {
		return
	}
	seq = vals[0].(core.Seq)
	clientID = vals[1].(grp.ClientID)
	n = vals[2].(int)
	err, _ = vals[3].(error)
	return
}

func (m Group[T]) Has(seq core.Seq, clientID grp.ClientID) (ok bool) {
	vals, err := m.Call("Has", seq, clientID)
	if err != nil {
		return
	}
	ok = vals[0].(bool)
	return
}

func (m Group[T]) Forget(seq core.Seq, clientID grp.ClientID) {
	_, _ = m.Call("Forget", seq, clientID)
}

func (m Group[T]) Done() <-chan struct{} {
	vals, err := m.Call("Done")
	if err != nil {
		return nil
	}
	return vals[0].(<-chan struct{})
}

func (m Group[T]) Error() (err error) {
	vals, err := m.Call("Error")
	if err != nil {
		return
	}
	err, _ = vals[0].(error)
	return
}

func (m Group[T]) Close() (err error) {
	vals, err := m.Call("Close")
	if err != nil {
		return
	}
	err, _ = vals[0].(error)
	return
}

func (m Group[T]) RegisterSend(
	fn func(cmd core.Cmd[T], results chan<- core.AsyncResult) (
		seq core.Seq, clientID grp.ClientID, n int, err error),
) Group[T] {
	m.Register("Send", fn)
	return m
}

func (m Group[T]) RegisterSendWithDeadline(
	fn func(deadline time.Time, cmd core.Cmd[T],
		results chan<- core.AsyncResult,
	) (seq core.Seq, clientID grp.ClientID, n int, err error),
) Group[T] {
	m.Register("SendWithDeadline", fn)
	return m
}

func (m Group[T]) RegisterHas(fn func(seq core.Seq,
	clientID grp.ClientID) (ok bool),
) Group[T] {
	m.Register("Has", fn)
	return m
}

func (m Group[T]) RegisterForget(fn func(seq core.Seq,
	clientID grp.ClientID),
) Group[T] {
	m.Register("Forget", fn)
	return m
}

func (m Group[T]) RegisterDone(fn func() <-chan struct{}) Group[T] {
	m.Register("Done", fn)
	return m
}

func (m Group[T]) RegisterError(fn func() (err error)) Group[T] {
	m.Register("Error", fn)
	return m
}

func (m Group[T]) RegisterClose(fn func() (err error)) Group[T] {
	m.Register("Close", fn)
	return m
}
