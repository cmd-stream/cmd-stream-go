package server

import (
	base_server "github.com/cmd-stream/base-go/server"
	delegate_server "github.com/cmd-stream/delegate-go/server"
	handler "github.com/cmd-stream/handler-go"
	transport_common "github.com/cmd-stream/transport-go/common"
)

// Conf is a Server configuration.
//
// Represents a union of the Transport, Handler, Delegate and Base configurations.
type Conf struct {
	Transport transport_common.Conf
	Handler   handler.Conf
	Delegate  delegate_server.Conf
	Base      base_server.Conf
}
