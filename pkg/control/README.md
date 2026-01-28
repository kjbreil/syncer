# Control

The control package defines the gRPC protocol and data types used for network synchronization in Syncer.

## Overview

This package contains:

- **Protocol Buffer definitions** (`proto/control.proto`) for the data structures transmitted between peers
- **Generated Go and gRPC code** for client/server communication
- **`control.Entry`** — the fundamental unit of change, representing a single field update extracted from a struct
- **`control.Entries`** — a collection of `Entry` objects representing a set of changes

The extractor package produces `control.Entries` from struct diffs, and the injector package consumes them to apply changes to a target struct.

## gRPC Services

The `Control` service defines four RPC methods:

| Method | Type | Description |
|--------|------|-------------|
| `Pull` | Server streaming | Client requests data; server streams `Entry` objects back |
| `Push` | Client streaming | Client streams `Entry` objects to the server |
| `PushPull` | Bidirectional streaming | Both sides can send and receive `Entry` objects simultaneously |
| `Control` | Unary | Client sends a `Message` and receives a `Response` |

## Architecture

The design is **client-centric**: clients initiate all connections and services. The server acts as a passive endpoint that responds to client requests. The `PushPull` service is particularly useful because it allows the server to send data back to the client when changes are detected, without the client needing to poll for updates.

## Regenerating Protobuf Code

```bash
# Using Make
make proto

# Or manually — Go protobuf and gRPC code
protoc -I=pkg/control/proto --go_out=. --go-grpc_out=. pkg/control/proto/*.proto

# JavaScript/gRPC-Web code
protoc -I=pkg/control/proto \
    --js_out=import_style=commonjs,binary:pkg/control/web \
    --grpc-web_out=import_style=commonjs,mode=grpcwebtext:pkg/control/web \
    pkg/control/proto/*.proto
```
