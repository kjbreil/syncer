package control

import (
	"reflect"
	"testing"
)

type TestInt int

const (
	TestInt1 TestInt = iota
	TestInt2
)

func TestNewObject(t *testing.T) {
	tests := []struct {
		name string
		v    any
		want *Object
	}{
		{
			name: "string",
			v:    "test",
			want: &Object{String_: MakePtr("test")},
		},
		{
			name: "pointer to string",
			v:    MakePtr("test"),
			want: &Object{String_: MakePtr("test")},
		},
		{
			name: "int",
			v:    1,
			want: &Object{Int64: MakePtr(int64(1))},
		},
		{
			name: "int8",
			v:    int8(1),
			want: &Object{Int64: MakePtr(int64(1))},
		},
		{
			name: "int16",
			v:    int16(1),
			want: &Object{Int64: MakePtr(int64(1))},
		},
		{
			name: "int32",
			v:    int32(1),
			want: &Object{Int64: MakePtr(int64(1))},
		},
		{
			name: "int64",
			v:    int64(1),
			want: &Object{Int64: MakePtr(int64(1))},
		},
		{
			name: "uint",
			v:    uint(1),
			want: &Object{Uint64: MakePtr(uint64(1))},
		},
		{
			name: "uint8",
			v:    uint8(1),
			want: &Object{Uint64: MakePtr(uint64(1))},
		},
		{
			name: "uint16",
			v:    uint16(1),
			want: &Object{Uint64: MakePtr(uint64(1))},
		},
		{
			name: "uint32",
			v:    uint32(1),
			want: &Object{Uint64: MakePtr(uint64(1))},
		},
		{
			name: "uint64",
			v:    uint64(1),
			want: &Object{Uint64: MakePtr(uint64(1))},
		},
		{
			name: "float32",
			v:    float32(1),
			want: &Object{Float32: MakePtr(float32(1))},
		},
		{
			name: "float64",
			v:    float64(1),
			want: &Object{Float64: MakePtr(float64(1))},
		},
		{
			name: "bool",
			v:    true,
			want: &Object{Bool: MakePtr(true)},
		},
		{
			name: "[]byte",
			v:    []byte("test"),
			want: &Object{Bytes: []byte("test")},
		},
		{
			name: "Object",
			v:    Object{},
			want: &Object{},
		},
		{
			name: "nil",
			v:    nil,
			want: nil,
		},
		{
			name: "Test Int Type",
			v:    TestInt1,
			want: &Object{Int64: MakePtr(int64(TestInt1))},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewObject(tt.v)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
