package equal

import (
	"reflect"
	"testing"

	"github.com/kjbreil/syncer/pkg/test"
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

func BenchmarkEqualPrimitive(b *testing.B) {
	val1 := reflect.ValueOf(42)
	val2 := reflect.ValueOf(42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(val1, val2)
	}
}

func BenchmarkEqualPrimitiveDifferent(b *testing.B) {
	val1 := reflect.ValueOf(42)
	val2 := reflect.ValueOf(43)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(val1, val2)
	}
}

func BenchmarkEqualString(b *testing.B) {
	val1 := reflect.ValueOf("test string for benchmarking performance")
	val2 := reflect.ValueOf("test string for benchmarking performance")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(val1, val2)
	}
}

func BenchmarkEqualStringDifferent(b *testing.B) {
	val1 := reflect.ValueOf("test string for benchmarking performance")
	val2 := reflect.ValueOf("different test string for benchmarking")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(val1, val2)
	}
}

func BenchmarkEqualSimpleStruct(b *testing.B) {
	data1 := simpleData{
		String: "test string",
		Int:    42,
		Bool:   true,
		Float:  3.14,
	}
	data2 := simpleData{
		String: "test string",
		Int:    42,
		Bool:   true,
		Float:  3.14,
	}
	val1 := reflect.ValueOf(data1)
	val2 := reflect.ValueOf(data2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(val1, val2)
	}
}

func BenchmarkEqualSimpleStructDifferent(b *testing.B) {
	data1 := simpleData{
		String: "test string",
		Int:    42,
		Bool:   true,
		Float:  3.14,
	}
	data2 := simpleData{
		String: "different string",
		Int:    42,
		Bool:   true,
		Float:  3.14,
	}
	val1 := reflect.ValueOf(data1)
	val2 := reflect.ValueOf(data2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(val1, val2)
	}
}

func BenchmarkEqualSliceInt(b *testing.B) {
	slice1 := make([]int, 100)
	slice2 := make([]int, 100)
	for i := range slice1 {
		slice1[i] = i
		slice2[i] = i
	}
	val1 := reflect.ValueOf(slice1)
	val2 := reflect.ValueOf(slice2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(val1, val2)
	}
}

func BenchmarkEqualSliceIntDifferent(b *testing.B) {
	slice1 := make([]int, 100)
	slice2 := make([]int, 100)
	for i := range slice1 {
		slice1[i] = i
		slice2[i] = i + 1
	}
	val1 := reflect.ValueOf(slice1)
	val2 := reflect.ValueOf(slice2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(val1, val2)
	}
}

func BenchmarkEqualSliceStruct(b *testing.B) {
	slice1 := make([]simpleData, 10)
	slice2 := make([]simpleData, 10)
	for i := range slice1 {
		data := simpleData{
			String: "test",
			Int:    i,
			Bool:   i%2 == 0,
			Float:  float64(i) * 1.5,
		}
		slice1[i] = data
		slice2[i] = data
	}
	val1 := reflect.ValueOf(slice1)
	val2 := reflect.ValueOf(slice2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(val1, val2)
	}
}

func BenchmarkEqualMap(b *testing.B) {
	m1 := make(map[string]int, 50)
	m2 := make(map[string]int, 50)
	for i := 0; i < 50; i++ {
		key := string(rune('a' + i))
		m1[key] = i
		m2[key] = i
	}
	val1 := reflect.ValueOf(m1)
	val2 := reflect.ValueOf(m2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(val1, val2)
	}
}

func BenchmarkEqualMapDifferent(b *testing.B) {
	m1 := make(map[string]int, 50)
	m2 := make(map[string]int, 50)
	for i := 0; i < 50; i++ {
		key := string(rune('a' + i))
		m1[key] = i
		m2[key] = i + 1
	}
	val1 := reflect.ValueOf(m1)
	val2 := reflect.ValueOf(m2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(val1, val2)
	}
}

func BenchmarkEqualComplexStruct(b *testing.B) {
	subData := simpleData{
		String: "nested",
		Int:    100,
		Bool:   true,
		Float:  2.71,
	}
	ptrData := &simpleData{
		String: "pointer data",
		Int:    200,
		Bool:   false,
		Float:  1.41,
	}
	data1 := complexData{
		String: "complex test string",
		Int:    42,
		Slice:  []int{1, 2, 3, 4, 5},
		Map: map[string]int{
			"one":   1,
			"two":   2,
			"three": 3,
		},
		SubStruct: subData,
		Pointer:   ptrData,
	}
	data2 := complexData{
		String: "complex test string",
		Int:    42,
		Slice:  []int{1, 2, 3, 4, 5},
		Map: map[string]int{
			"one":   1,
			"two":   2,
			"three": 3,
		},
		SubStruct: subData,
		Pointer:   ptrData,
	}
	val1 := reflect.ValueOf(data1)
	val2 := reflect.ValueOf(data2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(val1, val2)
	}
}

func BenchmarkEqualTestStruct(b *testing.B) {
	data1 := test.MakeBaseTestStruct()
	data2 := test.MakeBaseTestStruct()
	val1 := reflect.ValueOf(data1)
	val2 := reflect.ValueOf(data2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(val1, val2)
	}
}

func BenchmarkEqualTestStructDifferent(b *testing.B) {
	data1 := test.MakeBaseTestStruct()
	data2 := test.MakeChangeTestStruct()
	val1 := reflect.ValueOf(data1)
	val2 := reflect.ValueOf(data2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(val1, val2)
	}
}

func BenchmarkEqualPointer(b *testing.B) {
	data1 := &simpleData{
		String: "pointer test",
		Int:    42,
		Bool:   true,
		Float:  3.14,
	}
	data2 := &simpleData{
		String: "pointer test",
		Int:    42,
		Bool:   true,
		Float:  3.14,
	}
	val1 := reflect.ValueOf(data1)
	val2 := reflect.ValueOf(data2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(val1, val2)
	}
}

