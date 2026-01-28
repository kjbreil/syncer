package test

type TestStruct struct {
	String         string
	Int            int
	Int8           int8
	Int16          int16
	Int32          int32
	Int64          int64
	Uint           uint
	Uint8          uint8
	Uint16         uint16
	Uint32         uint32
	Uint64         uint64
	Float32        float32
	Float64        float64
	Complex64      complex64
	Complex128     complex128
	Bool           bool
	Byte           byte
	Bytes          []byte
	Interface      TestInterface
	Slice          []int
	SliceStruct    []SD
	SlicePtr       []*int
	SlicePtrStruct []*TestStruct
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
	MapKeyBool     map[bool]int
	MapKeyUint     map[uint]int
	MapKeyFloat    map[float64]int
	MapKeyFloat32  map[float32]int
	MapStruct      map[string]TestStruct
	MapMapType     map[string]MapType
	MapPtr         map[string]*int
	MapPtrStruct   map[string]*TestStruct
	MapInterface   map[string]TestInterface
	MapMap         map[string]map[string]int
	MapSlice       map[string][]int
	MapStructSlice map[int64]TestSub
	SubStruct      TestSub
	SubStructPtr   *TestStruct
	unexported     string
	Function       func()
}

type MapType map[string]TestSub

type TestSub struct {
	MapPtrStruct map[int64]*SD
	Slice        []SD
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

type Tool int

const (
	ToolDns Tool = iota
	ToolDhcp
)
