# cmd-stream-go
cmd-stream-go is a high-performance client-server library that implements the 
Command pattern.

# Brief cmd-stream-go Description
- Can work over TCP, TLS or mutual TLS.
- The client is asynchronous and can be used by several goroutines.
- Only one connection is used per client.
- Supports a deadline for sending a command or result.
- Supports timeouts.
- Supports the server streaming, i.e. a command can send back multiple results 
  (client streaming is not directly supported, but can also be implemented).
- The server can configure the client.
- Supports reconnect feature.
- Supports keepalive feature.
- Has a flexible architecture.

# Contents
- [cmd-stream-go](#cmd-stream-go)
- [Brief cmd-stream-go Description](#brief-cmd-stream-go-description)
- [Contents](#contents)
- [Tests](#tests)
- [Benchmarks](#benchmarks)
- [The Command Pattern as an Alternative to RPC](#the-command-pattern-as-an-alternative-to-rpc)
- [High-performance Communication Channel](#high-performance-communication-channel)
- [Network Protocols Support](#network-protocols-support)
- [Client](#client)
  - [Configuration](#configuration)
  - [Waiting for the Result with a Timeout](#waiting-for-the-result-with-a-timeout)
  - [Reconect](#reconect)
- [Server](#server)
  - [Configuration](#configuration-1)
  - [Command Size Restriction](#command-size-restriction)
- [How To Use](#how-to-use)
- [Architecture](#architecture)

# Tests
Test coverage of each submodule is over 90%.

# Benchmarks
[github.com/ymz-ncnk/go-client-server-communication-benchmarks](https://github.com/ymz-ncnk/go-client-server-communication-benchmarks)

# The Command Pattern as an Alternative to RPC
[https://medium.com/p/b08b3b2bba35](https://medium.com/p/b08b3b2bba35)

# High-performance Communication Channel
To build a high-performance communication channel between two services:
1. Use N connections. Because several connections can transmit much more 
   information than one. The number N depends on your system and can indicate 
   the number of connections after which adding another one will not provide any 
   benefits.
2. To improve the responsiveness of the system use all available connections
   from the start, instead of creating new ones as needed.
3. Use keepalive connections.

# Network Protocols Support
cmd-stream-go is built on top of the standard Golang net package, and supports 
connection-oriented protocols like TCP, TLS or mutual TLS (for client
authentication).

# Client
The client is asynchronous and can be used from different goroutines 
simultaneously. It uses only one connection to send commands and receive 
results. Commands sent from a single goroutine will be delivered to the server 
in order.

## Configuration
Client configuration options include (and not only):
- KeepaliveTime and KeepaliveIntvl - if both of them != 0, client will try to
  keep the connection alive. If there are no commands to send, it starts 
  Ping-Pong with the server - sends a Ping command and receives a Pong result, 
  both of which are transfered as a 0 (like a ball) byte.
- SysDataReceiveTimeout - determines how long the client will wait for system 
  data from the server.

## Waiting for the Result with a Timeout
```go
...
results := make(chan base.AsyncResult, 1)
seq, err := client.Send(cmd, results) // Where seq is the sequence number of the command.
...
select {
case <-time.NewTimer(3 * time.Second).C:
    client.Forget(seq)
  // Handle timeout.
case result := <-results:
  // Handle result.
}
```

## Reconect
`client.NewDefReconnect()` will create a client that attempts to reconnect to 
the server if the connection is lost. This can happen while sending a command -
we'll get an error, or while waiting for the result - we will not be sure 
whether the command was executed on the server or not.

In both cases, after a while, we can try to send the command again (idempotent
command). If the connection is restored, normal operation will continue, 
otherwise we will get the error again.

Regarding the time interval before retrying, it is better to choose it randomly 
for each goroutine to avoid overloading the server with a large number of 
simultaneously sent commands.

# Server
Before starting to receive commands, the server sends to the client system 
data: `ServerInfo` and `ServerSettings`. With `ServerInfo`, the client can 
determine  its compatibility with the server, for example, whether it and the 
server support the same set of commands. `ServerSettings`, in turn, contains the 
desired settings for interacting with the server.

A few words about commands execution:
- Each command is executed by a single `Invoker` (it should be thread-safe) in 
  a separete goroutine.
- There is a default `Invoker`, but you can specify your own.
- Command can send back several results. They all will be delivered to the 
  client in order.

## Configuration
Server configuration options include (and not only):
- FirstConnTimeout - the server will close if it does not receive the first 
  connection during this time.
- WorkersCount - each connection to the client is processed on the server by one 
  `Worker`.	That is, this parameter sets the maximum number of simultaneous 
  clients of the server.
- LostConnCallback - called when the server loses connection with the client.
- ReceiveTimeout - if the server has not received any commands from the client 
  during this time, it closes the connection.

## Command Size Restriction
The server may ask the client not to send too large commands.

To enable this, simply set `Conf.ServerSettings.MaxCmdSize` in bytes and 
implement the client codec's `Size()` method, it will be used to verify the 
command size.

Please note that even with this feature, the server must protect itself from
receiving too large commands. This can be done while decoding a command - the 
server codec's `Decode()` method may return an error, which will close the 
connection to the client.

# How To Use
All you need to do is implement the Command pattern and codecs - one for the 
client and one for the server:
1. First of all define the Receiver. In this case it will be a `Calculator` with
   two methods `Add()` and `Sub()`:
```go
type Calculator struct{}

func (c Calculator) Add(n1, n2 int) int {...}

func (c Calculator) Sub(n1, n2 int) int {...}
```

2. Define the Command and Result.
```go
// Eq1Cmd is an equation that we want to calculate on the server. It implements 
// base.Cmd[T] interface, where T is a Receiver.
type Eq1Cmd struct {...}

// Exec method will be called by the Invoker on the server.
func (c Eq1Cmd) Exec(ctx context.Context, 
  at time.Time, // If the server was configured with Conf.Handler.At == true, 
  // contains the command receiving time.
  seq base.Seq, // The sequence number of the command. It is used to send back 
  // results.
  receiver Calculator, // Receiver.
  proxy base.Proxy, // Allows command to send back results. Contains only
  // two methods: Send() and SendWithDeadline().
) error {
  // It uses Receiver here.
  result := Result(receiver.Add(...))
  // And sends back the result. In general, a command can send back several 
  // results, which will be received by the client in order.
  // If an error was encountered during execution, the command can send it back 
  // to the client as a result, or it can simply return it. In the latter case, 
  // the connection to the client will be closed.
  return proxy.Send(seq, result)
}

// Result is the result of calculating the equation on the server. It implements 
// the base.Result interface.
type Result int

// LastOne determines whether the result is the last one. If it returns false, 
// the client will wait for the next one.
func (r Result) LastOne() bool {
  return true
}
```

3. Define the client Codec.
```go
// ClientCodec encodes commands to the Writer and decodes results from the 
// Reader. It should implement the client.Codec[T] interface (from the 
// cmd-stream-go module), where T is a Receiver.
type ClientCodec struct{}

// Encode is used by the client to send commands to the server. If it fails with
// an error, the Client.Send() method will return it.
func (c ClientCodec) Encode(cmd base.Cmd[Calculator], w transport.Writer) (
  err error) {...}

// Decode is used by the client to receive results from the server. If it fails
// with an error, the client will be closed.
func (c ClientCodec) Decode(r transport.Reader) (result base.Result, 
err error) {...}

// Size returns the size of the command in bytes. If the server imposes any
// restrictions on the command size, the client will use this method to
// check it before sending.
func (c ClientCodec) Size(cmd base.Cmd[Calculator]) (size int) {...}
```

4. Define the server Codec.
```go
// ServerCodec encodes results to the Writer and decodes commands from the 
// Reader. It should implement the server.Codec[T] interface (from the 
// cmd-stream-go module), where T is a Receiver.
type ServerCodec struct{}

// Encode is used by the server to send results to the client. If it fails with 
// an error, the server closes the connection.
func (c ServerCodec) Encode(result base.Result, w transport.Writer) (
  err error) {...}

// Decode is used by the server to receive commands from the client. If it fails
// with an error, the server closes the connection.
func (c ServerCodec) Decode(r transport.Reader) (cmd base.Cmd[Calculator],
  err error) {...}
```

6. Create the server.
```go
server := cs_server.NewDef[Calculator](ServerCodec{}, Calculator{})
// Make the listener.
listener, err := net.Listen("tcp", Addr)
...
go func() {
  ...
  // Start the server.
  server.Serve(listener.(*net.TCPListener))
}()
```

7. Create the client.
```go
// Connect to the server.
conn, err := net.Dial("tcp", Addr)
...
// Create the client.
client, err := cs_client.NewDef[Calculator](ClientCodec{}, conn, nil)
...
```

8. Send a command and get the result.
```go
...
results := make(chan base.AsyncResult, 1)
_, err := client.Send(cmd, results)
...
result := (<-results).Result
...
```
The full code of this example, called [standard](https://github.com/cmd-stream/cmd-stream-examples-go/tree/main/standard) 
and several other examples of using cmd-stream-go can be found in 
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