func BenchmarkEqualArray(b *testing.B) {
	var arr1, arr2 [10]simpleData
	for i := range arr1 {
		data := simpleData{
			String: "array element",
			Int:    i,
			Bool:   i%2 == 0,
			Float:  float64(i),
		}
		arr1[i] = data
		arr2[i] = data
	}
	val1 := reflect.ValueOf(arr1)
	val2 := reflect.ValueOf(arr2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(val1, val2)
	}
}

func BenchmarkAnySimpleStruct(b *testing.B) {
	data1 := simpleData{
		String: "test string",
		Int:    42,
		Bool:   true,
		Float:  3.14,
	}
	data2 := simpleData{
		String: "test string",
		Int:    42,
		Bool:   true,
		Float:  3.14,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Any(data1, data2)
	}
}

func BenchmarkAnyComplexStruct(b *testing.B) {
	subData := simpleData{
		String: "nested",
		Int:    100,
		Bool:   true,
		Float:  2.71,
	}
	ptrData := &simpleData{
		String: "pointer data",
		Int:    200,
		Bool:   false,
		Float:  1.41,
	}
	data1 := complexData{
		String: "complex test string",
		Int:    42,
		Slice:  []int{1, 2, 3, 4, 5},
		Map: map[string]int{
			"one":   1,
			"two":   2,
			"three": 3,
		},
		SubStruct: subData,
		Pointer:   ptrData,
	}
	data2 := complexData{
		String: "complex test string",
		Int:    42,
		Slice:  []int{1, 2, 3, 4, 5},
		Map: map[string]int{
			"one":   1,
			"two":   2,
			"three": 3,
		},
		SubStruct: subData,
		Pointer:   ptrData,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Any(data1, data2)
	}
}

func BenchmarkAnyTestStruct(b *testing.B) {
	data1 := test.MakeBaseTestStruct()
	data2 := test.MakeBaseTestStruct()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Any(data1, data2)
	}
}

// Benchmark different sized slices
func BenchmarkEqualSliceSmall(b *testing.B) {
	slice1 := make([]int, 10)
	slice2 := make([]int, 10)
	for i := range slice1 {
		slice1[i] = i
		slice2[i] = i
	}
	val1 := reflect.ValueOf(slice1)
	val2 := reflect.ValueOf(slice2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(val1, val2)
	}
}

func BenchmarkEqualSliceMedium(b *testing.B) {
	slice1 := make([]int, 1000)
	slice2 := make([]int, 1000)
	for i := range slice1 {
		slice1[i] = i
		slice2[i] = i
	}
	val1 := reflect.ValueOf(slice1)
	val2 := reflect.ValueOf(slice2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(val1, val2)
	}
}

func BenchmarkEqualSliceLarge(b *testing.B) {
	slice1 := make([]int, 10000)
	slice2 := make([]int, 10000)
	for i := range slice1 {
		slice1[i] = i
		slice2[i] = i
	}
	val1 := reflect.ValueOf(slice1)
	val2 := reflect.ValueOf(slice2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(val1, val2)
	}
}

// Benchmark different sized maps
func BenchmarkEqualMapSmall(b *testing.B) {
	m1 := make(map[string]int, 10)
	m2 := make(map[string]int, 10)
	for i := 0; i < 10; i++ {
		key := string(rune('a' + i))
		m1[key] = i
		m2[key] = i
	}
	val1 := reflect.ValueOf(m1)
	val2 := reflect.ValueOf(m2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(val1, val2)
	}
}

func BenchmarkEqualMapMedium(b *testing.B) {
	m1 := make(map[string]int, 100)
	m2 := make(map[string]int, 100)
	for i := 0; i < 100; i++ {
		key := string(rune('a'+(i%26))) + string(rune('0'+(i/26)))
		m1[key] = i
		m2[key] = i
	}
	val1 := reflect.ValueOf(m1)
	val2 := reflect.ValueOf(m2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(val1, val2)
	}
}

func BenchmarkEqualMapLarge(b *testing.B) {
	m1 := make(map[string]int, 1000)
	m2 := make(map[string]int, 1000)
	for i := 0; i < 1000; i++ {
		key := string(rune('a'+(i%26))) + string(rune('0'+((i/26)%10))) + string(rune('0'+((i/260)%10)))
		m1[key] = i
		m2[key] = i
	}
	val1 := reflect.ValueOf(m1)
	val2 := reflect.ValueOf(m2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(val1, val2)
	}
}

// Early exit performance tests
func BenchmarkEqualSliceEarlyExit(b *testing.B) {
	slice1 := make([]int, 1000)
	slice2 := make([]int, 1000)
	for i := range slice1 {
		slice1[i] = i
		slice2[i] = i
	}
	// Make first element different for early exit
	slice2[0] = 999
	val1 := reflect.ValueOf(slice1)
	val2 := reflect.ValueOf(slice2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(val1, val2)
	}
}

func BenchmarkEqualStructEarlyExit(b *testing.B) {
	data1 := complexData{
		String: "complex test string",
		Int:    42,
		Slice:  []int{1, 2, 3, 4, 5},
		Map: map[string]int{
			"one":   1,
			"two":   2,
			"three": 3,
		},
	}
	data2 := complexData{
		String: "different string", // Different first field for early exit
		Int:    42,
		Slice:  []int{1, 2, 3, 4, 5},
		Map: map[string]int{
			"one":   1,
			"two":   2,
			"three": 3,
		},
	}
	val1 := reflect.ValueOf(data1)
	val2 := reflect.ValueOf(data2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(val1, val2)
	}
}