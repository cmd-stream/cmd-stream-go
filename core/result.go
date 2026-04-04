package core

// Result represents the outcome of the Сommand execution.
//
// LastOne method indicates the final Result of the Command.
type Result interface {
	LastOne() bool
}

// AsyncResult represents an asynchronous result.
//
// Seq is the sequence number of the Command, Error != nil if something went
// wrong with the connection.
type AsyncResult struct {
	Seq       Seq
	BytesRead int
	Result    Result
	Error     error
}
