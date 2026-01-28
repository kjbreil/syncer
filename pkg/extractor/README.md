# Extractor

The extractor package monitors Go structs for changes and produces `control.Entry` objects representing the differences.

## How It Works

When initialized with a struct, the extractor stores a deep copy as a baseline. Each call to `Entries()` compares the current struct state against the baseline, returns entries for any changed fields, and updates the baseline. The first call after initialization returns entries for the full struct state (since the baseline starts as a zero-value copy).

## Struct Tags

Use the `extractor:"-"` tag to exclude fields from change detection:

```go
type MyData struct {
    Name       string              // Tracked for changes
    Count      int                 // Tracked for changes
    LocalOnly  string `extractor:"-"` // Excluded from sync
    unexported string              // Ignored (unexported fields are not tracked)
}
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/kjbreil/syncer/pkg/extractor"
)

type data struct {
    String string
}

func main() {
    t := data{
        String: "test",
    }
    ext, err := extractor.New(&t)
    if err != nil {
        panic(err)
    }

    // First call returns the full current state (String = "test")
    entries, err := ext.Entries(&t)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Initial entries: %d\n", len(entries))

    // Modify the struct
    t.String = "new test"

    // Second call returns only the change (String: "test" -> "new test")
    entries, err = ext.Entries(&t)
    if err != nil {
        panic(err)
    }
    for _, entry := range entries {
        fmt.Printf("Changed: %s\n", entry.String())
    }

    // Reset clears the baseline, so the next Entries() call
    // will again return the full struct state
    ext.Reset()
}
```
