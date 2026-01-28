package control_test

import (
	"testing"

	"github.com/kjbreil/syncer/pkg/control"
	"google.golang.org/protobuf/proto"
)

// TestEntry_ProtobufRoundtrip verifies that Entry objects survive protobuf
// marshal/unmarshal, ensuring all Go types can be transmitted over gRPC.
func TestEntry_ProtobufRoundtrip(t *testing.T) {
	tests := []struct {
		name  string
		entry *control.Entry
	}{
		{
			name: "string value",
			entry: &control.Entry{
				Key: []*control.Key{
					{Key: "TestStruct"},
					{Key: "String"},
				},
				Value: control.NewObject(control.MakePtr("hello")),
			},
		},
		{
			name: "int64 value",
			entry: &control.Entry{
				Key: []*control.Key{
					{Key: "TestStruct"},
					{Key: "Int"},
				},
				Value: control.NewObject(control.MakePtr(int64(42))),
			},
		},
		{
			name: "int8 as int64",
			entry: &control.Entry{
				Key: []*control.Key{
					{Key: "TestStruct"},
					{Key: "Int8"},
				},
				Value: control.NewObject(control.MakePtr(int64(-42))),
			},
		},
		{
			name: "uint64 value",
			entry: &control.Entry{
				Key: []*control.Key{
					{Key: "TestStruct"},
					{Key: "Uint64"},
				},
				Value: control.NewObject(control.MakePtr(uint64(12345))),
			},
		},
		{
			name: "float32 value",
			entry: &control.Entry{
				Key: []*control.Key{
					{Key: "TestStruct"},
					{Key: "Float32"},
				},
				Value: control.NewObject(control.MakePtr(float32(3.14))),
			},
		},
		{
			name: "float64 value",
			entry: &control.Entry{
				Key: []*control.Key{
					{Key: "TestStruct"},
					{Key: "Float64"},
				},
				Value: control.NewObject(control.MakePtr(float64(2.71828))),
			},
		},
		{
			name: "bool value",
			entry: &control.Entry{
				Key: []*control.Key{
					{Key: "TestStruct"},
					{Key: "Bool"},
				},
				Value: control.NewObject(control.MakePtr(true)),
			},
		},
		{
			name: "bytes value (for []byte or complex)",
			entry: &control.Entry{
				Key: []*control.Key{
					{Key: "TestStruct"},
					{Key: "Bytes"},
				},
				Value: &control.Object{Bytes: []byte{0xDE, 0xAD, 0xBE, 0xEF}},
			},
		},
		{
			name: "complex64 as bytes",
			entry: &control.Entry{
				Key: []*control.Key{
					{Key: "TestStruct"},
					{Key: "Complex64"},
				},
				Value: control.NewObject(control.MakePtr(complex64(complex(1.5, 2.5)))),
			},
		},
		{
			name: "complex128 as bytes",
			entry: &control.Entry{
				Key: []*control.Key{
					{Key: "TestStruct"},
					{Key: "Complex128"},
				},
				Value: control.NewObject(control.MakePtr(complex(3.14, 2.71))),
			},
		},
		{
			name: "remove entry",
			entry: &control.Entry{
				Key: []*control.Key{
					{Key: "TestStruct"},
					{Key: "Ptr"},
				},
				Remove: true,
			},
		},
		{
			name: "slice index entry",
			entry: &control.Entry{
				Key: []*control.Key{
					{Key: "TestStruct"},
					{Key: "Slice", Index: control.NewObjects(control.NewObject(control.MakePtr(int64(2))))},
				},
				Value: control.NewObject(control.MakePtr(int64(42))),
			},
		},
		{
			name: "map index entry",
			entry: &control.Entry{
				Key: []*control.Key{
					{Key: "TestStruct"},
					{Key: "Map", Index: control.NewObjects(control.NewObject(control.MakePtr("mykey")))},
				},
				Value: control.NewObject(control.MakePtr(int64(100))),
			},
		},
		{
			name: "nested map entry",
			entry: &control.Entry{
				Key: []*control.Key{
					{Key: "TestStruct"},
					{Key: "MapMap", Index: control.NewObjects(
						control.NewObject(control.MakePtr("outer")),
						control.NewObject(control.MakePtr("inner")),
					)},
				},
				Value: control.NewObject(control.MakePtr(int64(5))),
			},
		},
		{
			name: "uint map key entry",
			entry: &control.Entry{
				Key: []*control.Key{
					{Key: "TestStruct"},
					{Key: "MapKeyUint", Index: control.NewObjects(control.NewObject(control.MakePtr(uint64(42))))},
				},
				Value: control.NewObject(control.MakePtr(int64(1))),
			},
		},
		{
			name: "float64 map key entry",
			entry: &control.Entry{
				Key: []*control.Key{
					{Key: "TestStruct"},
					{Key: "MapKeyFloat", Index: control.NewObjects(control.NewObject(control.MakePtr(float64(2.5))))},
				},
				Value: control.NewObject(control.MakePtr(int64(1))),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to protobuf bytes
			data, err := proto.Marshal(tt.entry)
			if err != nil {
				t.Fatalf("proto.Marshal() error: %v", err)
			}

			// Unmarshal back
			decoded := &control.Entry{}
			err = proto.Unmarshal(data, decoded)
			if err != nil {
				t.Fatalf("proto.Unmarshal() error: %v", err)
			}

			// Verify equality
			if !tt.entry.Equals(decoded) {
				t.Errorf("entry != decoded after protobuf roundtrip\n  original: %v\n  decoded:  %v", tt.entry, decoded)
			}
		})
	}
}

