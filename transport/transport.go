// Package transport provides the core abstractions and implementations required
// to deliver Commands and Results between a cmd-stream client and server.
//
// # Implementations
//
// It provides concrete implementations of the ClientTransport and
// ServerTransport interfaces defined in the delegate package.
//
// # Buffered I/O and Encoding
//
// All transport implementations use efficient buffered I/O via bufio.Reader and
// bufio.Writer, and rely on user-defined codecs to handle the serialization of
// Commands and Results.
package transport
