// Package sender provides a high-level abstraction for sending Commands and
// processing Results across a group of clients, with built-in support for hooks,
// deadlines, and multi-result handling.
package sender

import (
	"context"
	"errors"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	grp "github.com/cmd-stream/cmd-stream-go/group"
	hks "github.com/cmd-stream/cmd-stream-go/sender/hooks"
)

// Sender provides a high-level abstraction for sending Commands and processing
// Results. It uses a Group to communicate with the server and a
// hooks factory to customize behavior.
type Sender[T any] struct {
	group   Group[T]
	options Options[T]
}

// New creates a new Sender with the given client group and options.
func New[T any](group Group[T], opts ...SetOption[T]) Sender[T] {
	o := DefaultOptions[T]()
	Apply(&o, opts...)
	return Sender[T]{
		group:   group,
		options: o,
	}
}

// Send sends a Command and waits for the Result.
func (s Sender[T]) Send(ctx context.Context, cmd core.Cmd[T]) (
	result core.Result, err error,
) {
	var (
		results = make(chan core.AsyncResult, 1)
		hooks   = s.options.HooksFactory.New()
	)
	ctx, err = hooks.BeforeSend(ctx, cmd)
	if err != nil {
		return
	}
	seq, clientID, n, err := s.group.Send(cmd, results)
	sentCmd := hks.SentCmd[T]{
		Seq:  seq,
		Size: n,
		Cmd:  cmd,
	}
	if err != nil {
		hooks.OnError(ctx, sentCmd, err)
		return
	}
	return s.receive(ctx, sentCmd, results, clientID, hooks)
}

// SendWithDeadline sends a Command with the specified deadline and waits for
// the Result.
func (s Sender[T]) SendWithDeadline(ctx context.Context, deadline time.Time,
	cmd core.Cmd[T],
) (result core.Result, err error) {
	var (
		results = make(chan core.AsyncResult, 1)
		hooks   = s.options.HooksFactory.New()
	)
	ctx, err = hooks.BeforeSend(ctx, cmd)
	if err != nil {
		return
	}
	seq, clientID, n, err := s.group.SendWithDeadline(deadline, cmd, results)
	sentCmd := hks.SentCmd[T]{
		Seq:  seq,
		Size: n,
		Cmd:  cmd,
	}
	if err != nil {
		hooks.OnError(ctx, sentCmd, err)
		return
	}
	return s.receive(ctx, sentCmd, results, clientID, hooks)
}

// SendMulti sends a Command and waits for multiple Results. Each Result is
// passed to the provided handler.
func (s Sender[T]) SendMulti(ctx context.Context, cmd core.Cmd[T],
	resultsCount int, handler ResultHandler,
) (err error) {
	var (
		results = make(chan core.AsyncResult, resultsCount)
		hooks   = s.options.HooksFactory.New()
	)
	ctx, err = hooks.BeforeSend(ctx, cmd)
	if err != nil {
		return
	}
	seq, clientID, n, err := s.group.Send(cmd, results)
	sentCmd := hks.SentCmd[T]{
		Seq:  seq,
		Size: n,
		Cmd:  cmd,
	}
	if err != nil {
		hooks.OnError(ctx, sentCmd, err)
		return
	}
	s.receiveMulti(ctx, sentCmd, results, clientID, hooks, handler)
	return
}

// SendMultiWithDeadline sends a Command with the specified deadline and waits
// for multiple Results. Each Result is passed to the provided handler.
func (s Sender[T]) SendMultiWithDeadline(ctx context.Context, deadline time.Time,
	cmd core.Cmd[T],
	resultsCount int,
	handler ResultHandler,
) (err error) {
	var (
		results = make(chan core.AsyncResult, resultsCount)
		hooks   = s.options.HooksFactory.New()
	)
	ctx, err = hooks.BeforeSend(ctx, cmd)
	if err != nil {
		return
	}
	seq, clientID, n, err := s.group.SendWithDeadline(deadline, cmd, results)
	sentCmd := hks.SentCmd[T]{
		Seq:  seq,
		Size: n,
		Cmd:  cmd,
	}
	if err != nil {
		hooks.OnError(ctx, sentCmd, err)
		return
	}
	s.receiveMulti(ctx, sentCmd, results, clientID, hooks, handler)
	return
}

// CloseAndWait closes the sender and waits for all processing to complete or
// until the timeout is exceeded.
func (s Sender[T]) CloseAndWait(timeout time.Duration) (err error) {
	err = s.Close()
	if err != nil {
		return
	}
	select {
	case <-time.NewTimer(timeout).C:
		return errors.New("timeout exceeded")
	case <-s.Done():
		return
	}
}

// Close closes the underlying client group.
func (s Sender[T]) Close() error {
	return s.group.Close()
}

// Done returns a channel that is closed when the sender is closed and all
// processing is complete.
func (s Sender[T]) Done() <-chan struct{} {
	return s.group.Done()
}

func (s Sender[T]) receive(ctx context.Context, sentCmd hks.SentCmd[T],
	results <-chan core.AsyncResult,
	clientID grp.ClientID,
	hooks hks.Hooks[T],
) (result core.Result, err error) {
	select {
	case <-ctx.Done():
		err = ErrTimeout
		hooks.OnTimeout(ctx, sentCmd, err)
		s.group.Forget(sentCmd.Seq, clientID)
	case asyncResult := <-results:
		recvResult := hks.ReceivedResult{
			Seq:    core.Seq(1),
			Size:   asyncResult.BytesRead,
			Result: asyncResult.Result,
		}
		hooks.OnResult(ctx, sentCmd, recvResult, asyncResult.Error)
		result = asyncResult.Result
		err = asyncResult.Error
	}
	return
}

func (s Sender[T]) receiveMulti(ctx context.Context, sentCmd hks.SentCmd[T],
	results <-chan core.AsyncResult,
	clientID grp.ClientID,
	hooks hks.Hooks[T],
	handler ResultHandler,
) {
	var (
		result    core.Result
		handleErr error
		err       error
		i         = 1
	)
	for {
		select {
		case <-ctx.Done():
			err = ErrTimeout
			hooks.OnTimeout(ctx, sentCmd, err)
			s.group.Forget(sentCmd.Seq, clientID)
		case asyncResult := <-results:
			recvResult := hks.ReceivedResult{
				Seq:    core.Seq(i),
				Size:   asyncResult.BytesRead,
				Result: asyncResult.Result,
			}
			hooks.OnResult(ctx, sentCmd, recvResult, asyncResult.Error)
			result = asyncResult.Result
			err = asyncResult.Error
		}
		handleErr = handler.Handle(result, err)
		if handleErr != nil {
			err = handleErr
		}
		if err != nil || result.LastOne() {
			return
		}
		i++
	}
}
