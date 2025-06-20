# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Syncer is a Go-based tool for synchronizing struct data between programs over a network using gRPC. The project implements a client-server architecture where structs are automatically synchronized across network boundaries.

## Development Commands

### Building and Running
```bash
# Build the main demo application
go build -o syncer cmd/syncer/main.go

# Run the demo application
./syncer

# Run all tests
go test ./...

# Run tests for a specific package
go test ./extractor
go test ./injector
go test ./control
```

### Protocol Buffer Generation
```bash
# Generate Go protobuf and gRPC code
protoc -I=control/proto --go_out=. --go-grpc_out=. control/proto/*.proto

# Generate JavaScript/gRPC-Web code
protoc -I=control/proto --go_out=. --js_out=import_style=commonjs,binary:control/web --grpc-web_out=import_style=commonjs,mode=grpcwebtext:control/web control/proto/*.proto
```

## Architecture

### Core Components

1. **Extractor** (`extractor/`): Analyzes structs and generates change entries by comparing current state with previous snapshots. Uses reflection to traverse struct fields and create `control.Entry` objects representing changes.

2. **Injector** (`injector/`): Applies changes from `control.Entry` objects to target structs. Handles type conversion and field updates across different Go types.

3. **Control** (`control/`): Contains the gRPC protocol definition and generated code. Implements streaming services for Push, Pull, and bidirectional Push/Pull operations.

4. **Endpoint** (`endpoint/`): High-level API that combines extractor, injector, and control components into a single interface. Handles client/server roles and network communication.

### Key Files

- `control/proto/control.proto`: Protocol buffer definitions for network communication
- `endpoint/endpoint.go`: Main endpoint implementation with client/server logic
- `cmd/syncer/main.go`: Demo application showing dual-endpoint synchronization with TUI
- `combined/combined.go`: Utilities for combining multiple structs for synchronization

### Data Flow

1. Extractor monitors a struct for changes and generates `control.Entry` objects
2. Entries are transmitted via gRPC streams (Push/Pull/PushPull services)
3. Receiving endpoint uses Injector to apply changes to its local struct copy
4. Bidirectional synchronization keeps both endpoints in sync

### Struct Tagging

Use `extractor:"-"` struct tags to exclude fields from synchronization:

```go
type MyStruct struct {
    SyncedField   string
    LocalField    string `extractor:"-"`  // Won't be synchronized
    privateField  string                  // Also won't be synchronized (unexported)
}
```

### Testing Patterns

- Tests are organized by component (extractor, injector, control, helpers)
- Use `helpers/test/` for shared test utilities and common test structures
- Deep equality testing is handled by `helpers/equal/`
- Deep copying for test setup is in `helpers/deepcopy/`

### Network Architecture

- Client-centric design: clients initiate all connections and services
- Server acts as a passive endpoint responding to client requests
- Push/Pull service enables server-initiated updates when changes are detected
- Auto-discovery and peer management through `endpoint/settings`