package deepcopy

import (
	"fmt"
	"reflect"
	"testing"
)

type TestStruct struct {
	String         string
	Int            int
	Interface      TestInterface
	Slice          []int
	SliceStruct    []SD
	SlicePtr       []*int
	SlicePtrStruct []*SD
	SliceInterface []TestInterface
	SliceSlice     [][]int
	SliceMap       []map[string]int
	Array          [10]int
	ArrayStruct    [10]SD
	ArrayPtr       [10]*int
	ArrayPtrStruct [10]*SD
	ArrayInterface [10]TestInterface
	ArrayArray     [10][10]int
	Map            map[string]int
	MapKeyInt      map[int]int
	MapStruct      map[string]TestStruct
	MapPtr         map[string]*int
	MapPtrStruct   map[string]*TestStruct
	MapInterface   map[string]TestInterface
	MapMap         map[string]map[string]int
	MapSlice       map[string][]int
	SubStruct      TestSub
	SubStructPtr   *TestStruct
}

type SD struct {
	Name string
	Data string
}

type TestSub struct {
	MapPtrStruct map[int64]*SD
	S            string
}

type TestInterface interface {
	String() string
}

type TestInterfaceImpl struct {
	S string
}

func (t *TestInterfaceImpl) String() string {
	return t.S
}

func Test_Any(t *testing.T) {
	tests := []struct {
		name   string
		dst    any
		src    any
		wantFn func(src, dst any) (bool, string)
	}{
		{
			name: "int",
			dst:  3,
			src:  1,
			wantFn: func(src, dst any) (bool, string) {
				return reflect.DeepEqual(src, dst), fmt.Sprintf("src: %v, dst: %v", src, dst)
			},
		},

		{
			name: "int pointer",
			dst:  makePtr(3),
			src:  makePtr(1),
			wantFn: func(src, dst any) (bool, string) {
				if src == dst {
					return false, "pointers pointing to same"
				}
				return reflect.DeepEqual(src, dst), ""
			},
		},

		{
			name: "nil val in struct",
			dst:  nil,
			src: struct {
				Val *int
			}{
				Val: nil,
			},
			wantFn: func(src, dst any) (bool, string) {
				return reflect.DeepEqual(src, dst), ""
			},
		},

		{
			name: "nil slice in struct",
			dst:  nil,
			src: struct {
				Val []int
			}{
				Val: nil,
			},
			wantFn: func(src, dst any) (bool, string) {
				dstS := dst.(struct {
					Val []int
				})
				if dstS.Val != nil {
					return false, "dst slice is not nil"
				}
				return reflect.DeepEqual(src, dst), ""
			},
		},
		{
			name: "nil map in struct",
			dst:  nil,
			src: struct {
				Val map[string]int
			}{
				Val: nil,
			},
			wantFn: func(src, dst any) (bool, string) {
				dstS := dst.(struct {
					Val map[string]int
				})
				if dstS.Val != nil {
					return false, "dst slice is not nil"
				}
				return reflect.DeepEqual(src, dst), ""
			},
		},
		{
			name: "map ptr val",
			dst:  nil,
			src: map[string]*int{
				"test": makePtr(1),
			},
			wantFn: func(src, dst any) (bool, string) {
				s := src.(map[string]*int)
				d := dst.(map[string]*int)

				if s["test"] == d["test"] {
					return false, "pointers pointing to same"
				}
				return reflect.DeepEqual(src, dst), ""
			},
		},
		{
			name: "slice",
			dst:  nil,
			src: []int{
				1,
			},
			wantFn: func(src, dst any) (bool, string) {
				return reflect.DeepEqual(src, dst), fmt.Sprintf("slice src: %v, dst: %v", src, dst)
			},
		},
		{
			name: "slice ptr val",
			dst:  nil,
			src: []*int{
				makePtr(1),
			},
			wantFn: func(src, dst any) (bool, string) {
				s := src.([]*int)
				d := dst.([]*int)

				if s[0] == d[0] {
					return false, "pointers pointing to same"
				}
				return reflect.DeepEqual(src, dst), ""
			},
		},
		{
			name: "array",
			dst:  nil,
			src: [1]int{
				1,
			},
			wantFn: func(src, dst any) (bool, string) {
				return reflect.DeepEqual(src, dst), fmt.Sprintf("slice src: %v, dst: %v", src, dst)
			},
		},
		{
			name: "array ptr val",
			dst:  nil,
			src: [1]*int{
				makePtr(1),
			},
			wantFn: func(src, dst any) (bool, string) {
				s := src.([1]*int)
				d := dst.([1]*int)

				if s[0] == d[0] {
					return false, "pointers pointing to same"
				}
				return reflect.DeepEqual(src, dst), ""
			},
		},
		{
			name: "interface",
			dst:  nil,
			src: map[string]TestInterface{
				"test": &TestInterfaceImpl{S: "test"},
			},
			wantFn: func(src, dst any) (bool, string) {
				s := src.(map[string]TestInterface)
				d := dst.(map[string]TestInterface)

				if s["test"] == d["test"] {
					return false, "pointers pointing to same"
				}
				return reflect.DeepEqual(src, dst), ""
			},
		},
		{
			name: "unexported type",
			dst:  nil,
			src: struct {
				String     string
				unexported string
			}{
				String:     "String",
				unexported: "unexported",
			},
			wantFn: func(src, dst any) (bool, string) {
				s := src.(struct {
					String     string
					unexported string
				})
				d := dst.(struct {
					String     string
					unexported string
				})
				if s.String != d.String {
					return false, fmt.Sprintf("s String: %v != d String: %v", s.String, d.String)
				}

				if s.unexported == d.unexported {
					return false, fmt.Sprintf("s unexported: %v == d unexported: %v", s.unexported, d.unexported)
				}

				return true, ""
			},
		},
		{
			name: "lots of types",
			dst:  nil,
			src:  makeBaseTestStruct(),
			wantFn: func(src, dst any) (bool, string) {
				return reflect.DeepEqual(src, dst), ""
			},
		},
		{
			name: "lots of types zero val",
			dst:  nil,
			src:  TestStruct{},
			wantFn: func(src, dst any) (bool, string) {
				return reflect.DeepEqual(src, dst), ""
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.dst = Any(tt.src)
			if ok, errS := tt.wantFn(tt.src, tt.dst); !ok {
				t.Errorf(errS)
			}
		})
	}
}

