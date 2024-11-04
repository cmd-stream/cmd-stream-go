# cmd-stream-go
cmd-stream-go is a high-performance client-server library that implements the 
Command pattern, supports reconnect and keepalive features.

# Tests
Test coverage of each submodule is over 90%.

# Benchmarks
[github.com/ymz-ncnk/go-client-server-communication-benchmarks](https://github.com/ymz-ncnk/go-client-server-communication-benchmarks)

# Command Pattern Over Network vs RPC
By mapping commands to remote procedure calls, it is quite easy to implement the
RPC approach. Thus, if you are already using one of the RPC products, you can 
switch to cmd-stream-go even without changing your interfaces.

# Network Protocols Support
cmd-stream-go is built on top of the standard Golang net package, and supports 
connection-oriented protocols like TCP, TLS or mutual TLS (for client
authentication).

# Client
The client is asynchronous and can be used from different gorountines 
simultaneously. Also it uses only one connection to send commands and receive 
results.

With `client.NewReconnect()`, you can create a client that tries to reconnect to
the server if it loses the connection.

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

A small example:
```go
// 1. First of all we have to define Receiver.
type Calculator struct{}

func (c Calculator) Add(n1, n2 int) int {...}

func (c Calculator) Sub(n1, n2 int) int {...}

// 2. Than a command. All commands should implement base.Cmd[T] interface.
type Eq1Cmd struct {...}

func (c Eq1Cmd) Exec(ctx context.Context, at time.Time, seq base.Seq,
  receiver Calculator,
  proxy base.Proxy,
) error {
  // It uses Receiver here.
  result := Result(receiver.Add(...))
  // And sends back result.
  return proxy.Send(seq, result)
}

// 3. Than a result. All results should implement the base.Result interface. The 
// client will wait for more command results if the LastOne method of the 
// received result returns false.
type Result int

func (r Result) LastOne() bool {
  return true
}

// 4. Than a client codec, which should implement the cs_client.Codec[T] 
// interface.
type ClientCodec struct{}

// Encode is used by the client to send commands to the server. If Encode fails
// with an error, the Client.Send method will return it.
func (c ClientCodec) Encode(cmd base.Cmd[Calculator], w transport.Writer) (
  err error) {...}

// Decode is used by the client to receive resulsts from the server. If Decode
// fails with an error, the client will be closed.
func (c ClientCodec) Decode(r transport.Reader) (result base.Result, 
err error) {...}

// Size returns the size of the command in bytes. If the server imposes any
// restrictions on the command size, the client will use this method to
// check it before sending.
func (c ClientCodec) Size(cmd base.Cmd[Calculator]) (size int) {...}

// 5. Than a server codec, which should implement the cs_server.Codec[T] 
// interface.
type ServerCodec struct{}

// Encode is used by the server to send results to the client. If Encode fails
// with an error, the server closes the connection.
func (c ServerCodec) Encode(result base.Result, w transport.Writer) (
  err error) {...}

// Decode is used by the server to receive commands from the client. If Decode
// fails with an error, the server closes the connection.
func (c ServerCodec) Decode(r transport.Reader) (cmd base.Cmd[Calculator],
  err error) {...}

// 6. And that's it, the only thing left to do is to create a server and client.
// Create the server.
server := cs_server.NewDef[Calculator](ServerCodec{}, Calculator{})
// Start the server.
listener, err := net.Listen("tcp", Addr)
...
go func() {
  ...
  server.Serve(listener.(*net.TCPListener))
}()
// Connect to the server.
conn, err := net.Dial("tcp", Addr)
...
// Create the client.
client, err := cs_client.NewDef[Calculator](ClientCodec{}, conn, nil)
...
```
You can find the full code of this example, called 
[standard](https://github.com/cmd-stream/cmd-stream-examples-go/tree/main/standard) 
and several other examples of using cmd-stream-go in 
[cmd-stream-examples-go](https://github.com/cmd-stream/cmd-stream-examples-go).

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