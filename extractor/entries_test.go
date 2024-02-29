package extractor

import (
	"testing"

	"github.com/kjbreil/syncer/control"
	. "github.com/kjbreil/syncer/helpers/test"
)

func TestExtractor_Entries(t *testing.T) {
	tests := []struct {
		name string
		// structure must be a &struct otherwise reflect only sees an interfaces when dereferencing the pointer
		structure any
		change    any
		want      []*control.Entry
	}{
		{
			name: "string",
			structure: &struct {
				String string
			}{
				String: "test",
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "",
						},
						{
							Key: "String",
						},
					},
					Value: control.NewObject(control.MakePtr("test")),
				},
			},
		},
		{
			name: "pointer",
			structure: &struct {
				Ptr *struct {
					String string
				}
			}{
				Ptr: &struct {
					String string
				}{
					String: "test",
				},
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "",
						},
						{
							Key: "Ptr",
						},
						{
							Key: "String",
						},
					},
					Value: control.NewObject(control.MakePtr("test")),
				},
			},
		},
		{
			name: "nil pointer",
			structure: &struct {
				Ptr *struct {
					String string
				}
			}{
				Ptr: &struct {
					String string
				}{
					String: "test",
				},
			},
			change: &struct {
				Ptr *struct {
					String string
				}
			}{
				Ptr: nil,
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "",
						},
						{
							Key: "Ptr",
						},
					},
					Remove: true,
				},
			},
		},
		{
			name: "slice",
			structure: &struct {
				Slice []int
			}{
				Slice: []int{1, 2},
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "",
						},
						{
							Key:   "Slice",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(1))),
				},
				{
					Key: []*control.Key{
						{
							Key: "",
						},
						{
							Key:   "Slice",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(2))),
				},
			},
		},
		{
			name: "make slice nil",
			structure: &struct {
				Slice []int
			}{
				Slice: []int{1, 2},
			},
			change: &struct {
				Slice []int
			}{
				Slice: nil,
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "",
						},
						{
							Key: "Slice",
						},
					},
					Remove: true,
				},
			},
		},
		{
			name: "remove from slice",
			structure: &struct {
				Slice []int
			}{
				Slice: []int{1, 2},
			},
			change: &struct {
				Slice []int
			}{
				Slice: []int{1},
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "",
						},
						{
							Key:   "Slice",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
						},
					},
					Remove: true,
				},
			},
		},
		{
			name: "Slice Slice",
			structure: &struct {
				Slice [][]int
			}{
				Slice: [][]int{{1}},
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "",
						},
						{
							Key:   "Slice",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0))), control.NewObject(control.MakePtr(int64(0)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(1))),
				},
			},
		},
		{
			name: "Slice Struct Slice",
			structure: &struct {
				Slice []struct {
					Slice []int
				}
			}{
				Slice: []struct {
					Slice []int
				}{
					{
						Slice: []int{1, 2},
					},
				},
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "",
						},
						{
							Key:   "Slice",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
						{
							Key:   "Slice",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(1))),
				},
				{
					Key: []*control.Key{
						{
							Key: "",
						},
						{
							Key:   "Slice",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
						{
							Key:   "Slice",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(2))),
				},
			},
		},
		{
			name: "Nil Slice -> Slice",
			structure: &struct {
				Slice []int
			}{
				Slice: nil,
			},
			change: &struct {
				Slice []int
			}{
				Slice: []int{1, 2},
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "",
						},
						{
							Key:   "Slice",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(1))),
				},
				{
					Key: []*control.Key{
						{
							Key: "",
						},
						{
							Key:   "Slice",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(2))),
				},
			},
		},
		{
			name: "map",
			structure: &struct {
				Map map[string]int
			}{
				Map: map[string]int{"test": 1},
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "",
						},
						{
							Key:   "Map",
							Index: control.NewObjects(control.NewObject(control.MakePtr("test"))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(1))),
				},
			},
		},
		{
			name: "nil map",
			structure: &struct {
				Map map[string]int
			}{
				Map: map[string]int{"test": 1},
			},
			change: &struct {
				Map map[string]int
			}{
				Map: nil,
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "",
						},
						{
							Key: "Map",
						},
					},
					Remove: true,
				},
			},
		},
		{
			name: "remove from map",
			structure: &struct {
				Map map[string]int
			}{
				Map: map[string]int{"one": 1, "two": 2},
			},
			change: &struct {
				Map map[string]int
			}{
				Map: map[string]int{"one": 1},
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "",
						},
						{
							Key:   "Map",
							Index: control.NewObjects(control.NewObject(control.MakePtr("two"))),
						},
					},
					Remove: true,
				},
			},
		},
		{
			name: "array",
			structure: &struct {
				Array [5]int
			}{
				Array: [5]int{1, 2},
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "",
						},
						{
							Key:   "Array",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(1))),
				},
				{
					Key: []*control.Key{
						{
							Key: "",
						},
						{
							Key:   "Array",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(2))),
				},
			},
		},
		{
			name: "array of pointers",
			structure: &struct {
				Array [5]*int
			}{
				Array: [5]*int{control.MakePtr(1), control.MakePtr(2)},
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "",
						},
						{
							Key:   "Array",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(1))),
				},
				{
					Key: []*control.Key{
						{
							Key: "",
						},
						{
							Key:   "Array",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(2))),
				},
			},
		},
		{
			name: "remove from array of pointers",
			structure: &struct {
				Array [5]*int
			}{
				Array: [5]*int{control.MakePtr(1)},
			},
			change: &struct {
				Array [5]*int
			}{
				Array: [5]*int{nil},
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "",
						},
						{
							Key:   "Array",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
					},
					Remove: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ext, err := New(tt.structure)
			if err != nil {
				t.Fatalf("could not create extractor: %v", err)
			}
			got, err := ext.Entries(tt.structure)
			if err != nil {
				t.Fatalf("could not get diff: %v", err)
			}
			if tt.change != nil {
				got, err = ext.Entries(tt.change)
				if err != nil {
					t.Fatalf("could not get diff: %v", err)
				}
			}
			if !got.Equals(tt.want) {
				t.Errorf("Entries() = %v, want %v", got, tt.want)
				// t.Logf("##########\n\n%s\n\n##########", got.Diff(tt.want).Struct())
				t.Logf("##########\n\n%s\n\n##########", got.Struct())
			}
			if tt.change != nil {
				got, err = ext.Entries(tt.change)
				if err != nil {
					t.Fatalf("could not get diff: %v", err)
				}
			}
		})
	}
}

