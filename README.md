# Syncer

Syncer is a Go tool for synchronizing struct data between programs over a network using gRPC. It automatically detects changes in structs and propagates them to connected peers in real-time.

## WARNING

Do not use this project directly at the moment, the API is changing and all the needed tests have not been created. Parts of this project, like DeepCopy and Equal can be considered finished or at least API stable.

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

## Core Packages

The Syncer project provides several standalone utility packages that can be used independently:

### pkg/deepcopy

A high-performance deep copy library for Go that uses reflection to create independent copies of any Go value.

**Key Features:**
- Handles all Go types (structs, slices, maps, pointers, interfaces, arrays)  
- Optimized performance with primitive type detection
- Uses `reflect.Copy` for primitive arrays/slices
- Thread-safe operations
- Generic `Any[T]` function for type-safe copying

**Usage:**
```go
package main

import (
    "reflect"
    "github.com/kjbreil/syncer/pkg/deepcopy"
)

// Using the generic Any function (recommended)
original := MyStruct{Name: "test", Data: []int{1, 2, 3}}
copied := deepcopy.Any(original)

// Using with reflect.Value directly
originalValue := reflect.ValueOf(original)
copiedValue := deepcopy.DeepCopy(originalValue)
copied := copiedValue.Interface().(MyStruct)
```

### pkg/equal

A specialized equality comparison library that provides more flexible comparison than Go's `reflect.DeepEqual`.

**Key Features:**
- Cross-type comparison for numeric types (int8 vs int32, float32 vs float64)
- Pointer dereferencing for value comparison
- Interface value comparison
- Handles all Go collection types (arrays, slices, maps)
- Function signature comparison

**Comparison with reflect.DeepEqual:**
- `reflect.DeepEqual`: Strict type matching, `int32(5) != int64(5)`
- `pkg/equal`: Flexible numeric comparison, `int32(5) == int64(5)`
- `reflect.DeepEqual`: Pointer address comparison  
- `pkg/equal`: Dereferences pointers to compare values

**Usage:**
```go
package main

import (
    "reflect"
    "github.com/kjbreil/syncer/pkg/equal"
)

// Using the generic Any function (recommended)
a := int32(42)
b := int64(42)
isEqual := equal.Any(a, b) // true (unlike reflect.DeepEqual)

// Using with reflect.Value directly
valueA := reflect.ValueOf(&MyStruct{Name: "test"})
valueB := reflect.ValueOf(&MyStruct{Name: "test"})
isEqual := equal.Equal(valueA, valueB) // true (compares values, not addresses)
```

## Manual Usage of Core Components

### pkg/extractor

The extractor monitors structs for changes and generates `control.Entry` objects representing the differences.

**Manual Usage:**
```go
package main

import (
    "github.com/kjbreil/syncer/pkg/extractor"
)

type MyData struct {
    Name        string
    Count       int
    LocalField  string `extractor:"-"` // Excluded from sync
}

func main() {
    // Create extractor with initial data
    data := &MyData{Name: "initial", Count: 0}
    ext, err := extractor.New(data)
    if err != nil {
        panic(err)
    }

    // Modify your data
    data.Name = "modified"
    data.Count = 42

    // Extract changes
    entries, err := ext.Entries(data)
    if err != nil {
        panic(err)
    }

    // entries now contains the differences
    for _, entry := range entries {
        // Process each change entry
        fmt.Printf("Changed: %s\n", entry.String())
    }

    // Reset to clean state
    ext.Reset()
}
```

### pkg/injector

The injector applies `control.Entry` changes to target structs.

**Manual Usage:**
```go
package main

import (
    "github.com/kjbreil/syncer/pkg/injector"
    "github.com/kjbreil/syncer/pkg/control"
)

func main() {
    // Target data structure (must be pointer)
    targetData := &MyData{Name: "original", Count: 0}
    
    // Create injector
    inj, err := injector.New(targetData)
    if err != nil {
        panic(err)
    }

    // Apply single entry
    entry := &control.Entry{/* entry data */}
    err = inj.Add(entry)
    if err != nil {
        panic(err)
    }

    // Apply multiple entries at once
    entries := control.Entries{/* multiple entries */}
    err = inj.AddAll(entries)
    if err != nil {
        panic(err)
    }

    // targetData is now updated with the injected changes
}
```

### pkg/combined

The combined package provides a higher-level interface that manages both extraction and injection with change debouncing.

**Usage:**
```go
package main

import (
    "context"
    "time"
    "github.com/kjbreil/syncer/pkg/combined"
)

func main() {
    ctx := context.Background()
    data := &MyData{Name: "test", Count: 0}
    
    // Create combined extractor/injector
    combo, err := combined.New(ctx, data)
    if err != nil {
        panic(err)
    }
    defer combo.Close()

    // Set custom debounce duration (default: 2 seconds)
    combo.Debounce = time.Millisecond * 500

    // Set callback for when extractor detects changes
    combo.ExtractorChanges(func() error {
        // Called when data changes are detected
        fmt.Println("Data changed!")
        return nil
    })

    // Set callback for when injector applies changes  
    combo.InjectorChanges(func() error {
        // Called when changes are applied
        fmt.Println("Changes applied!")
        return nil
    })

    // Extract changes from modified data
    data.Name = "modified"
    entries, err := combo.Entries(data)
    if err != nil {
        panic(err)
    }

    // Apply changes from external source
    err = combo.Add(entry)
    if err != nil {
        panic(err)
    }

    // Reset to clean state
    combo.Reset()
}