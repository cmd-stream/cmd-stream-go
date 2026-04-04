package sender

import (
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	grp "github.com/cmd-stream/cmd-stream-go/group"
	"github.com/ymz-ncnk/mok"
)

type SenderGroup[T any] struct {
	*mok.Mock
}

func NewSenderGroup[T any]() SenderGroup[T] {
	return SenderGroup[T]{mok.New("SenderGroup")}
}

func (m SenderGroup[T]) Send(cmd core.Cmd[T], results chan<- core.AsyncResult) (
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

func (m SenderGroup[T]) SendWithDeadline(deadline time.Time, cmd core.Cmd[T],
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

func (m SenderGroup[T]) Has(seq core.Seq, clientID grp.ClientID) (ok bool) {
	vals, err := m.Call("Has", seq, clientID)
	if err != nil {
		return
	}
	ok = vals[0].(bool)
	return
}

func (m SenderGroup[T]) Forget(seq core.Seq, clientID grp.ClientID) {
	_, _ = m.Call("Forget", seq, clientID)
}

func (m SenderGroup[T]) Done() <-chan struct{} {
	vals, err := m.Call("Done")
	if err != nil {
		return nil
	}
	return vals[0].(<-chan struct{})
}

func (m SenderGroup[T]) Error() (err error) {
	vals, err := m.Call("Error")
	if err != nil {
		return
	}
	err, _ = vals[0].(error)
	return
}

func (m SenderGroup[T]) Close() (err error) {
	vals, err := m.Call("Close")
	if err != nil {
		return
	}
	err, _ = vals[0].(error)
	return
}

func (m SenderGroup[T]) RegisterSend(
	fn func(cmd core.Cmd[T], results chan<- core.AsyncResult) (
		seq core.Seq, clientID grp.ClientID, n int, err error),
) SenderGroup[T] {
	m.Register("Send", fn)
	return m
}

func (m SenderGroup[T]) RegisterSendWithDeadline(
	fn func(deadline time.Time, cmd core.Cmd[T],
		results chan<- core.AsyncResult,
	) (seq core.Seq, clientID grp.ClientID, n int, err error),
) SenderGroup[T] {
	m.Register("SendWithDeadline", fn)
	return m
}

func (m SenderGroup[T]) RegisterHas(fn func(seq core.Seq,
	clientID grp.ClientID) (ok bool),
) SenderGroup[T] {
	m.Register("Has", fn)
	return m
}

func (m SenderGroup[T]) RegisterForget(fn func(seq core.Seq,
	clientID grp.ClientID),
) SenderGroup[T] {
	m.Register("Forget", fn)
	return m
}

func (m SenderGroup[T]) RegisterDone(fn func() <-chan struct{}) SenderGroup[T] {
	m.Register("Done", fn)
	return m
}

func (m SenderGroup[T]) RegisterError(fn func() (err error)) SenderGroup[T] {
	m.Register("Error", fn)
	return m
}

func (m SenderGroup[T]) RegisterClose(fn func() (err error)) SenderGroup[T] {
	m.Register("Close", fn)
	return m
}
