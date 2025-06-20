# Syncer

Syncer is a Go tool for synchronizing struct data between programs over a network using gRPC. It automatically detects changes in structs and propagates them to connected peers in real-time.

## WARNING

Do not use this project directly at the moment, the API is changing and all the needed tests have not been created.

## Features

- Real-time struct synchronization over gRPC
- Automatic change detection using reflection
- Bidirectional client-server communication
- Support for all Go types except interfaces (which require pre-existing instances)
- Field-level control with struct tags
- Built-in demo application with TUI

## Quick Start

### Running the Demo

```bash
# Build and run the interactive demo
make run

# Or manually:
go run cmd/syncer/main.go
```

The demo creates two synchronized endpoints with a terminal UI showing real-time sync status and logs.

### Basic Usage

```go
package main

import (
    "net"
    "github.com/kjbreil/syncer/endpoint"
    "github.com/kjbreil/syncer/endpoint/settings"
)

// Define your data structure
type data struct {
    Name        string
    Count       int
    SyncedField string
    LocalField  string `extractor:"-"` // Won't be synchronized
}

func main() {
    // Initialize your data
    myData := &data{
        Name:        "Example",
        Count:       42,
        SyncedField: "This will sync",
        LocalField:  "This won't sync",
    }

    // Configure the endpoint
    settings := &settings.Settings{
        Port: 45012,
        Peers: []net.TCPAddr{{
            IP:   net.ParseIP("10.0.2.2"),
            Port: 45012,
        }},
        AutoUpdate: true, // Automatically sync changes
    }

    // Create and start the endpoint
    ep, err := endpoint.New(myData, settings)
    if err != nil {
        panic(err)
    }
    
    // Run as server (blocking)
    ep.Run(true)
}
```

## Struct Tags

Use the `extractor:"-"` tag to exclude fields from synchronization:

```go
type MyStruct struct {
    PublicSynced  string              // Will be synchronized
    privatefield  string              // Won't sync (unexported)
    ExcludedField string `extractor:"-"` // Explicitly excluded
}
```

## Development

### Building

```bash
# Build the demo application
make build

# Run tests
make test

# Generate protobuf files
make proto
```

### Protocol Buffers

The project uses Protocol Buffers for network communication. To regenerate the protobuf files:

```bash
make proto
```

This generates both Go and JavaScript/gRPC-Web bindings from `control/proto/control.proto`.