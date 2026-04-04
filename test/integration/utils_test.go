package integration

import (
	"errors"
	"testing"
	"time"

	"github.com/cmd-stream/cmd-stream-go/core"
	"github.com/cmd-stream/cmd-stream-go/testkit"
	asserterror "github.com/ymz-ncnk/assert/error"
)

func assertResult(t *testing.T, asyncResult core.AsyncResult, seq core.Seq,
	expectedResult testkit.Result) {
	asserterror.Equal(t, asyncResult.Seq, seq)
	asserterror.Equal(t, asyncResult.BytesRead, testkit.CalcResultSize(seq,
		expectedResult))
	asserterror.EqualDeep(t, asyncResult.Result, expectedResult)
}

func receiveAndAssert(t *testing.T, results <-chan core.AsyncResult, seq core.Seq,
	expectedResult testkit.Result) {
	var (
		asyncResult core.AsyncResult
		err         error
		timer       = time.NewTimer(time.Second)
	)
	defer timer.Stop()
	select {
	case <-timer.C:
		err = errors.New("test lasts too long")
	case asyncResult = <-results:
	}
	asserterror.EqualError(t, err, nil)
	assertResult(t, asyncResult, seq, expectedResult)
}
