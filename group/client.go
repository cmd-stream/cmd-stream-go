package cgrp

import (
	"time"

	"github.com/cmd-stream/base-go"
)

// ClientID identifies a specific client within a Group.
type ClientID int

// Client represents a client used by the Group for sending commands and
// receiving results.
type Client[T any] interface {
	Send(cmd base.Cmd[T], results chan<- base.AsyncResult) (seq base.Seq, n int,
		err error)
	SendWithDeadline(cmd base.Cmd[T], results chan<- base.AsyncResult,
		deadline time.Time) (seq base.Seq, n int, err error)
	Has(seq base.Seq) bool
	Forget(seq base.Seq)
	Err() (err error)
	Close() (err error)
	Done() <-chan struct{}
}
