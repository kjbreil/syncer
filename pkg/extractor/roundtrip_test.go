package extractor_test

import (
	"testing"

	"github.com/kjbreil/syncer/pkg/extractor"
	"github.com/kjbreil/syncer/pkg/injector"
	. "github.com/kjbreil/syncer/pkg/test"
)

// TestRoundtrip_AllTypes verifies that extracting entries from a populated struct
// and injecting them into an empty struct produces an equivalent result.
// This covers all Go types: all int/uint sizes, floats, complex, bool, string,
// byte, []byte, slices, arrays, maps, structs, pointers, interfaces.
func TestRoundtrip_AllTypes(t *testing.T) {
	base := MakeBaseTestStruct()

	// Create extractor from zero-value struct (so all fields are "changed")
	ext, err := extractor.New(&base)
	if err != nil {
		t.Fatalf("extractor.New() error: %v", err)
	}

	// Extract all entries (all fields are new relative to zero)
	entries, err := ext.Entries(&base)
	if err != nil {
		t.Fatalf("Entries() error: %v", err)
	}

	if len(entries) == 0 {
		t.Fatal("expected entries from base struct, got none")
	}

	// Create target struct with interface pre-initialized (interfaces need concrete type)
	target := TestStruct{
		Interface: &TestInterfaceImpl{},
		SliceInterface: []TestInterface{&TestInterfaceImpl{}},
		ArrayInterface: [10]TestInterface{&TestInterfaceImpl{}},
		MapInterface:   map[string]TestInterface{"one": &TestInterfaceImpl{}},
	}
	inj, err := injector.New(&target)
	if err != nil {
		t.Fatalf("injector.New() error: %v", err)
	}

	// Inject all entries
	err = inj.AddAll(entries)
	if err != nil {
		t.Fatalf("AddAll() error: %v", err)
	}

	// Verify all primitive types
	if target.String != base.String {
		t.Errorf("String: got %q, want %q", target.String, base.String)
	}
	if target.Int != base.Int {
		t.Errorf("Int: got %d, want %d", target.Int, base.Int)
	}
	if target.Int8 != base.Int8 {
		t.Errorf("Int8: got %d, want %d", target.Int8, base.Int8)
	}
	if target.Int16 != base.Int16 {
		t.Errorf("Int16: got %d, want %d", target.Int16, base.Int16)
	}
	if target.Int32 != base.Int32 {
		t.Errorf("Int32: got %d, want %d", target.Int32, base.Int32)
	}
	if target.Int64 != base.Int64 {
		t.Errorf("Int64: got %d, want %d", target.Int64, base.Int64)
	}
	if target.Uint != base.Uint {
		t.Errorf("Uint: got %d, want %d", target.Uint, base.Uint)
	}
	if target.Uint8 != base.Uint8 {
		t.Errorf("Uint8: got %d, want %d", target.Uint8, base.Uint8)
	}
	if target.Uint16 != base.Uint16 {
		t.Errorf("Uint16: got %d, want %d", target.Uint16, base.Uint16)
	}
	if target.Uint32 != base.Uint32 {
		t.Errorf("Uint32: got %d, want %d", target.Uint32, base.Uint32)
	}
	if target.Uint64 != base.Uint64 {
		t.Errorf("Uint64: got %d, want %d", target.Uint64, base.Uint64)
	}
	if target.Float32 != base.Float32 {
		t.Errorf("Float32: got %f, want %f", target.Float32, base.Float32)
	}
	if target.Float64 != base.Float64 {
		t.Errorf("Float64: got %f, want %f", target.Float64, base.Float64)
	}
	if target.Complex64 != base.Complex64 {
		t.Errorf("Complex64: got %v, want %v", target.Complex64, base.Complex64)
	}
	if target.Complex128 != base.Complex128 {
		t.Errorf("Complex128: got %v, want %v", target.Complex128, base.Complex128)
	}
	if target.Bool != base.Bool {
		t.Errorf("Bool: got %v, want %v", target.Bool, base.Bool)
	}
	if target.Byte != base.Byte {
		t.Errorf("Byte: got %d, want %d", target.Byte, base.Byte)
	}
	if len(target.Bytes) != len(base.Bytes) {
		t.Errorf("Bytes length: got %d, want %d", len(target.Bytes), len(base.Bytes))
	} else {
		for i := range base.Bytes {
			if target.Bytes[i] != base.Bytes[i] {
				t.Errorf("Bytes[%d]: got %d, want %d", i, target.Bytes[i], base.Bytes[i])
			}
		}
	}

	// Verify slices
	if len(target.Slice) != len(base.Slice) {
		t.Errorf("Slice length: got %d, want %d", len(target.Slice), len(base.Slice))
	} else {
		for i, v := range base.Slice {
			if target.Slice[i] != v {
				t.Errorf("Slice[%d]: got %d, want %d", i, target.Slice[i], v)
			}
		}
	}
	if len(target.SliceStruct) != len(base.SliceStruct) {
		t.Errorf("SliceStruct length: got %d, want %d", len(target.SliceStruct), len(base.SliceStruct))
	}

	// Verify maps
	if len(target.Map) != len(base.Map) {
		t.Errorf("Map length: got %d, want %d", len(target.Map), len(base.Map))
	}
	for k, v := range base.Map {
		if tv, ok := target.Map[k]; !ok || tv != v {
			t.Errorf("Map[%q]: got %d, want %d", k, tv, v)
		}
	}

	// Verify sub struct
	if target.SubStruct.S != base.SubStruct.S {
		t.Errorf("SubStruct.S: got %q, want %q", target.SubStruct.S, base.SubStruct.S)
	}

	// Verify sub struct ptr
	if target.SubStructPtr == nil {
		t.Error("SubStructPtr: got nil, want non-nil")
	} else if target.SubStructPtr.String != base.SubStructPtr.String {
		t.Errorf("SubStructPtr.String: got %q, want %q", target.SubStructPtr.String, base.SubStructPtr.String)
	}

	// Verify interface
	if target.Interface == nil {
		t.Error("Interface: got nil, want non-nil")
	} else if target.Interface.String() != base.Interface.String() {
		t.Errorf("Interface.String(): got %q, want %q", target.Interface.String(), base.Interface.String())
	}

	// Verify arrays
	if target.Array != base.Array {
		t.Errorf("Array: got %v, want %v", target.Array, base.Array)
	}
	if target.ArrayArray != base.ArrayArray {
		t.Errorf("ArrayArray: got %v, want %v", target.ArrayArray, base.ArrayArray)
	}
}

