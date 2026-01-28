package endpoint

import (
    "net"
    "testing"

    "github.com/kjbreil/syncer/pkg/endpoint/settings"
    "github.com/kjbreil/syncer/pkg/endpoint/server"
)

type stub struct{ V int }

func defaultSettings() *settings.Settings {
    return &settings.Settings{
        Port:       12345,
        Peers:      []net.TCPAddr{},
        AutoUpdate: false,
    }
}

// TestNewValidation checks constructor validation paths.
func TestNewValidation(t *testing.T) {
    ptrData := &stub{}
    st := defaultSettings()

    if _, err := New(ptrData, st); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if _, err := New(nil, st); err == nil {
        t.Fatalf("expected error when data is nil")
    }
    if _, err := New(stub{}, st); err != ErrNotPointer {
        t.Fatalf("expected ErrNotPointer, got %v", err)
    }
    if _, err := New(ptrData, nil); err == nil {
        t.Fatalf("expected error when settings is nil")
    }
}

// TestRandomIntRange ensures randomInt respects boundaries and fallback path.
func TestRandomIntRange(t *testing.T) {
    low, high := 1, 10
    for i := 0; i < 100; i++ {
        v := randomInt(low, high)
        if v < low || v > high {
            t.Fatalf("randomInt produced %d outside range %d-%d", v, low, high)
        }
    }
}

// TestEndpointRunningStop verifies Running/Stop bookkeeping without starting goroutines.
func TestEndpointRunningStop(t *testing.T) {
    ep, _ := New(&stub{}, defaultSettings())
    if ep.Running() {
        t.Fatalf("expected not running initially")
    }

    // Manually set server to non-nil to mimic run
    ep.server = &server.Server{}
    if !ep.Running() {
        t.Fatalf("expected running when server set")
    }
    ep.Stop()
    if ep.Running() {
        t.Fatalf("expected not running after Stop")
    }
}
