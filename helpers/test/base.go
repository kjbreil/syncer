package test

import "github.com/kjbreil/syncer/control"

func MakeBaseTestStruct() TestStruct {
	return TestStruct{
		String:    "Base String",
		Int:       1,
		Interface: &TestInterfaceImpl{S: "TestInterface"},
		Slice:     []int{1, 2, 3},
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
		MapKeyInt:   map[int]int{1: 1},
		MapKeyBool:  map[bool]int{true: 1},
		MapKeyUint:  map[uint]int{1: 1},
		MapKeyFloat: map[float64]int{1: 1},
		MapKeyPtr:   map[*int]int{control.MakePtr(1): 1},
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
