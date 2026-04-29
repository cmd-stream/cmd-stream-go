# cmd-stream-go

[![Go Reference](https://pkg.go.dev/badge/github.com/cmd-stream/cmd-stream-go.svg)](https://pkg.go.dev/github.com/cmd-stream/cmd-stream-go)
[![GoReportCard](https://goreportcard.com/badge/cmd-stream/cmd-stream-go)](https://goreportcard.com/report/github.com/cmd-stream/cmd-stream-go)
[![codecov](https://codecov.io/gh/cmd-stream/cmd-stream-go/graph/badge.svg?token=RXPJ6ZIPK7)](https://codecov.io/gh/cmd-stream/cmd-stream-go)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/12510/badge)](https://www.bestpractices.dev/projects/12510)
[![Follow on X](https://img.shields.io/twitter/url?url=https%3A%2F%2Fx.com%2Fcmdstream_lib)](https://x.com/cmdstream_lib)


**cmd-stream** is a high-performance networking library that implements the 
Distributed [Command Pattern](https://en.wikipedia.org/wiki/Command_pattern) (DCP) 
for Go. Designed for low-latency communication over TCP/TLS, it provides a 
flexible, decoupled alternative to traditional RPC by treating requests as 
first-class Command objects.

The architecture is straightforward: a client sends Commands to the server, 
where an Invoker executes them, and a Receiver provides the actual server-side
functionality.

*Want to learn more about how the Command Pattern applies to network
communication?  Check out [this series of posts](https://medium.com/p/f9e53442c85d)*.

## Contents

- [cmd-stream-go](#cmd-stream-go)
  - [Contents](#contents)
  - [Why cmd-stream?](#why-cmd-stream)
  - [Overview](#overview)
  - [Benchmarks](#benchmarks)
  - [Installation](#installation)
  - [How To](#how-to)
    - [Quick Look](#quick-look)
    - [Additional Resources](#additional-resources)
  - [Network Protocols Support](#network-protocols-support)
  - [High-performance Communication Channel](#high-performance-communication-channel)
  - [cmd-stream and RPC](#cmd-stream-and-rpc)
  - [Architecture](#architecture)
  - [Contributing \& Security](#contributing--security)
  - [Version Compatibility](#version-compatibility)

## Why cmd-stream?

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

## Benchmarks

![QPS Benchmark](https://github.com/ymz-ncnk/go-client-server-benchmarks/blob/main/results/qps/img/qps.png)

See [go-client-server-benchmarks](https://github.com/ymz-ncnk/go-client-server-benchmarks)
for detailed performance comparisons.

## Installation

To obtain the library, use: 

```bash
go get github.com/cmd-stream/cmd-stream-go
```

## How To

Getting started is easy:

1. Implement the Command Pattern.
2. Use one of the codecs:
   - [codec-json-go](https://github.com/cmd-stream/codec-json-go) - simple and
   easy-to-use JSON codec (ideal for prototyping).
   - [codec-protobuf-go](https://github.com/cmd-stream/codec-protobuf-go) -
   Protobuf-based codec (requires code generation).
   - [codec-mus-stream-go](https://github.com/cmd-stream/codec-mus-stream-go) -
   high-performance MUS codec (requires code generation).

**Tip:** Start with JSON for simplicity, and switch to MUS later for maximum
performance.

### Quick Look

Here's a minimal end-to-end example showing how Commands can be defined, sent,
and executed over the network:

```go
// Calc handles arithmetic logic.
type Calc struct{}

func (c Calc) Add(a, b int) int { return a + b }
func (c Calc) Sub(a, b int) int { return a - b }

// AddCmd executes addition via Calc.
type AddCmd struct{ A, B int }

func (c AddCmd) Exec(ctx context.Context, seq core.Seq, _ time.Time, calc Calc,
  proxy core.Proxy) error {
  _, err := proxy.Send(seq, CalcResult(calc.Add(c.A, c.B)))
  return err
}

// SubCmd executes subtraction via Calc.
type SubCmd struct{ A, B int }

func (c SubCmd) Exec(ctx context.Context, seq core.Seq, _ time.Time, calc Calc,
  proxy core.Proxy) error {
  _, err := proxy.Send(seq, CalcResult(calc.Sub(c.A, c.B)))
  return err
}

// CalcResult represents the Command output.
type CalcResult int

func (r CalcResult) LastOne() bool { return true }

func main() {
  const addr = "127.0.0.1:9000"

  // 1. Setup codecs.
  reg := cdcjson.NewRegistry(
    cdcjson.WithCmd[Calc, AddCmd](),
    cdcjson.WithCmd[Calc, SubCmd](),
    cdcjson.WithResult[Calc, CalcResult](),
  )
  serverCodec := cdcjson.NewServerCodecWith(reg)
  clientCodec := cdcjson.NewClientCodecWith(reg)

  // 2. Start server.
  server, _ := cmdstream.NewServer(Calc{}, serverCodec)
  go server.ListenAndServe(addr)
  time.Sleep(100 * time.Millisecond)

  // 3. Create sender.
  sender, _ := cmdstream.NewSender(addr, clientCodec)

  // 4. Send commands.
  sum, _ := sender.Send(context.Background(), AddCmd{A: 2, B: 3})
  fmt.Println(sum) // Output: 5

  diff, _ := sender.Send(context.Background(), SubCmd{A: 8, B: 4})
  fmt.Println(diff) // Output: 4
}
```

The full, runnable example is available in the [calc_json](https://github.com/cmd-stream/examples-go/tree/main/calc_json).

### Additional Resources

- [Examples](https://github.com/cmd-stream/examples-go)
- [OpenTelemetry Instrumentation](https://github.com/cmd-stream/otelcmd-stream-go)

## Network Protocols Support

Built on Go’s standard `net` package, `cmd-stream` supports
connection-oriented protocols, such as TCP, TLS, and mutual TLS (for client
authentication).

## High-performance Communication Channel

To maximize performance between services:

1. Use N parallel connections. More connections typically improve throughput,
   until a saturation point.
2. Pre-establish all connections instead of opening them on-demand.
3. Keep connections alive to avoid the overhead from reconnections.

These practices, implemented via the [Group](group/), can significantly
enhance throughput and reduce latency between your services.

## cmd-stream and RPC

Already using RPC? You can use `cmd-stream` as a faster transport layer. See
the [RPC example](https://github.com/cmd-stream/examples-go/tree/main/rpc).

## Architecture

`cmd-stream` is built on a layered architecture that ensures clear separation of 
concerns while maintaining maximum performance:

- [core](core/): The core client and server definitions.
- [delegate](delegate/): All communication-related tasks and connection
  initialization.
- [handler](handler/): Server-side Command processing.
- [transport](transport/): Delivery of Commands and Results over the network.
- [sender](sender/): High-level sender implementation.
- [testkit](testkit/): Data and foundations for integration tests.

`cmd-stream` was designed in such a way that you can easily replace any part
of it.

## Contributing & Security

We welcome contributions of all kinds! Please see [CONTRIBUTING.md](CONTRIBUTING.md) 
for details on how to get involved.

If you find a security vulnerability, please refer to 
[Security Policy](SECURITY.md) for instructions on how to report it privately.

For bugs, feedback, or feature requests, please open an issue!

## Version Compatibility

For a complete list of compatible module versions, see [VERSIONS.md](VERSIONS.md).