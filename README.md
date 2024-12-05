# cmd-stream-go
cmd-stream-go is a high-performance client-server library that implements the 
Command pattern.

# Brief cmd-stream-go Description
- Can work over TCP, TLS or mutual TLS.
- Has an asynchronous client that uses only one connection for both sending 
  commands and receiving results.
- Supports the server streaming, i.e. a command can send back multiple results 
  (client streaming is not directly supported, but can also be implemented).
- Supports deadlines for sending commands and results.
- Supports timeouts.
- Supports reconnect feature.
- Supports keepalive feature.
- Can work with various serialization formats ([here](https://github.com/cmd-stream/cmd-stream-examples-go/tree/main/standard_protobuf) is an example using the Protobuf serializer).
- Has a flexible architecture.

# Contents
- [cmd-stream-go](#cmd-stream-go)
- [Brief cmd-stream-go Description](#brief-cmd-stream-go-description)
- [Contents](#contents)
- [Tests](#tests)
- [Benchmarks](#benchmarks)
- [High-performance Communication Channel](#high-performance-communication-channel)
- [Network Protocols Support](#network-protocols-support)
- [Client](#client)
  - [Configuration](#configuration)
  - [Waiting for the Result with a Timeout](#waiting-for-the-result-with-a-timeout)
  - [Reconect](#reconect)
- [Server](#server)
  - [Configuration](#configuration-1)
  - [Command Size Restriction](#command-size-restriction)
  - [Close and Shutdown](#close-and-shutdown)
- [How To Use](#how-to-use)
- [Architecture](#architecture)

# Tests
Test coverage of each submodule is over 90%.

# Benchmarks
[github.com/ymz-ncnk/go-client-server-communication-benchmarks](https://github.com/ymz-ncnk/go-client-server-communication-benchmarks)

# High-performance Communication Channel
To build a high-performance communication channel between two services:
1. Use N connections. Multiple connections can transfer significantly more 
   information than a single one. The number N depends on your system and can 
   represents the number of connections after which adding another one will not 
   provide additional benefits.
2. To minimize system latency, use all available connections from the start 
   instead of creating new ones on demand.
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
- KeepaliveTime and KeepaliveIntvl - If both != 0, client will try to keep the 
  connection alive. When there are no commands to send, it starts Ping-Pong with
  the server - sends a Ping command and receives a Pong result, both of which 
  are transfered as a 0 (like a ball) byte.
- SysDataReceiveTimeout - Specifies how long the client will wait to receive 
  system data from the server.

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
`client.NewDefReconnect()` creates a client that attempts to reconnect to the 
server if the connection is lost. This may occur while sending a command -
we'll get an error, or while waiting for the result - we will not be sure 
whether the command was executed on the server or not.

In both cases, after some time, we can try to resend the command (idempotent
command). If the connection is restored, normal operation will continue, 
otherwise we will get the error again.

When deciding on the retry interval, it is advisable to randomize it for each 
goroutine to prevent overloading the server with a large number of 
simultaneously sent commands.

# Server
Before receiving commands, the server sends system data to the client: 
`ServerInfo` and `ServerSettings`. Using `ServerInfo`, the client can 
determine  its compatibility with the server, for instance, whether they both 
support the same set of commands. `ServerSettings`, on the other hand, contains 
the preferred configuration for interacting with the server.

A few words about command execution:
- Each command is executed by a single `Invoker` (it should be thread-safe) in 
  a separete goroutine.
- There is a default `Invoker`, but you can provide your own.
- A command can send back multiple results, all of which will be delivered to 
  the client in order.

## Configuration
Server configuration options include (and not only):
- FirstConnTimeout - Specifies the time limit for receiving the first 
  connection. If the server does not receive a connection within this time, it 
  will close.
- WorkersCount - Each connection is processed by one `Worker`, so this parameter
  determines the maximum number of simultaneous clients that the server can 
  handle.
- LostConnCallback - A callback function triggered when the server loses its 
  connection with a client.
- ReceiveTimeout - Specifies the maximum duration the server will wait for a
  command from a client. If no one command is received within this time, the 
  server will close the connection.

## Command Size Restriction
The server may ask the client not to send too large commands.

To enable this, simply set `Conf.ServerSettings.MaxCmdSize` in bytes and 
implement the client codec's `Size()` method, it will be used to verify the 
command size.

Please note that even with this feature, the server must protect itself from
receiving too large commands. This can be achieved during the command decoding 
process - the server codec's `Decode()` method may return an error, which will 
close the connection to the client.

## Close and Shutdown
`Server.Close()` terminates all connections and immediately stops the server. 
`Server.Shutdown()` allows the server to complete processing of already 
established connections while rejecting new ones.

# How To Use
All you need to do is implement the Command pattern and codecs (one for the 
client and one for the server):
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
  at time.Time, // If Conf.Handler.At == true on the server, this param will 
  // contain the time when the command was received.
  seq base.Seq, // The sequence number of the command. It is used to send back 
  // results.
  receiver Calculator, // Receiver.
  proxy base.Proxy, // Allows command to send back results. Contains only
  // two methods: Send() and SendWithDeadline().
) error {
  // It uses Receiver here.
  result := Result(receiver.Add(...))
  // And sends back the result. 
  // In general, a command can send back multiple results, which will be 
  // received by the client in order.
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
// One ServerCodec will be used by all server Workers, so it must be thread-safe.
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
asyncResults := make(chan base.AsyncResult, 1)
_, err := client.Send(cmd, asyncResults)
...
asyncResult := <-asyncResults
if asyncResult.Error != nil {
  // asyncResult.Error != nil if something is wrong with the connection.
  ...
}
// The result sent by the command.
result := asyncResult.Result.(Result)
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