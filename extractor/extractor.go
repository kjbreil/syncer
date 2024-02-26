package extractor

import (
	"errors"
	"fmt"
	"github.com/kjbreil/syncer/control"
	"github.com/kjbreil/syncer/helpers/deepcopy"
	"github.com/kjbreil/syncer/helpers/equal"
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
func New(data any) (*Extractor, error) {
	if data == nil {
		return nil, errors.New("data is nil")
	}
	t := reflect.Indirect(reflect.ValueOf(data)).Type()
	dataStruct := reflect.New(t)
	aStruct := deepcopy.Any(dataStruct.Interface())
	return &Extractor{
		data:    aStruct,
		history: make([]*control.Diff, 0, historySize),
		mut:     new(sync.Mutex),
	}, nil
}

func NewEntries(data any) (control.Entries, error) {
	if data == nil {
		return nil, errors.New("data is nil")
	}
	t := reflect.Indirect(reflect.ValueOf(data)).Type()
	dataStruct := reflect.New(t)
	aStruct := deepcopy.Any(dataStruct.Interface())
	e := &Extractor{
		data:    aStruct,
		history: make([]*control.Diff, 0, historySize),
		mut:     new(sync.Mutex),
	}
	return e.Entries(&data), nil
}

func (ext *Extractor) addHistory(head *control.Diff) {
	if len(head.GetChildren()) == 0 {
		return
	}
	// if length of history Equal to capacity drop first item and move everything down one
	if len(ext.history) == cap(ext.history) {
		for i := 0; i < len(ext.history)-1; i++ {
			ext.history[i] = ext.history[i+1]
		}
		ext.history[len(ext.history)-1] = head
	} else {
		ext.history = append(ext.history, head)
	}
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

// Entries returns the Entries from the Extractor.
//
// The function takes in a parameter 'data' of type 'any' and returns a value of type 'control.Entries'.
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

	// deep copy the current data as a point in time
	pitData := deepcopy.Any(data)

	newValue := reflect.ValueOf(pitData)
	if newValue.Kind() != reflect.Ptr {
		return nil, ErrNotPointer
	}

	newValue = reflect.Indirect(newValue)
	// newValue = reflect.ValueOf(newValue.Interface())

	oldValue := reflect.Indirect(reflect.ValueOf(ext.data))

	// children, err := extract(newValue, oldValue, newValue.Type())
	//
	// if err != nil && !errors.Is(err, ErrUnsupportedType) {
	// 	return nil, err
	// }
	//
	// ext.data = pitData
	// if len(children) > 0 {
	// 	head := children[0]
	// 	head.Timestamp()
	// 	return head, nil
	// }
	// return nil, nil

	head := extractObject(newValue, oldValue, newValue.Type().Name())
	// TODO: This needs to be moved withing the head != nil but maps are messed up
	// use the pitData for ext.data since it is already a deep copy of data
	ext.data = pitData
	if head != nil {
		head.Timestamp()
	}

	return head, nil
}

// func extract(newValue reflect.Value, oldValue reflect.Value, upperType reflect.Type) ([]*control.Diff, error) {
// 	if iFn, ok := extFns[newValue.Kind()]; ok {
// 		head, err := iFn(newValue, oldValue, upperType)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return head, nil
// 	}
// 	return nil, fmt.Errorf("%w: %s", ErrUnsupportedType, newValue.Kind())
// }

// extractObject performs a deep comparison of two reflect Values and creates a control.Diff
//
// newValue, oldValue reflect.Value, keyName string
// *control.Diff
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
		// skip if extractor tag is set
		etag := newValue.Type().Field(i).Tag.Get("extractor")
		if etag == "-" {
			continue
		}
		if !newValue.Field(i).CanInterface() {
			continue
		}

		newValueField := newValue.Field(i)
		newValueFieldKind := newValueField.Kind()

		oldValueField := oldValue.Field(i)
		newValueTypeField := newValue.Type().Field(i)
		switch newValueFieldKind {
		case reflect.Struct:
			child := extractObject(newValueField, oldValueField, newValueTypeField.Name)
			if child != nil {
				current.Children = append(current.Children, child)
			}
		case reflect.Pointer:
			child := extractObjectPtr(newValueField, oldValueField, newValueTypeField.Name)
			if child != nil {
				current.Children = append(current.Children, child)
			}
		case reflect.Interface:
			child := extractObjectInterface(newValueField, oldValueField, newValueTypeField.Name)
			if child != nil {
				current.Children = append(current.Children, child)
			}
		case reflect.Map:
			children := extractMapOld(newValueField, oldValueField, newValueTypeField.Type, newValueTypeField.Name)
			if len(children) > 0 {
				current.Children = append(current.Children, children...)
			}
		case reflect.Slice, reflect.Array:
			children := extractSliceOld(newValueField, oldValueField, newValueTypeField.Name)
			if len(children) > 0 {
				current.Children = append(current.Children, children...)
			}
		case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			child := extractPrimitiveOld(newValueField, oldValueField, newValueTypeField.Name)
			if child != nil {
				current.Children = append(current.Children, child)
			}
		case reflect.Complex64, reflect.Complex128:
			panic(fmt.Sprintf("complex types not yet supported"))
		case reflect.Invalid, reflect.Uintptr, reflect.Chan, reflect.Func, reflect.UnsafePointer:

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

func extractObjectInterface(newValue, oldValue reflect.Value, keyName string) *control.Diff {
	// check if the newValue is null
	if newValue.IsNil() {
		// if newValue is null and oldValue is not null then create delete
		if !oldValue.IsValid() || !oldValue.IsNil() {
			return control.NewDelDiff(&control.Key{
				Key: keyName,
			})
		}
		return nil
	}

	if !oldValue.IsValid() {
		oldValue = reflect.New(newValue.Elem().Type().Elem())
	}

	// if the old value is null then generate a blank type to compare against in oldValue
	if !oldValue.IsValid() || oldValue.IsNil() {
		oldValue = reflect.New(newValue.Elem().Type())
		// t := newValue.Elem().Type()
		// if t.Kind() == reflect.Ptr {
		// 	t = t.Elem()
		// }
		// if oldValue.CanSet() {
		// 	oldValue.Set(reflect.New(t))
		// }
	}
	child := extractObject(newValue.Elem(), oldValue.Elem(), keyName)

	return child
}

func extractPrimitiveOld(newValue reflect.Value, oldValue reflect.Value, keyName string) *control.Diff {
	if !equal.Equal(newValue, oldValue) {
		child := control.NewDiff(&control.Key{
			Key: keyName,
		})
		err := child.SetValue(reflect.Indirect(newValue))
		if err != nil {
			panic(err)
		}
		return child
	}
	return nil
}
func extractIndexBuiltIn(newValue reflect.Value, oldValue reflect.Value, keyName string, index any) []*control.Diff {
	switch newValue.Kind() {
	case reflect.Struct:
		child := extractObject(newValue, oldValue, keyName)
		if child != nil {
			child.Key.Index = control.NewObjects(index)
			return []*control.Diff{child}
		}
	case reflect.Interface:

		child := extractObjectInterface(newValue, oldValue, keyName)
		if child != nil {
			child.Key.Index = control.NewObjects(index)
			return []*control.Diff{child}
		}
	case reflect.Map:
		if !oldValue.IsValid() {
			keyType := newValue.Type().Key()
			valueType := newValue.Type().Elem()
			mapType := reflect.MapOf(keyType, valueType)
			oldValue = reflect.MakeMapWithSize(mapType, 0)
		}
		children := extractMapOld(newValue, oldValue, newValue.Type(), keyName)
		if len(children) > 0 {
			for _, c := range children {
				c.Key.Index = control.NewObjects(index, c.GetKey().GetIndex()...)
			}
			return children
		}
	case reflect.Slice, reflect.Array:
		children := extractSliceOld(newValue, oldValue, keyName)
		if len(children) > 0 {
			for _, c := range children {
				c.Key.Index = control.NewObjects(index, c.GetKey().GetIndex()...)
			}
			return children
		}
	default:
		if !equal.Equal(newValue, oldValue) {
			child := control.NewDiff(&control.Key{
				Key:   keyName,
				Index: control.NewObjects(index),
			})
			err := child.SetValue(reflect.Indirect(newValue))
			if err != nil {
				panic(err)
			}
			return []*control.Diff{child}
		}
	}
	return nil
}

func deleteIndexBuiltIn(keyName string, index any) *control.Diff {
	return control.NewDelDiff(&control.Key{
		Key:   keyName,
		Index: control.NewObjects(index),
	})
}

func extractSliceOld(newValue, oldValue reflect.Value, keyName string) []*control.Diff {
	var children []*control.Diff

	// if newValue is a slice and is null while oldValue is not null then create remove
	if newValue.Kind() == reflect.Slice && newValue.IsNil() && !oldValue.IsNil() {
		children = append(children, control.NewDelDiff(&control.Key{
			Key: keyName,
		}))
		return children
	}

	if !oldValue.IsValid() {
		oldValue = reflect.MakeSlice(newValue.Type(), newValue.Len(), newValue.Cap())
	}
	// make the old slice match the new slice
	// oldValue is shorter, add the extra entries and just run compare
	if oldValue.Len() < newValue.Len() {
		// make a new slice for oldValue of capacity the newValue Slice
		newOldSlice := reflect.MakeSlice(newValue.Type(), newValue.Len(), newValue.Cap())
		// copy the values from the oldSlice into the newOldSlice
		reflect.Copy(newOldSlice, oldValue)
		// set the oldSlice to the newOldSlice
		oldValue = newOldSlice
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

		if equal.Equal(newIndexValue, oldIndexValue) {
			continue
		}

		childs := extractIndexBuiltIn(newIndexValue, oldIndexValue, keyName, i)
		if len(childs) > 0 {
			children = append(children, childs...)
		}
	}

	return children
}

func extractMapOld(newValue, oldValue reflect.Value, newUpperType reflect.Type, keyName string) []*control.Diff {
	var children []*control.Diff
	if newValue.IsNil() && !oldValue.IsNil() {
		children = append(children, control.NewDelDiff(&control.Key{
			Key: keyName,
		}))
		return children
	}
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
		newMapIndexValue := newValue.MapIndex(k)
		oldMapIndexValue := oldValue.MapIndex(k)

		if !oldMapIndexValue.IsValid() {
			// create a dataStruct of the type in the slice to append to the oldValue slice
			dataStruct := reflect.New(newMapIndexValue.Type()).Elem()
			oldValue.SetMapIndex(k, dataStruct)
		}

		newMapIndexValue = reflect.Indirect(newMapIndexValue)
		oldMapIndexValue = reflect.Indirect(oldMapIndexValue)

		childs := extractIndexBuiltIn(newMapIndexValue, oldMapIndexValue, keyName, k.Interface())
		if len(childs) > 0 {
			children = append(children, childs...)
		}
	}

	for _, k := range oldValue.MapKeys() {
		newMapIndexValue := newValue.MapIndex(k)

		if !newMapIndexValue.IsValid() {
			children = append(children, deleteIndexBuiltIn(keyName, k.Interface()))
		}
	}

	return children
}

func copyData[T any](data T) T {
	t := reflect.TypeOf(data).Elem()
	dataStruct := reflect.New(t)

	currData := reflect.ValueOf(data)

	dataStruct.Elem().Set(currData.Elem())

	return dataStruct.Interface().(T)
}
