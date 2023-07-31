package server

import (
	"context"
	"errors"
	"testing"
	"time"

	base "github.com/cmd-stream/base-go"
	base_mock "github.com/cmd-stream/base-go/testdata/mock"
)

func TestInvoker(t *testing.T) {

	t.Run("Invoke should execute cmd", func(t *testing.T) {
		var (
			wantCtx                 = context.TODO()
			wantAt                  = time.Time{}
			wantSeq                 = base.Seq(10)
			wantReceiver            = struct{}{}
			wantProxy    base.Proxy = nil
			cmd                     = base_mock.NewCmd().RegisterExec(
				func(ctx context.Context, at time.Time, seq base.Seq, receiver any, proxy base.Proxy) (err error) {
					if ctx != wantCtx {
						t.Errorf("unexpected ctx, want '%v' actual '%v'", wantCtx, ctx)
					}
					if at != wantAt {
						t.Errorf("unexpected at, want '%v' actual '%v'", wantAt, at)
					}
					if seq != wantSeq {
						t.Errorf("unexpected seq, want '%v' actual '%v'", wantSeq, seq)
					}
					if receiver != wantReceiver {
						t.Errorf("unexpected receiver, want '%v' actual '%v'", wantReceiver,
							receiver)
					}
					return nil
				},
			)
			invoker = DefInvoker[any]{wantReceiver}
			err     = invoker.Invoke(wantCtx, wantAt, wantSeq, cmd, wantProxy)
		)
		if err != nil {
			t.Errorf("unexpected er, want '%v' actual '%v'", nil, err)
		}
	})

	t.Run("If cmd execution failed with an error, Invoke should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("cmd execution error")
				cmd     = base_mock.NewCmd().RegisterExec(
					func(ctx context.Context, at time.Time, seq base.Seq, receiver any,
						proxy base.Proxy) (err error) {
						return wantErr
					},
				)
				invoker = DefInvoker[any]{nil}
				err     = invoker.Invoke(context.TODO(), time.Time{}, base.Seq(1), cmd, nil)
			)
			if err != wantErr {
				t.Errorf("unexpected er, want '%v' actual '%v'", wantErr, err)
			}
		})

}