func TestExtractor_GetDiff_Big(t *testing.T) {
	// ts := MakeBaseTestStruct()
	ts := TestStruct{}
	newTS := func() {
		ts = TestStruct{}
	}

	ext, err := New(ts)
	if err != nil {
		t.Fatalf("could not create extractor: %v", err)
	}
	tests := []struct {
		name  string
		preFn func()
		modFn func()
		want  []*control.Entry
	}{
		{
			name: "change string",
			modFn: func() {
				ts.String = "change string"
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key: "String",
						},
					},
					Value: &control.Object{String_: control.MakePtr("change string")},
				},
			},
		},
		{
			name: "add to int slice",
			preFn: func() {
				ts.Slice = []int{1, 2, 3}
			},
			modFn: func() {
				ts.Slice = append(ts.Slice, 4)
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "Slice",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(3)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(4))),
				},
			},
		},
		{
			name: "remove from int slice",
			modFn: func() {
				ts.Slice = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
				_, err = ext.Entries(&ts)
				if err != nil {
					t.Fatalf("could not get Entries in modFn: %v", err)
				}
				ts.Slice = ts.Slice[:len(ts.Slice)-5]
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "Slice",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(4)))),
						},
					},
					Remove: true,
				},
			},
		},
		{
			name: "remove from int slice check update old",
			preFn: func() {
				ts.Slice = []int{1, 2, 3}
			},
			modFn: func() {
				ts.Slice = ts.Slice[:len(ts.Slice)-1]
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "Slice",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(2)))),
						},
					},
					Remove: true,
				},
			},
		},
		{
			name: "make slice nil",
			preFn: func() {
				ts.Slice = []int{1, 2, 3}
			},
			modFn: func() {
				ts.Slice = nil
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key: "Slice",
						},
					},
					Remove: true,
				},
			},
		},
		{
			name: "empty slice",
			preFn: func() {
				ts.Slice = []int{1, 2, 3}
			},
			modFn: func() {
				ts.Slice = []int{}
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "Slice",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
					},
					Remove: true,
				},
			},
		},
		{
			name: "empty slice",
			preFn: func() {
				ts.Slice = []int{1, 2, 3}
			},
			modFn: func() {
				ts.Slice = []int{}
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "Slice",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
					},
					Remove: true,
				},
			},
		},
		{
			name: "MapStructSlice",
			preFn: func() {

			},
			modFn: func() {
				ts.MapStructSlice = map[int64]TestSub{
					1: {
						Slice: []SD{
							{
								Name: "test",
							},
						},
					},
				}
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "MapStructSlice",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
						},
						{
							Key:   "Slice",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
						{
							Key: "Name",
						},
					},
					Value: control.NewObject(control.MakePtr("test")),
				},
			},
		},
		{
			name: "add to SliceStruct",
			modFn: func() {
				ts.SliceStruct = []SD{
					{Name: "SliceStruct Name 1", Data: "SliceStruct Data 1"},
					{Name: "SliceStruct Name 2", Data: "SliceStruct Data 2"},
				}
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "SliceStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
						{
							Key: "Name",
						},
					},
					Value: control.NewObject(control.MakePtr("SliceStruct Name 1")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "SliceStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
						{
							Key: "Data",
						},
					},
					Value: control.NewObject(control.MakePtr("SliceStruct Data 1")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "SliceStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
						},
						{
							Key: "Name",
						},
					},
					Value: control.NewObject(control.MakePtr("SliceStruct Name 2")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "SliceStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
						},
						{
							Key: "Data",
						},
					},
					Value: control.NewObject(control.MakePtr("SliceStruct Data 2")),
				},
			},
		},
		{
			name: "add to SlicePtrStruct",
			modFn: func() {
				ts.SlicePtrStruct = []*TestStruct{
					{String: "SlicePtrStruct String 1"},
					{String: "SlicePtrStruct String 2"},
				}
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "SlicePtrStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
						{
							Key: "String",
						},
					},
					Value: control.NewObject(control.MakePtr("SlicePtrStruct String 1")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "SlicePtrStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
						},
						{
							Key: "String",
						},
					},
					Value: control.NewObject(control.MakePtr("SlicePtrStruct String 2")),
				},
			},
		},
		{
			name: "remove from map",
			preFn: func() {
				ts.Map = map[string]int{"Base First": 1}
			},
			modFn: func() {
				delete(ts.Map, "Base First")
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "Map",
							Index: control.NewObjects(control.NewObject(control.MakePtr("Base First"))),
						},
					},
					Remove: true,
				},
			},
		},
		{
			name: "make map nil",
			preFn: func() {
				ts.Map = map[string]int{"Base First": 1}
			},
			modFn: func() {
				ts.Map = nil
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key: "Map",
						},
					},
					Remove: true,
				},
			},
		},
		{
			name: "add to empty MapStruct",
			preFn: func() {
			},
			modFn: func() {
				ts.MapStruct = map[string]TestStruct{
					"Base First": {String: "Base First Test Struct"},
				}
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "MapStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr("Base First"))),
						},
						{
							Key: "String",
						},
					},
					Value: control.NewObject(control.MakePtr("Base First Test Struct")),
				},
			},
		},

		{
			name: "add to empty MapMapType",
			preFn: func() {
			},
			modFn: func() {
				ts.MapMapType = map[string]MapType{
					"Base First": {
						"Base Second": {
							S: "inside",
						},
					},
				}
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "MapMapType",
							Index: control.NewObjects(control.NewObject(control.MakePtr("Base First")), control.NewObject(control.MakePtr("Base Second"))),
						},
						{
							Key: "S",
						},
					},
					Value: control.NewObject(control.MakePtr("inside")),
				},
			},
		},
		{
			name: "remove from MapStruct",
			preFn: func() {
				ts.MapStruct = map[string]TestStruct{
					"Base First": {String: "Base First Test Struct"},
				}
			},
			modFn: func() {
				delete(ts.MapStruct, "Base First")
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "MapStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr("Base First"))),
						},
					},
					Remove: true,
				},
			},
		},
		{
			name: "add pointer struct",
			modFn: func() {
				ts.SubStructPtr = &TestStruct{
					String: "SubStructPtr String",
				}
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key: "SubStructPtr",
						},
						{
							Key: "String",
						},
					},
					Value: control.NewObject(control.MakePtr("SubStructPtr String")),
				},
			},
		},
		{
			name: "remove pointer struct",
			preFn: func() {
				ts.SubStructPtr = &TestStruct{
					String: "SubStructPtr String",
				}
			},
			modFn: func() {
				ts.SubStructPtr = nil
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key: "SubStructPtr",
						},
					},
					Remove: true,
				},
			},
		},
		{
			name: "MapKeyInt",
			preFn: func() {
				ts.MapKeyInt = map[int]int{0: 0}
			},
			modFn: func() {
				ts.MapKeyInt[0] = 2
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "MapKeyInt",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(2))),
				},
			},
		},
		{
			name: "MapKeyUint",
			modFn: func() {
				ts.MapKeyUint = map[uint]int{0: 2}
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "MapKeyUint",
							Index: control.NewObjects(control.NewObject(control.MakePtr(uint64(0)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(2))),
				},
			},
		},
		{
			name: "MapKeyFloat",
			modFn: func() {
				ts.MapKeyFloat = map[float64]int{1: 2}
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "MapKeyFloat",
							Index: control.NewObjects(control.NewObject(control.MakePtr(float64(1.000000)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(2))),
				},
			},
		},
		{
			name: "Add Interface",
			modFn: func() {
				ts.Interface = &TestInterfaceImpl{S: "TestInterface"}
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key: "Interface",
						},
						{
							Key: "S",
						},
					},
					Value: control.NewObject(control.MakePtr("TestInterface")),
				},
			},
		},
		{
			name: "Remove Interface",
			preFn: func() {
				ts.Interface = &TestInterfaceImpl{S: "TestInterface"}
			},
			modFn: func() {
				ts.Interface = nil
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key: "Interface",
						},
					},
					Remove: true,
				},
			},
		},
		{
			name: "Add to Array",
			preFn: func() {
				ts.Array = [10]int{1, 2, 3, 4}
			},
			modFn: func() {
				ts.Array[4] = 5
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "Array",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(4)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(5))),
				},
			},
		},
		{
			name: "All Test Types in Pre Should have no entries",
			preFn: func() {
				ts = MakeBaseTestStruct()
			},
			want: []*control.Entry{},
		},

		{
			name: "All Test Types All Entries",
			modFn: func() {
				ts = MakeBaseTestStruct()
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key: "String",
						},
					},
					Value: control.NewObject(control.MakePtr("Base String")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key: "Int",
						},
					},
					Value: control.NewObject(control.MakePtr(int64(1))),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key: "Interface",
						},
						{
							Key: "S",
						},
					},
					Value: control.NewObject(control.MakePtr("TestInterface")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "Slice",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(1))),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "Slice",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(2))),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "Slice",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(2)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(3))),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "SliceStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
						{
							Key: "Name",
						},
					},
					Value: control.NewObject(control.MakePtr("Base SliceStruct Name 1")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "SliceStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
						{
							Key: "Data",
						},
					},
					Value: control.NewObject(control.MakePtr("Base SliceStruct Data 1")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "SliceStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
						},
						{
							Key: "Name",
						},
					},
					Value: control.NewObject(control.MakePtr("Base SliceStruct Name 2")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "SliceStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
						},
						{
							Key: "Data",
						},
					},
					Value: control.NewObject(control.MakePtr("Base SliceStruct Data 2")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "SlicePtr",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(1))),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "SlicePtr",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(2))),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "SlicePtrStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
						{
							Key: "String",
						},
					},
					Value: control.NewObject(control.MakePtr("Base SlicePtrStruct Name 1")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "SlicePtrStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
						},
						{
							Key: "String",
						},
					},
					Value: control.NewObject(control.MakePtr("Base SlicePtrStruct Name 2")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "SliceInterface",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
						{
							Key: "S",
						},
					},
					Value: control.NewObject(control.MakePtr("TestInterface")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "SliceSlice",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0))), control.NewObject(control.MakePtr(int64(0)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(1))),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "Array",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(1))),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "ArrayStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
						{
							Key: "Name",
						},
					},
					Value: control.NewObject(control.MakePtr("Base SliceStruct Name 1")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "ArrayStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
						{
							Key: "Data",
						},
					},
					Value: control.NewObject(control.MakePtr("Base SliceStruct Data 1")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "ArrayStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
						},
						{
							Key: "Name",
						},
					},
					Value: control.NewObject(control.MakePtr("Base SliceStruct Name 2")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "ArrayStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
						},
						{
							Key: "Data",
						},
					},
					Value: control.NewObject(control.MakePtr("Base SliceStruct Data 2")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "ArrayPtr",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(1))),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "ArrayPtr",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(2))),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "ArrayPtrStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
						{
							Key: "Name",
						},
					},
					Value: control.NewObject(control.MakePtr("Base SlicePtrStruct Name 1")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "ArrayPtrStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
						{
							Key: "Data",
						},
					},
					Value: control.NewObject(control.MakePtr("Base SlicePtrStruct Data 1")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "ArrayPtrStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
						},
						{
							Key: "Name",
						},
					},
					Value: control.NewObject(control.MakePtr("Base SlicePtrStruct Name 2")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "ArrayPtrStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
						},
						{
							Key: "Data",
						},
					},
					Value: control.NewObject(control.MakePtr("Base SlicePtrStruct Data 2")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "ArrayInterface",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0)))),
						},
						{
							Key: "S",
						},
					},
					Value: control.NewObject(control.MakePtr("TestInterface")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "ArrayArray",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0))), control.NewObject(control.MakePtr(int64(0)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(1))),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "ArrayArray",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(0))), control.NewObject(control.MakePtr(int64(1)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(2))),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "ArrayArray",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1))), control.NewObject(control.MakePtr(int64(0)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(3))),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "ArrayArray",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1))), control.NewObject(control.MakePtr(int64(1)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(4))),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "Map",
							Index: control.NewObjects(control.NewObject(control.MakePtr("Base First"))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(1))),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "Map",
							Index: control.NewObjects(control.NewObject(control.MakePtr("Base Second"))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(2))),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "MapKeyInt",
							Index: control.NewObjects(control.NewObject(control.MakePtr(int64(1)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(1))),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "MapKeyBool",
							Index: control.NewObjects(control.NewObject(control.MakePtr(true))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(1))),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "MapKeyUint",
							Index: control.NewObjects(control.NewObject(control.MakePtr(uint64(1)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(1))),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "MapKeyFloat",
							Index: control.NewObjects(control.NewObject(control.MakePtr(float64(1.000000)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(1))),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "MapStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr("Base First"))),
						},
						{
							Key: "String",
						},
					},
					Value: control.NewObject(control.MakePtr("Base First Test Struct")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "MapPtr",
							Index: control.NewObjects(control.NewObject(control.MakePtr("Base First"))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(1))),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "MapPtrStruct",
							Index: control.NewObjects(control.NewObject(control.MakePtr("Base First"))),
						},
						{
							Key: "String",
						},
					},
					Value: control.NewObject(control.MakePtr("Base First Test Struct")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "MapInterface",
							Index: control.NewObjects(control.NewObject(control.MakePtr("one"))),
						},
						{
							Key: "S",
						},
					},
					Value: control.NewObject(control.MakePtr("TestInterface")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "MapMap",
							Index: control.NewObjects(control.NewObject(control.MakePtr("top")), control.NewObject(control.MakePtr("second"))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(3))),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "MapSlice",
							Index: control.NewObjects(control.NewObject(control.MakePtr("one")), control.NewObject(control.MakePtr(int64(0)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(1))),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key:   "MapSlice",
							Index: control.NewObjects(control.NewObject(control.MakePtr("one")), control.NewObject(control.MakePtr(int64(1)))),
						},
					},
					Value: control.NewObject(control.MakePtr(int64(2))),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key: "SubStruct",
						},
						{
							Key: "S",
						},
					},
					Value: control.NewObject(control.MakePtr("Base Test Sub")),
				},
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key: "SubStructPtr",
						},
						{
							Key: "String",
						},
					},
					Value: control.NewObject(control.MakePtr("Base Sub Test Struct")),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// make sure a base slate is set before running the modFn
			newTS()
			if tt.preFn != nil {
				tt.preFn()
			}
			_, err = ext.Entries(&ts)
			if err != nil {
				t.Errorf("Entries() error = %v", err)
			}
			if tt.modFn != nil {
				tt.modFn()
			}

			got, err := ext.Entries(&ts)
			if err != nil {
				t.Fatalf("Entries() error = %v", err)
			}
			if !got.Equals(tt.want) {
				t.Errorf("Entries() = %v, want %v", got, tt.want)
				// t.Logf("##########\n\n%s\n\n##########", got.Diff(tt.want).Struct())
				t.Logf("##########\n\n%s\n\n##########", got.Struct())
			}
		})
	}
}

func BenchmarkExtractor_GetDiff(b *testing.B) {
	cs := MakeChangeTestStruct()
	ts := MakeBaseTestStruct()

	sliceEntries := 100000
	ts.Slice = make([]int, 0, sliceEntries)
	for ii := 0; ii < sliceEntries; ii++ {
		ts.Slice = append(ts.Slice, ii)
	}
	ext, err := New(ts)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		_, err = ext.Entries(&ts)
		if err != nil {
			b.Fatal(err)
		}
		_, err = ext.Entries(&cs)
		if err != nil {
			b.Fatal(err)
		}
	}
}
