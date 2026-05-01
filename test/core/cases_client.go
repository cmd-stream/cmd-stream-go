package core

import (
	"errors"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/cmd-stream/cmd-stream-go/core/cln"
	"github.com/cmd-stream/cmd-stream-go/test"
	mock "github.com/cmd-stream/cmd-stream-go/test/mock/core"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

func ReconnectTestCase() ClientTestCase[any] {
	name := "If the client has lost a connection it should try to reconnect"

	var (
		reconnectDone = make(chan struct{})
		delegate      = mock.NewReconnectDelegate()
	)

	delegate.RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			return 0, nil, 0, net.ErrClosed
		},
	).RegisterReconnect(
		func() error {
			close(reconnectDone)
			return nil
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			<-reconnectDone
			return 0, nil, 0, errors.New("receive error")
		},
	).RegisterClose(
		func() (err error) { return nil },
	)
	return ClientTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{},
		},
		Mocks: []*mok.Mock{delegate.Mock},
	}
}

// -----------------------------------------------------------------------------

func ReconnectOnEOFTestCase() ClientTestCase[any] {
	name := "If the client received EOF it should try to reconnect"

	var (
		reconnectDone = make(chan struct{})
		delegate      = mock.NewReconnectDelegate()
	)

	delegate.RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			return 0, nil, 0, io.EOF
		},
	).RegisterReconnect(
		func() error {
			close(reconnectDone)
			return nil
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			<-reconnectDone
			return 0, nil, 0, errors.New("receive error")
		},
	).RegisterClose(
		func() (err error) { return nil },
	)
	return ClientTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{},
		},
		Mocks: []*mok.Mock{delegate.Mock},
	}
}

// -----------------------------------------------------------------------------

func NoReconnectOnCloseTestCase() ClientTestCase[any] {
	name := "If the client is closed it should not reconnect"

	var (
		receiveDone = make(chan struct{})
		delegate    = mock.NewReconnectDelegate()
	)

	delegate.RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			<-receiveDone
			return 0, nil, 0, cln.ErrClosed
		},
	).RegisterClose(
		func() (err error) {
			close(receiveDone)
			return nil
		},
	)
	return ClientTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{},
		},
		Action: func(t *testing.T, client *cln.Client[any], results chan core.AsyncResult) {
			err := client.Close()
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{delegate.Mock},
	}
}

// -----------------------------------------------------------------------------

func ReconnectFailTestCase() ClientTestCase[any] {
	name := "If reconnection fails with an error, it should become the client error"

	var (
		wantErr  = errors.New("reconnection error")
		delegate = mock.NewReconnectDelegate()
	)
	delegate.RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			return 0, nil, 0, net.ErrClosed
		},
	).RegisterReconnect(
		func() error { return wantErr },
	).RegisterClose(
		func() (err error) { return nil },
	)
	return ClientTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{},
		},
		Action: func(t *testing.T, client *cln.Client[any], results chan core.AsyncResult) {
			<-client.Done()
			asserterror.Equal(t, client.Error(), wantErr)
		},
		Mocks: []*mok.Mock{delegate.Mock},
	}
}

// -----------------------------------------------------------------------------

func KeepaliveTestCase() ClientTestCase[any] {
	name := "Upon creation, the client should call KeepaliveDelegate.Keepalive()"

	var (
		keepaliveDone = make(chan struct{})
		delegate      = mock.NewKeepaliveDelegate()
	)

	delegate.RegisterKeepalive(
		func(muSn *sync.Mutex) {
			close(keepaliveDone)
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			<-keepaliveDone
			return 0, nil, 0, errors.New("receive error")
		},
	).RegisterClose(
		func() (err error) { return nil },
	)
	return ClientTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{},
		},
		Mocks: []*mok.Mock{delegate.Mock},
	}
}

// Send ------------------------------------------------------------------------

