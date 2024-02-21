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

type mapKey interface {
	int | ~string
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
	aStruct := dataStruct.Interface()
	return &Extractor{
		data:    aStruct,
		history: make([]*control.Diff, 0, historySize),
		mut:     new(sync.Mutex),
	}, nil
}

func (ext *Extractor) addHistory(head *control.Diff) {
	// if length of history equal to capacity drop first item and move everything down one
	// if len(ext.history) == cap(ext.history) {
	//     ext.history = ext.history[1:]
	// }
	// ext.history = append(ext.history, head)

	if len(head.GetChildren()) == 0 {
		return
	}
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

func (ext *Extractor) Diff(currData any) (*control.Diff, error) {
	// force single threaded access
	ext.mut.Lock()
	defer ext.mut.Unlock()

	// copy the current data as a point in time
	data := copyData(currData)

	newValue := reflect.ValueOf(data)
	if newValue.Kind() != reflect.Ptr {
		return nil, ErrNotPointer
	}

	newValue = reflect.Indirect(newValue)
	oldValue := reflect.Indirect(reflect.ValueOf(ext.data))

	head := extractObject(newValue, oldValue, newValue.Type().Name())

	head.Timestamp()
	return head, nil
}

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
			children := extractMap(newValueField, oldValueField, newValueTypeField.Type, newValueTypeField.Name)
			if len(children) > 0 {
				current.Children = append(current.Children, children...)
			}
		case reflect.Slice, reflect.Array:
			children := extractSlice(newValueField, oldValueField, newValueTypeField.Name)
			if len(children) > 0 {
				current.Children = append(current.Children, children...)
			}
		case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			child := extractBuiltIn(newValueField, oldValueField, newValueTypeField.Name)
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

func extractObjectInterface(newValue, oldValue reflect.Value, keyName string) *control.Diff {
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

	if !oldValue.IsValid() {
		oldValue = reflect.New(newValue.Elem().Type().Elem())
	}

	// if the old value is null then generate a blank type to compare against in oldValue
	if oldValue.IsNil() {
		t := newValue.Elem().Type()
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		oldValue.Set(reflect.New(t))
	}
	child := extractObject(newValue.Elem(), oldValue.Elem(), keyName)

	return child
}

func extractBuiltIn(newValue reflect.Value, oldValue reflect.Value, keyName string) *control.Diff {
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
		children := extractMap(newValue, oldValue, newValue.Type(), keyName)
		if len(children) > 0 {
			for _, c := range children {
				c.Key.Index = control.NewObjects(index, c.GetKey().GetIndex()...)
			}
			return children
		}
	case reflect.Slice, reflect.Array:
		children := extractSlice(newValue, oldValue, keyName)
		if len(children) > 0 {
			for _, c := range children {
				c.Key.Index = control.NewObjects(index, c.GetKey().GetIndex()...)
			}
			return children
		}
	default:
		if !equal(newValue, oldValue) {
			child := control.NewDiff(&control.Key{
				Key:   keyName,
				Index: control.NewObjects(index),
			})
			err := child.SetValue(reflect.Indirect(newValue))
			if err != nil {
				panic(err)
			}
			if oldValue.CanSet() {
				oldValue.Set(newValue)
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

		childs := extractIndexBuiltIn(newIndexValue, oldIndexValue, keyName, i)
		if len(childs) > 0 {
			children = append(children, childs...)
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

		childs := extractIndexBuiltIn(newMapIndexValue, oldMapIndexValue, keyName, k.Interface())
		if len(childs) > 0 {
			children = append(children, childs...)
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

func equal(n, o reflect.Value) bool {
	if n.Kind() != o.Kind() {
		return false
	}
	switch n.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return n.Int() == o.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return n.Uint() == o.Uint()
	case reflect.String:
		newS := n.String()
		oldS := o.String()
		return newS == oldS
		// return n.String() == o.String()
	case reflect.Bool:
		return n.Bool() == o.Bool()
	case reflect.Float32, reflect.Float64:
		return n.Float() == o.Float()
	case reflect.Complex64, reflect.Complex128:
		return n.Complex() == o.Complex()
	default:
		return false
	}
}

func copyData[T any](data T) T {
	t := reflect.TypeOf(data).Elem()
	dataStruct := reflect.New(t)

	currData := reflect.ValueOf(data)

	dataStruct.Elem().Set(currData.Elem())

	return dataStruct.Interface().(T)
}
