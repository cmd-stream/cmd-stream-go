package cgrp_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/cmd-stream/base-go"
	bmock "github.com/cmd-stream/base-go/testdata/mock"
	cgrp "github.com/cmd-stream/cmd-stream-go/group"
	"github.com/cmd-stream/cmd-stream-go/group/testdata/mock"
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
				func() []cgrp.Client[any] {
					return []cgrp.Client[any]{client1, client2}
				},
			)
			mocks = []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock}
		)
		group := cgrp.NewClientGroup(strategy)

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
			wantErr      error         = nil
			wantSeq      base.Seq      = 10
			wantN        int           = 1
			wantClientID cgrp.ClientID = 1
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
				func(cmd base.Cmd[any], results chan<- base.AsyncResult) (seq base.Seq, n int, err error) {
					asserterror.Equal[any](cmd, wantCmd, t)
					asserterror.Equal(results, wantResults, t)
					return wantSeq, wantN, nil
				},
			).RegisterDone(func() <-chan struct{} {
				defer wg.Done()
				return done2
			})
			strategy = mock.NewDispatchStrategy[any]().RegisterNext(
				func() (cgrp.Client[any], int64) {
					return client2, 1
				},
			).RegisterSlice(
				func() []cgrp.Client[any] {
					return []cgrp.Client[any]{client1, client2}
				},
			)
			mocks = []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock}
		)
		group := cgrp.NewClientGroup(strategy)

		seq, clientID, n, err := group.Send(wantCmd, wantResults)
		asserterror.EqualError(err, wantErr, t)
		asserterror.Equal(seq, wantSeq, t)
		asserterror.Equal(clientID, wantClientID, t)
		asserterror.Equal(n, wantN, t)

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
			wantN        int           = 2
			wantClientID cgrp.ClientID = 1
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
					results chan<- base.AsyncResult) (seq base.Seq, n int, err error) {
					asserterror.Equal(deadline, wantDeadline, t)
					asserterror.Equal[any](cmd, wantCmd, t)
					asserterror.Equal(results, wantResults, t)
					return wantSeq, wantN, nil
				},
			).RegisterDone(func() <-chan struct{} {
				defer wg.Done()
				return done2
			})
			strategy = mock.NewDispatchStrategy[any]().RegisterNext(
				func() (cgrp.Client[any], int64) {
					return client2, 1
				},
			).RegisterSlice(
				func() []cgrp.Client[any] {
					return []cgrp.Client[any]{client1, client2}
				},
			)
			mocks = []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock}
		)
		group := cgrp.NewClientGroup(strategy)

		seq, clientID, n, err := group.SendWithDeadline(wantCmd, wantResults,
			wantDeadline)
		asserterror.EqualError(err, wantErr, t)
		asserterror.Equal(seq, wantSeq, t)
		asserterror.Equal(clientID, wantClientID, t)
		asserterror.Equal(n, wantN, t)

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
			wantClientID cgrp.ClientID = 1

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
				func() []cgrp.Client[any] {
					return []cgrp.Client[any]{client1, client2}
				},
			)
			mocks = []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock}
		)
		group := cgrp.NewClientGroup(strategy)

		result := group.Has(wantSeq, wantClientID)
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
			wantClientID cgrp.ClientID = 1

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
				func() []cgrp.Client[any] {
					return []cgrp.Client[any]{client1, client2}
				},
			)
			mocks = []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock}
		)
		group := cgrp.NewClientGroup(strategy)

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
				func() []cgrp.Client[any] {
					return []cgrp.Client[any]{client1, client2}
				},
			)
			mocks = []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock}
		)
		group := cgrp.NewClientGroup(strategy)

		err := group.Err()
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
				func() []cgrp.Client[any] {
					return []cgrp.Client[any]{client1, client2}
				},
			)
			mocks = []*mok.Mock{client1.Mock, client2.Mock, strategy.Mock}
		)
		group := cgrp.NewClientGroup(strategy)

		close(done1)
		close(done2)
		wg.Wait()

		err := group.Close()
		asserterror.EqualError(err, wantErr, t)

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
	// 					return bmock.NewConn().RegisterSetReadDeadline(
	// 						func(deadline time.Time) (err error) { return },
	// 					).RegisterRead(
	// 						func(b []byte) (n int, err error) {
	// 							w := &bytes.Buffer{}
	// 							n, err = ord.ByteSlice.Marshal(cser.ServerInfo, w)
	// 							copy(b, w.Bytes())
	// 							return
	// 						},
	// 					).RegisterSetReadDeadline(
	// 						func(deadline time.Time) (err error) { return },
	// 					), nil
	// 				},
	// 			).RegisterNew(
	// 				func() (net.Conn, error) {
	// 					return bmock.NewConn().RegisterSetReadDeadline(
	// 						func(deadline time.Time) (err error) { return },
	// 					).RegisterRead(
	// 						func(b []byte) (n int, err error) {
	// 							w := &bytes.Buffer{}
	// 							n, err = ord.ByteSlice.Marshal(cser.ServerInfo, w)
	// 							copy(b, w.Bytes())
	// 							return
	// 						},
	// 					).RegisterSetReadDeadline(
	// 						func(deadline time.Time) (err error) { return },
	// 					), nil
	// 				},
	// 			)
	// 		)
	// 		cgrp, err := cgrp.MakeGroup(2, cdmock.NewCodec[base.Cmd[any], base.Result](),
	// 			connFactory)
	// 		asserterror.EqualError(err, wantErr, t)
	// 		asserterror.Equal(cgrp.Size(), 2, t)
	// 	})

	// })

}
