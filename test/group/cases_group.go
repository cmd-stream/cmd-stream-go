package group

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

func GroupSendTestCase(t *testing.T) GroupTestCase {
	name := "Should successfully send command using Group"

	var (
		client1         = gmock.NewClient[any]()
		client1DoneChan = make(chan struct{})
		client2         = gmock.NewClient[any]()
		client2DoneChan = make(chan struct{})
		strategy        = gmock.NewDispatchStrategy[group.Client[any]]()
		wantCmd         = cmock.NewCmd[any]()
		results         = make(chan core.AsyncResult, 1)
		wantSeq         = core.Seq(1)
		wantN           = 10
	)
	strategy.RegisterNext(
		func() (t any, index int64) { return client1, 0 },
	).RegisterSliceN(2,
		func() any { return []group.Client[any]{client1, client2} },
	)
	// client1 is first in round-robin
	client1.RegisterSend(
		func(cmd core.Cmd[any], r chan<- core.AsyncResult) (core.Seq, int, error) {
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
	return GroupTestCase{
		Name:     name,
		Strategy: strategy,
		Action: func(t *testing.T, g group.Group[any]) {
			seq, clientID, n, err := g.Send(wantCmd, results)
			asserterror.Equal(t, seq, wantSeq)
			asserterror.Equal(t, int(clientID), 0)
			asserterror.Equal(t, n, wantN)
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock},
	}
}

func GroupSendWithDeadlineTestCase(t *testing.T) GroupTestCase {
	name := "Should successfully send command with deadline using Group"

	var (
		client1         = gmock.NewClient[any]()
		client1DoneChan = make(chan struct{})
		client2         = gmock.NewClient[any]()
		client2DoneChan = make(chan struct{})
		strategy        = gmock.NewDispatchStrategy[group.Client[any]]()
		wantCmd         = cmock.NewCmd[any]()
		results         = make(chan core.AsyncResult, 1)
		wantSeq         = core.Seq(1)
		wantN           = 10
		wantDeadline    = time.Now().Add(time.Second)
	)
	strategy.RegisterNext(
		func() (t any, index int64) { return client1, 0 },
	).RegisterSliceN(2,
		func() any { return []group.Client[any]{client1, client2} },
	)
	// client1 is first in round-robin
	client1.RegisterSendWithDeadline(
		func(deadline time.Time, cmd core.Cmd[any], r chan<- core.AsyncResult) (core.Seq, int, error) {
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
	return GroupTestCase{
		Name:     name,
		Strategy: strategy,
		Action: func(t *testing.T, g group.Group[any]) {
			seq, clientID, n, err := g.SendWithDeadline(wantDeadline, wantCmd, results)
			asserterror.Equal(t, seq, wantSeq)
			asserterror.Equal(t, int(clientID), 0)
			asserterror.Equal(t, n, wantN)
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock},
	}
}

func GroupHasTestCase(t *testing.T) GroupTestCase {
	name := "Group.Has should return true when client has cmd"

	var (
		client1         = gmock.NewClient[any]()
		client1DoneChan = make(chan struct{})
		client2         = gmock.NewClient[any]()
		client2DoneChan = make(chan struct{})
		strategy        = gmock.NewDispatchStrategy[group.Client[any]]()
		seq             = core.Seq(1)
	)
	strategy.RegisterSliceN(3,
		func() any { return []group.Client[any]{client1, client2} },
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
	return GroupTestCase{
		Name:     name,
		Strategy: strategy,
		Action: func(t *testing.T, g group.Group[any]) {
			asserterror.Equal(t, g.Has(seq, 0), true)
		},
		Mocks: []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock},
	}
}

func GroupForgetTestCase(t *testing.T) GroupTestCase {
	name := "Group.Forget should call Forget on the correct client"

	var (
		client1         = gmock.NewClient[any]()
		client1DoneChan = make(chan struct{})
		client2         = gmock.NewClient[any]()
		client2DoneChan = make(chan struct{})
		strategy        = gmock.NewDispatchStrategy[group.Client[any]]()
		seq             = core.Seq(1)
	)
	strategy.RegisterSliceN(3,
		func() any { return []group.Client[any]{client1, client2} },
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
	return GroupTestCase{
		Name:     name,
		Strategy: strategy,
		Action: func(t *testing.T, g group.Group[any]) {
			g.Forget(seq, 0)
		},
		Mocks: []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock},
	}
}

func GroupErrorTestCase(t *testing.T) GroupTestCase {
	name := "Group.Error should return joint error of all clients"

	var (
		err1            = errors.New("error 1")
		err2            = errors.New("error 2")
		client1         = gmock.NewClient[any]()
		client1DoneChan = make(chan struct{})
		client2         = gmock.NewClient[any]()
		client2DoneChan = make(chan struct{})
		strategy        = gmock.NewDispatchStrategy[group.Client[any]]()
	)
	strategy.RegisterSliceN(3,
		func() any { return []group.Client[any]{client1, client2} },
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
	return GroupTestCase{
		Name:     name,
		Strategy: strategy,
		Action: func(t *testing.T, g group.Group[any]) {
			err := g.Error()
			asserterror.EqualError(t, err, errors.Join(err1, err2))
		},
		Mocks: []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock},
	}
}

func GroupCloseTestCase(t *testing.T) GroupTestCase {
	name := "Should successfully close all clients in Group"

	var (
		client1         = gmock.NewClient[any]()
		client1DoneChan = make(chan struct{})
		client2         = gmock.NewClient[any]()
		client2DoneChan = make(chan struct{})
		strategy        = gmock.NewDispatchStrategy[group.Client[any]]()
	)
	strategy.RegisterSliceN(2,
		func() any { return []group.Client[any]{client1, client2} },
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
	return GroupTestCase{
		Name:     name,
		Strategy: strategy,
		Action:   func(t *testing.T, g group.Group[any]) {},
		Mocks:    []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock},
	}
}
