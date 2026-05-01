package test

import (
	"errors"
	"testing"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/cmd-stream/cmd-stream-go/group"
	cmock "github.com/cmd-stream/cmd-stream-go/test/mock/core"
	gmock "github.com/cmd-stream/cmd-stream-go/test/mock/group"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

type GroupTestCase[T any] struct {
	Name     string
	Strategy group.DispatchStrategy[group.Client[T]]
	Action   func(t *testing.T, g group.Group[T])
	Mocks    []*mok.Mock
}

func RunGroupTestCase[T any](t *testing.T, tc GroupTestCase[T]) {
	t.Run(tc.Name, func(t *testing.T) {
		g := group.New[T](tc.Strategy)
		tc.Action(t, g)
		_ = g.Close()
		select {
		case <-g.Done():
		case <-time.After(time.Second):
			t.Fatal("timeout waiting for group to be done")
		}
		asserterror.EqualDeep(t, mok.CheckCalls(tc.Mocks), mok.EmptyInfomap)
	})
}

type GroupSuite[T any] struct{}

// -----------------------------------------------------------------------------
// Test Cases
// -----------------------------------------------------------------------------

func (GroupSuite[T]) Send(t *testing.T) GroupTestCase[T] {
	name := "Should successfully send command using Group"

	var (
		client1         = gmock.NewClient[T]()
		client1DoneChan = make(chan struct{})
		client2         = gmock.NewClient[T]()
		client2DoneChan = make(chan struct{})
		strategy        = gmock.NewDispatchStrategy[group.Client[T]]()
		wantCmd         = cmock.NewCmd[T]()
		results         = make(chan core.AsyncResult, 1)
		wantSeq         = core.Seq(1)
		wantN           = 10
	)
	strategy.RegisterNext(
		func() (t group.Client[T], index int64) { return client1, 0 },
	).RegisterSliceN(2,
		func() []group.Client[T] { return []group.Client[T]{client1, client2} },
	)
	// client1 is first in round-robin
	client1.RegisterSend(
		func(cmd core.Cmd[T], r chan<- core.AsyncResult) (core.Seq, int, error) {
			asserterror.EqualDeep(t, cmd, wantCmd)
			return wantSeq, wantN, nil
		},
	).RegisterClose(
		func() error {
			close(client1DoneChan)
			return nil
		},
	).RegisterDone(
		func() <-chan struct{} {
			return client1DoneChan
		},
	)
	client2.RegisterClose(
		func() error {
			close(client2DoneChan)
			return nil
		},
	).RegisterDone(
		func() <-chan struct{} {
			return client2DoneChan
		},
	)
	return GroupTestCase[T]{
		Name:     name,
		Strategy: strategy,
		Action: func(t *testing.T, g group.Group[T]) {
			seq, clientID, n, err := g.Send(wantCmd, results)
			asserterror.Equal(t, seq, wantSeq)
			asserterror.Equal(t, int(clientID), 0)
			asserterror.Equal(t, n, wantN)
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock},
	}
}

func (GroupSuite[T]) SendWithDeadline(t *testing.T) GroupTestCase[T] {
	name := "Should successfully send command with deadline using Group"

	var (
		client1         = gmock.NewClient[T]()
		client1DoneChan = make(chan struct{})
		client2         = gmock.NewClient[T]()
		client2DoneChan = make(chan struct{})
		strategy        = gmock.NewDispatchStrategy[group.Client[T]]()
		wantCmd         = cmock.NewCmd[T]()
		results         = make(chan core.AsyncResult, 1)
		wantSeq         = core.Seq(1)
		wantN           = 10
		wantDeadline    = time.Now().Add(time.Second)
	)
	strategy.RegisterNext(
		func() (t group.Client[T], index int64) { return client1, 0 },
	).RegisterSliceN(2,
		func() []group.Client[T] { return []group.Client[T]{client1, client2} },
	)
	// client1 is first in round-robin
	client1.RegisterSendWithDeadline(
		func(deadline time.Time, cmd core.Cmd[T], r chan<- core.AsyncResult) (core.Seq, int, error) {
			asserterror.EqualDeep(t, cmd, wantCmd)
			return wantSeq, wantN, nil
		},
	).RegisterClose(
		func() error {
			close(client1DoneChan)
			return nil
		},
	).RegisterDone(
		func() <-chan struct{} { return client1DoneChan },
	)
	client2.RegisterClose(func() error {
		close(client2DoneChan)
		return nil
	}).RegisterDone(func() <-chan struct{} {
		return client2DoneChan
	})
	return GroupTestCase[T]{
		Name:     name,
		Strategy: strategy,
		Action: func(t *testing.T, g group.Group[T]) {
			seq, clientID, n, err := g.SendWithDeadline(wantDeadline, wantCmd, results)
			asserterror.Equal(t, seq, wantSeq)
			asserterror.Equal(t, int(clientID), 0)
			asserterror.Equal(t, n, wantN)
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock},
	}
}

func (GroupSuite[T]) Has(t *testing.T) GroupTestCase[T] {
	name := "Group.Has should return true when client has cmd"

	var (
		client1         = gmock.NewClient[T]()
		client1DoneChan = make(chan struct{})
		client2         = gmock.NewClient[T]()
		client2DoneChan = make(chan struct{})
		strategy        = gmock.NewDispatchStrategy[group.Client[T]]()
		seq             = core.Seq(1)
	)
	strategy.RegisterSliceN(3,
		func() []group.Client[T] { return []group.Client[T]{client1, client2} },
	)
	client1.RegisterHas(
		func(s core.Seq) bool { return true },
	).RegisterClose(
		func() error {
			close(client1DoneChan)
			return nil
		},
	).RegisterDone(
		func() <-chan struct{} { return client1DoneChan },
	)
	client2.RegisterClose(
		func() error {
			close(client2DoneChan)
			return nil
		},
	).RegisterDone(
		func() <-chan struct{} { return client2DoneChan },
	)
	return GroupTestCase[T]{
		Name:     name,
		Strategy: strategy,
		Action: func(t *testing.T, g group.Group[T]) {
			asserterror.Equal(t, g.Has(seq, 0), true)
		},
		Mocks: []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock},
	}
}