// TestRoundtrip_Changes verifies that changes between states are correctly
// extracted and injected, simulating a real sync scenario.
func TestRoundtrip_Changes(t *testing.T) {
	base := MakeBaseTestStruct()
	change := MakeChangeTestStruct()

	// Set up interface on change too
	change.Interface = &TestInterfaceImpl{S: "ChangedInterface"}

	// Create extractor with base state
	ext, err := extractor.New(&base)
	if err != nil {
		t.Fatalf("extractor.New() error: %v", err)
	}

	// First extraction establishes baseline
	_, err = ext.Entries(&base)
	if err != nil {
		t.Fatalf("Entries() baseline error: %v", err)
	}

	// Second extraction gets the changes
	entries, err := ext.Entries(&change)
	if err != nil {
		t.Fatalf("Entries() changes error: %v", err)
	}

	if len(entries) == 0 {
		t.Fatal("expected change entries, got none")
	}

	// Inject changes into a copy of base
	target := MakeBaseTestStruct()
	target.Interface = &TestInterfaceImpl{} // pre-init for injection
	inj, err := injector.New(&target)
	if err != nil {
		t.Fatalf("injector.New() error: %v", err)
	}

	err = inj.AddAll(entries)
	if err != nil {
		t.Fatalf("AddAll() error: %v", err)
	}

	// Verify changed fields
	if target.String != change.String {
		t.Errorf("String: got %q, want %q", target.String, change.String)
	}
	if target.Int != change.Int {
		t.Errorf("Int: got %d, want %d", target.Int, change.Int)
	}
	if target.Int8 != change.Int8 {
		t.Errorf("Int8: got %d, want %d", target.Int8, change.Int8)
	}
	if target.Float32 != change.Float32 {
		t.Errorf("Float32: got %f, want %f", target.Float32, change.Float32)
	}
	if target.Float64 != change.Float64 {
		t.Errorf("Float64: got %f, want %f", target.Float64, change.Float64)
	}
	if target.Complex64 != change.Complex64 {
		t.Errorf("Complex64: got %v, want %v", target.Complex64, change.Complex64)
	}
	if target.Complex128 != change.Complex128 {
		t.Errorf("Complex128: got %v, want %v", target.Complex128, change.Complex128)
	}
	if target.Bool != change.Bool {
		t.Errorf("Bool: got %v, want %v", target.Bool, change.Bool)
	}
	if target.Byte != change.Byte {
		t.Errorf("Byte: got %d, want %d", target.Byte, change.Byte)
	}
	if len(target.Bytes) != len(change.Bytes) {
		t.Errorf("Bytes length: got %d, want %d", len(target.Bytes), len(change.Bytes))
	} else {
		for i := range change.Bytes {
			if target.Bytes[i] != change.Bytes[i] {
				t.Errorf("Bytes[%d]: got %d, want %d", i, target.Bytes[i], change.Bytes[i])
			}
		}
	}
	if target.SubStruct.S != change.SubStruct.S {
		t.Errorf("SubStruct.S: got %q, want %q", target.SubStruct.S, change.SubStruct.S)
	}
}

