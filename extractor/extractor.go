package extractor

import (
	"errors"
	"github.com/kjbreil/syncer/control"
	"reflect"
	"sync"
)

type Extractor struct {
	data    any
	history []*control.Diff
	mut     *sync.Mutex
}

var (
	ErrNotPointer         = errors.New("data is not a pointer")
	ErrDataStructMisMatch = errors.New("data structs do not match")
	ErrUnsupportedType    = errors.New("unsupported type")
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
func New(data any) *Extractor {
	// Get the type of the data
	t := reflect.TypeOf(data)
	// Iterate through pointer types until we reach the base type
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Create a new value of the base type
	dataStruct := reflect.New(t)
	// Convert the new value to an interface{} so we can access it
	aStruct := dataStruct.Interface()
	// Create a new instance of the Extractor struct and return it
	return &Extractor{
		data:    aStruct,
		history: make([]*control.Diff, 0, historySize),
		mut:     &sync.Mutex{},
	}
}

// Reset resets the data to its initial state.
func (ext *Extractor) Reset() {
	t := reflect.TypeOf(ext.data)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	dataStruct := reflect.New(t)
	aStruct := dataStruct.Interface()

	ext.data = aStruct
}

func (ext *Extractor) Entries(data any) control.Entries {
	h, err := ext.Diff(data)
	if err != nil {
		panic(err)
	}
	return h.Entries()
}

func (ext *Extractor) Diff(data any) (*control.Diff, error) {
	// force single threaded access
	ext.mut.Lock()
	defer ext.mut.Unlock()

	newValue := reflect.ValueOf(data)
	if newValue.Kind() != reflect.Ptr {
		return nil, ErrNotPointer
	}

	// follow the pointer to the actual value
	newValue = reflect.Indirect(newValue)
	oldValue := reflect.Indirect(reflect.ValueOf(ext.data))

	head := extractObject(newValue, oldValue, newValue.Type().Name())

	head.Timestamp()

	return head, nil
}

func extractObject(newValue, oldValue reflect.Value, keyName string) *control.Diff {

	// TODO: This should check if oldValue is valid and return a delete if it is
	if !newValue.IsValid() {
		return nil
	}
	// check if the oldValue is valid (exists) and create it if it does not
	if !oldValue.IsValid() {
		oldValue = reflect.New(newValue.Type()).Elem()
	}

	current := control.NewDiff(&control.Key{Key: keyName})

	// indirect old and new values
	newValue = reflect.Indirect(newValue)
	oldValue = reflect.Indirect(oldValue)
	// newValue, oldValue = indirect(newValue, oldValue)

	// loop over the fields of the newValue finding the relevant matching field in the old value
	numFields := newValue.NumField()
	for i := 0; i < numFields; i++ {
		newValueFieldKind := newValue.Field(i).Kind()

		switch newValueFieldKind {
		case reflect.Struct, reflect.Interface:
			child := extractObject(newValue.Field(i), oldValue.Field(i), newValue.Type().Field(i).Name)
			if child != nil {
				current.Children = append(current.Children, child)
			}
		case reflect.Pointer:
			child := extractObjectPtr(newValue.Field(i), oldValue.Field(i), newValue.Type().Field(i).Name)
			if child != nil {
				current.Children = append(current.Children, child)
			}
		case reflect.Map:
			children := extractMap(newValue.Field(i), oldValue.Field(i), oldValue.Type().Field(i).Type, newValue.Type().Field(i).Name)
			if len(children) > 0 {
				current.Children = append(current.Children, children...)
			}
		case reflect.Slice, reflect.Array:
			children := extractSlice(newValue.Field(i), oldValue.Field(i), newValue.Type().Field(i).Name)
			if len(children) > 0 {
				current.Children = append(current.Children, children...)
			}
		case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			child := extractConcrete(newValue.Field(i), oldValue.Field(i), newValue.Type().Field(i).Name)
			if child != nil {
				current.Children = append(current.Children, child)
			}
		case reflect.Invalid:
		case reflect.Uintptr:
		case reflect.Complex64, reflect.Complex128:
		case reflect.Chan:
		case reflect.Func:
		case reflect.UnsafePointer:
		}
	}
	if len(current.Children) == 0 {
		return nil
	}

	return current
}

func extractObjectPtr(newValue, oldValue reflect.Value, keyName string) *control.Diff {
	// check if the newValue is null
	if newValue.IsNil() {
		// if newValue is null and oldValue is not null then create delete
		if !oldValue.IsNil() {
			return control.NewDelDiff(&control.Key{
				Key: keyName,
			})
		}
		return nil
	}
	// if the old value is null then generate a blank type to compare against in oldValue
	if oldValue.IsNil() {
		oldValue.Set(reflect.New(newValue.Elem().Type()))
	}
	child := extractObject(newValue.Elem(), oldValue.Elem(), keyName)

	return child

}

func extractConcrete(newValue reflect.Value, oldValue reflect.Value, keyName string) *control.Diff {
	if !equal(newValue, oldValue) {
		child := control.NewDiff(&control.Key{
			Key: keyName,
		})
		err := child.SetValue(reflect.Indirect(newValue))
		if err != nil {
			panic(err)
		}

		if oldValue.CanSet() {
			oldValue.Set(newValue)
		}
		return child
	}
	return nil
}
func extractIndexBuiltIn[T int | string](newValue reflect.Value, oldValue reflect.Value, keyName string, index T) *control.Diff {
	switch newValue.Kind() {
	case reflect.Struct:
		child := extractObject(newValue, oldValue, keyName)
		if child != nil {
			child.Key.Index = control.NewObject(index)
			return child
		}
	case reflect.Map:
		panic("map not supported here")
	case reflect.Slice, reflect.Array:
		panic("slice not supported here")
	default:
		if !equal(newValue, oldValue) {
			child := control.NewDiff(&control.Key{
				Key:   keyName,
				Index: control.NewObject(index),
			})
			err := child.SetValue(reflect.Indirect(newValue))
			if err != nil {
				panic(err)
			}
			if oldValue.CanSet() {
				oldValue.Set(newValue)
			}
			return child
		}
	}
	return nil
}

func deleteIndexBuiltIn[T int | string](keyName string, index T) *control.Diff {
	return control.NewDelDiff(&control.Key{
		Key:   keyName,
		Index: control.NewObject(index),
	})
}

func extractSlice(newValue, oldValue reflect.Value, keyName string) []*control.Diff {
	var children []*control.Diff

	// make the old slice match the new slice
	// oldValue is shorter, add the extra entries and just run compare
	if oldValue.Len() < newValue.Len() {
		// make a new slice for oldValue of capacity the newValue Slice
		newOldSlice := reflect.MakeSlice(newValue.Type(), newValue.Len(), newValue.Len())
		// copy the values from the oldSlice into the newOldSlice
		reflect.Copy(newOldSlice, oldValue)
		// set the oldSlice to the newOldSlice
		oldValue.Set(newOldSlice)
		// if value is a pointer loop over and create a zero value entry for each element
	}

	// newValue is shorter, set a delete starting at the index of the difference
	if newValue.Len() < oldValue.Len() {
		children = append(children, deleteIndexBuiltIn(keyName, newValue.Len()))
		// now set the length of the oldValue slice to the newValue slice length
		oldValue.SetLen(newValue.Len())
	}
	for i := 0; i < newValue.Len(); i++ {
		newIndexValue, oldIndexValue := reflect.Indirect(newValue.Index(i)), reflect.Indirect(oldValue.Index(i))

		if equal(newIndexValue, oldIndexValue) {
			continue
		}

		child := extractIndexBuiltIn(newIndexValue, oldIndexValue, keyName, i)
		if child != nil {
			children = append(children, child)
		}
	}

	if len(children) > 0 {
		reflect.Copy(oldValue, newValue)
	}

	return children
}

func extractMap(newValue, oldValue reflect.Value, newUpperType reflect.Type, keyName string) []*control.Diff {
	var children []*control.Diff
	// indirect the values to get at the concrete values
	oldValue = reflect.Indirect(oldValue)
	newValue = reflect.Indirect(newValue)

	// if the oldValue length is 0 then the map needs to be created
	if oldValue.Len() == 0 {
		keyType := newUpperType.Key()
		valueType := newUpperType.Elem()
		mapType := reflect.MapOf(keyType, valueType)
		if oldValue.CanSet() {
			oldValue.Set(reflect.MakeMapWithSize(mapType, 0))
		}
	}

	for _, k := range newValue.MapKeys() {
		// append that value to the oldValue slice
		newMapIndexValue := newValue.MapIndex(k)
		oldMapIndexValue := oldValue.MapIndex(k)

		if !oldMapIndexValue.IsValid() {
			// create a dataStruct of the type in the slice to append to the oldValue slice
			dataStruct := reflect.New(newMapIndexValue.Type()).Elem()
			oldValue.SetMapIndex(k, dataStruct)
		}

		newMapIndexValue = reflect.Indirect(newMapIndexValue)
		oldMapIndexValue = reflect.Indirect(oldMapIndexValue)

		child := extractIndexBuiltIn(newMapIndexValue, oldMapIndexValue, keyName, makeString(k))
		if child != nil {
			children = append(children, child)
		}

		if newUpperType.Elem().Kind() == reflect.Ptr {
			oldValue.SetMapIndex(k, newMapIndexValue.Addr())
		} else {
			oldValue.SetMapIndex(k, newMapIndexValue)
		}

	}

	return children
}

func indirect(newValue, oldValue reflect.Value) (reflect.Value, reflect.Value) {
	for newValue.Kind() == reflect.Ptr {
		newValue = newValue.Elem()
	}
	for oldValue.Kind() == reflect.Ptr {
		oldValue = oldValue.Elem()
	}
	return newValue, oldValue
}
