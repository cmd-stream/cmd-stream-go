# cmd-stream-go

[![Go Reference](https://pkg.go.dev/badge/github.com/cmd-stream/cmd-stream-go.svg)](https://pkg.go.dev/github.com/cmd-stream/cmd-stream-go)
[![GoReportCard](https://goreportcard.com/badge/cmd-stream/cmd-stream-go)](https://goreportcard.com/report/github.com/cmd-stream/cmd-stream-go)
[![Follow on X](https://img.shields.io/twitter/url?url=https%3A%2F%2Fx.com%2Fcmdstream_lib)](https://x.com/cmdstream_lib)

**cmd-stream-go** is a high-performance, modular client–server library for Go,
based on the [Command Pattern](https://en.wikipedia.org/wiki/Command_pattern).
It's designed for efficient, low-latency communication over TCP/TLS, with
support for streaming and observability.

In short, the concept is simple: a client sends Commands to the server, where an
Invoker executes them, and a Receiver provides the actual server-side
functionality.

Want to learn more about how the Command Pattern applies to network
communication?  Check out [this series of posts](https://medium.com/p/f9e53442c85d).

## Contents

- [cmd-stream-go](#cmd-stream-go)
  - [Contents](#contents)
  - [Why cmd-stream-go?](#why-cmd-stream-go)
  - [Overview](#overview)
  - [Tests](#tests)
  - [Benchmarks](#benchmarks)
  - [How To](#how-to)
  - [Network Protocols Support](#network-protocols-support)
  - [High-performance Communication Channel](#high-performance-communication-channel)
  - [cmd-stream-go and RPC](#cmd-stream-go-and-rpc)
  - [Architecture](#architecture)

## Why cmd-stream-go?

It delivers high-performance and resource efficiency, helping reduce
infrastructure costs and scale more effectively.

## Overview

- Works over TCP, TLS or mutual TLS.
- Has an asynchronous client that uses only one connection for both sending
  Commands and receiving Results.
- Supports the server streaming, i.e. a Command can send back multiple Results.
- Provides reconnect and keepalive features.
- Supports the Circuit Breaker pattern.
- Has OpenTelemetry integration.
- Can work with various serialization formats.
- Follows a modular design.

## Tests

The main `cmd-stream-go` module contains basic integration tests, while each
submodule (see [Architecture](#architecture)) has approximately 90% code
coverage.

## Benchmarks

![QPS Benchmark](https://github.com/ymz-ncnk/go-client-server-benchmarks/blob/main/results/qps/img/qps.png)

See [go-client-server-benchmarks](https://github.com/ymz-ncnk/go-client-server-benchmarks)
for detailed performance comparisons.

## How To

Getting started is easy:

1. Implement the Command Pattern.
2. Generate the serialization code.

### Quick Look

Here's a minimal end-to-end example showing how Commands can be defined, sent,
and executed over the network:

```go
// Calc represents the Receiver (application layer).
type Calc struct{}
func (c Calc) Add(a, b int) int { return a + b }
func (c Calc) Sub(a, b int) int { return a - b }

// AddCmd is a Command that uses Calc to perform addition.
type AddCmd struct {A, B int}
func (c AddCmd) Exec(ctx context.Context, seq core.Seq, at time.Time,
  calc rcvr.Calc, proxy core.Proxy,
) (err error) {
  result := CalcResult(calc.Add(c.A, c.B))
  _, err = proxy.Send(seq, result) // send result back
  return
}

// SubCmd is a Command that uses Calc to perform subtraction.
type SubCmd struct {A, B int}
func (c SubCmd) Exec(ctx context.Context, seq core.Seq, at time.Time,
  calc rcvr.Calc, proxy core.Proxy,
) (err error) {
  result := CalcResult(calc.Sub(c.A, c.B))
  _, err = proxy.Send(seq, result)  // send result back
  return
}

// CalcResult is the Result returned by Commands.
type CalcResult int
func (r CalcResult) LastOne() bool { return true }

func main() {
  const addr = "127.0.0.1:9000"
  var (
    invoker     = srv.NewInvoker(Calc{})
    serverCodec = ...
    clientCodec = ...
  )
  // Start server.
  go func() {
    server := cmdstream.MakeServer(serverCodec, invoker)
    server.ListenAndServe(addr)
  }()
  // Make sender.
  sender, _ := sndr.Make(addr, clientCodec)

  // Send AddCmd.
  sum, _ := sender.Send(context.Background(), AddCmd{A: 2, B: 3})
  fmt.Println(sum) // Output: 5
  // Send SubCmd.
  diff, _ := sender.Send(context.Background(), SubCmd{A: 8, B: 4})
  fmt.Println(diff) // Output: 4
}
```

The full, runnable example, including codec definitions, is available in the
[calc](https://github.com/cmd-stream/examples-go/tree/main/calc).

For further learning, see the additional resources below.

### Additional Resources

- [Tutorial](https://ymz-ncnk.medium.com/cmd-stream-go-tutorial-0276d39c91e8)
- [Examples](https://github.com/cmd-stream/examples-go)
- [OpenTelemetry Instrumentation](https://ymz-ncnk.medium.com/cmd-stream-go-with-opentelemetry-adeecfbe7987)

## Network Protocols Support

Built on Go’s standard `net` package, `cmd-stream-go` supports
connection-oriented protocols, such as TCP, TLS, and mutual TLS (for client
authentication).

## High-performance Communication Channel

To maximize performance between services:

1. Use N parallel connections. More connections typically improve throughput,
   until a saturation point.
2. Pre-establish all connections instead of opening them on-demand.
3. Keep connections alive to avoid the overhead from reconnections.

These practices, implemented via the [ClientGroup](group/group.go), can
significantly enhance throughput and reduce latency between your services.

## cmd-stream-go and RPC

Already using RPC? You can use `cmd-stream-go` as a faster transport layer. See
the [RPC example](https://github.com/cmd-stream/examples-go/tree/main/rpc).

## Architecture

`cmd-stream-go` is split into the following submodules:

- [core-go](https://github.com/cmd-stream/core-go): The core module that includes
  client and server definitions.
- [delegate-go](https://github.com/cmd-stream/delegate-go): The client delegates
  all communication-related tasks to its delegate, the server follows the same
  approach. The connection is also initialized at this level.
- [handler-go](https://github.com/cmd-stream/handler-go): The server delegate
  uses a handler to receive and process Commands.
- [transport-go](https://github.com/cmd-stream/transport-go): Responsible for
  Commands/Results delivery.

`cmd-stream-go` was designed in such a way that you can easily replace any part
of it.