func SendSuccessTestCase() ClientTestCase[any] {
	name := "Should successfully send cmd and receive result"

	var (
		wantSeq    core.Seq = 1
		wantN               = 1
		wantResult          = mock.NewResult()
		cmd                 = mock.NewCmd[any]()
		sendDone            = make(chan struct{})
		delegate            = mock.NewClientDelegate()
	)

	delegate.RegisterSend(
		func(seq core.Seq, cmd core.Cmd[any]) (n int, err error) {
			return wantN, nil
		},
	).RegisterFlush(
		func() (err error) {
			close(sendDone)
			return nil
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			<-sendDone
			return wantSeq, wantResult, 10, nil
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			return 0, nil, 0, errors.New("receive error")
		},
	).RegisterClose(
		func() (err error) {
			return nil
		},
	)
	wantResult.RegisterLastOne(func() bool { return true })

	return ClientTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{},
		},
		Action: func(t *testing.T, client *cln.Client[any], results chan core.AsyncResult) {
			AssertSend(t, client, cmd, results, wantSeq, wantN, nil)
			AssertResults(t, results, core.AsyncResult{
				Seq:       wantSeq,
				Result:    wantResult,
				BytesRead: 10,
			})
		},
		Mocks: []*mok.Mock{delegate.Mock, cmd.Mock, wantResult.Mock},
	}
}

// -----------------------------------------------------------------------------

func HasTestCase() ClientTestCase[any] {
	name := "Client.Has should return true when client has cmd"

	var (
		wantSeq   core.Seq = 1
		wantN              = 1
		cmd                = mock.NewCmd[any]()
		checkDone          = make(chan struct{})
		delegate           = mock.NewClientDelegate()
	)
	delegate.RegisterSend(
		func(seq core.Seq, cmd core.Cmd[any]) (n int, err error) {
			return wantN, nil
		},
	).RegisterFlush(
		func() (err error) {
			return nil
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			<-checkDone
			return 0, nil, 0, errors.New("receive error")
		},
	).RegisterClose(
		func() (err error) {
			return nil
		},
	)
	return ClientTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{},
		},
		Action: func(t *testing.T, client *cln.Client[any], results chan core.AsyncResult) {
			seq := AssertSend(t, client, cmd, results, wantSeq, wantN, nil)
			AssertHas(t, client, seq, true)
			close(checkDone)
		},
		Mocks: []*mok.Mock{delegate.Mock, cmd.Mock},
	}
}

// -----------------------------------------------------------------------------
func ForgetTestCase() ClientTestCase[any] {
	name := "Client.Forget should remove cmd from waiting map"

	var (
		wantSeq   core.Seq = 1
		wantN              = 1
		cmd                = mock.NewCmd[any]()
		checkDone          = make(chan struct{})
		delegate           = mock.NewClientDelegate()
	)
	delegate.RegisterSend(
		func(seq core.Seq, cmd core.Cmd[any]) (n int, err error) {
			return wantN, nil
		},
	).RegisterFlush(
		func() (err error) {
			return nil
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			<-checkDone
			return 0, nil, 0, errors.New("receive error")
		},
	).RegisterClose(
		func() (err error) {
			return nil
		},
	)
	return ClientTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{},
		},
		Action: func(t *testing.T, client *cln.Client[any], results chan core.AsyncResult) {
			seq := AssertSend(t, client, cmd, results, wantSeq, wantN, nil)
			AssertHas(t, client, seq, true)
			client.Forget(seq)
			AssertHas(t, client, seq, false)
			close(checkDone)
		},
		Mocks: []*mok.Mock{delegate.Mock, cmd.Mock},
	}
}

// -----------------------------------------------------------------------------

func ForgetOnFailTestCase() ClientTestCase[any] {
	name := "The cmd should be forgoten if send fails"

	var (
		wantSeq  core.Seq = 1
		wantErr           = errors.New("send error")
		cmd               = mock.NewCmd[any]()
		sendDone          = make(chan struct{})
		delegate          = mock.NewClientDelegate()
	)
	delegate.RegisterSend(
		func(seq core.Seq, cmd core.Cmd[any]) (n int, err error) {
			close(sendDone)
			return 0, wantErr
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			<-sendDone
			return 0, nil, 0, errors.New("receive error")
		},
	).RegisterClose(
		func() (err error) {
			return nil
		},
	)
	return ClientTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{},
		},
		Action: func(t *testing.T, client *cln.Client[any], results chan core.AsyncResult) {
			seq := AssertSend(t, client, cmd, results, wantSeq, 0, cln.NewClientError(wantErr))
			AssertHas(t, client, seq, false)
		},
		Mocks: []*mok.Mock{delegate.Mock, cmd.Mock},
	}
}

// SendWithDeadline ------------------------------------------------------------

