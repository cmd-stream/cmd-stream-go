# cmd-stream-go
cmd-stream-go is a high-performance RCX (**R**emote **C**ommand e**X**ecution) 
library, which applies the Command design pattern to the client-server 
architecture and supports reconnect and keepalive features.

# Tests
Test coverage of each submodule is over 90%.

# Benchmarks
[github.com/ymz-ncnk/go-client-server-communication-benchmarks](https://github.com/ymz-ncnk/go-client-server-communication-benchmarks)

# RCX vs RPC
By mapping commands to remote procedure calls, you can simply implement RPC with 
RCX. This way, if you are already using one of the RPC products, you can switch 
to RCX even without modifying your interfaces.

# Network Protocols Support
cmd-stream-go is built on top of the standard Golang net package, and supports 
connection-oriented protocols like TCP or TLS (if you need to authenticate 
clients as well, you can use mutual TLS).

# Client
The client is asynchronous and can be used from different gorountines 
simultaneously. Also it uses only one connection to send commands and receive 
results.

You can create a regular client or a "reconnet" client. The last one tries to 
reconnect to the server if it has lost the connection.

Among the client configuration, you can find (and not only):
- KeepaliveTime and KeepaliveIntvl - if both of them != 0, client will try to
  keep the connection alive. If there are no commands to send, it starts 
  Ping-Pong with the server - sends a Ping command and receives a Pong result, 
  both of which are transfered as a 0 (like a ball) byte.
- SysDataReceiveTimeout - determines how long the client will wait for system 
  data from the server.

# Server
Before starting to receive commands from the client, the server sends it system 
data: `ServerInfo` and `ServerSettings`. With `ServerInfo`, the client can 
determine  its compatibility with the server, for example, whether it and the 
server support the same set of commands. `ServerSettings`, in turn, contains the 
desired settings for interacting with the server.

It should also be noted that the number of simultaneous server clients is 
limited (it can be configured). And each command on the server is executed in a
separate gorountine, with help of user-defined `Invoker` and `Receiver`. Also, 
a command can have more than one result.

Among the server configuration, you can find (and not only):
- FirstConnTimeout - the server will close if it does not receive the first 
  connection during this time.
- WorkersCount - each connection to the client is processed on the server by one 
  `Worker`.	That is, this parameter sets the number of simultaneous clients 
	of the server.
- LostConnCallback - called when the server loses connection with the client.
- ReceiveTimeout - if the server has not received any commands from the client 
  during this time, it closes the connection.

# How To Use
All we need to do is define `Receiver`, commands, results, and codecs for 
the client and server.

The client codec encodes commands and decodes results from the connection.
The server codec does the same thing, but in reverse. cms-stream-go was designed
with [mus-stream-go](https://github.com/mus-format/mus-stream-go) in mind,
but you can use any other serializer with it.

Thanks to the super simple [MUS format](https://github.com/mus-format/specification), 
the mus-stream-go serializer uses a small number of bytes to encode the data.
Also, with mus-stream-go there is no need to put the length of the data before
the data itself. This all can have a positive impact on your bandwidth.

In [cmd-stream-examples-go](https://github.com/cmd-stream/cmd-stream-examples-go)
you can find examples of using cmd-stream-go.

# Architecture
There are the following cmd-stream-go submodules:
- [base-go](https://github.com/cmd-stream/base-go) - basic module for creating 
  the client and server.
- [delegate-go](https://github.com/cmd-stream/delegate-go) - the client entrusts 
  all its communication-related work to the delegate. The server does the same. 
  The connection is also initialized at this level.
- [handler-go](https://github.com/cmd-stream/handler-go) - the server delegate 
  uses a handler to receive commands, execute them, and return results. Here you 
  can find a `Proxy` definition (the proxy of the server transport), which 
  allows commands to send back results.
- [transport-go](https://github.com/cmd-stream/transport-go) is resposible for 
  commands/results delivery. Here you can find a `Codec` definition.

cmd-stream-go was designed in such a way that you can easily replace any part of 
it.