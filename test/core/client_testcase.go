package core

import (
	"sync"
	"testing"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/cmd-stream/cmd-stream-go/core/cln"
	asserterror "github.com/ymz-ncnk/assert/error"
	"github.com/ymz-ncnk/mok"
)

type ClientTestCase[T any] struct {
	Name   string
	Setup  ClientSetup[T]
	Action func(t *testing.T, client *cln.Client[T], results chan core.AsyncResult)
	Mocks  []*mok.Mock
}

type ClientSetup[T any] struct {
	Delegate core.ClientDelegate[T]
	Opts     []cln.SetOption
}

func RunClientTestCase[T any](t *testing.T, tc ClientTestCase[T]) {
	t.Run(tc.Name, func(t *testing.T) {
		var (
			results = make(chan core.AsyncResult, 10)
			client  = cln.New(tc.Setup.Delegate, tc.Setup.Opts...)
		)
		if tc.Action != nil {
			tc.Action(t, client, results)
		}
		select {
		case <-client.Done():
		case <-time.After(time.Second):
			t.Fatal("timeout waiting for client to be done")
		}
		asserterror.EqualDeep(t, mok.CheckCalls(tc.Mocks), mok.EmptyInfomap)
	})
}

// -----------------------------------------------------------------------------

type MultiSendTestCase[T any] struct {
	Name        string
	Setup       ClientSetup[T]
	Params      MultiSendParams[T]
	Want        MultiSendWant
	CheckDuring func(t *testing.T, client *cln.Client[T], seqs []core.Seq)
	Concurrent  bool
}

type MultiSendParams[T any] struct {
	Cmds      []core.Cmd[T]
	Results   chan core.AsyncResult
	Deadlines []time.Time
}

type MultiSendWant struct {
	Seqs    []core.Seq
	Ns      []int
	Errs    []error
	Results []core.AsyncResult
	Mocks   []*mok.Mock
	Has     bool
}

func RunMultiSendTestCase[T any](t *testing.T, tc MultiSendTestCase[T]) {
	t.Run(tc.Name, func(t *testing.T) {
		var (
			client = cln.New(tc.Setup.Delegate, tc.Setup.Opts...)
			seqs   = make([]core.Seq, len(tc.Params.Cmds))
		)
		if tc.Concurrent {
			wg := sync.WaitGroup{}
			for i, cmd := range tc.Params.Cmds {
				wg.Add(1)
				go func() {
					defer wg.Done()
					_, _, err := client.Send(cmd, tc.Params.Results)
					asserterror.EqualError(t, err, tc.Want.Errs[i])
				}()
			}
			wg.Wait()
		} else {
			for i, cmd := range tc.Params.Cmds {
				var (
					seq core.Seq
					n   int
					err error
				)
				if i < len(tc.Params.Deadlines) && !tc.Params.Deadlines[i].IsZero() {
					seq, n, err = client.SendWithDeadline(tc.Params.Deadlines[i], cmd, tc.Params.Results)
				} else {
					seq, n, err = client.Send(cmd, tc.Params.Results)
				}
				seqs[i] = seq
				asserterror.Equal(t, seq, tc.Want.Seqs[i])
				asserterror.Equal(t, n, tc.Want.Ns[i])
				asserterror.EqualError(t, err, tc.Want.Errs[i])
			}
		}

		if tc.CheckDuring != nil {
			tc.CheckDuring(t, client, seqs)
		}

		for i := 0; i < len(tc.Want.Results); i++ {
			select {
			case result := <-tc.Params.Results:
				asserterror.EqualDeep(t, result, tc.Want.Results[i])
			case <-time.After(time.Second):
				t.Fatalf("timeout waiting for result %d", i)
			}
		}

		select {
		case <-client.Done():
		case <-time.After(time.Second):
			t.Fatal("timeout waiting for client to be done")
		}

		for _, seq := range tc.Want.Seqs {
			asserterror.Equal(t, client.Has(seq), tc.Want.Has)
		}
		asserterror.EqualDeep(t, mok.CheckCalls(tc.Want.Mocks), mok.EmptyInfomap)
	})
}

// -----------------------------------------------------------------------------

func AssertSend[T any](t *testing.T, client *cln.Client[T],
	cmd core.Cmd[T], results chan core.AsyncResult, wantSeq core.Seq,
	wantN int, wantErr error) (seq core.Seq) {
	t.Helper()
	seq, n, err := client.Send(cmd, results)
	asserterror.Equal(t, seq, wantSeq)
	asserterror.Equal(t, n, wantN)
	asserterror.EqualError(t, err, wantErr)
	return
}

func AssertSendWithDeadline[T any](t *testing.T, client *cln.Client[T],
	deadline time.Time, cmd core.Cmd[T], results chan core.AsyncResult,
	wantSeq core.Seq, wantN int, wantErr error) (seq core.Seq) {
	t.Helper()
	seq, n, err := client.SendWithDeadline(deadline, cmd, results)
	asserterror.Equal(t, seq, wantSeq)
	asserterror.Equal(t, n, wantN)
	asserterror.EqualError(t, err, wantErr)
	return
}

func AssertResults(t *testing.T, results <-chan core.AsyncResult,
	wantResults ...core.AsyncResult) {
	t.Helper()
	for i, want := range wantResults {
		select {
		case got := <-results:
			asserterror.EqualDeep(t, got, want)
		case <-time.After(time.Second):
			t.Fatalf("timeout waiting for result %d", i)
		}
	}
}

func AssertDone[T any](t *testing.T, client *cln.Client[T]) {
	t.Helper()
	select {
	case <-client.Done():
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for client to be done")
	}
}

func AssertHas[T any](t *testing.T, client *cln.Client[T], seq core.Seq,
	wantHas bool) {
	t.Helper()
	asserterror.Equal(t, client.Has(seq), wantHas)
}