func (GroupSuite[T]) Forget(t *testing.T) GroupTestCase[T] {
	name := "Group.Forget should call Forget on the correct client"

	var (
		client1         = gmock.NewClient[T]()
		client1DoneChan = make(chan struct{})
		client2         = gmock.NewClient[T]()
		client2DoneChan = make(chan struct{})
		strategy        = gmock.NewDispatchStrategy[group.Client[T]]()
		seq             = core.Seq(1)
	)
	strategy.RegisterSliceN(3,
		func() []group.Client[T] { return []group.Client[T]{client1, client2} },
	)
	client1.RegisterForget(
		func(s core.Seq) {},
	).RegisterClose(
		func() error {
			close(client1DoneChan)
			return nil
		},
	).RegisterDone(
		func() <-chan struct{} { return client1DoneChan },
	)
	client2.RegisterClose(
		func() error {
			close(client2DoneChan)
			return nil
		},
	).RegisterDone(
		func() <-chan struct{} { return client2DoneChan },
	)
	return GroupTestCase[T]{
		Name:     name,
		Strategy: strategy,
		Action: func(t *testing.T, g group.Group[T]) {
			g.Forget(seq, 0)
		},
		Mocks: []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock},
	}
}

func (GroupSuite[T]) Error(t *testing.T) GroupTestCase[T] {
	name := "Group.Error should return joint error of all clients"

	var (
		err1            = errors.New("error 1")
		err2            = errors.New("error 2")
		client1         = gmock.NewClient[T]()
		client1DoneChan = make(chan struct{})
		client2         = gmock.NewClient[T]()
		client2DoneChan = make(chan struct{})
		strategy        = gmock.NewDispatchStrategy[group.Client[T]]()
	)
	strategy.RegisterSliceN(3,
		func() []group.Client[T] { return []group.Client[T]{client1, client2} },
	)
	client1.RegisterError(
		func() error { return err1 },
	).RegisterClose(
		func() error {
			close(client1DoneChan)
			return nil
		},
	).RegisterDone(
		func() <-chan struct{} { return client1DoneChan },
	)
	client2.RegisterError(
		func() error { return err2 },
	).RegisterClose(
		func() error {
			close(client2DoneChan)
			return nil
		},
	).RegisterDone(
		func() <-chan struct{} { return client2DoneChan },
	)
	return GroupTestCase[T]{
		Name:     name,
		Strategy: strategy,
		Action: func(t *testing.T, g group.Group[T]) {
			err := g.Error()
			asserterror.EqualError(t, err, errors.Join(err1, err2))
		},
		Mocks: []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock},
	}
}

func (GroupSuite[T]) Close(t *testing.T) GroupTestCase[T] {
	name := "Should successfully close all clients in Group"

	var (
		client1         = gmock.NewClient[T]()
		client1DoneChan = make(chan struct{})
		client2         = gmock.NewClient[T]()
		client2DoneChan = make(chan struct{})
		strategy        = gmock.NewDispatchStrategy[group.Client[T]]()
	)
	strategy.RegisterSliceN(2,
		func() []group.Client[T] { return []group.Client[T]{client1, client2} },
	)
	client1.RegisterClose(
		func() error {
			close(client1DoneChan)
			return nil
		},
	).RegisterDone(
		func() <-chan struct{} {
			return client1DoneChan
		},
	)
	client2.RegisterClose(
		func() error {
			close(client2DoneChan)
			return nil
		},
	).RegisterDone(
		func() <-chan struct{} {
			return client2DoneChan
		},
	)
	return GroupTestCase[T]{
		Name:     name,
		Strategy: strategy,
		Action:   func(t *testing.T, g group.Group[T]) {},
		Mocks:    []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock},
	}
}
