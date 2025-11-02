package integration_test

import (
	"errors"
	"time"

	"github.com/cmd-stream/core-go"
)

func receiveResult(results <-chan core.AsyncResult) (result core.AsyncResult,
	err error,
) {
	select {
	case <-time.NewTimer(time.Second).C:
		err = errors.New("test lasts too long")
	case result = <-results:
	}
	return
}