func Test_DeepCopy(t *testing.T) {
	tests := []struct {
		name   string
		dst    any
		src    any
		wantFn func(src, dst any) (bool, string)
	}{
		{
			name: "lots of types",
			dst:  nil,
			src:  makeBaseTestStruct(),
			wantFn: func(src, dst any) (bool, string) {
				return reflect.DeepEqual(src, dst), ""
			},
		},
		{
			name: "lots of types zero val",
			dst:  nil,
			src:  TestStruct{},
			wantFn: func(src, dst any) (bool, string) {
				return reflect.DeepEqual(src, dst), ""
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srcV := reflect.ValueOf(tt.src)
			tt.dst = DeepCopy(srcV).Interface()
			if ok, errS := tt.wantFn(tt.src, tt.dst); !ok {
				t.Errorf(errS)
			}
		})
	}
}

func Benchmark_Any(b *testing.B) {
	ts := makeBaseTestStruct()
	for i := 0; i < b.N; i++ {
		_ = Any(ts)
	}
}

func makeBaseTestStruct() TestStruct {
	return TestStruct{
		String:    "Base String",
		Int:       1,
		Interface: &TestInterfaceImpl{S: "TestInterface"},
		Slice:     []int{1, 2, 3},
		SliceStruct: []SD{
			{Name: "Base SliceStruct Name 1", Data: "Base SliceStruct Data 1"},
			{Name: "Base SliceStruct Name 2", Data: "Base SliceStruct Data 2"},
		},
		SlicePtr: []*int{makePtr(1), makePtr(2)},
		SlicePtrStruct: []*SD{
			{Name: "Base SlicePtrStruct Name 1", Data: "Base SlicePtrStruct Data 1"},
			{Name: "Base SlicePtrStruct Name 2", Data: "Base SlicePtrStruct Data 2"},
		},
		SliceInterface: []TestInterface{&TestInterfaceImpl{S: "TestInterface"}},
		SliceSlice:     [][]int{{1}},
		Array:          [10]int{1},
		ArrayStruct: [10]SD{
			{Name: "Base SliceStruct Name 1", Data: "Base SliceStruct Data 1"},
			{Name: "Base SliceStruct Name 2", Data: "Base SliceStruct Data 2"},
		},
		ArrayPtr: [10]*int{makePtr(1), makePtr(2)},
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
		MapStruct: map[string]TestStruct{
			"Base First": {String: "Base First Test Struct"},
		},
		MapPtr: map[string]*int{
			"Base First": makePtr(1),
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
	}
}

func makePtr[V any](v V) *V {
	return &v
}
