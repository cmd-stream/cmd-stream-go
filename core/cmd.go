package core

import (
	"context"
	"time"
)

// Cmd defines the general Command interface.
//
// The Exec method is invoked by the server's Invoker and is responsible for
// executing the Command with the following parameters:
//   - ctx: Execution context.
//   - seq: Sequence number uniquely identifying the Command.
//   - at: Timestamp when the server received the Command.
//   - receiver: The Receiver of type T, which handles the Command's execution logic.
//   - proxy: A server transport proxy used to send Results back to the client.
type Cmd[T any] interface {
	Exec(ctx context.Context, seq Seq, at time.Time, receiver T, proxy Proxy) error
}
