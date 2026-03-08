package group_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	grp "github.com/cmd-stream/cmd-stream-go/group"
	grpmock "github.com/cmd-stream/cmd-stream-go/test/mock/group"
	"github.com/cmd-stream/core-go"
	cmock "github.com/cmd-stream/core-go/test/mock"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

func TestGroup(t *testing.T) {
	t.Run("NewGroup", func(t *testing.T) {
		wg := sync.WaitGroup{}
		wg.Add(2)

		var (
			done1   = make(chan struct{})
			done2   = make(chan struct{})
			client1 = grpmock.NewClient[any]().RegisterDone(
				func() <-chan struct{} {
					defer wg.Done()
					return done1
				},
			)
			client2 = grpmock.NewClient[any]().RegisterDone(
				func() <-chan struct{} {
					defer wg.Done()
					return done2
				},
			)
			strategy = grpmock.NewDispatchStrategy[any]().RegisterSlice(
				func() []grp.Client[any] {
					return []grp.Client[any]{client1, client2}
				},
			)
			mocks = []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock}
		)
		group := grp.NewClientGroup(strategy)

		close(done1)
		close(done2)
		wg.Wait()

		select {
		case <-time.NewTimer(time.Second).C:
			t.Error("timeout")
		case <-group.Done():
		}

		if infomap := mok.CheckCalls(mocks); len(infomap) > 0 {
			t.Error(infomap)
		}
	})

	t.Run("Send", func(t *testing.T) {
		wg := sync.WaitGroup{}
		wg.Add(2)

		var (
			wantErr      error        = nil
			wantSeq      core.Seq     = 10
			wantN                     = 1
			wantClientID grp.ClientID = 1
			wantCmd                   = cmock.NewCmd()
			wantResults               = make(chan core.AsyncResult)

			done1   = make(chan struct{})
			done2   = make(chan struct{})
			client1 = grpmock.NewClient[any]().RegisterDone(
				func() <-chan struct{} {
					defer wg.Done()
					return done1
				})
			client2 = grpmock.NewClient[any]().RegisterSend(
				func(cmd core.Cmd[any], results chan<- core.AsyncResult) (seq core.Seq, n int, err error) {
					asserterror.Equal[any](t, cmd, wantCmd)
					asserterror.Equal(t, results, wantResults)
					return wantSeq, wantN, nil
				},
			).RegisterDone(func() <-chan struct{} {
				defer wg.Done()
				return done2
			})
			strategy = grpmock.NewDispatchStrategy[any]().RegisterNext(
				func() (grp.Client[any], int64) {
					return client2, 1
				},
			).RegisterSlice(
				func() []grp.Client[any] {
					return []grp.Client[any]{client1, client2}
				},
			)
			mocks = []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock}
		)
		group := grp.NewClientGroup(strategy)

		seq, clientID, n, err := group.Send(wantCmd, wantResults)
		asserterror.EqualError(t, err, wantErr)
		asserterror.Equal(t, seq, wantSeq)
		asserterror.Equal(t, clientID, wantClientID)
		asserterror.Equal(t, n, wantN)

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
			wantErr      error        = nil
			wantSeq      core.Seq     = 10
			wantN                     = 2
			wantClientID grp.ClientID = 1
			wantDeadline              = time.Now().Add(time.Second)
			wantCmd                   = cmock.NewCmd()
			wantResults               = make(chan core.AsyncResult)

			done1   = make(chan struct{})
			done2   = make(chan struct{})
			client1 = grpmock.NewClient[any]().RegisterDone(
				func() <-chan struct{} {
					defer wg.Done()
					return done1
				})
			client2 = grpmock.NewClient[any]().RegisterSendWithDeadline(
				func(deadline time.Time, cmd core.Cmd[any],
					results chan<- core.AsyncResult,
				) (seq core.Seq, n int, err error) {
					asserterror.Equal(t, deadline, wantDeadline)
					asserterror.Equal[any](t, cmd, wantCmd)
					asserterror.Equal(t, results, wantResults)
					return wantSeq, wantN, nil
				},
			).RegisterDone(func() <-chan struct{} {
				defer wg.Done()
				return done2
			})
			strategy = grpmock.NewDispatchStrategy[any]().RegisterNext(
				func() (grp.Client[any], int64) {
					return client2, 1
				},
			).RegisterSlice(
				func() []grp.Client[any] {
					return []grp.Client[any]{client1, client2}
				},
			)
			mocks = []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock}
		)
		group := grp.NewClientGroup(strategy)

		seq, clientID, n, err := group.SendWithDeadline(wantCmd, wantResults,
			wantDeadline)
		asserterror.EqualError(t, err, wantErr)
		asserterror.Equal(t, seq, wantSeq)
		asserterror.Equal(t, clientID, wantClientID)
		asserterror.Equal(t, n, wantN)

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
			wantResult                = true
			wantSeq      core.Seq     = 10
			wantClientID grp.ClientID = 1

			done1   = make(chan struct{})
			done2   = make(chan struct{})
			client1 = grpmock.NewClient[any]().RegisterDone(
				func() <-chan struct{} {
					defer wg.Done()
					return done1
				})
			client2 = grpmock.NewClient[any]().RegisterHas(
				func(seq core.Seq) bool {
					asserterror.Equal(t, seq, wantSeq)
					return wantResult
				},
			).RegisterDone(func() <-chan struct{} {
				defer wg.Done()
				return done2
			})
			strategy = grpmock.NewDispatchStrategy[any]().RegisterNSlice(2,
				func() []grp.Client[any] {
					return []grp.Client[any]{client1, client2}
				},
			)
			mocks = []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock}
		)
		group := grp.NewClientGroup(strategy)

		result := group.Has(wantSeq, wantClientID)
		asserterror.Equal(t, result, wantResult)

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
			wantSeq      core.Seq     = 10
			wantClientID grp.ClientID = 1

			done1   = make(chan struct{})
			done2   = make(chan struct{})
			client1 = grpmock.NewClient[any]().RegisterDone(
				func() <-chan struct{} {
					defer wg.Done()
					return done1
				})
			client2 = grpmock.NewClient[any]().RegisterForget(
				func(seq core.Seq) {
					asserterror.Equal(t, seq, wantSeq)
				},
			).RegisterDone(func() <-chan struct{} {
				defer wg.Done()
				return done2
			})
			strategy = grpmock.NewDispatchStrategy[any]().RegisterNSlice(2,
				func() []grp.Client[any] {
					return []grp.Client[any]{client1, client2}
				},
			)
			mocks = []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock}
		)
		group := grp.NewClientGroup(strategy)

		group.Forget(wantSeq, wantClientID)

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
			wantErr1 = errors.New("client1 error")
			wantErr2 = errors.New("client2 error")
			wantErr  = errors.Join(wantErr1, wantErr2)

			done1   = make(chan struct{})
			done2   = make(chan struct{})
			client1 = grpmock.NewClient[any]().RegisterErr(
				func() (err error) {
					return wantErr1
				},
			).RegisterDone(
				func() <-chan struct{} {
					defer wg.Done()
					return done1
				},
			)
			client2 = grpmock.NewClient[any]().RegisterErr(
				func() (err error) {
					return wantErr2
				},
			).RegisterDone(func() <-chan struct{} {
				defer wg.Done()
				return done2
			})
			strategy = grpmock.NewDispatchStrategy[any]().RegisterNSlice(2,
				func() []grp.Client[any] {
					return []grp.Client[any]{client1, client2}
				},
			)
			mocks = []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock}
		)
		group := grp.NewClientGroup(strategy)

		err := group.Err()
		asserterror.EqualError(t, err, wantErr)

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
			wantErr1 = errors.New("client1 error")
			wantErr2 = errors.New("client2 error")
			wantErr  = grp.NewGroupError(errors.Join(wantErr1, wantErr2))

			done1   = make(chan struct{})
			done2   = make(chan struct{})
			client1 = grpmock.NewClient[any]().RegisterClose(
				func() (err error) {
					return wantErr1
				},
			).RegisterDone(
				func() <-chan struct{} {
					defer wg.Done()
					return done1
				},
			)
			client2 = grpmock.NewClient[any]().RegisterClose(
				func() (err error) {
					return wantErr2
				},
			).RegisterDone(func() <-chan struct{} {
				defer wg.Done()
				return done2
			})
			strategy = grpmock.NewDispatchStrategy[any]().RegisterNSlice(2,
				func() []grp.Client[any] {
					return []grp.Client[any]{client1, client2}
				},
			)
			mocks = []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock}
		)
		group := grp.NewClientGroup(strategy)

		close(done1)
		close(done2)
		wg.Wait()

		err := group.Close()
		asserterror.EqualError(t, err, wantErr)

		if infomap := mok.CheckCalls(mocks); len(infomap) > 0 {
			t.Error(infomap)
		}
	})

	// t.Run("MakeGroup", func(t *testing.T) {

	// 	t.Run("Should work", func(t *testing.T) {
	// 		var (
	// 			wantErr     error = nil
	// 			connFactory       = cmock.NewConnFactory[any]().RegisterNew(
	// 				func() (net.Conn, error) {
	// 					return cmock.NewConn().RegisterSetReadDeadline(
	// 						func(deadline time.Time) (err error) { return },
	// 					).RegisterRead(
	// 						func(b []byte) (n int, err error) {
	// 							w := &bytes.Buffer{}
	// 							n, err = ord.ByteSlice.Marshal(csrv.ServerInfo, w)
	// 							copy(b, w.Bytes())
	// 							return
	// 						},
	// 					).RegisterSetReadDeadline(
	// 						func(deadline time.Time) (err error) { return },
	// 					), nil
	// 				},
	// 			).RegisterNew(
	// 				func() (net.Conn, error) {
	// 					return cmock.NewConn().RegisterSetReadDeadline(
	// 						func(deadline time.Time) (err error) { return },
	// 					).RegisterRead(
	// 						func(b []byte) (n int, err error) {
	// 							w := &bytes.Buffer{}
	// 							n, err = ord.ByteSlice.Marshal(csrv.ServerInfo, w)
	// 							copy(b, w.Bytes())
	// 							return
	// 						},
	// 					).RegisterSetReadDeadline(
	// 						func(deadline time.Time) (err error) { return },
	// 					), nil
	// 				},
	// 			)
	// 		)
	// 		grp, err := grp.MakeGroup(2, cdmock.NewCodec[core.Cmd[any], core.Result](),
	// 			connFactory)
	// 		asserterror.EqualError(t, err, wantErr)
	// 		asserterror.Equal(t, grp.Size(), 2)
	// 	})

	// })
}
