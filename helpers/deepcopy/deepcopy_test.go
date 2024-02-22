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
	SliceSlice     [][]int
	SliceInterface []TestInterface
	Map            map[string]int
	MapKeyInt      map[int]int
	MapStruct      map[string]TestStruct
	MapPtr         map[string]*int
	MapMap         map[string]map[string]int
	MapPtrStruct   map[string]*TestStruct
	MapInterface   map[string]TestInterface
	Sub            TestSub
	SubPtr         *TestStruct
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

func (t *TestSub) String() string {
	return t.S
}

func Test_DeepCopy(t *testing.T) {
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
		// {
		// 	name: "interface",
		// 	dst:  nil,
		// 	src: map[string]IFace{
		// 		"test": &IFaceImpl{S: "test"},
		// 	},
		// 	wantFn: func(src, dst any) (bool, string) {
		// 		s := src.(map[string]IFace)
		// 		d := dst.(map[string]IFace)
		//
		// 		if s["test"] == d["test"] {
		// 			return false, "pointers pointing to same"
		// 		}
		// 		return reflect.DeepEqual(src, dst), ""
		// 	},
		// },
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
			name: "array",
			dst:  nil,
			src: struct {
				String     string
				unexported string
			}{
				String:     "String",
				unexported: "unexported",
			},
			wantFn: func(src, dst any) (bool, string) {

				return true, ""
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

	sliceEntries := 100000
	ts.Slice = make([]int, 0, sliceEntries)
	for ii := 0; ii < sliceEntries; ii++ {
		ts.Slice = append(ts.Slice, ii)
	}

	for i := 0; i < b.N; i++ {
		_ = Any(ts)
	}
}

func makeBaseTestStruct() TestStruct {
	return TestStruct{
		String:    "Base String",
		Int:       1,
		Interface: &TestSub{S: "TestInterface"},
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
		Sub: TestSub{
			S: "Base Test Sub",
		},
		SubPtr: &TestStruct{String: "Base Sub Test Struct"},
	}
}

func makePtr[V any](v V) *V {
	return &v
}
