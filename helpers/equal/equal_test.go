package equal

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/kjbreil/syncer/control"
)

type TestInterface interface {
	String() string
}

type TestInterfaceImpl struct {
	S string
}

func (t *TestInterfaceImpl) String() string {
	return t.S
}

func TestEqual(t *testing.T) {
	type args struct {
		newValue any
		oldValue any
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "int match",
			args: args{
				newValue: 1,
				oldValue: 1,
			},
			want: true,
		},
		{
			name: "int not match",
			args: args{
				newValue: 1,
				oldValue: 2,
			},
			want: false,
		},
		{
			name: "uint match",
			args: args{
				newValue: uint(1),
				oldValue: uint(1),
			},
			want: true,
		},
		{
			name: "uint not match",
			args: args{
				newValue: uint(1),
				oldValue: uint(2),
			},
			want: false,
		},
		{
			name: "bool match",
			args: args{
				newValue: true,
				oldValue: true,
			},
			want: true,
		},
		{
			name: "bool not match",
			args: args{
				newValue: true,
				oldValue: false,
			},
			want: false,
		},

		{
			name: "float32 match",
			args: args{
				newValue: float32(1.2),
				oldValue: float32(1.2),
			},
			want: true,
		},
		{
			name: "float32 not match",
			args: args{
				newValue: float32(1.2),
				oldValue: float32(2.1),
			},
			want: false,
		},
		{
			name: "float32 == float64 not match, this is a trick",
			args: args{
				newValue: float32(1.2),
				oldValue: float64(1.2),
			},
			want: false,
		},
		{
			name: "float32 == float64  match",
			args: args{
				newValue: float32(1),
				oldValue: float64(1),
			},
			want: true,
		},
		{
			name: "float64 match",
			args: args{
				newValue: float64(1.2),
				oldValue: float64(1.2),
			},
			want: true,
		},
		{
			name: "float64 not match",
			args: args{
				newValue: float64(1.2),
				oldValue: float64(2.1),
			},
			want: false,
		},
		{
			name: "float64 not match",
			args: args{
				newValue: float64(1.2),
				oldValue: float64(2.1),
			},
			want: false,
		},
		{
			name: "complex64 match",
			args: args{
				newValue: complex64(1.2),
				oldValue: complex64(1.2),
			},
			want: true,
		},
		{
			name: "complex128 match",
			args: args{
				newValue: complex128(1.2),
				oldValue: complex128(1.2),
			},
			want: true,
		},
		{
			name: "pointers",
			args: args{
				newValue: control.MakePtr(1),
				oldValue: control.MakePtr(1),
			},
			want: true,
		},
		{
			name: "nil",
			args: args{
				newValue: nil,
				oldValue: nil,
			},
			want: true,
		},
		{
			name: "struct",
			args: args{
				newValue: struct {
					String string
				}{
					String: "test",
				},
				oldValue: struct {
					String string
				}{
					String: "test",
				},
			},
			want: true,
		},
		{
			name: "struct not match",
			args: args{
				newValue: struct {
					String string
				}{
					String: "test",
				},
				oldValue: struct {
					String string
				}{
					String: "not",
				},
			},
			want: false,
		},
		{
			name: "struct structure not match",
			args: args{
				newValue: struct {
					String string
				}{
					String: "test",
				},
				oldValue: struct {
					String  string
					String2 string
				}{
					String: "test",
				},
			},
			want: false,
		},
		{
			name: "interface",
			args: args{
				newValue: struct {
					IFace TestInterface
				}{
					IFace: &TestInterfaceImpl{},
				},
				oldValue: struct {
					IFace TestInterface
				}{
					IFace: &TestInterfaceImpl{},
				},
			},
			want: true,
		},
		{
			name: "slice",
			args: args{
				newValue: []int{1},
				oldValue: []int{1},
			},
			want: true,
		},
		{
			name: "slice not matched",
			args: args{
				newValue: []int{1},
				oldValue: []int{2},
			},
			want: false,
		},
		{
			name: "slice not matched length",
			args: args{
				newValue: []int{1, 2},
				oldValue: []int{1},
			},
			want: false,
		},
		{
			name: "array",
			args: args{
				newValue: [1]int{1},
				oldValue: [1]int{1},
			},
			want: true,
		},
		{
			name: "array not matched",
			args: args{
				newValue: [1]int{1},
				oldValue: [1]int{2},
			},
			want: false,
		},
		{
			name: "array of ptr matched with nil entries",
			args: args{
				newValue: [5]*int{control.MakePtr(1)},
				oldValue: [5]*int{control.MakePtr(1)},
			},
			want: true,
		},
		{
			name: "function",
			args: args{
				newValue: func() {},
				oldValue: func() {},
			},
			want: true,
		},
		{
			name: "function not matched",
			args: args{
				newValue: func() {},
				oldValue: func(bool) {},
			},
			want: false,
		},
		{
			name: "channel",
			args: args{
				newValue: make(chan int, 2),
				oldValue: make(chan int, 2),
			},
			want: true,
		},
		{
			name: "channel not match",
			args: args{
				newValue: make(chan int, 2),
				oldValue: make(chan string, 2),
			},
			want: false,
		},
		{
			name: "map",
			args: args{
				newValue: map[string]int{"one": 1},
				oldValue: map[string]int{"one": 1},
			},
			want: true,
		},
		{
			name: "map different length",
			args: args{
				newValue: map[string]int{"one": 1},
				oldValue: map[string]int{"one": 1, "two": 2},
			},
			want: false,
		},
		{
			name: "map different",
			args: args{
				newValue: map[string]int{"one": 1},
				oldValue: map[string]int{"two": 2},
			},
			want: false,
		},
		{
			name: "UnsafePointer ",
			args: args{
				newValue: unsafe.Pointer(&struct {
					String string
				}{
					String: "test",
				}),
				oldValue: unsafe.Pointer(&struct {
					String string
				}{
					String: "test",
				}),
			},
			want: false,
		},
		{
			name: "UnsafePointer ",
			args: args{
				newValue: struct {
					String string
				}{
					String: "test",
				},
				oldValue: nil,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newValue, oldValue := reflect.ValueOf(tt.args.newValue), reflect.ValueOf(tt.args.oldValue)
			if got := Equal(newValue, oldValue); got != tt.want {
				t.Errorf("Equal() = %v, want %v, newValue = %v, oldValue = %v", got, tt.want, newValue, oldValue)
			}
		})
	}
}
