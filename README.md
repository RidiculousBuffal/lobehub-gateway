# LobeHub Gateway Go

[简体中文](./README.zh-CN.md)

`lobehub-gateway-go` is a Go reimplementation of the LobeHub gateway services. It is intended for self-hosted deployments that need a simple, single-instance gateway without Cloudflare platform bindings or external infrastructure dependencies.

The original gateway implementations remain the protocol reference. This repository focuses on keeping the public HTTP and WebSocket behavior compatible while making the runtime easier to run as a normal Go service.

## Goals

- Go implementation of LobeHub gateway protocols.
- Single-instance design for straightforward self-hosting.
- No Cloudflare Workers or Durable Objects dependency.
- No Redis, PostgreSQL, NATS, or other external runtime services.
- Keep protocol behavior compatible with the reference implementations where practical.

## Alignment Target

The protocol baseline for this repository is the pair of standalone reference implementations:

- `device-gateway` at `d9df5542f06df6f5555356efa0735f767f0e5ccf`
- `agent-gateway` at `82dc31148e270866356a24c3b4fbefa912479ebb`

For implemented gateway services, the target is compatibility with the public HTTP routes, WebSocket messages, request and response shapes, authentication behavior, routing semantics, and timeout behavior exposed by those reference implementations. This target intentionally keeps the Go runtime self-hosted and single-instance oriented.

The alignment target does not include deliberately excluded platform surfaces: admin APIs, metrics, Durable Object storage or WebSocket hibernation, multi-instance coordination, geo reporting, and admin reporting.

## Projects

| Project | Status | Description |
| --- | --- | --- |
| [`device-gateway-go`](./device-gateway-go/README.md) | Implemented | Device gateway for routing LobeHub device connections, status queries, tool calls, system information, and agent-run requests. |
| [`agent-gateway-go`](./agent-gateway-go/README.md) | Implemented | Agent gateway for relaying browser WebSocket sessions and server-pushed agent events. |

## Unified Binary

The root module (`go.mod`) provides a unified entry point that starts both gateway services in a single process:

```text
lobehub-gateway-go/
  go.mod                          # umbrella module with replace directives
  cmd/gateway/main.go             # unified entry: runs agent + device simultaneously
  agent-gateway-go/
    gateway.go                    # facade re-exporting internal/gateway symbols
    go.mod                        # standalone module (unchanged)
    internal/gateway/
  device-gateway-go/
    gateway.go                    # facade re-exporting internal/gateway symbols
    go.mod                        # standalone module (unchanged)
    internal/gateway/
```

Each sub-project remains an independent Go module and can be built/tested on its own. The root module uses `replace` directives to import the sub-modules via facade files (`gateway.go`) that re-export the internal package symbols.

### Build & Run

```bash
# Build the unified binary
 go build -o gateway ./cmd/gateway

# Run both services
 SERVICE_TOKEN=your-token ./gateway

# Defaults: agent on :8787, device on :8788
# Override ports via env vars
 AGENT_PORT=9000 DEVICE_PORT=9001 SERVICE_TOKEN=your-token ./gateway
```

### Environment Variables

| Variable | Default | Shared | Description |
| --- | --- | --- | --- |
| `SERVICE_TOKEN` | — | Yes | Service auth token (required) |
| `JWKS_PUBLIC_KEY` | — | Yes | JWKS public key for JWT verification |
| `LOBE_API_BASE_URL` | `https://app.lobehub.com` | Agent only | LobeHub API base URL |
| `AGENT_PORT` | `8787` | Agent only | Agent gateway listen port |
| `DEVICE_PORT` | `8788` | Device only | Device gateway listen port |
| `READ_TIMEOUT` | `30s` | Yes | HTTP read timeout |
| `WRITE_TIMEOUT` | `30s` | Yes | HTTP write timeout |
| `SHUTDOWN_TIMEOUT` | `10s` | Yes | Graceful shutdown timeout |

## Architecture

This repository is organized as a collection of gateway services. Each gateway lives in its own directory so more gateway implementations can be added without coupling them to the existing services. The root module provides a unified binary that runs all gateways together.

The current implementation style is intentionally minimal: state is held in memory, services run as ordinary HTTP/WebSocket servers, and deployment can be handled by systemd, Docker, Kubernetes, or any other process manager.

## Relationship To LobeHub Gateway

Compared with the original LobeHub gateway projects, this repository has several deliberate differences:

| Area | Original gateway | This repository |
| --- | --- | --- |
| Language | TypeScript | Go |
| Runtime | Cloudflare Workers | Standard Go process |
| Coordination | Durable Objects | In-memory single instance |
| Platform binding | Cloudflare-specific | Platform-neutral |
| Runtime dependencies | Cloudflare services | No external services required |

Because this repository targets a single gateway instance, it does not provide distributed session coordination by default. If you need horizontal scaling, run one active gateway instance per routing domain or add a shared coordination layer explicitly.

## Development

### Unified binary

```bash
# Build
 go build -o gateway ./cmd/gateway

# Test root module (unified binary)
 go test ./...
```

### Individual gateways

Each gateway remains a standalone module and can be built/tested independently:

```bash
cd agent-gateway-go && go test ./...
cd device-gateway-go && go test ./...
```

See each gateway README for service-specific configuration, API routes, and local run instructions.
