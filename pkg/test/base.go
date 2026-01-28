package test

import "github.com/kjbreil/syncer/pkg/control"

func MakeBaseTestStruct() TestStruct {
	return TestStruct{
		String:     "Base String",
		Int:        1,
		Int8:       8,
		Int16:      16,
		Int32:      32,
		Int64:      64,
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
		Bytes:      []byte{0x01, 0x02, 0x03},
		Interface:  &TestInterfaceImpl{S: "TestInterface"},
		Slice:      []int{1, 2, 3},
		SliceStruct: []SD{
			{Name: "Base SliceStruct Name 1", Data: "Base SliceStruct Data 1"},
			{Name: "Base SliceStruct Name 2", Data: "Base SliceStruct Data 2"},
		},
		SlicePtr: []*int{control.MakePtr(1), control.MakePtr(2)},
		SlicePtrStruct: []*TestStruct{
			{String: "Base SlicePtrStruct Name 1"},
			{String: "Base SlicePtrStruct Name 2"},
		},
		SliceInterface: []TestInterface{&TestInterfaceImpl{S: "TestInterface"}},
		SliceSlice:     [][]int{{1}},
		Array:          [10]int{1},
		ArrayStruct: [10]SD{
			{Name: "Base SliceStruct Name 1", Data: "Base SliceStruct Data 1"},
			{Name: "Base SliceStruct Name 2", Data: "Base SliceStruct Data 2"},
		},
		ArrayPtr: [10]*int{control.MakePtr(1), control.MakePtr(2)},
		ArrayPtrStruct: [10]*SD{
			{Name: "Base SlicePtrStruct Name 1", Data: "Base SlicePtrStruct Data 1"},
			{Name: "Base SlicePtrStruct Name 2", Data: "Base SlicePtrStruct Data 2"},
		},
		ArrayInterface: [10]TestInterface{&TestInterfaceImpl{S: "TestInterface"}},
		ArrayArray:     [10][10]int{{1, 2}, {3, 4}},
		Map: map[string]int{
			"Base First":  1,
			"Base Second": 2,
		},
		MapKeyInt:     map[int]int{1: 1},
		MapKeyBool:    map[bool]int{true: 1},
		MapKeyUint:    map[uint]int{1: 1},
		MapKeyFloat:   map[float64]int{1: 1},
		MapKeyFloat32: map[float32]int{1.5: 1},
		MapStruct: map[string]TestStruct{
			"Base First": {String: "Base First Test Struct"},
		},
		MapPtr: map[string]*int{
			"Base First": control.MakePtr(1),
		},
		MapPtrStruct: map[string]*TestStruct{
			"Base First": {String: "Base First Test Struct"},
		},
		MapInterface: map[string]TestInterface{"one": &TestInterfaceImpl{S: "TestInterface"}},
		MapMap:       map[string]map[string]int{"top": {"second": 3}},
		MapSlice:     map[string][]int{"one": {1, 2}},
		SubStruct: TestSub{
			S: "Base Test Sub",
		},
		SubStructPtr: &TestStruct{String: "Base Sub Test Struct"},
		unexported:   "not exported",
		Function: func() {
		},
	}
}

func MakeChangeTestStruct() TestStruct {
	return TestStruct{
		String:     "Change String",
		Int:        2,
		Int8:       -8,
		Int16:      -16,
		Int32:      -32,
		Int64:      -64,
		Uint:       20,
		Uint8:      28,
		Uint16:     216,
		Uint32:     232,
		Uint64:     264,
		Float32:    6.28,
		Float64:    1.41421,
		Complex64:  complex(10.5, 20.5),
		Complex128: complex(30.14, 20.71),
		Bool:       false,
		Byte:       0xAA,
		Bytes:      []byte{0xDE, 0xAD},
		Slice:      []int{4, 5, 6, 7},
		SliceStruct: []SD{
			{Name: "Change SliceStruct Name 1", Data: "Change SliceStruct Data 1"},
			{Name: "Change SliceStruct Name 2", Data: "Change SliceStruct Data 2"},
		},
		SlicePtr:       []*int{control.MakePtr(1), control.MakePtr(2)},
		SlicePtrStruct: []*TestStruct{},
		Map: map[string]int{
			"Change Third":  3,
			"Change Fourth": 4,
		},
		MapStruct: map[string]TestStruct{
			"Change First": {String: "Change First Test Struct"},
		},
		MapPtr: map[string]*int{
			"Change First": control.MakePtr(1),
		},
		MapPtrStruct: map[string]*TestStruct{
			"Change First": {String: "Change First Test Struct"},
		},
		SubStruct: TestSub{
			S: "Change Test Sub",
		},
		SubStructPtr: &TestStruct{String: "Change Sub Test Struct"},
	}
}
