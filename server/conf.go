package server

import (
	bser "github.com/cmd-stream/base-go/server"
	dser "github.com/cmd-stream/delegate-go/server"
	"github.com/cmd-stream/handler-go"
	tcom "github.com/cmd-stream/transport-go/common"
)

// Conf configures the server.
//
// Represents a union of the Transport, Handler, Delegate and Base configurations.
type Conf struct {
	Transport tcom.Conf
	Handler   handler.Conf
	Delegate  dser.Conf
	Base      bser.Conf
}