func SendWDTestCase() ClientTestCase[any] {
	name := "Should successfully send cmd by SendWithDeadline"

	var (
		wantSeq      core.Seq = 1
		wantDeadline          = time.Now()
		wantCmd               = mock.NewCmd[any]()
		receiveDone           = make(chan struct{})
		delegate              = mock.NewClientDelegate()
	)

	delegate.RegisterSetSendDeadline(
		func(deadline time.Time) (err error) {
			return nil
		},
	).RegisterSend(
		func(seq core.Seq, cmd core.Cmd[any]) (n int, err error) {
			return 1, nil
		},
	).RegisterFlush(
		func() (err error) {
			return nil
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			<-receiveDone
			return 0, nil, 0, errors.New("receive error")
		},
	).RegisterClose(
		func() (err error) {
			return nil
		},
	)

	return ClientTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{},
		},
		Action: func(t *testing.T, client *cln.Client[any], results chan core.AsyncResult) {
			seq := AssertSendWithDeadline(t, client, wantDeadline, wantCmd, results, wantSeq, 1, nil)
			AssertHas(t, client, seq, true)
			close(receiveDone)
		},
		Mocks: []*mok.Mock{wantCmd.Mock, delegate.Mock},
	}
}

// -----------------------------------------------------------------------------

func SendWDFailSetDeadlineTestCase() ClientTestCase[any] {
	name := "If Delegate.SetSendDeadline fails with an error, SendWithDeadline should return it"

	var (
		wantSeq     core.Seq = 1
		delegateErr          = errors.New("Delegate.SetSendDeadline error")
		wantErr              = cln.NewClientError(delegateErr)
		sendDone             = make(chan struct{})
		delegate             = mock.NewClientDelegate()
	)

	delegate.RegisterSetSendDeadline(
		func(deadline time.Time) (err error) {
			close(sendDone)
			return delegateErr
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			<-sendDone
			return 0, nil, 0, errors.New("receive error")
		},
	).RegisterClose(
		func() (err error) {
			return nil
		},
	)
	return ClientTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{},
		},
		Action: func(t *testing.T, client *cln.Client[any], results chan core.AsyncResult) {
			AssertSendWithDeadline(t, client, time.Unix(1, 0), nil, results, wantSeq, 0, wantErr)
		},
		Mocks: []*mok.Mock{delegate.Mock},
	}
}

// -----------------------------------------------------------------------------

func SendWDFailTestCase() ClientTestCase[any] {
	name := "If Delegate.Send fails with an error, SendWithDeadline should return it"

	var (
		wantSeq     core.Seq = 1
		delegateErr          = errors.New("Delegate.Send error")
		wantErr              = cln.NewClientError(delegateErr)
		sendDone             = make(chan struct{})
		delegate             = mock.NewClientDelegate()
	)
	delegate.RegisterSetSendDeadline(
		func(deadline time.Time) (err error) {
			return nil
		},
	).RegisterSend(
		func(seq core.Seq, cmd core.Cmd[any]) (n int, err error) {
			close(sendDone)
			return 0, delegateErr
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			<-sendDone
			return 0, nil, 0, errors.New("receive error")
		},
	).RegisterClose(
		func() (err error) {
			return nil
		},
	)
	return ClientTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{},
		},
		Action: func(t *testing.T, client *cln.Client[any], results chan core.AsyncResult) {
			AssertSendWithDeadline(t, client, time.Unix(1, 0), nil, results, wantSeq, 0, wantErr)
		},
		Mocks: []*mok.Mock{delegate.Mock},
	}
}

// -----------------------------------------------------------------------------

func ClosedOnReceiveErrorTestCase() ClientTestCase[any] {
	name := "If Receive fails with an error, further Send calls should return ErrClosed"

	var (
		cmd         = mock.NewCmd[any]()
		receiveDone = make(chan struct{})
		delegate    = mock.NewClientDelegate()
	)
	delegate.RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			close(receiveDone)
			return 0, nil, 0, errors.New("receive error")
		},
	).RegisterClose(
		func() (err error) {
			return nil
		},
	)
	return ClientTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{},
		},
		Action: func(t *testing.T, client *cln.Client[any], results chan core.AsyncResult) {
			<-receiveDone
			AssertDone(t, client)
			AssertSend(t, client, cmd, results, 0, 0, cln.ErrClosed)
		},
		Mocks: []*mok.Mock{delegate.Mock},
	}
}

