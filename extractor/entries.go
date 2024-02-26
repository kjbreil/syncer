package extractor

import (
	"errors"
	"fmt"
	"github.com/kjbreil/syncer/control"
	"github.com/kjbreil/syncer/helpers/deepcopy"
	"reflect"
)

type extFn func(newValue, oldValue reflect.Value, upperType reflect.StructField, level int) (control.Entries, error)

var extFns map[reflect.Kind]extFn

func init() {
	extFns = map[reflect.Kind]extFn{
		reflect.Bool:       extractPrimitive,
		reflect.Int:        extractPrimitive,
		reflect.Int8:       extractPrimitive,
		reflect.Int16:      extractPrimitive,
		reflect.Int32:      extractPrimitive,
		reflect.Int64:      extractPrimitive,
		reflect.Uint:       extractPrimitive,
		reflect.Uint8:      extractPrimitive,
		reflect.Uint16:     extractPrimitive,
		reflect.Uint32:     extractPrimitive,
		reflect.Uint64:     extractPrimitive,
		reflect.Uintptr:    extractPrimitive,
		reflect.Float32:    extractPrimitive,
		reflect.Float64:    extractPrimitive,
		reflect.Complex64:  extractPrimitive,
		reflect.Complex128: extractPrimitive,
		reflect.String:     extractPrimitive,
		reflect.Struct:     extractStruct,
		// reflect.Array:      extractArray,
		// reflect.Map:        extractMap,
		// reflect.Ptr:        extractPointer,
		reflect.Interface: extractInterface,
		// reflect.Slice:      extractSlice,
	}
}

func (ext *Extractor) GetDiff(data any) (control.Entries, error) {
	ext.mut.Lock()
	defer ext.mut.Unlock()

	// deep copy the current data as a point in time
	pitData := deepcopy.Any(data)

	newValue := reflect.ValueOf(pitData)
	if newValue.Kind() != reflect.Ptr {
		return nil, ErrNotPointer
	}

	newValue = reflect.Indirect(newValue)
	// newValue = reflect.ValueOf(newValue.Interface())

	oldValue := reflect.Indirect(reflect.ValueOf(ext.data))

	entries, err := extract(newValue, oldValue, reflect.StructField{
		Name: newValue.Type().Name(),
	}, 0)

	if err != nil && !errors.Is(err, ErrUnsupportedType) {
		return nil, err
	}

	ext.data = pitData

	return entries, nil
}

func extract(newValue reflect.Value, oldValue reflect.Value, upperType reflect.StructField, level int) (control.Entries, error) {
	if iFn, ok := extFns[newValue.Kind()]; ok {
		head, err := iFn(newValue, oldValue, upperType, level)
		if err != nil {
			return nil, err
		}
		head.AddKey(upperType.Name)
		return head, nil
	}
	return nil, fmt.Errorf("%w: %s", ErrUnsupportedType, newValue.Kind())
}
