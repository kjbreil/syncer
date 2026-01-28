package combined

import (
    "context"
    "testing"
    "time"

    "github.com/kjbreil/syncer/pkg/control"
)

// TestAdd verifies that Add injects an entry and mutates the target data.
func TestAdd(t *testing.T) {
    ctx := context.Background()
    data := &simpleStruct{}
    c, err := New(ctx, data)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    entry := control.NewEntry(2, "Bob")
    // Build key path: simpleStruct.Name
    entry.Key = append(entry.Key, &control.Key{Key: "simpleStruct"})
    entry.Key = append(entry.Key, &control.Key{Key: "Name"})

    if err := c.Add(entry); err != nil {
        t.Fatalf("Add() error = %v", err)
    }
    if data.Name != "Bob" {
        t.Fatalf("expected Name to be Bob, got %s", data.Name)
    }
}

// TestEntries ensures Entries produces a diff and signals change.
func TestEntries(t *testing.T) {
    ctx := context.Background()
    base := &simpleStruct{Name: "Alice", Age: 20}
    c, _ := New(ctx, base)

    modified := &simpleStruct{Name: "Alice", Age: 21}
    entries, err := c.Entries(modified)
    if err != nil {
        t.Fatalf("Entries() error = %v", err)
    }
    if len(entries) == 0 {
        t.Fatalf("expected non-empty diff entries")
    }

    select {
    case <-c.extractorChgChan:
        // success
    case <-time.After(50 * time.Millisecond):
        t.Fatalf("expected extractorChgChan to receive a signal")
    }
}

// TestReset verifies that Reset clears extractor state.
func TestReset(t *testing.T) {
    ctx := context.Background()
    st := &simpleStruct{Name: "Z", Age: 42}
    c, _ := New(ctx, st)

    c.Reset()

    zero := &simpleStruct{}
    entries, err := c.Entries(zero)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(entries) != 0 {
        t.Fatalf("expected no diff after reset; got %d entries", len(entries))
    }
}

// TestClose ensures Close cancels context and closes channels.
func TestClose(t *testing.T) {
    ctx := context.Background()
    c, _ := New(ctx, &simpleStruct{})

    if err := c.Close(); err != nil {
        t.Fatalf("Close() error = %v", err)
    }

    select {
    case <-c.ctx.Done():
        // ok
    default:
        t.Fatalf("context not cancelled after Close")
    }

    // Closed channel should panic on send
    defer func() {
        if r := recover(); r == nil {
            t.Fatalf("expected panic when sending to closed channel")
        }
    }()
    c.extractorChgChan <- struct{}{}
}