// -----------------------------------------------------------------------------

func ForgetOnSendWDFailSetDeadlineTestCase() ClientTestCase[any] {
	name := "Should forget the cmd if SendWithDeadline failed to Delegate.SetSendDeadline"

	var (
		wantSeq     core.Seq = 1
		delegateErr          = errors.New("Delegate.SetSendDeadline error")
		wantErr              = cln.NewClientError(delegateErr)
		sendDone             = make(chan struct{})
		delegate             = mock.NewClientDelegate()
	)
	delegate.RegisterSetSendDeadline(
		func(deadline time.Time) (err error) {
			close(sendDone)
			return delegateErr
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			<-sendDone
			return 0, nil, 0, errors.New("receive error")
		},
	).RegisterClose(
		func() (err error) {
			return nil
		},
	)
	return ClientTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{},
		},
		Action: func(t *testing.T, client *cln.Client[any], results chan core.AsyncResult) {
			seq := AssertSendWithDeadline(t, client, time.Unix(1, 0), nil, results, wantSeq, 0, wantErr)
			AssertHas(t, client, seq, false)
		},
		Mocks: []*mok.Mock{delegate.Mock},
	}
}

// -----------------------------------------------------------------------------

func ForgetOnSendWDFailSendTestCase() ClientTestCase[any] {
	name := "Should forget the cmd if SendWithDeadline failed to Delegate.Send"

	var (
		wantSeq     core.Seq = 1
		delegateErr          = errors.New("Delegate.Send error")
		wantErr              = cln.NewClientError(delegateErr)
		sendDone             = make(chan struct{})
		delegate             = mock.NewClientDelegate()
	)
	delegate.RegisterSetSendDeadline(
		func(deadline time.Time) (err error) {
			return nil
		},
	).RegisterSend(
		func(seq core.Seq, cmd core.Cmd[any]) (n int, err error) {
			close(sendDone)
			return 0, delegateErr
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			<-sendDone
			return 0, nil, 0, errors.New("receive error")
		},
	).RegisterClose(
		func() (err error) {
			return nil
		},
	)
	return ClientTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{},
		},
		Action: func(t *testing.T, client *cln.Client[any], results chan core.AsyncResult) {
			seq := AssertSendWithDeadline(t, client, time.Unix(1, 0), nil, results, wantSeq, 0, wantErr)
			AssertHas(t, client, seq, false)
		},
		Mocks: []*mok.Mock{delegate.Mock},
	}
}

// Send Multi ------------------------------------------------------------------

func IncrementSeqOnSendWDFailTestCase() MultiSendTestCase[any] {
	name := "Should increment seq even after SendWithDeadline fail"

	var (
		wantSeq1    core.Seq = 1
		wantSeq2    core.Seq = 2
		wantErr              = errors.New("SetSendDeadline error")
		cmd1                 = mock.NewCmd[any]()
		cmd2                 = mock.NewCmd[any]()
		receiveDone          = make(chan struct{})
		delegate             = mock.NewClientDelegate()
	)
	delegate.RegisterSetSendDeadline(
		func(deadline time.Time) (err error) {
			return wantErr
		},
	).RegisterSend(
		func(seq core.Seq, cmd core.Cmd[any]) (n int, err error) {
			return 1, nil
		},
	).RegisterFlush(
		func() (err error) {
			return nil
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			<-receiveDone
			return 0, nil, 0, errors.New("receive error")
		},
	).RegisterClose(
		func() (err error) {
			return nil
		},
	)
	return MultiSendTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{},
		},
		Params: MultiSendParams[any]{
			Cmds:      []core.Cmd[any]{cmd1, cmd2},
			Deadlines: []time.Time{time.Unix(1, 0)},
			Results:   make(chan core.AsyncResult, 1),
		},
		Want: MultiSendWant{
			Seqs:    []core.Seq{wantSeq1, wantSeq2},
			Ns:      []int{0, 1},
			Errs:    []error{cln.NewClientError(wantErr), nil},
			Results: []core.AsyncResult{},
			Mocks:   []*mok.Mock{delegate.Mock, cmd1.Mock, cmd2.Mock},
			Has:     false,
		},
		CheckDuring: func(t *testing.T, client *cln.Client[any], seqs []core.Seq) {
			close(receiveDone)
		},
	}
}

