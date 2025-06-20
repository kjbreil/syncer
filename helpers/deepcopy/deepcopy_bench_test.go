package deepcopy

import (
	"reflect"
	"testing"

	"github.com/kjbreil/syncer/helpers/test"
)

// Simple test data for basic benchmarks
type simpleData struct {
	String string
	Int    int
	Bool   bool
	Float  float64
}

type complexData struct {
	String    string
	Int       int
	Slice     []int
	Map       map[string]int
	SubStruct simpleData
	Pointer   *simpleData
}

func BenchmarkDeepCopyPrimitive(b *testing.B) {
	val := reflect.ValueOf(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DeepCopy(val)
	}
}

func BenchmarkDeepCopyString(b *testing.B) {
	val := reflect.ValueOf("test string for benchmarking performance")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DeepCopy(val)
	}
}

func BenchmarkDeepCopySimpleStruct(b *testing.B) {
	data := simpleData{
		String: "test string",
		Int:    42,
		Bool:   true,
		Float:  3.14,
	}
	val := reflect.ValueOf(data)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DeepCopy(val)
	}
}

func BenchmarkDeepCopySliceInt(b *testing.B) {
	slice := make([]int, 100)
	for i := range slice {
		slice[i] = i
	}
	val := reflect.ValueOf(slice)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DeepCopy(val)
	}
}

func BenchmarkDeepCopySliceStruct(b *testing.B) {
	slice := make([]simpleData, 10)
	for i := range slice {
		slice[i] = simpleData{
			String: "test",
			Int:    i,
			Bool:   i%2 == 0,
			Float:  float64(i) * 1.5,
		}
	}
	val := reflect.ValueOf(slice)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DeepCopy(val)
	}
}

func BenchmarkDeepCopyMap(b *testing.B) {
	m := make(map[string]int, 50)
	for i := 0; i < 50; i++ {
		m[string(rune('a'+i))] = i
	}
	val := reflect.ValueOf(m)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DeepCopy(val)
	}
}

func BenchmarkDeepCopyComplexStruct(b *testing.B) {
	data := complexData{
		String: "complex test string",
		Int:    42,
		Slice:  []int{1, 2, 3, 4, 5},
		Map: map[string]int{
			"one":   1,
			"two":   2,
			"three": 3,
		},
		SubStruct: simpleData{
			String: "nested",
			Int:    100,
			Bool:   true,
			Float:  2.71,
		},
		Pointer: &simpleData{
			String: "pointer data",
			Int:    200,
			Bool:   false,
			Float:  1.41,
		},
	}
	val := reflect.ValueOf(data)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DeepCopy(val)
	}
}

func BenchmarkDeepCopyTestStruct(b *testing.B) {
	data := test.MakeBaseTestStruct()
	val := reflect.ValueOf(data)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DeepCopy(val)
	}
}

func BenchmarkDeepCopyPointer(b *testing.B) {
	data := &simpleData{
		String: "pointer test",
		Int:    42,
		Bool:   true,
		Float:  3.14,
	}
	val := reflect.ValueOf(data)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DeepCopy(val)
	}
}

func BenchmarkDeepCopyArray(b *testing.B) {
	var arr [10]simpleData
	for i := range arr {
		arr[i] = simpleData{
			String: "array element",
			Int:    i,
			Bool:   i%2 == 0,
			Float:  float64(i),
		}
	}
	val := reflect.ValueOf(arr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DeepCopy(val)
	}
}

func BenchmarkAnySimpleStruct(b *testing.B) {
	data := simpleData{
		String: "test string",
		Int:    42,
		Bool:   true,
		Float:  3.14,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Any(data)
	}
}

func BenchmarkAnyComplexStruct(b *testing.B) {
	data := complexData{
		String: "complex test string",
		Int:    42,
		Slice:  []int{1, 2, 3, 4, 5},
		Map: map[string]int{
			"one":   1,
			"two":   2,
			"three": 3,
		},
		SubStruct: simpleData{
			String: "nested",
			Int:    100,
			Bool:   true,
			Float:  2.71,
		},
		Pointer: &simpleData{
			String: "pointer data",
			Int:    200,
			Bool:   false,
			Float:  1.41,
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Any(data)
	}
}

func BenchmarkAnyTestStruct(b *testing.B) {
	data := test.MakeBaseTestStruct()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Any(data)
	}
}

// Benchmark different sized slices
func BenchmarkDeepCopySliceSmall(b *testing.B) {
	slice := make([]int, 10)
	for i := range slice {
		slice[i] = i
	}
	val := reflect.ValueOf(slice)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DeepCopy(val)
	}
}

func BenchmarkDeepCopySliceMedium(b *testing.B) {
	slice := make([]int, 1000)
	for i := range slice {
		slice[i] = i
	}
	val := reflect.ValueOf(slice)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DeepCopy(val)
	}
}

func BenchmarkDeepCopySliceLarge(b *testing.B) {
	slice := make([]int, 10000)
	for i := range slice {
		slice[i] = i
	}
	val := reflect.ValueOf(slice)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DeepCopy(val)
	}
}

// Benchmark different sized maps
func BenchmarkDeepCopyMapSmall(b *testing.B) {
	m := make(map[string]int, 10)
	for i := 0; i < 10; i++ {
		m[string(rune('a'+i))] = i
	}
	val := reflect.ValueOf(m)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DeepCopy(val)
	}
}

func BenchmarkDeepCopyMapMedium(b *testing.B) {
	m := make(map[string]int, 100)
	for i := 0; i < 100; i++ {
		m[string(rune('a'+(i%26)))+string(rune('0'+(i/26)))] = i
	}
	val := reflect.ValueOf(m)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DeepCopy(val)
	}
}

func BenchmarkDeepCopyMapLarge(b *testing.B) {
	m := make(map[string]int, 1000)
	for i := 0; i < 1000; i++ {
		m[string(rune('a'+(i%26)))+string(rune('0'+((i/26)%10)))+string(rune('0'+((i/260)%10)))] = i
	}
	val := reflect.ValueOf(m)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DeepCopy(val)
	}
}