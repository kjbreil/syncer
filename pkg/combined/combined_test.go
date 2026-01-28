package combined

import (
    "context"
    "sync/atomic"
    "testing"
    "time"
)

// simpleStruct is a trivial struct used for testing.
type simpleStruct struct {
    Name string
    Age  int
}

// TestNewSuccess ensures that New returns a valid *Combined when provided
// with a non-nil context and pointer to data.
func TestNewSuccess(t *testing.T) {
    ctx := context.Background()
    data := &simpleStruct{}

    c, err := New(ctx, data)
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if c == nil {
        t.Fatalf("expected combined instance, got nil")
    }
}

// TestNewNilContext expects an error when a nil context is supplied.
func TestNewNilContext(t *testing.T) {
    data := &simpleStruct{}
    _, err := New(nil, data)
    if err == nil {
        t.Fatalf("expected error when context is nil")
    }
}

// TestNewNilData expects an error when data is nil.
func TestNewNilData(t *testing.T) {
    ctx := context.Background()
    _, err := New(ctx, nil)
    if err == nil {
        t.Fatalf("expected error when data is nil")
    }
}

// TestInjectorDebounce verifies that the InjectorChanges callback is executed
// after a change signal respecting the debounce duration.
func TestInjectorDebounce(t *testing.T) {
    ctx := context.Background()
    data := &simpleStruct{}

    c, err := New(ctx, data)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    // shorten debounce to speed up the test
    c.Debounce = 10 * time.Millisecond

    var called int32
    c.InjectorChanges(func() error {
        atomic.StoreInt32(&called, 1)
        return nil
    })

    // trigger change
    c.injectorChgChan <- struct{}{}

    // wait longer than debounce duration
    time.Sleep(20 * time.Millisecond)

    if atomic.LoadInt32(&called) == 0 {
        t.Fatalf("expected InjectorChanges callback to be invoked")
    }
}

// TestExtractorDebounce verifies that the ExtractorChanges callback is executed
// after a change signal respecting the debounce duration.
func TestExtractorDebounce(t *testing.T) {
    ctx := context.Background()
    data := &simpleStruct{}

    c, err := New(ctx, data)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    // shorten debounce to speed up the test
    c.Debounce = 10 * time.Millisecond

    var called int32
    c.ExtractorChanges(func() error {
        atomic.StoreInt32(&called, 1)
        return nil
    })

    // trigger change
    c.extractorChgChan <- struct{}{}

    // wait longer than debounce duration
    time.Sleep(20 * time.Millisecond)

    if atomic.LoadInt32(&called) == 0 {
        t.Fatalf("expected ExtractorChanges callback to be invoked")
    }
}