// -----------------------------------------------------------------------------

func MultiSuccessTestCase() MultiSendTestCase[any] {
	name := "Should successfully send multiple cmds and receive results"

	var (
		wantSeq1    core.Seq = 1
		wantSeq2    core.Seq = 2
		wantResult1          = mock.NewResult()
		wantResult2          = mock.NewResult()
		cmd1                 = mock.NewCmd[any]()
		cmd2                 = mock.NewCmd[any]()
		sendDone             = make(chan struct{}, 2)
		delegate             = mock.NewClientDelegate()
	)
	delegate.RegisterSend(
		func(seq core.Seq, cmd core.Cmd[any]) (n int, err error) {
			return 1, nil
		},
	).RegisterSend(
		func(seq core.Seq, cmd core.Cmd[any]) (n int, err error) {
			return 1, nil
		},
	).RegisterFlush(
		func() (err error) {
			sendDone <- struct{}{}
			return nil
		},
	).RegisterFlush(
		func() (err error) {
			sendDone <- struct{}{}
			return nil
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			<-sendDone
			return wantSeq1, wantResult1, 10, nil
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			return wantSeq2, wantResult2, 20, nil
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			return 0, nil, 0, errors.New("receive error")
		},
	).RegisterClose(
		func() (err error) { return nil },
	)
	wantResult1.RegisterLastOne(func() bool { return true })
	wantResult2.RegisterLastOne(func() bool { return true })
	return MultiSendTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{},
		},
		Params: MultiSendParams[any]{
			Cmds:    []core.Cmd[any]{cmd1, cmd2},
			Results: make(chan core.AsyncResult, 2),
		},
		Want: MultiSendWant{
			Seqs: []core.Seq{wantSeq1, wantSeq2},
			Ns:   []int{1, 1},
			Errs: []error{nil, nil},
			Results: []core.AsyncResult{
				{Seq: wantSeq1, Result: wantResult1, BytesRead: 10},
				{Seq: wantSeq2, Result: wantResult2, BytesRead: 20},
			},
			Mocks: []*mok.Mock{delegate.Mock, cmd1.Mock, cmd2.Mock, wantResult1.Mock, wantResult2.Mock},
			Has:   false,
		},
	}
}

// -----------------------------------------------------------------------------

func IncrementSeqTestCase() MultiSendTestCase[any] {
	name := "Client sequence should be incremented"

	var (
		wantSeq1  core.Seq = 1
		wantSeq2  core.Seq = 2
		cmd1               = mock.NewCmd[any]()
		cmd2               = mock.NewCmd[any]()
		checkDone          = make(chan struct{})
		delegate           = mock.NewClientDelegate()
	)
	delegate.RegisterSend(
		func(seq core.Seq, cmd core.Cmd[any]) (n int, err error) { return 1, nil },
	).RegisterSend(
		func(seq core.Seq, cmd core.Cmd[any]) (n int, err error) { return 1, nil },
	).RegisterFlush(
		func() (err error) { return nil },
	).RegisterFlush(
		func() (err error) { return nil },
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			<-checkDone
			return 0, nil, 0, errors.New("receive error")
		},
	).RegisterClose(
		func() (err error) { return nil },
	)
	return MultiSendTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{},
		},
		Params: MultiSendParams[any]{
			Cmds:    []core.Cmd[any]{cmd1, cmd2},
			Results: make(chan core.AsyncResult, 1),
		},
		Want: MultiSendWant{
			Seqs:    []core.Seq{wantSeq1, wantSeq2},
			Ns:      []int{1, 1},
			Errs:    []error{nil, nil},
			Results: []core.AsyncResult{},
			Mocks:   []*mok.Mock{delegate.Mock, cmd1.Mock, cmd2.Mock},
			Has:     false,
		},
		CheckDuring: func(t *testing.T, client *cln.Client[any], seqs []core.Seq) {
			asserterror.Equal(t, client.Has(seqs[0]), true)
			asserterror.Equal(t, client.Has(seqs[1]), true)
			close(checkDone)
		},
	}
}

// -----------------------------------------------------------------------------

