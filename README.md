# cmd-stream-go

[![Go Reference](https://pkg.go.dev/badge/github.com/cmd-stream/cmd-stream-go.svg)](https://pkg.go.dev/github.com/cmd-stream/cmd-stream-go)
[![GoReportCard](https://goreportcard.com/badge/cmd-stream/cmd-stream-go)](https://goreportcard.com/report/github.com/cmd-stream/cmd-stream-go)
[![Follow on X](https://img.shields.io/twitter/url?url=https%3A%2F%2Fx.com%2Fcmdstream_lib)](https://x.com/cmdstream_lib)

**cmd-stream-go** is a high-performance, modular client–server library for Go, 
based on the [Command Pattern](https://en.wikipedia.org/wiki/Command_pattern). 
It's designed for efficient, low-latency communication over TCP/TLS, with 
support for streaming and observability.

Want to learn how the Command Pattern applies to network communication? Check 
out [this series of posts](https://medium.com/p/f9e53442c85d).

# Why cmd-stream-go?
It delivers high-performance and resource efficiency, helping reduce 
infrastructure costs and scale more effectively.

# Overview
- Works over TCP, TLS or mutual TLS.
- Has an asynchronous client that uses only one connection for both sending 
  Commands and receiving Results.
- Supports the server streaming, i.e. a Command can send back multiple Results.
- Provides reconnect and keepalive features.
- Supports the Circuit Breaker pattern.
- Has OpenTelemetry integration.
- Can work with various serialization formats.
- Follows a modular design.

# Contents
- [cmd-stream-go](#cmd-stream-go)
- [Why cmd-stream-go?](#why-cmd-stream-go)
- [Overview](#overview)
- [Contents](#contents)
- [Tests](#tests)
- [Benchmarks](#benchmarks)
- [How To](#how-to)
- [Network Protocols Support](#network-protocols-support)
- [Client \& Server](#client--server)
- [High-performance Communication Channel](#high-performance-communication-channel)
- [cmd-stream-go and RPC](#cmd-stream-go-and-rpc)
- [Architecture](#architecture)

# Tests
The main `cmd-stream-go` module contains basic integration tests, while each 
submodule (see [Architecture](#architecture)) has ~90% code coverage.

# Benchmarks
[github.com/ymz-ncnk/go-client-server-benchmarks](https://github.com/ymz-ncnk/go-client-server-benchmarks)

# How To
Getting started is easy:
1. Implement the Command Pattern.
2. Generate the serialization code. 

For more details, explore the following resources:
- [Tutorial](https://ymz-ncnk.medium.com/cmd-stream-go-tutorial-0276d39c91e8)
- [OpenTelemetry Instrumentation](https://ymz-ncnk.medium.com/cmd-stream-go-with-opentelemetry-adeecfbe7987)
- [Examples](https://github.com/cmd-stream/examples-go)

# Network Protocols Support
Built on Go’s standard `net` package, `cmd-stream-go` supports 
connection-oriented protocols, such as TCP, TLS, and mutual TLS (for client 
authentication).

# Client & Server
The client operates asynchronously, sending Commands to the server. On the 
server side, the Invoker executes the Commands, while the Receiver provides the 
underlying server functionality.

# High-performance Communication Channel
To maximize performance between services:
1. Use N parallel connections. More connections typically improve throughput, 
   until a saturation point.
2. Pre-establish all connections instead of opening them on-demand.
3. Keep connections alive to avoid the overhead from reconnections.

These practices, implemented via the [ClientGroup](group/group.go), can 
significantly enhance throughput and reduce latency between your services.

# cmd-stream-go and RPC
Already using RPC? You can use `cmd-stream-go` as a faster transport layer. See 
the [RPC example](https://github.com/cmd-stream/examples-go/tree/main/rpc).

# Architecture
`cmd-stream-go` is split into the following submodules:
- [core-go](https://github.com/cmd-stream/core-go): The core module that includes 
  client and server definitions.
- [delegate-go](https://github.com/cmd-stream/delegate-go): The client delegates
  all communication-related tasks to its delegate, the server follows the same 
  approach. The connection is also initialized at this level.
- [handler-go](https://github.com/cmd-stream/handler-go): The server delegate 
  uses a handler to receive and process Commands.
- [transport-go](https://github.com/cmd-stream/transport-go): Resposible for 
  Commands/Results delivery.

cmd-stream-go was designed in such a way that you can easily replace any part of 
it.