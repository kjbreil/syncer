package extractor

import (
	"fmt"
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
type TestStruct2 struct {
	String    string
	MapStruct map[int64]TestSub
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

func (t *TestSub) String() string {
	return t.S
}

func TestNew(t *testing.T) {
	ts := TestStruct2{
		MapStruct: map[int64]TestSub{
			1: {
				S: "test",
				MapPtrStruct: map[int64]*SD{
					1: &SD{
						Name: "test",
					},
				},
			},
		},
	}
	ext, err := New(&ts)
	if err != nil {
		t.Fatal(err)
	}
	entries := ext.Entries(&ts)
	fmt.Println(entries.Struct())
	ts.MapStruct[1].MapPtrStruct[1].Name = "new"
	entries = ext.Entries(&ts)
	fmt.Println(entries.Struct())

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
		SlicePtr: []*int{control.MakePtr(1), control.MakePtr(2)},
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
			"Base First": control.MakePtr(1),
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
		Sub: TestSub{
			S: "Change Test Sub",
		},
		SubPtr: &TestStruct{String: "Change Sub Test Struct"},
	}
}

func Test_equal(t *testing.T) {
	type args struct {
		newValue reflect.Value
		oldValue reflect.Value
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				newValue: reflect.ValueOf(1),
				oldValue: reflect.ValueOf(1),
			},
			want: true,
		},
		{
			name: "",
			args: args{
				newValue: reflect.ValueOf(1),
				oldValue: reflect.ValueOf(2),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := equal(tt.args.newValue, tt.args.oldValue); got != tt.want {
				t.Errorf("equal() = %v, want %v", got, tt.want)
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
		equal(o, n)
	}
}

func TestExtractor_Entries(t *testing.T) {

	ts := makeBaseTestStruct()
	// ts := TestStruct{
	// 	String: "test",
	// }
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
			name: "remove pointer struct",
			modFn: func() {
				ts.SubPtr = nil
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key: "TestStruct",
						},
						{
							Key: "SubPtr",
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
