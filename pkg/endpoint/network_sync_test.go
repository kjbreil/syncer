package endpoint

import (
	"net"
	"testing"
	"time"

	"github.com/kjbreil/syncer/pkg/endpoint/settings"
)

type syncStruct struct {
	String     string
	Int        int
	Int8       int8
	Int16      int16
	Int32      int32
	Int64      int64
	Uint       uint
	Uint8      uint8
	Uint16     uint16
	Uint32     uint32
	Uint64     uint64
	Float32    float32
	Float64    float64
	Complex64  complex64
	Complex128 complex128
	Bool       bool
	Byte       byte
	Bytes      []byte
	Slice      []int
	Map        map[string]int
	Sub        subStruct
	SubPtr     *syncStruct
}

type subStruct struct {
	Name string
}

// TestNetworkSync_AllTypes tests end-to-end synchronization of all Go types
// over a real gRPC connection between server and client endpoints.
func TestNetworkSync_AllTypes(t *testing.T) {
	port := findFreePort(t)

	serverData := &syncStruct{
		String:     "hello",
		Int:        42,
		Int8:       -8,
		Int16:      -16,
		Int32:      -32,
		Int64:      -64,
		Uint:       10,
		Uint8:      18,
		Uint16:     116,
		Uint32:     132,
		Uint64:     164,
		Float32:    3.14,
		Float64:    2.71828,
		Complex64:  complex(1.5, 2.5),
		Complex128: complex(3.14, 2.71),
		Bool:       true,
		Byte:       0xFF,
		Bytes:      []byte{0xDE, 0xAD, 0xBE, 0xEF},
		Slice:      []int{1, 2, 3},
		Map:        map[string]int{"a": 1, "b": 2},
		Sub:        subStruct{Name: "sub"},
		SubPtr:     &syncStruct{String: "ptr"},
	}

	clientData := &syncStruct{}

	// Server has no peers (it will start as a server)
	serverEP, err := New(serverData, &settings.Settings{
		Port:       port,
		Peers:      []net.TCPAddr{},
		AutoUpdate: true,
	})
	if err != nil {
		t.Fatalf("server New() error: %v", err)
	}

	// Start server first and wait for it to be a server
	serverEP.Run(false)
	waitForServer(t, serverEP)
	// Give the HTTP/gRPC server time to start accepting connections
	time.Sleep(500 * time.Millisecond)

	// Client connects to the server
	clientEP, err := New(clientData, &settings.Settings{
		Port: port + 1,
		Peers: []net.TCPAddr{
			{
				IP:   net.ParseIP("127.0.0.1"),
				Port: port,
			},
		},
		AutoUpdate: true,
	})
	if err != nil {
		t.Fatalf("client New() error: %v", err)
	}

	// Start client with onlyClient=true so it doesn't fall back to server mode
	clientEP.Run(true)
	waitForRunning2(t, clientEP)

	// Wait for data to sync (Init is called inside tryPeers on successful connection)
	time.Sleep(2 * time.Second)

	// Verify synced data
	if clientData.String != serverData.String {
		t.Errorf("String: got %q, want %q", clientData.String, serverData.String)
	}
	if clientData.Int != serverData.Int {
		t.Errorf("Int: got %d, want %d", clientData.Int, serverData.Int)
	}
	if clientData.Int8 != serverData.Int8 {
		t.Errorf("Int8: got %d, want %d", clientData.Int8, serverData.Int8)
	}
	if clientData.Int16 != serverData.Int16 {
		t.Errorf("Int16: got %d, want %d", clientData.Int16, serverData.Int16)
	}
	if clientData.Int32 != serverData.Int32 {
		t.Errorf("Int32: got %d, want %d", clientData.Int32, serverData.Int32)
	}
	if clientData.Int64 != serverData.Int64 {
		t.Errorf("Int64: got %d, want %d", clientData.Int64, serverData.Int64)
	}
	if clientData.Uint != serverData.Uint {
		t.Errorf("Uint: got %d, want %d", clientData.Uint, serverData.Uint)
	}
	if clientData.Uint8 != serverData.Uint8 {
		t.Errorf("Uint8: got %d, want %d", clientData.Uint8, serverData.Uint8)
	}
	if clientData.Uint16 != serverData.Uint16 {
		t.Errorf("Uint16: got %d, want %d", clientData.Uint16, serverData.Uint16)
	}
	if clientData.Uint32 != serverData.Uint32 {
		t.Errorf("Uint32: got %d, want %d", clientData.Uint32, serverData.Uint32)
	}
	if clientData.Uint64 != serverData.Uint64 {
		t.Errorf("Uint64: got %d, want %d", clientData.Uint64, serverData.Uint64)
	}
	if clientData.Float32 != serverData.Float32 {
		t.Errorf("Float32: got %f, want %f", clientData.Float32, serverData.Float32)
	}
	if clientData.Float64 != serverData.Float64 {
		t.Errorf("Float64: got %f, want %f", clientData.Float64, serverData.Float64)
	}
	if clientData.Complex64 != serverData.Complex64 {
		t.Errorf("Complex64: got %v, want %v", clientData.Complex64, serverData.Complex64)
	}
	if clientData.Complex128 != serverData.Complex128 {
		t.Errorf("Complex128: got %v, want %v", clientData.Complex128, serverData.Complex128)
	}
	if clientData.Bool != serverData.Bool {
		t.Errorf("Bool: got %v, want %v", clientData.Bool, serverData.Bool)
	}
	if clientData.Byte != serverData.Byte {
		t.Errorf("Byte: got %d, want %d", clientData.Byte, serverData.Byte)
	}
	if len(clientData.Bytes) != len(serverData.Bytes) {
		t.Errorf("Bytes length: got %d, want %d", len(clientData.Bytes), len(serverData.Bytes))
	} else {
		for i := range serverData.Bytes {
			if clientData.Bytes[i] != serverData.Bytes[i] {
				t.Errorf("Bytes[%d]: got %d, want %d", i, clientData.Bytes[i], serverData.Bytes[i])
			}
		}
	}
	if len(clientData.Slice) != len(serverData.Slice) {
		t.Errorf("Slice length: got %d, want %d", len(clientData.Slice), len(serverData.Slice))
	}
	if clientData.Sub.Name != serverData.Sub.Name {
		t.Errorf("Sub.Name: got %q, want %q", clientData.Sub.Name, serverData.Sub.Name)
	}

	// Cleanup
	serverEP.Stop()
	clientEP.Stop()
}

func findFreePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("findFreePort: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return port
}

func waitForServer(t *testing.T, ep *Endpoint) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for !ep.IsServer() {
		if time.Now().After(deadline) {
			t.Fatal("endpoint did not start as server in time")
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func waitForRunning2(t *testing.T, ep *Endpoint) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for !ep.Running() {
		if time.Now().After(deadline) {
			t.Fatal("endpoint did not start in time")
		}
		time.Sleep(100 * time.Millisecond)
	}
}
