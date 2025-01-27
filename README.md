# cmd-stream-go

[![Go Reference](https://pkg.go.dev/badge/github.com/cmd-stream/cmd-stream-go.svg)](https://pkg.go.dev/github.com/cmd-stream/cmd-stream-go)
[![GoReportCard](https://goreportcard.com/badge/cmd-stream/cmd-stream-go)](https://goreportcard.com/report/github.com/cmd-stream/cmd-stream-go)

cmd-stream-go allows execution of Commands on the server using the 
[Command pattern](https://en.wikipedia.org/wiki/Command_pattern).

# Why cmd-stream-go?
It provides an extremely fast and flexible communication mechanism.

# Command Pattern as an API Architecture Style
[ymz-ncnk.medium.com/command-pattern-as-an-api-architecture-style-be9ac25d6d94](https://ymz-ncnk.medium.com/command-pattern-as-an-api-architecture-style-be9ac25d6d94)

# Brief cmd-stream-go Description
- Can work over TCP, TLS or mutual TLS.
- Has an asynchronous client that uses only one connection for both sending 
  Commands and receiving Results.
- Supports the server streaming, i.e. a Command can send back multiple Results
  (client streaming is not directly supported, but can also be implemented).
- Supports deadlines for sending Commands and Results.
- Supports timeouts.
- Supports reconnect feature.
- Supports keepalive feature.
- Can work with various serialization formats ([here](https://github.com/cmd-stream/cmd-stream-examples-go/tree/main/standard_protobuf) is an example using the Protobuf serializer).
- Has a modular architecture.

# Contents
- [cmd-stream-go](#cmd-stream-go)
- [Why cmd-stream-go?](#why-cmd-stream-go)
- [Command Pattern as an API Architecture Style](#command-pattern-as-an-api-architecture-style)
- [Brief cmd-stream-go Description](#brief-cmd-stream-go-description)
- [Contents](#contents)
- [Tests](#tests)
- [Benchmarks](#benchmarks)
- [High-performance Communication Channel](#high-performance-communication-channel)
- [cmd-stream-go and RPC](#cmd-stream-go-and-rpc)
- [Network Protocols Support](#network-protocols-support)
- [Client](#client)
  - [Reconect](#reconect)
- [Server](#server)
  - [Command Size Restriction](#command-size-restriction)
  - [Close and Shutdown](#close-and-shutdown)
- [How To Use](#how-to-use)
- [Architecture](#architecture)

# Tests
The cmd-stream-go module includes only a few integration tests, while each 
submodule (see the [Architecture](#architecture) section) has approximately 90% 
test coverage.

# Benchmarks
[github.com/ymz-ncnk/go-client-server-communication-benchmarks](https://github.com/ymz-ncnk/go-client-server-communication-benchmarks)

# High-performance Communication Channel
To build a high-performance communication channel between two services:
1. Use N connections. Several connections can transfer significantly more
   data than a single one. The number N, depends on your system and represents 
   the point after which adding more connections will not provide additional 
   benefits.
2. To minimize system latency, use all available connections from the start 
   instead of creating new ones on demand.
3. Use keepalive connections.

# cmd-stream-go and RPC
If you are already using RPC, cmd-stream-go can boost its performance by 
providing a faster communication tool. [Here's](https://github.com/cmd-stream/cmd-stream-examples-go/tree/main/rpc) 
an example.

# Network Protocols Support
cmd-stream-go is built on top of the standard Golang net package, and supports 
connection-oriented protocols like TCP, TLS or mutual TLS (for client
authentication).

# Client
The client is asynchronous and can be used from different goroutines 
simultaneously. It uses only one connection to send Commands and receive 
Results. Commands sent from a single goroutine will be delivered to the server 
in order.

## Reconect
The client may lose connection to the server in the following cases:
- While sending a Command – Client.Send() will return an error.
- While waiting for a Result – whether the Command was executed on the server 
  remains unknown.

In both cases, if the `client.NewReconnect()` constructor is used, we can try to
resend the Command (idempotent Command) after a while. If the client has already
successfully reconnected, normal operation will continue, otherwise 
`Client.Send()` will return an error again.

An example of using the "reconnect" client can be found [here](https://github.com/cmd-stream/cmd-stream-examples-go/tree/main/reconnect).

# Server
The server initializes the connection to the client by sending system data: 
`ServerInfo` and `ServerSettings`. Using `ServerInfo`, the client can 
determine its compatibility with the server, for instance, whether they both 
support the same set of Commands. `ServerSettings` specifies the configuration 
parameters for interacting with the server.

A few words about Command execution:
- Each Command is executed by a single `Invoker` (it should be thread-safe) in 
  a separete goroutine.
- A Command can send multiple Results, all of which will be delivered to 
  the client in order. See [this example](https://github.com/cmd-stream/cmd-stream-examples-go/tree/main/multi_result).

## Command Size Restriction
The server may ask the client not to send too large Commands - simply set 
`Conf.ServerSettings.MaxCmdSize` in bytes and implement the client codec's 
`Size()` method, it will be used to verify the Command size.

Even with this feature, the server must protect itself against excessively large
Commands. This can be handled in the `Codec.Decode()` method on the server. See 
[this example](https://github.com/cmd-stream/cmd-stream-examples-go/tree/main/max_cmd_size).

## Close and Shutdown
`Server.Close()` terminates all connections and immediately stops the server. 
`Server.Shutdown()` allows the server to complete processing of already 
established connections without accepting new ones.

# How To Use
[Here](https://ymz-ncnk.medium.com/cmd-stream-go-tutorial-0276d39c91e8) is a 
detailed tutorial and [here](https://github.com/cmd-stream/cmd-stream-examples-go/tree/main/echo) 
is a short example.

# Architecture
There are the following cmd-stream-go submodules:
- [base-go](https://github.com/cmd-stream/base-go): Basic module, that contains 
  the client and server definitions.
- [delegate-go](https://github.com/cmd-stream/delegate-go): The client delegates
  all communication-related tasks to its delegate, the server follows the same 
  approach. The connection is also initialized at this level.
- [handler-go](https://github.com/cmd-stream/handler-go): The server delegate 
  uses a handler to receive and execute Commands.
- [transport-go](https://github.com/cmd-stream/transport-go): Resposible for 
  Commands/Results delivery.

cmd-stream-go was designed in such a way that you can easily replace any part of 
it.