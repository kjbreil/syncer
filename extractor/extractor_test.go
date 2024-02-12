package extractor

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
	"testing"

	"github.com/kjbreil/syncer/control"
)

type TestStruct struct {
	String         string
	Int            int
	Slice          []int
	SliceStruct    []SD
	SlicePtr       []*int
	SlicePtrStruct []*SD
	Map            map[string]int
	MapStruct      map[string]TestStruct
	MapPtr         map[string]*int
	MapPtrStruct   map[string]*TestStruct
	Sub            TestSub
	SubPtr         *TestStruct
}

type TestSub struct {
	String string
}

type SD struct {
	Name string
	Data string
}

func TestNew(t *testing.T) {
	sliceEntries := 1

	ts := TestStruct{
		Slice: make([]int, 0, sliceEntries),
	}
	ext := New(ts)
	// add in slice entries
	for ii := 0; ii < sliceEntries; ii++ {
		ts.Slice = append(ts.Slice, ii)
	}
	ext.Diff(ts)
	fmt.Println("here")
}

func makeBaseTestStruct() TestStruct {
	return TestStruct{
		String: "Base String",
		Int:    1,
		Slice:  []int{1, 2, 3},
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
			String: "Base Test Sub",
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
			String: "Change Test Sub",
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

func TestExtractor_Diff(t *testing.T) {
	ts := TestStruct{
		Map: make(map[string]int),
	}
	ext := New(ts)

	tests := []struct {
		name    string
		addFunc func()
		want    []*control.Entry
	}{
		{
			name: "TestAddString",
			addFunc: func() {
				ts.String = "TestAddString"
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{Key: "TestStruct"},
						{Key: "String"},
					},
					Value: &control.Object{
						String_: control.MakePtr("TestAddString"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.addFunc()
			head, _ := ext.Diff(ts)
			molds := head.Entries()
			if len(molds) != len(tt.want) {
				t.Errorf("got molds length not matching expected")
			}
			for i := range molds {
				if !reflect.DeepEqual(molds[i], tt.want[i]) {
					t.Fatalf("mold not match: %v != %v", molds[i], tt.want[i])
				}
			}
			head, _ = ext.Diff(ts)
			molds = head.Entries()
			if len(molds) > 0 {
				t.Fatal("changes detected when they should not have been")
			}
		})
	}
}

func BenchmarkExtractor_Diff(b *testing.B) {
	sliceEntries := 100000
	cs := makeChangeTestStruct()
	ts := makeBaseTestStruct()

	ts.Slice = make([]int, 0, sliceEntries)
	for ii := 0; ii < sliceEntries; ii++ {
		ts.Slice = append(ts.Slice, ii)
	}
	ext := New(ts)
	for i := 0; i < b.N; i++ {
		_ = ext.Entries(&ts)
		_ = ext.Entries(&cs)
	}
}

func BenchmarkCompareToGob(b *testing.B) {
	sliceEntries := 100000

	ts := TestStruct{
		Slice: make([]int, 0, sliceEntries),
	}
	for i := 0; i < b.N; i++ {
		// add in slice entries
		for ii := 0; ii < sliceEntries; ii++ {
			ts.Slice = append(ts.Slice, ii)
		}

		var b bytes.Buffer
		err := gob.NewEncoder(&b).Encode(ts)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkExtractor_Diff2(b *testing.B) {
	sliceEntries := 100000

	ts := TestStruct{
		Slice: make([]int, 0, sliceEntries),
	}
	ext := New(ts)
	for i := 0; i < b.N; i++ {
		// add in slice entries
		for ii := 0; ii < sliceEntries; ii++ {
			ts.Slice = append(ts.Slice, ii)
		}
		head, _ := ext.Diff(&ts)
		_ = head.Entries()
		ts.Slice = make([]int, 0, sliceEntries)
		head, _ = ext.Diff(&ts)
		_ = head.Entries()
		for ii := 0; ii < sliceEntries; ii++ {
			ts.Slice = append(ts.Slice, ii)
		}
		head, _ = ext.Diff(&ts)
		_ = head.Entries()
	}
}

func BenchmarkExtractor_Diff3(b *testing.B) {
	ts := TestStruct{
		String: "Test",
	}
	ext := New(ts)
	for i := 0; i < b.N; i++ {
		head, _ := ext.Diff(&ts)
		_ = head.Entries()
		head, _ = ext.Diff(&ts)
		_ = head.Entries()
	}
}

// reflect: 5.58 ns/op
// stringEqual: 112.3 ns/op.
func Benchmark_equal1(b *testing.B) {
	var o, n reflect.Value
	for i := 0; i < b.N; i++ {
		o, n = reflect.ValueOf(i), reflect.ValueOf(i)
		equal(o, n)
	}
}

func Test_equal1(t *testing.T) {
	type args struct {
		n reflect.Value
		o reflect.Value
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				n: reflect.ValueOf(1),
				o: reflect.ValueOf(1),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := equal(tt.args.n, tt.args.o); got != tt.want {
				t.Errorf("equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractor_Entries(t *testing.T) {

	ts := makeBaseTestStruct()
	// ts := TestStruct{
	// 	String: "test",
	// }
	ext := New(ts)

	tests := []struct {
		name  string
		modFn func()
		want  []*control.Entry
	}{
		{
			name: "base extract all",
			modFn: func() {
				ext.Reset()
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "String",
							Index: &control.Object{},
						},
					},
					Value: &control.Object{String_: control.MakePtr("Base String")},
				},
				{
					Key: []*control.Key{
						{
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "Int",
							Index: &control.Object{},
						},
					},
					Value: &control.Object{Int64: control.MakePtr(int64(1))},
				},
				{
					Key: []*control.Key{
						{
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "Slice",
							Index: &control.Object{Int64: control.MakePtr(int64(0))},
						},
					},
					Value: &control.Object{Int64: control.MakePtr(int64(1))},
				},
				{
					Key: []*control.Key{
						{
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "Slice",
							Index: &control.Object{Int64: control.MakePtr(int64(1))},
						},
					},
					Value: &control.Object{Int64: control.MakePtr(int64(2))},
				},
				{
					Key: []*control.Key{
						{
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "Slice",
							Index: &control.Object{Int64: control.MakePtr(int64(2))},
						},
					},
					Value: &control.Object{Int64: control.MakePtr(int64(3))},
				},
				{
					Key: []*control.Key{
						{
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "SliceStruct",
							Index: &control.Object{},
						},
						{
							Key:   "Name",
							Index: &control.Object{},
						},
					},
					Value: &control.Object{String_: control.MakePtr("Base SliceStruct Name 1")},
				},
				{
					Key: []*control.Key{
						{
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "SliceStruct",
							Index: &control.Object{},
						},
						{
							Key:   "Data",
							Index: &control.Object{},
						},
					},
					Value: &control.Object{String_: control.MakePtr("Base SliceStruct Data 1")},
				},
				{
					Key: []*control.Key{
						{
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "SliceStruct",
							Index: &control.Object{},
						},
						{
							Key:   "Name",
							Index: &control.Object{},
						},
					},
					Value: &control.Object{String_: control.MakePtr("Base SliceStruct Name 2")},
				},
				{
					Key: []*control.Key{
						{
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "SliceStruct",
							Index: &control.Object{},
						},
						{
							Key:   "Data",
							Index: &control.Object{},
						},
					},
					Value: &control.Object{String_: control.MakePtr("Base SliceStruct Data 2")},
				},
				{
					Key: []*control.Key{
						{
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "SlicePtr",
							Index: &control.Object{Int64: control.MakePtr(int64(0))},
						},
					},
					Value: &control.Object{Int64: control.MakePtr(int64(1))},
				},
				{
					Key: []*control.Key{
						{
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "SlicePtr",
							Index: &control.Object{Int64: control.MakePtr(int64(1))},
						},
					},
					Value: &control.Object{Int64: control.MakePtr(int64(2))},
				},
				{
					Key: []*control.Key{
						{
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "SlicePtrStruct",
							Index: &control.Object{},
						},
						{
							Key:   "Name",
							Index: &control.Object{},
						},
					},
					Value: &control.Object{String_: control.MakePtr("Base SlicePtrStruct Name 1")},
				},
				{
					Key: []*control.Key{
						{
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "SlicePtrStruct",
							Index: &control.Object{},
						},
						{
							Key:   "Data",
							Index: &control.Object{},
						},
					},
					Value: &control.Object{String_: control.MakePtr("Base SlicePtrStruct Data 1")},
				},
				{
					Key: []*control.Key{
						{
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "SlicePtrStruct",
							Index: &control.Object{},
						},
						{
							Key:   "Name",
							Index: &control.Object{},
						},
					},
					Value: &control.Object{String_: control.MakePtr("Base SlicePtrStruct Name 2")},
				},
				{
					Key: []*control.Key{
						{
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "SlicePtrStruct",
							Index: &control.Object{},
						},
						{
							Key:   "Data",
							Index: &control.Object{},
						},
					},
					Value: &control.Object{String_: control.MakePtr("Base SlicePtrStruct Data 2")},
				},
				{
					Key: []*control.Key{
						{
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "Map",
							Index: &control.Object{String_: control.MakePtr("Base First")},
						},
					},
					Value: &control.Object{Int64: control.MakePtr(int64(1))},
				},
				{
					Key: []*control.Key{
						{
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "Map",
							Index: &control.Object{String_: control.MakePtr("Base Second")},
						},
					},
					Value: &control.Object{Int64: control.MakePtr(int64(2))},
				},
				{
					Key: []*control.Key{
						{
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "MapStruct",
							Index: &control.Object{},
						},
						{
							Key:   "String",
							Index: &control.Object{},
						},
					},
					Value: &control.Object{String_: control.MakePtr("Base First Test Struct")},
				},
				{
					Key: []*control.Key{
						{
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "MapPtr",
							Index: &control.Object{String_: control.MakePtr("Base First")},
						},
					},
					Value: &control.Object{Int64: control.MakePtr(int64(1))},
				},
				{
					Key: []*control.Key{
						{
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "MapPtrStruct",
							Index: &control.Object{},
						},
						{
							Key:   "String",
							Index: &control.Object{},
						},
					},
					Value: &control.Object{String_: control.MakePtr("Base First Test Struct")},
				},
				{
					Key: []*control.Key{
						{
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "Sub",
							Index: &control.Object{},
						},
						{
							Key:   "String",
							Index: &control.Object{},
						},
					},
					Value: &control.Object{String_: control.MakePtr("Base Test Sub")},
				},
				{
					Key: []*control.Key{
						{
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "Sub",
							Index: &control.Object{},
						},
						{
							Key:   "String",
							Index: &control.Object{},
						},
					},
					Value: &control.Object{String_: control.MakePtr("Base Test Sub")},
				},
				{
					Key: []*control.Key{
						{
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "SubPtr",
							Index: &control.Object{},
						},
						{
							Key:   "String",
							Index: &control.Object{},
						},
					},
					Value: &control.Object{String_: control.MakePtr("Base Sub Test Struct")},
				},
			},
		},
		{
			name: "change string",
			modFn: func() {
				ts.String = "change string"
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "String",
							Index: &control.Object{},
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
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "Slice",
							Index: &control.Object{Int64: control.MakePtr(int64(3))},
						},
					},
					Value: &control.Object{Int64: control.MakePtr(int64(4))},
				},
			},
		},
		{
			name: "remove from int slice",
			modFn: func() {
				ts.Slice = ts.Slice[:len(ts.Slice)-1]
			},
			want: []*control.Entry{
				{
					Key: []*control.Key{
						{
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "Slice",
							Index: &control.Object{Int64: control.MakePtr(int64(2))},
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
							Key:   "TestStruct",
							Index: &control.Object{},
						},
						{
							Key:   "SubPtr",
							Index: &control.Object{},
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
			ext.Entries(&ts)
			tt.modFn()

			got := ext.Entries(&ts)
			if !got.Equals(tt.want) {
				t.Errorf("Entries() = %v, want %v", got, tt.want)
				t.Logf("##########\n\n%s\n\n##########", got.Struct())
			}
		})
	}
}
