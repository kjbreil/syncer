package extractor

import (
	"errors"
	"reflect"
	"sync"

	"github.com/kjbreil/syncer/helpers/deepcopy"
)

type Extractor struct {
	data any
	mut  *sync.Mutex
}

var (
	ErrNotPointer      = errors.New("data is not a pointer")
	ErrUnsupportedType = errors.New("unsupported type")
)

const (
	historySize = 100
)

// New creates a new instance of the Extractor struct.
//
// data: the data to be extracted from
//
// Returns:
// *Extractor: a new instance of the Extractor struct.
func New(data any) (*Extractor, error) {
	if data == nil {
		return nil, errors.New("data is nil")
	}
	t := reflect.Indirect(reflect.ValueOf(data)).Type()
	dataStruct := reflect.New(t)
	aStruct := deepcopy.Any(dataStruct.Interface())
	return &Extractor{
		data: aStruct,
		mut:  new(sync.Mutex),
	}, nil
}

// Reset resets the data to its initial state.
func (ext *Extractor) Reset() {
	ext.mut.Lock()
	defer ext.mut.Unlock()

	if ext.data == nil {
		return
	}

	t := reflect.TypeOf(ext.data)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	dataStruct := reflect.New(t)
	aStruct := dataStruct.Interface()

	ext.data = aStruct
}
