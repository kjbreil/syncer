package extractor

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/kjbreil/syncer/pkg/control"
	"github.com/kjbreil/syncer/pkg/deepcopy"
)

type extFn func(newValue, oldValue reflect.Value, upperType reflect.StructField, level int) (control.Entries, error)

// extFns is a map of reflect.Kind to their respective extraction function.
var extFns map[reflect.Kind]extFn

func init() {
	extFns = map[reflect.Kind]extFn{
		reflect.Invalid:       extractInvalid,
		reflect.Bool:          extractPrimitive,
		reflect.Int:           extractPrimitive,
		reflect.Int8:          extractPrimitive,
		reflect.Int16:         extractPrimitive,
		reflect.Int32:         extractPrimitive,
		reflect.Int64:         extractPrimitive,
		reflect.Uint:          extractPrimitive,
		reflect.Uint8:         extractPrimitive,
		reflect.Uint16:        extractPrimitive,
		reflect.Uint32:        extractPrimitive,
		reflect.Uint64:        extractPrimitive,
		reflect.Uintptr:       extractPrimitive,
		reflect.Float32:       extractPrimitive,
		reflect.Float64:       extractPrimitive,
		reflect.Complex64:     extractPrimitive,
		reflect.Complex128:    extractPrimitive,
		reflect.Array:         extractArray,
		reflect.Chan:          extractUnsupported,
		reflect.Func:          extractUnsupported,
		reflect.Interface:     extractInterface,
		reflect.Map:           extractMap,
		reflect.Pointer:       extractPointer,
		reflect.Slice:         extractSlice,
		reflect.String:        extractPrimitive,
		reflect.Struct:        extractStruct,
		reflect.UnsafePointer: extractUnsupported,
	}
}

// Entries returns a list of changes between the current and previous states of the data.
// The previous state is taken at the time of the last call to Entries.
// If the data is not a pointer, an error is returned.
// The returned list of changes is a tree structure, where each node represents a single change.
// The node's key is the name of the field that was changed, and its value is a list of changes
// to the field's value.
// If the value of a field is a struct, array, or map, Entries is recursively called on that value.
// If the value of a field is a pointer, and the pointer is not nil, Entries is recursively called
// on the value that the pointer points to.
// If the value of a field is a slice, Entries is only called on the first N elements (where N is
// the length of the slice), as slices are considered immutable.
// If the value of a field is an interface, and the interface is not nil, Entries is recursively
// called on the value that the interface contains.
// If the value of a field is unsupported, an error is returned.
// The returned list of changes is thread-safe and can be modified concurrently.
func (ext *Extractor) Entries(data any) (control.Entries, error) {
	ext.mut.Lock()
	defer ext.mut.Unlock()

	if data == nil {
		return nil, errors.New("data is nil")
	}

	// deep copy the current data as a point in time
	pitData := deepcopy.Any(data)

	// check if ext.data is nil before proceeding
	if ext.data != nil {
		// get the reflect values of the current and previous states
		newValue := reflect.ValueOf(pitData)
		if newValue.Kind() != reflect.Ptr {
			return nil, ErrNotPointer
		}
		newValue = reflect.Indirect(newValue)

		oldValue := reflect.Indirect(reflect.ValueOf(ext.data))

		// recursively extract the changes between the current and previous states
		entries, err := extract(newValue, oldValue, reflect.StructField{
			Name: newValue.Type().Name(),
		}, 0, true)

		if err != nil && !errors.Is(err, ErrUnsupportedType) {
			return nil, err
		}

		// set the current state to the point in time data
		ext.data = pitData

		return entries, nil
	}

	return nil, nil
}

// extract recursively compares the current and previous states of a value and returns a list of
// changes between them.
func extract(newValue, oldValue reflect.Value, upperType reflect.StructField, level int, makeKey bool) (control.Entries, error) {
	if !oldValue.IsValid() {
		oldValue = reflect.New(newValue.Type()).Elem()
	}

	if iFn, ok := extFns[newValue.Kind()]; ok {
		// if the value kind has a registered extraction function, use it
		head, err := iFn(newValue, oldValue, upperType, level)
		if err != nil {
			return nil, err
		}
		if !makeKey {
			return head, nil
		}
		// add the field name as a key to the list of changes
		head.AddKey(upperType.Name)
		return head, nil
	}

	// if the value kind does not have a registered extraction function, return an error
	return nil, fmt.Errorf("%w: %s", ErrUnsupportedType, newValue.Kind())
}
