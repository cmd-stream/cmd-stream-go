package ccln_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/cmd-stream/base-go"
	bmock "github.com/cmd-stream/base-go/testdata/mock"
	ccln "github.com/cmd-stream/cmd-stream-go/client"
	"github.com/cmd-stream/cmd-stream-go/client/testdata/mock"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

func TestStreamline(t *testing.T) {

	t.Run("NewStreamline", func(t *testing.T) {
		wg := sync.WaitGroup{}
		wg.Add(2)

		var (
			done1   = make(chan struct{})
			done2   = make(chan struct{})
			client1 = mock.NewClient[any]().RegisterDone(
				func() <-chan struct{} {
					defer wg.Done()
					return done1
				},
			)
			client2 = mock.NewClient[any]().RegisterDone(
				func() <-chan struct{} {
					defer wg.Done()
					return done2
				},
			)
			strategy = mock.NewDispatchStrategy[any]().RegisterSlice(
				func() []ccln.Client[any] {
					return []ccln.Client[any]{client1, client2}
				},
			)
			mocks = []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock}
		)
		grp := ccln.NewGroup(strategy)

		close(done1)
		close(done2)
		wg.Wait()

		select {
		case <-time.NewTimer(time.Second).C:
			t.Error("timeout")
		case <-grp.Done():
		}

		if infomap := mok.CheckCalls(mocks); len(infomap) > 0 {
			t.Error(infomap)
		}
	})

	t.Run("Send", func(t *testing.T) {
		wg := sync.WaitGroup{}
		wg.Add(2)

		var (
			wantErr      error         = nil
			wantSeq      base.Seq      = 10
			wantClientID ccln.ClientID = 1
			wantCmd                    = bmock.NewCmd()
			wantResults                = make(chan base.AsyncResult)

			done1   = make(chan struct{})
			done2   = make(chan struct{})
			client1 = mock.NewClient[any]().RegisterDone(
				func() <-chan struct{} {
					defer wg.Done()
					return done1
				})
			client2 = mock.NewClient[any]().RegisterSend(
				func(cmd base.Cmd[any], results chan<- base.AsyncResult) (seq base.Seq, err error) {
					asserterror.Equal[any](cmd, wantCmd, t)
					asserterror.Equal(results, wantResults, t)
					return wantSeq, nil
				},
			).RegisterDone(func() <-chan struct{} {
				defer wg.Done()
				return done2
			})
			strategy = mock.NewDispatchStrategy[any]().RegisterNext(
				func() (ccln.Client[any], int64) {
					return client2, 1
				},
			).RegisterSlice(
				func() []ccln.Client[any] {
					return []ccln.Client[any]{client1, client2}
				},
			)
			mocks = []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock}
		)
		grp := ccln.NewGroup(strategy)

		seq, clientID, err := grp.Send(wantCmd, wantResults)
		asserterror.EqualError(err, wantErr, t)
		asserterror.Equal(seq, wantSeq, t)
		asserterror.Equal(clientID, wantClientID, t)

		close(done1)
		close(done2)
		wg.Wait()

		if infomap := mok.CheckCalls(mocks); len(infomap) > 0 {
			t.Error(infomap)
		}
	})

	t.Run("SendWithDeadline", func(t *testing.T) {
		wg := sync.WaitGroup{}
		wg.Add(2)

		var (
			wantErr      error         = nil
			wantSeq      base.Seq      = 10
			wantClientID ccln.ClientID = 1
			wantDeadline               = time.Now().Add(time.Second)
			wantCmd                    = bmock.NewCmd()
			wantResults                = make(chan base.AsyncResult)

			done1   = make(chan struct{})
			done2   = make(chan struct{})
			client1 = mock.NewClient[any]().RegisterDone(
				func() <-chan struct{} {
					defer wg.Done()
					return done1
				})
			client2 = mock.NewClient[any]().RegisterSendWithDeadline(
				func(deadline time.Time, cmd base.Cmd[any],
					results chan<- base.AsyncResult) (seq base.Seq, err error) {
					asserterror.Equal(deadline, wantDeadline, t)
					asserterror.Equal[any](cmd, wantCmd, t)
					asserterror.Equal(results, wantResults, t)
					return wantSeq, nil
				},
			).RegisterDone(func() <-chan struct{} {
				defer wg.Done()
				return done2
			})
			strategy = mock.NewDispatchStrategy[any]().RegisterNext(
				func() (ccln.Client[any], int64) {
					return client2, 1
				},
			).RegisterSlice(
				func() []ccln.Client[any] {
					return []ccln.Client[any]{client1, client2}
				},
			)
			mocks = []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock}
		)
		grp := ccln.NewGroup(strategy)

		seq, clientID, err := grp.SendWithDeadline(wantDeadline, wantCmd, wantResults)
		asserterror.EqualError(err, wantErr, t)
		asserterror.Equal(seq, wantSeq, t)
		asserterror.Equal(clientID, wantClientID, t)

		close(done1)
		close(done2)
		wg.Wait()

		if infomap := mok.CheckCalls(mocks); len(infomap) > 0 {
			t.Error(infomap)
		}
	})

	t.Run("Has", func(t *testing.T) {
		wg := sync.WaitGroup{}
		wg.Add(2)

		var (
			wantResult                 = true
			wantSeq      base.Seq      = 10
			wantClientID ccln.ClientID = 1

			done1   = make(chan struct{})
			done2   = make(chan struct{})
			client1 = mock.NewClient[any]().RegisterDone(
				func() <-chan struct{} {
					defer wg.Done()
					return done1
				})
			client2 = mock.NewClient[any]().RegisterHas(
				func(seq base.Seq) bool {
					asserterror.Equal(seq, wantSeq, t)
					return wantResult
				},
			).RegisterDone(func() <-chan struct{} {
				defer wg.Done()
				return done2
			})
			strategy = mock.NewDispatchStrategy[any]().RegisterNSlice(2,
				func() []ccln.Client[any] {
					return []ccln.Client[any]{client1, client2}
				},
			)
			mocks = []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock}
		)
		grp := ccln.NewGroup(strategy)

		result := grp.Has(wantSeq, wantClientID)
		asserterror.Equal(result, wantResult, t)

		close(done1)
		close(done2)
		wg.Wait()

		if infomap := mok.CheckCalls(mocks); len(infomap) > 0 {
			t.Error(infomap)
		}
	})

	t.Run("Forget", func(t *testing.T) {
		wg := sync.WaitGroup{}
		wg.Add(2)

		var (
			wantSeq      base.Seq      = 10
			wantClientID ccln.ClientID = 1

			done1   = make(chan struct{})
			done2   = make(chan struct{})
			client1 = mock.NewClient[any]().RegisterDone(
				func() <-chan struct{} {
					defer wg.Done()
					return done1
				})
			client2 = mock.NewClient[any]().RegisterForget(
				func(seq base.Seq) {
					asserterror.Equal(seq, wantSeq, t)
				},
			).RegisterDone(func() <-chan struct{} {
				defer wg.Done()
				return done2
			})
			strategy = mock.NewDispatchStrategy[any]().RegisterNSlice(2,
				func() []ccln.Client[any] {
					return []ccln.Client[any]{client1, client2}
				},
			)
			mocks = []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock}
		)
		grp := ccln.NewGroup(strategy)

		grp.Forget(wantSeq, wantClientID)

		close(done1)
		close(done2)
		wg.Wait()

		if infomap := mok.CheckCalls(mocks); len(infomap) > 0 {
			t.Error(infomap)
		}
	})

	t.Run("Error", func(t *testing.T) {
		wg := sync.WaitGroup{}
		wg.Add(2)

		var (
			wantErr1 error = errors.New("client1 error")
			wantErr2 error = errors.New("client2 error")
			wantErr        = errors.Join(wantErr1, wantErr2)

			done1   = make(chan struct{})
			done2   = make(chan struct{})
			client1 = mock.NewClient[any]().RegisterErr(
				func() (err error) {
					return wantErr1
				},
			).RegisterDone(
				func() <-chan struct{} {
					defer wg.Done()
					return done1
				},
			)
			client2 = mock.NewClient[any]().RegisterErr(
				func() (err error) {
					return wantErr2
				},
			).RegisterDone(func() <-chan struct{} {
				defer wg.Done()
				return done2
			})
			strategy = mock.NewDispatchStrategy[any]().RegisterNSlice(2,
				func() []ccln.Client[any] {
					return []ccln.Client[any]{client1, client2}
				},
			)
			mocks = []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock}
		)
		grp := ccln.NewGroup(strategy)

		err := grp.Err()
		asserterror.EqualError(err, wantErr, t)

		close(done1)
		close(done2)
		wg.Wait()

		if infomap := mok.CheckCalls(mocks); len(infomap) > 0 {
			t.Error(infomap)
		}
	})

	t.Run("Close", func(t *testing.T) {
		wg := sync.WaitGroup{}
		wg.Add(2)

		var (
			wantErr1 error = errors.New("client1 error")
			wantErr2 error = errors.New("client2 error")
			wantErr        = errors.Join(wantErr1, wantErr2)

			done1   = make(chan struct{})
			done2   = make(chan struct{})
			client1 = mock.NewClient[any]().RegisterClose(
				func() (err error) {
					return wantErr1
				},
			).RegisterDone(
				func() <-chan struct{} {
					defer wg.Done()
					return done1
				},
			)
			client2 = mock.NewClient[any]().RegisterClose(
				func() (err error) {
					return wantErr2
				},
			).RegisterDone(func() <-chan struct{} {
				defer wg.Done()
				return done2
			})
			strategy = mock.NewDispatchStrategy[any]().RegisterNSlice(2,
				func() []ccln.Client[any] {
					return []ccln.Client[any]{client1, client2}
				},
			)
			mocks = []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock}
		)
		grp := ccln.NewGroup(strategy)

		close(done1)
		close(done2)
		wg.Wait()

		err := grp.Close()
		asserterror.EqualError(err, wantErr, t)

		if infomap := mok.CheckCalls(mocks); len(infomap) > 0 {
			t.Error(infomap)
		}
	})

}
