// Package core defines the fundamental interfaces and types for the cmd-stream
// protocol, including Commands, Results, and sequence numbers (Seq).
//
// It contains definitions for both the client and server components of the
// library.
//
// # Internal Workflow
//
// In cmd-stream, the Client delegates all communication-related tasks,
// such as sending Commands, receiving Results, and managing the connection
// lifecycle, to a ClientDelegate.
//
// Similarly, the Server handles client connections through a ServerDelegate,
// using a configurable number of background workers to process incoming
// connections and execute Commands.
package core
