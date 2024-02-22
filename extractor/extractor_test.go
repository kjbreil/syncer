package extractor

import (
	"github.com/kjbreil/syncer/helpers/equal"
	"reflect"
	"testing"

	"github.com/kjbreil/syncer/control"
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
	unexported     string
	Function       func()
}

type TestSub struct {
	MapPtrStruct map[int64]*SD
	S            string
}

type SD struct {
	Name string
	Data string
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

func makeChangeTestStruct() TestStruct {
	return TestStruct{
		String: "Change String",
		Int:    2,
		Slice:  []int{4, 5, 6, 7},
		SliceStruct: []SD{
			{Name: "Change SliceStruct Name 1", Data: "Change SliceStruct Data 1"},
			{Name: "Change SliceStruct Name 2", Data: "Change SliceStruct Data 2"},
		},
		SlicePtr: []*int{control.MakePtr(1), control.MakePtr(2)},
		SlicePtrStruct: []*SD{
			{Name: "Change SlicePtrStruct Name 1", Data: "Change SlicePtrStruct Data 1"},
			{Name: "Change SlicePtrStruct Name 2", Data: "Change SlicePtrStruct Data 2"},
		},
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

func TestExtractor_Entries(t *testing.T) {

	ts := makeBaseTestStruct()
	ext, _ := New(ts)

	tests := []struct {
		name  string
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
							// Index: control.NewObjects(),
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
				_ = ext.Entries(&ts)
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
			name: "remove from map",
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
			name: "remove pointer struct",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// make sure a base slate is set before running the modFn
			ts = makeBaseTestStruct()
			ext.Entries(&ts)
			tt.modFn()

			got := ext.Entries(&ts)
			if !got.Equals(tt.want) {
				t.Errorf("Entries() = %v, want %v", got, tt.want)
				// t.Logf("##########\n\n%s\n\n##########", got.Diff(tt.want).Struct())
				t.Logf("##########\n\n%s\n\n##########", got.Struct())
			}
		})
	}
}

func BenchmarkExtractor_Diff(b *testing.B) {
	cs := makeChangeTestStruct()
	ts := makeBaseTestStruct()

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
		_ = ext.Entries(&ts)
		_ = ext.Entries(&cs)
	}
}

func Benchmark_equal1(b *testing.B) {
	var o, n reflect.Value
	for i := 0; i < b.N; i++ {
		o, n = reflect.ValueOf(i), reflect.ValueOf(i)
		equal.Equal(o, n)
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
		SlicePtr: []*int{control.MakePtr(1), control.MakePtr(2)},
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