// TestEntries_ProtobufRoundtrip tests that a full set of entries survives serialization.
func TestEntries_ProtobufRoundtrip(t *testing.T) {
	entries := control.Entries{
		{
			Key: []*control.Key{
				{Key: "TestStruct"},
				{Key: "String"},
			},
			Value: control.NewObject(control.MakePtr("test")),
		},
		{
			Key: []*control.Key{
				{Key: "TestStruct"},
				{Key: "Int"},
			},
			Value: control.NewObject(control.MakePtr(int64(42))),
		},
		{
			Key: []*control.Key{
				{Key: "TestStruct"},
				{Key: "Bool"},
			},
			Value: control.NewObject(control.MakePtr(true)),
		},
		{
			Key: []*control.Key{
				{Key: "TestStruct"},
				{Key: "Float64"},
			},
			Value: control.NewObject(control.MakePtr(float64(3.14))),
		},
		{
			Key: []*control.Key{
				{Key: "TestStruct"},
				{Key: "Uint64"},
			},
			Value: control.NewObject(control.MakePtr(uint64(999))),
		},
		{
			Key: []*control.Key{
				{Key: "TestStruct"},
				{Key: "Bytes"},
			},
			Value: &control.Object{Bytes: []byte{1, 2, 3, 4}},
		},
	}

	// Marshal and unmarshal each entry
	decoded := make(control.Entries, len(entries))
	for i, e := range entries {
		data, err := proto.Marshal(e)
		if err != nil {
			t.Fatalf("Marshal entry %d error: %v", i, err)
		}
		decoded[i] = &control.Entry{}
		err = proto.Unmarshal(data, decoded[i])
		if err != nil {
			t.Fatalf("Unmarshal entry %d error: %v", i, err)
		}
	}

	if !entries.Equals(decoded) {
		t.Error("entries != decoded after protobuf roundtrip")
	}
}

// TestObject_BytesEquality verifies that the fixed bytes comparison works
// (previously panicked).
func TestObject_BytesEquality(t *testing.T) {
	o1 := &control.Object{Bytes: []byte{1, 2, 3}}
	o2 := &control.Object{Bytes: []byte{1, 2, 3}}
	o3 := &control.Object{Bytes: []byte{4, 5, 6}}
	o4 := &control.Object{Bytes: nil}

	if !o1.Equals(o2) {
		t.Error("identical bytes should be equal")
	}
	if o1.Equals(o3) {
		t.Error("different bytes should not be equal")
	}
	if o1.Equals(o4) {
		t.Error("bytes vs nil should not be equal")
	}
	if !o4.Equals(&control.Object{Bytes: nil}) {
		t.Error("both nil bytes should be equal")
	}
}

// TestObject_ComplexRoundtrip tests complex number encoding/decoding via Object.
func TestObject_ComplexRoundtrip(t *testing.T) {
	// Complex64
	c64 := complex64(complex(1.5, -2.5))
	obj64 := control.NewObject(control.MakePtr(c64))
	if obj64.GetBytes() == nil {
		t.Fatal("Complex64 should be encoded as bytes")
	}
	if len(obj64.GetBytes()) != 8 {
		t.Fatalf("Complex64 bytes length: got %d, want 8", len(obj64.GetBytes()))
	}

	// Marshal/unmarshal
	data, err := proto.Marshal(&control.Entry{
		Key:   []*control.Key{{Key: "test"}},
		Value: obj64,
	})
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	decoded := &control.Entry{}
	err = proto.Unmarshal(data, decoded)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if !obj64.Equals(decoded.GetValue()) {
		t.Error("Complex64 object not equal after protobuf roundtrip")
	}

	// Complex128
	c128 := complex(3.14159265358979, -2.71828182845904)
	obj128 := control.NewObject(control.MakePtr(c128))
	if obj128.GetBytes() == nil {
		t.Fatal("Complex128 should be encoded as bytes")
	}
	if len(obj128.GetBytes()) != 16 {
		t.Fatalf("Complex128 bytes length: got %d, want 16", len(obj128.GetBytes()))
	}

	data, err = proto.Marshal(&control.Entry{
		Key:   []*control.Key{{Key: "test"}},
		Value: obj128,
	})
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	decoded = &control.Entry{}
	err = proto.Unmarshal(data, decoded)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if !obj128.Equals(decoded.GetValue()) {
		t.Error("Complex128 object not equal after protobuf roundtrip")
	}
}
