// Package delegate provides standard client and server delegate implementations
// for the cmd-stream library, bridging core protocol logic with the transport
// layer.
//
// These delegates manage the connection lifecycle, including the initial
// handshake, protocol-level heartbeat (keepalive), and automated reconnection.
//
// # Connection Handshake & ServerInfo
//
// During initialization, the server typically transmits a ServerInfo message.
// This allows the client to verify version compatibility and supported command
// sets. To ensure security and prevent memory exhaustion attacks, a default 1KB
// limit is strictly enforced for ServerInfo payloads.
//
// # Helper Delegates
//
// Specialized client-side delegates are provided in the cln subpackage:
//
//   - KeepaliveDelegate: Maintains connection health when idle. It initiates a
//     heartbeat exchange by sending a Ping Command and awaiting a Pong Result -
//     both transmitted as a single zero-byte payload to minimize bandwidth
//     overhead.
//
//   - ReconnectDelegate: Implements resilient connectivity by providing a
//     standardized Reconnect method to restore operations following a transport
//     failure.
package delegate
