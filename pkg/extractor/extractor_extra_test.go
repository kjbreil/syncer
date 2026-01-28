package extractor

import (
    "testing"
)

type sample struct {
    Name string
    Age  int
}

// TestReset ensures Reset zeroes the internal snapshot so that subsequent Entries against a zero struct shows no diff.
func TestReset(t *testing.T) {
    original := &sample{Name: "Alice", Age: 30}
    ext, err := New(original)
    if err != nil {
        t.Fatalf("New() error = %v", err)
    }

    // First diff against non-zero struct should produce entries.
    if diff, err := ext.Entries(original); err != nil || len(diff) == 0 {
        t.Fatalf("expected diff entries, got %v err %v", diff, err)
    }

    // Reset internal snapshot.
    ext.Reset()

    // Now diff against zero struct should be nil (because internal snapshot also zero).
    zero := &sample{}
    if diff, err := ext.Entries(zero); err != nil || len(diff) != 0 {
        t.Fatalf("after Reset expected no diff, got %v err %v", diff, err)
    }
}

// TestEntriesNilData ensures Entries returns error on nil input.
func TestEntriesNilData(t *testing.T) {
    ext, _ := New(&sample{})
    if _, err := ext.Entries(nil); err == nil {
        t.Fatalf("expected error when data is nil")
    }
}
