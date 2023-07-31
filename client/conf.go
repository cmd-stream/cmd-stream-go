package client

import (
	delegate_client "github.com/cmd-stream/delegate-go/client"
	transport_common "github.com/cmd-stream/transport-go/common"
)

// Conf is a Сlient configuration.

// It is a union of the transport and delegate configurations.
type Conf struct {
	Transport transport_common.Conf
	Delegate  delegate_client.Conf
}

// KeepaliveOn checks if keepalive is enabled.
func (c Conf) KeepaliveOn() bool {
	return c.Delegate.KeepaliveTime != 0 && c.Delegate.KeepaliveIntvl != 0
}
