package extractor

import (
	"fmt"
	"github.com/kjbreil/syncer/control"
	"reflect"
	"testing"
)

type testStruct struct {
	String         string
	Int            int
	Slice          []int
	SliceStruct    []sd
	SlicePtr       []*int
	SlicePtrStruct []*sd
	Map            map[string]int
	MapStruct      map[string]testStruct
	MapPtr         map[string]*int
	MapPtrStruct   map[string]*testStruct
	Sub            testSub
	SubPtr         *testStruct
}

type testSub struct {
	String string
}

type sd struct {
	name string
	data string
}

func TestNew2(t *testing.T) {
	sliceEntries := 1

	ts := testStruct{
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

func TestNew(t *testing.T) {
	// one := 1

	ts := testStruct{
		String: "Test",
		// Int:    1,
		Slice: []int{1, 2},
		// SliceStruct: []sd{
		// 	{
		// 		name: "testName1",
		// 		data: "testData1",
		// 	},
		// 	{
		// 		name: "testName2",
		// 		data: "testData2",
		// 	},
		// },
		// SlicePtr: []*int{
		// 	&one,
		// },
		// SlicePtrStruct: []*sd{
		// 	{
		// 		name: "testName1",
		// 		data: "testData1",
		// 	},
		// 	{
		// 		name: "testName2",
		// 		data: "testData2",
		// 	},
		// },
		// Map: map[string]int{
		// 	"test1": 1,
		// 	"test2": 2,
		// },
		// MapPtr: map[string]*int{
		// 	"test1": &one,
		// },
		// MapStruct: map[string]testStruct{
		// 	"test1": {
		// 		String: "test1data",
		// 	},
		// },
		// MapPtrStruct: map[string]*testStruct{
		// 	"test1": {
		// 		String: "MapPtrStructTest1Data",
		// 	},
		// },
		// Sub: testSub{
		// 	String: "Sub1",
		// },
		// SubPtr: &testStruct{
		// 	String: "SubTest",
		// 	// Slice:  nil,
		// 	// Map:    nil,
		// },
	}

	ext := New(ts)

	head := ext.Diff(ts)
	moulds := head.Entries()

	fmt.Println(moulds)

	ts.Slice[0] = 2
	// delete(ts.Map, "test2")
	head = ext.Diff(ts)
	moulds = head.Entries()

	fmt.Println(moulds)
	ts.Slice[0] = 3

	head = ext.Diff(ts)
	moulds = head.Entries()

	fmt.Println(moulds)
	// 	var moulds []Diff
	// 	moulds = ext.Diff(ts)
	// 	fmt.Println(len(moulds))
	// 	// ts.Slice = append(ts.Slice, sd{
	// 	// 	name: "test3",
	// 	// 	data: "test3",
	// 	// })
	// 	//
	// 	// ts.Map["M2"] = testStruct{
	// 	// 	String: "M2",
	// 	// 	Int:    2,
	// 	// }
	//
	// 	ts.Sub.String = "TestSubUpdate"
	//
	// 	moulds = ext.Diff(ts)
	// 	fmt.Println(len(moulds))
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
	ts := testStruct{
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
						{Key: "testStruct"},
						{Key: "String"},
					},
					Value: &control.Object{
						String_: stringPtr("TestAddString"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.addFunc()
			head := ext.Diff(ts)
			moulds := head.Entries()
			if len(moulds) != len(tt.want) {
				t.Errorf("got moulds length not matching expected")
			}
			for i := range moulds {
				if !reflect.DeepEqual(moulds[i], tt.want[i]) {
					t.Fatalf("mould not match: %v != %v", moulds[i], tt.want[i])
				}
			}
			head = ext.Diff(ts)
			moulds = head.Entries()
			if len(moulds) > 0 {
				t.Fatal("changes detected when they should not have been")
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

// string equal: 175723234
// reflect equal: 90242146
// disable GC: 85727556
//

func BenchmarkExtractor_Diff(b *testing.B) {
	sliceEntries := 100000

	ts := testStruct{
		Slice: make([]int, 0, sliceEntries),
	}
	ext := New(ts)
	for i := 0; i < b.N; i++ {
		// add in slice entries
		for ii := 0; ii < sliceEntries; ii++ {
			ts.Slice = append(ts.Slice, ii)
		}
		head := ext.Diff(ts)
		_ = head.Entries()
		head = ext.Diff(ts)
		_ = head.Entries()
	}
}

func BenchmarkExtractor_Diff2(b *testing.B) {
	sliceEntries := 100000

	ts := testStruct{
		Slice: make([]int, 0, sliceEntries),
	}
	ext := New(ts)
	for i := 0; i < b.N; i++ {
		// add in slice entries
		for ii := 0; ii < sliceEntries; ii++ {
			ts.Slice = append(ts.Slice, ii)
		}
		head := ext.Diff(ts)
		_ = head.Entries()
		ts.Slice = make([]int, 0, sliceEntries)
		head = ext.Diff(ts)
		_ = head.Entries()
		for ii := 0; ii < sliceEntries; ii++ {
			ts.Slice = append(ts.Slice, ii)
		}
		head = ext.Diff(ts)
		_ = head.Entries()
	}
}

// reflect: 5.58 ns/op
// stringEqual: 112.3 ns/op
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