func MultiResultSuccessTestCase() MultiSendTestCase[any] {
	name := "Should successfully send one cmd and receive multiple results"

	var (
		wantSeq     core.Seq = 1
		wantResult1          = mock.NewResult()
		wantResult2          = mock.NewResult()
		cmd                  = mock.NewCmd[any]()
		sendDone             = make(chan struct{})
		delegate             = mock.NewClientDelegate()
	)

	delegate.RegisterSend(
		func(seq core.Seq, cmd core.Cmd[any]) (n int, err error) {
			return 1, nil
		},
	).RegisterFlush(
		func() (err error) {
			close(sendDone)
			return nil
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			<-sendDone
			return wantSeq, wantResult1, 10, nil
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			return wantSeq, wantResult2, 20, nil
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			return 0, nil, 0, errors.New("receive error")
		},
	).RegisterClose(
		func() (err error) {
			return nil
		},
	)
	wantResult1.RegisterLastOne(func() bool { return false })
	wantResult2.RegisterLastOne(func() bool { return true })

	return MultiSendTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{},
		},
		Params: MultiSendParams[any]{
			Cmds:    []core.Cmd[any]{cmd},
			Results: make(chan core.AsyncResult, 2),
		},
		Want: MultiSendWant{
			Seqs: []core.Seq{wantSeq},
			Ns:   []int{1},
			Errs: []error{nil},
			Results: []core.AsyncResult{
				{Seq: wantSeq, Result: wantResult1, BytesRead: 10},
				{Seq: wantSeq, Result: wantResult2, BytesRead: 20},
			},
			Mocks: []*mok.Mock{delegate.Mock, cmd.Mock, wantResult1.Mock, wantResult2.Mock},
			Has:   false,
		},
	}
}

// -----------------------------------------------------------------------------

func PartialResultsTestCase() MultiSendTestCase[any] {
	name := "Should remember the comand after partial results"

	var (
		wantSeq         core.Seq = 1
		wantResult1              = mock.NewResult()
		wantResult2              = mock.NewResult()
		cmd                      = mock.NewCmd[any]()
		delegate                 = mock.NewClientDelegate()
		firstResultSent          = make(chan struct{})
		checkDone                = make(chan struct{})
	)

	delegate.RegisterSend(
		func(seq core.Seq, cmd core.Cmd[any]) (n int, err error) {
			return 1, nil
		},
	).RegisterFlush(
		func() (err error) {
			return nil
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			seq, result, n, err = wantSeq, wantResult1, 10, nil
			close(firstResultSent)
			return
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			<-checkDone
			return wantSeq, wantResult2, 20, nil
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			return 0, nil, 0, errors.New("receive error")
		},
	).RegisterClose(
		func() (err error) {
			return nil
		},
	)
	wantResult1.RegisterLastOne(func() bool { return false })
	wantResult2.RegisterLastOne(func() bool { return true })

	return MultiSendTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{},
		},
		Params: MultiSendParams[any]{
			Cmds:    []core.Cmd[any]{cmd},
			Results: make(chan core.AsyncResult, 2),
		},
		Want: MultiSendWant{
			Seqs: []core.Seq{wantSeq},
			Ns:   []int{1},
			Errs: []error{nil},
			Results: []core.AsyncResult{
				{Seq: wantSeq, Result: wantResult1, BytesRead: 10},
				{Seq: wantSeq, Result: wantResult2, BytesRead: 20},
			},
			Mocks: []*mok.Mock{delegate.Mock, cmd.Mock, wantResult1.Mock, wantResult2.Mock},
			Has:   false,
		},
		CheckDuring: func(t *testing.T, client *cln.Client[any], seqs []core.Seq) {
			<-firstResultSent
			asserterror.Equal(t, client.Has(seqs[0]), true)
			close(checkDone)
		},
	}
}

// -----------------------------------------------------------------------------