// TestRoundtrip_PerType tests individual type roundtrips for focused coverage.
func TestRoundtrip_PerType(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ts *TestStruct)
		verify func(t *testing.T, ts *TestStruct)
	}{
		{
			name: "int8",
			setup: func(ts *TestStruct) {
				ts.Int8 = -42
			},
			verify: func(t *testing.T, ts *TestStruct) {
				if ts.Int8 != -42 {
					t.Errorf("Int8: got %d, want -42", ts.Int8)
				}
			},
		},
		{
			name: "int16",
			setup: func(ts *TestStruct) {
				ts.Int16 = -1000
			},
			verify: func(t *testing.T, ts *TestStruct) {
				if ts.Int16 != -1000 {
					t.Errorf("Int16: got %d, want -1000", ts.Int16)
				}
			},
		},
		{
			name: "int32",
			setup: func(ts *TestStruct) {
				ts.Int32 = -100000
			},
			verify: func(t *testing.T, ts *TestStruct) {
				if ts.Int32 != -100000 {
					t.Errorf("Int32: got %d, want -100000", ts.Int32)
				}
			},
		},
		{
			name: "int64",
			setup: func(ts *TestStruct) {
				ts.Int64 = -9999999999
			},
			verify: func(t *testing.T, ts *TestStruct) {
				if ts.Int64 != -9999999999 {
					t.Errorf("Int64: got %d, want -9999999999", ts.Int64)
				}
			},
		},
		{
			name: "uint8",
			setup: func(ts *TestStruct) {
				ts.Uint8 = 200
			},
			verify: func(t *testing.T, ts *TestStruct) {
				if ts.Uint8 != 200 {
					t.Errorf("Uint8: got %d, want 200", ts.Uint8)
				}
			},
		},
		{
			name: "uint16",
			setup: func(ts *TestStruct) {
				ts.Uint16 = 50000
			},
			verify: func(t *testing.T, ts *TestStruct) {
				if ts.Uint16 != 50000 {
					t.Errorf("Uint16: got %d, want 50000", ts.Uint16)
				}
			},
		},
		{
			name: "uint32",
			setup: func(ts *TestStruct) {
				ts.Uint32 = 3000000000
			},
			verify: func(t *testing.T, ts *TestStruct) {
				if ts.Uint32 != 3000000000 {
					t.Errorf("Uint32: got %d, want 3000000000", ts.Uint32)
				}
			},
		},
		{
			name: "uint64",
			setup: func(ts *TestStruct) {
				ts.Uint64 = 18446744073709551000
			},
			verify: func(t *testing.T, ts *TestStruct) {
				if ts.Uint64 != 18446744073709551000 {
					t.Errorf("Uint64: got %d, want 18446744073709551000", ts.Uint64)
				}
			},
		},
		{
			name: "float32",
			setup: func(ts *TestStruct) {
				ts.Float32 = 3.14
			},
			verify: func(t *testing.T, ts *TestStruct) {
				if ts.Float32 != 3.14 {
					t.Errorf("Float32: got %f, want 3.14", ts.Float32)
				}
			},
		},
		{
			name: "float64",
			setup: func(ts *TestStruct) {
				ts.Float64 = 2.718281828459045
			},
			verify: func(t *testing.T, ts *TestStruct) {
				if ts.Float64 != 2.718281828459045 {
					t.Errorf("Float64: got %f, want 2.718281828459045", ts.Float64)
				}
			},
		},
		{
			name: "complex64",
			setup: func(ts *TestStruct) {
				ts.Complex64 = complex(1.5, 2.5)
			},
			verify: func(t *testing.T, ts *TestStruct) {
				if ts.Complex64 != complex(1.5, 2.5) {
					t.Errorf("Complex64: got %v, want (1.5+2.5i)", ts.Complex64)
				}
			},
		},
		{
			name: "complex128",
			setup: func(ts *TestStruct) {
				ts.Complex128 = complex(3.14, 2.71)
			},
			verify: func(t *testing.T, ts *TestStruct) {
				if ts.Complex128 != complex(3.14, 2.71) {
					t.Errorf("Complex128: got %v, want (3.14+2.71i)", ts.Complex128)
				}
			},
		},
		{
			name: "bool true",
			setup: func(ts *TestStruct) {
				ts.Bool = true
			},
			verify: func(t *testing.T, ts *TestStruct) {
				if !ts.Bool {
					t.Errorf("Bool: got false, want true")
				}
			},
		},
		{
			name: "byte",
			setup: func(ts *TestStruct) {
				ts.Byte = 0xAB
			},
			verify: func(t *testing.T, ts *TestStruct) {
				if ts.Byte != 0xAB {
					t.Errorf("Byte: got %d, want 0xAB", ts.Byte)
				}
			},
		},
		{
			name: "bytes",
			setup: func(ts *TestStruct) {
				ts.Bytes = []byte{0xDE, 0xAD, 0xBE, 0xEF}
			},
			verify: func(t *testing.T, ts *TestStruct) {
				want := []byte{0xDE, 0xAD, 0xBE, 0xEF}
				if len(ts.Bytes) != len(want) {
					t.Errorf("Bytes length: got %d, want %d", len(ts.Bytes), len(want))
					return
				}
				for i := range want {
					if ts.Bytes[i] != want[i] {
						t.Errorf("Bytes[%d]: got %d, want %d", i, ts.Bytes[i], want[i])
					}
				}
			},
		},
		{
			name: "map with float32 key",
			setup: func(ts *TestStruct) {
				ts.MapKeyFloat32 = map[float32]int{1.5: 42}
			},
			verify: func(t *testing.T, ts *TestStruct) {
				v, ok := ts.MapKeyFloat32[1.5]
				if !ok || v != 42 {
					t.Errorf("MapKeyFloat32[1.5]: got %d (exists=%v), want 42", v, ok)
				}
			},
		},
		{
			name: "map with float64 key",
			setup: func(ts *TestStruct) {
				ts.MapKeyFloat = map[float64]int{2.5: 99}
			},
			verify: func(t *testing.T, ts *TestStruct) {
				v, ok := ts.MapKeyFloat[2.5]
				if !ok || v != 99 {
					t.Errorf("MapKeyFloat[2.5]: got %d (exists=%v), want 99", v, ok)
				}
			},
		},
		{
			name: "map with uint key",
			setup: func(ts *TestStruct) {
				ts.MapKeyUint = map[uint]int{42: 100}
			},
			verify: func(t *testing.T, ts *TestStruct) {
				v, ok := ts.MapKeyUint[42]
				if !ok || v != 100 {
					t.Errorf("MapKeyUint[42]: got %d (exists=%v), want 100", v, ok)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := TestStruct{}
			ext, err := extractor.New(&src)
			if err != nil {
				t.Fatalf("extractor.New() error: %v", err)
			}

			// Set the value
			tt.setup(&src)

			// Extract
			entries, err := ext.Entries(&src)
			if err != nil {
				t.Fatalf("Entries() error: %v", err)
			}

			if len(entries) == 0 {
				t.Fatal("expected entries, got none")
			}

			// Inject into empty struct
			target := TestStruct{}
			inj, err := injector.New(&target)
			if err != nil {
				t.Fatalf("injector.New() error: %v", err)
			}

			err = inj.AddAll(entries)
			if err != nil {
				t.Fatalf("AddAll() error: %v", err)
			}

			// Verify
			tt.verify(t, &target)
		})
	}
}
