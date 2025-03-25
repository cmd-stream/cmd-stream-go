# cmd-stream-go

[![Go Reference](https://pkg.go.dev/badge/github.com/cmd-stream/cmd-stream-go.svg)](https://pkg.go.dev/github.com/cmd-stream/cmd-stream-go)
[![GoReportCard](https://goreportcard.com/badge/cmd-stream/cmd-stream-go)](https://goreportcard.com/report/github.com/cmd-stream/cmd-stream-go)

cmd-stream-go allows execution of Commands on the server using the 
[Command pattern](https://en.wikipedia.org/wiki/Command_pattern).

# Why cmd-stream-go?
It provides an extremely fast and flexible communication mechanism.

# Command Pattern as an API Architecture Style
[Article Link](https://ymz-ncnk.medium.com/command-pattern-as-an-api-architecture-style-be9ac25d6d94)

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
- [Server](#server)
- [How To Use](#how-to-use)
- [Architecture](#architecture)

# Tests
The cmd-stream-go module includes only a few integration tests, while each 
submodule (see the [Architecture](#architecture) section) has approximately 90% 
test coverage.

# Benchmarks
[github.com/ymz-ncnk/go-client-server-communication-benchmarks](https://github.com/ymz-ncnk/go-client-server-communication-benchmarks)

# High-performance Communication Channel
To build a high-performance communication channel between two services, consider 
the following guidelines:
1. Use N connections, as several connections can transfer significantly more 
   data than a single one. The optimal number, N, depends on your system and 
   represents the point after which adding more connections does not improve 
   performance.
2. To minimize latency, open all available connections at the start rather than 
   creating new ones on demand.
3. Keep connections alive to avoid the overhead of frequent connection setup and 
   teardown.

Following these practices can significantly enhance throughput and reduce 
latency between your services.   

# cmd-stream-go and RPC
If you're already using RPC, cmd-stream-go can enhance performance by providing 
a faster communication toolб [here's](https://github.com/cmd-stream/cmd-stream-examples-go/tree/main/rpc) 
an example.

# Network Protocols Support
cmd-stream-go is built on top of Go's standard net package and supports 
connection-oriented protocols such as TCP, TLS, and mutual TLS (for client 
authentication).

# Client
The client is asynchronous and can be safely used from multiple goroutines 
concurrently. It uses a single connection to send Commands and receive Results.
Commands sent from a single goroutine are delivered to the server in order.

# Server
The server initiates the connection by sending ServerInfo to the client. The 
client uses this information to verify compatibility, such as ensuring both 
endpoints support the same set of Commands.

Each Command is executed by a single `Invoker` (it should be thread-safe) in a 
separete goroutine. Also a Command can send multiple Results back, all of which 
will be delivered to the client in order, [here's](https://github.com/cmd-stream/cmd-stream-examples-go/tree/main/multi_result) 
an example.

# How To Use
- [Tutorial](https://ymz-ncnk.medium.com/cmd-stream-go-tutorial-0276d39c91e8)
- [Examples](https://github.com/cmd-stream/cmd-stream-examples-go)

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