func IncrementSeqAfterFailTestCase() MultiSendTestCase[any] {
	name := "Should increment seq even after the cmd send has been failed"

	var (
		wantSeq1    core.Seq = 1
		wantSeq2    core.Seq = 2
		wantErr              = errors.New("send error")
		cmd1                 = mock.NewCmd[any]()
		cmd2                 = mock.NewCmd[any]()
		receiveDone          = make(chan struct{})
		delegate             = mock.NewClientDelegate()
	)

	delegate.RegisterSend(
		func(seq core.Seq, cmd core.Cmd[any]) (n int, err error) {
			return 0, wantErr
		},
	).RegisterSend(
		func(seq core.Seq, cmd core.Cmd[any]) (n int, err error) {
			return 1, nil
		},
	).RegisterFlush(
		func() (err error) {
			return nil
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			<-receiveDone
			return 0, nil, 0, errors.New("receive error")
		},
	).RegisterClose(
		func() (err error) {
			return nil
		},
	)

	return MultiSendTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{},
		},
		Params: MultiSendParams[any]{
			Cmds:    []core.Cmd[any]{cmd1, cmd2},
			Results: make(chan core.AsyncResult, 1),
		},
		Want: MultiSendWant{
			Seqs:    []core.Seq{wantSeq1, wantSeq2},
			Ns:      []int{0, 1},
			Errs:    []error{cln.NewClientError(wantErr), nil},
			Results: []core.AsyncResult{},
			Mocks:   []*mok.Mock{delegate.Mock, cmd1.Mock, cmd2.Mock},
			Has:     false,
		},
		CheckDuring: func(t *testing.T, client *cln.Client[any], seqs []core.Seq) {
			close(receiveDone)
		},
	}
}

// -----------------------------------------------------------------------------

func ErrForAllCmdsOnFlushFailTestCase() MultiSendTestCase[any] {
	name := "If Delegate.Flush fails with an error, Send of all involved Commands should return error"

	var (
		delegateErr = errors.New("flush error")
		wantErr     = cln.NewClientError(delegateErr)
		cmds        = make([]core.Cmd[any], 10)
		wantErrs    = make([]error, 10)
		delegate    = mock.NewClientDelegate()
	)
	for i := range 10 {
		cmds[i] = mock.NewCmd[any]()
		wantErrs[i] = wantErr
	}
	delegate.RegisterSendN(10,
		func(seq core.Seq, cmd core.Cmd[any]) (n int, err error) { return 0, nil },
	).RegisterFlushN(10,
		func() (err error) {
			time.Sleep(test.TimeDelta)
			return delegateErr
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			time.Sleep(test.TimeDelta)
			return 0, nil, 0, errors.New("receive error")
		},
	).RegisterClose(
		func() (err error) { return nil },
	)
	return MultiSendTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{},
		},
		Params: MultiSendParams[any]{
			Cmds:    cmds,
			Results: nil,
		},
		Want: MultiSendWant{
			// Mocks are intentionally omitted: Flush count is non-deterministic
			// under concurrency (multiple goroutines may batch under one Flush call).
			Errs: wantErrs,
			Has:  false,
		},
		Concurrent: true,
	}
}

// Close -----------------------------------------------------------------------

func CloseSuccessTestCase() ClientTestCase[any] {
	name := "After Close the done channel should be closed"

	var delegate = mock.NewClientDelegate()

	delegate.RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			return 0, nil, 0, errors.New("receive error")
		},
	).RegisterClose(
		func() (err error) { return nil },
	)
	return ClientTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{},
		},
		Action: func(t *testing.T, client *cln.Client[any], results chan core.AsyncResult) {
			err := client.Close()
			asserterror.EqualError(t, err, nil)
			AssertDone(t, client)
		},
		Mocks: []*mok.Mock{delegate.Mock},
	}
}

// -----------------------------------------------------------------------------

func CloseDuringQueueResultTestCase() ClientTestCase[any] {
	name := "Should be able to close while queuing the result"

	var (
		wantCmd    = mock.NewCmd[any]()
		wantResult = mock.NewResult()
		results    = make(chan core.AsyncResult) // unbuffered — delivery will block
		delegate   = mock.NewClientDelegate()
	)

	wantResult.RegisterLastOne(func() bool { return true })

	delegate.RegisterSend(
		func(seq core.Seq, cmd core.Cmd[any]) (n int, err error) {
			return 1, nil
		},
	).RegisterFlush(
		func() (err error) {
			return nil
		},
	).RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			return 1, wantResult, 0, nil
		},
	).RegisterClose(
		func() (err error) {
			return nil
		},
	)

	return ClientTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{},
		},
		Action: func(t *testing.T, client *cln.Client[any], res chan core.AsyncResult) {
			_, _, err := client.Send(wantCmd, results)
			asserterror.EqualError(t, err, nil)
			// give the receive goroutine time to block on sending to results
			time.Sleep(test.TimeDelta)

			err = client.Close()
			asserterror.EqualError(t, err, nil)
		},
		Mocks: []*mok.Mock{wantCmd.Mock, wantResult.Mock, delegate.Mock},
	}
}

