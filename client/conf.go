package client

import (
	dcln "github.com/cmd-stream/delegate-go/client"
	tcom "github.com/cmd-stream/transport-go/common"
)

// Conf configures the client.
//
// Represents a union of the Transport and Delegate configurations.
type Conf struct {
	Transport tcom.Conf
	Delegate  dcln.Conf
}

// KeepaliveOn returns true if keepalive is enabled.
func (c Conf) KeepaliveOn() bool {
	return c.Delegate.KeepaliveTime != 0 && c.Delegate.KeepaliveIntvl != 0
}
