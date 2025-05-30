package cgrp

import (
	ccln "github.com/cmd-stream/cmd-stream-go/client"
)

// Options defines the configuration settings for creating a ClientGroup.
type Options[T any] struct {
	Factory   DispatchStrategyFactory[T]
	Reconnect bool
	ClientOps []ccln.SetOption
}

type SetOption[T any] func(o *Options[T])

// WithFactory sets the dispatch strategy factory for the client group.
//
// The dispatch strategy determines how Commands are distributed among clients.
// For example, a round-robin strategy will rotate client usage evenly.
func WithFactory[T any](factory DispatchStrategyFactory[T]) SetOption[T] {
	return func(o *Options[T]) { o.Factory = factory }
}

// WithReconnect enables automatic reconnection for all clients in the group.
//
// When this option is set, reconnect-capable clients are created, which attempt
// to re-establish the connection if it's lost during communication.
func WithReconnect[T any]() SetOption[T] {
	return func(o *Options[T]) { o.Reconnect = true }
}

// WithClientOps sets client-specific options to be applied when initializing
// each client in the group.
func WithClientOps[T any](ops ...ccln.SetOption) SetOption[T] {
	return func(o *Options[T]) { o.ClientOps = ops }
}

func ApplyGroup[T any](ops []SetOption[T], o *Options[T]) {
	for i := range ops {
		if ops[i] != nil {
			ops[i](o)
		}
	}
}