// -----------------------------------------------------------------------------

func CloseDelegateFailTestCase() ClientTestCase[any] {
	name := "If Delegate.Close fails with an error, Close should return it"

	var (
		delegateErr = errors.New("Delegate.Close error")
		wantErr     = cln.NewClientError(delegateErr)
		receiveDone = make(chan struct{})
		delegate    = mock.NewClientDelegate()
	)
	delegate.RegisterReceive(
		func() (seq core.Seq, result core.Result, n int, err error) {
			<-receiveDone
			return 0, nil, 0, errors.New("receive error")
		},
	).RegisterClose(
		func() (err error) { return delegateErr },
	)
	return ClientTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{},
		},
		Action: func(t *testing.T, client *cln.Client[any], results chan core.AsyncResult) {
			err := client.Close()
			asserterror.EqualError(t, err, wantErr)
			close(receiveDone)
		},
		Mocks: []*mok.Mock{delegate.Mock},
	}
}

// Unexpected Results ----------------------------------------------------------

func UnexpectedResultTestCase() ClientTestCase[any] {
	name := "Should ignore unexpected results (results for unknown sequence numbers)"

	var (
		unexpectedSeq core.Seq = 100
		result                 = mock.NewResult()
		receiveDone            = make(chan struct{})
		delegate               = mock.NewClientDelegate()
	)

	result.RegisterLastOne(func() bool { return true })

	delegate.RegisterReceive(
		func() (seq core.Seq, res core.Result, n int, err error) {
			return unexpectedSeq, result, 10, nil
		},
	).RegisterReceive(
		func() (seq core.Seq, res core.Result, n int, err error) {
			<-receiveDone
			return 0, nil, 0, errors.New("receive error")
		},
	).RegisterClose(
		func() (err error) {
			close(receiveDone)
			return nil
		},
	)

	return ClientTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{},
		},
		Action: func(t *testing.T, client *cln.Client[any], results chan core.AsyncResult) {
			// Give the receive goroutine time to process the unexpected result.
			// If it hangs (deadlock), this test will eventually timeout.
			time.Sleep(test.TimeDelta)
			_ = client.Close()
		},
		Mocks: []*mok.Mock{delegate.Mock, result.Mock},
	}
}

// -----------------------------------------------------------------------------

func UnexpectedResultCallbackTestCase() ClientTestCase[any] {
	name := "Should invoke UnexpectedResultCallback when an unexpected result is received"

	var (
		unexpectedSeq core.Seq = 100
		wantResult             = mock.NewResult()
		receiveDone            = make(chan struct{})
		callbackDone           = make(chan struct{})
		delegate               = mock.NewClientDelegate()
		gotSeq        core.Seq
		gotResult     core.Result
	)

	wantResult.RegisterLastOne(func() bool { return true })

	delegate.RegisterReceive(
		func() (seq core.Seq, res core.Result, n int, err error) {
			return unexpectedSeq, wantResult, 10, nil
		},
	).RegisterReceive(
		func() (seq core.Seq, res core.Result, n int, err error) {
			<-receiveDone
			return 0, nil, 0, errors.New("receive error")
		},
	).RegisterClose(
		func() (err error) {
			close(receiveDone)
			return nil
		},
	)

	callback := func(seq core.Seq, result core.Result) {
		gotSeq = seq
		gotResult = result
		close(callbackDone)
	}

	return ClientTestCase[any]{
		Name: name,
		Setup: ClientSetup[any]{
			Delegate: delegate,
			Opts:     []cln.SetOption{cln.WithUnexpectedResultCallback(callback)},
		},
		Action: func(t *testing.T, client *cln.Client[any], results chan core.AsyncResult) {
			select {
			case <-callbackDone:
				asserterror.Equal(t, gotSeq, unexpectedSeq)
				asserterror.EqualDeep(t, gotResult, wantResult)
			case <-time.After(time.Second):
				t.Fatal("timeout waiting for UnexpectedResultCallback")
			}
			_ = client.Close()
		},
		Mocks: []*mok.Mock{delegate.Mock, wantResult.Mock},
	}
}
