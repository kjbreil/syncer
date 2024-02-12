package extractor

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"

	"github.com/kjbreil/syncer/control"
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

// Entries returns the entries extracted from the diff obtained by comparing the current data with the provided data.
// If an error occurs during the diff extraction, it will panic.
// The extracted entries are returned as an array of control.Entry objects.
func (ext *Extractor) Entries(data any) control.Entries {
	h, err := ext.Diff(data)
	if err != nil {
		panic(err)
	}
	return h.Entries()
}

func (ext *Extractor) Diff(data any) (*control.Diff, error) {
	ext.mut.Lock()
	defer ext.mut.Unlock()
	newValue := reflect.ValueOf(data)
	if newValue.Kind() != reflect.Ptr {
		return nil, ErrNotPointer
	}

	oldValue := reflect.ValueOf(ext.data)

	// if it's a pointer follow to the real data
	for newValue.Kind() == reflect.Ptr {
		newValue = newValue.Elem()
	}
	for oldValue.Kind() == reflect.Ptr {
		oldValue = oldValue.Elem()
	}

	newType := newValue.Type()
	oldType := oldValue.Type()

	if newType != oldType {
		panic("not same types")
	}

	head := control.NewDiff([]*control.Key{
		{
			Key: oldType.Name(),
		},
	},
	)

	err := extractLevel(head, newValue, oldValue)
	if err != nil {
		return nil, err
	}
	head.Timestamp()

	// ext.addHistory(head)

	return head, nil
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

func extractLevel(parent *control.Diff, newValue, oldValue reflect.Value) error {
	newType := newValue.Type()

	if !oldValue.IsValid() {
		dataStruct := reflect.New(newValue.Type()).Elem()
		oldValue = dataStruct
	}
	oldType := oldValue.Type()

	numFields := newValue.NumField()
	for i := 0; i < numFields; i++ {
		if newType.Field(i).Name != oldType.Field(i).Name {
			return fmt.Errorf("%w: %s != %s", ErrDataStructMisMatch, newType.Field(i).Name, oldType.Field(i).Name)
		}

		etag := newType.Field(i).Tag.Get("extractor")
		if etag == "-" {
			continue
		}

		key := oldType.Field(i).Name
		child := control.NewDiff(append(parent.Key, &control.Key{
			Key: key,
		}),
		)
		hasChildren := false
		newValueFieldKind := newValue.Field(i).Kind()

		switch newValueFieldKind {
		case reflect.Pointer:
			extractLevelPointer(parent, newValue, oldValue, i, child)
		case reflect.Map:
			err := extractLevelMap(parent, newValue, oldValue, i, oldType, child, key)
			if err != nil {
				return fmt.Errorf("extractLevel: %w", err)
			}
		case reflect.Slice, reflect.Array:
			err := extractLevelSlice(parent, newValue, oldValue, i, key, child)
			if err != nil {
				return fmt.Errorf("extractLevel: %w", err)
			}
		case reflect.Struct:
			err := extractLevel(child, newValue.Field(i), oldValue.Field(i))
			if err != nil {
				return fmt.Errorf("extractLevel: %w", err)
			}
			if child.Children != nil {
				parent.AddChild(child, numFields)
			}
			err = extractChildren(parent, child, newValue.Field(i), oldValue.Field(i), &hasChildren)
			if err != nil {
				return fmt.Errorf("extractLevel: %w", err)
			}
		default:
			if !equal(newValue.Field(i), oldValue.Field(i)) {
				child.Value = &control.Object{}
				err := setValue(newValue.Field(i), child)
				if err != nil {
					return fmt.Errorf("extractLevel: %w", err)
				}
				parent.AddChild(child, numFields)
				if oldValue.Field(i).CanSet() {
					oldValue.Field(i).Set(newValue.Field(i))
				}
			}
		}
	}
	return nil
}

func extractLevelSlice(parent *control.Diff, newValue, oldValue reflect.Value, i int, key string, child *control.Diff) error {
	newFieldValue, oldFieldValue := newValue.Field(i), oldValue.Field(i)
	shortest := min(newFieldValue.Len(), oldFieldValue.Len())
	var hasChildren bool
	for ii := 0; ii < shortest; ii++ {
		newIndexValue, oldIndexValue := newFieldValue.Index(ii), oldFieldValue.Index(ii)

		// indirect newIndexValue and oldIndexValue so we can set the value
		newIndexValue = reflect.Indirect(newIndexValue)
		oldIndexValue = reflect.Indirect(oldIndexValue)

		if equal(newIndexValue, oldIndexValue) {
			continue
		}
		indexNewValue := newIndexValue
		switch {
		case indexNewValue.Type().Kind() == reflect.Ptr:
			// TODO: Remove this case statement after testing
			panic("indirect above should prevent this case statement")
		case indexNewValue.Type().Kind() != reflect.Struct:
			child, err := extractNonStruct(parent, newIndexValue, oldIndexValue, ii, key)
			if err != nil {
				return fmt.Errorf("extractLevelSlice: %w", err)
			}
			if child != nil {
				parent.AddChild(child, shortest)
			}
		default:
			err := extractChildren(parent, child, newIndexValue, oldIndexValue, &hasChildren)
			if err != nil {
				return fmt.Errorf("extractLevelSlice: %w", err)
			}
		}
	}
	// new value has more data than the olddata
	if newFieldValue.Len() > oldFieldValue.Len() {
		for ii := shortest; ii < newFieldValue.Len(); ii++ {
			// create a dataStruct of the type in the slice to append to the oldValue slice
			newIndexValue := newFieldValue.Index(ii)

			dataStruct := reflect.New(newIndexValue.Type()).Elem()
			// append that value to the oldValue slice
			oldFieldValue.Set(reflect.Append(oldFieldValue, dataStruct))

			oldIndexValue := oldFieldValue.Index(ii)

			// indirect newIndexValue and oldIndexValue so we can set the value
			newIndexValue = reflect.Indirect(newIndexValue)
			oldIndexValue = reflect.Indirect(oldIndexValue)

			// now extract
			switch {
			case newIndexValue.Type().Kind() == reflect.Ptr:
				// TODO: Remove this case statement after testing
				panic("indirect above should prevent this case statement")
			case newIndexValue.Type().Kind() != reflect.Struct:
				child, err := extractNonStruct(parent, newIndexValue, oldIndexValue, ii, key)
				if err != nil {
					return fmt.Errorf("extractLevelSlice: %w", err)
				}
				if child != nil {
					parent.AddChild(child, newFieldValue.Len())
				}
			default:
				err := extractChildren(parent, child, newIndexValue, oldIndexValue, &hasChildren)
				if err != nil {
					return fmt.Errorf("extractLevelSlice: %w", err)
				}
			}
		}
	}
	// oldValue slice is longer than the newValue so items were deleted
	if oldFieldValue.Len() > newFieldValue.Len() {
		for ii := shortest; ii < oldFieldValue.Len(); ii++ {
			deleteNonStruct(parent, ii, key)
		}
	}

	// change the old slice to the same length of the new slice for copying
	if oldFieldValue.CanSet() {

		oldFieldValue.SetLen(newFieldValue.Len())
	}

	reflect.Copy(oldFieldValue, newFieldValue)
	return nil
}

func extractLevelMap(parent *control.Diff, newValue, oldValue reflect.Value, i int, oldType reflect.Type, child *control.Diff, key string) error {
	// Make the map for the oldValue if it doesn't exist
	oldValueField := reflect.Indirect(oldValue.Field(i))
	newValueField := reflect.Indirect(newValue.Field(i))
	if oldValueField.Len() == 0 {
		keyType := oldType.Field(i).Type.Key()
		valueType := oldType.Field(i).Type.Elem()
		mapType := reflect.MapOf(keyType, valueType)
		if oldValueField.CanSet() {
			oldValueField.Set(reflect.MakeMapWithSize(mapType, 0))
		}
	}
	for _, k := range newValueField.MapKeys() {
		// append that value to the oldValue slice
		newMapIndexValue := newValueField.MapIndex(k)
		oldMapIndexValue := oldValueField.MapIndex(k)

		if !oldMapIndexValue.IsValid() {
			// create a dataStruct of the type in the slice to append to the oldValue slice
			dataStruct := reflect.New(newMapIndexValue.Type()).Elem()
			oldValueField.SetMapIndex(k, dataStruct)
		}

		newMapIndexValue = reflect.Indirect(newMapIndexValue)
		oldMapIndexValue = reflect.Indirect(oldMapIndexValue)

		var hasChildren bool
		switch newMapIndexValue.Type().Kind() {
		case reflect.Ptr:
			// TODO: Remove this case statement after testing
			panic("indirect above should prevent this case statement")
		case reflect.Struct:
			err := extractChildren(parent, child, newMapIndexValue, oldMapIndexValue, &hasChildren)
			if err != nil {
				return fmt.Errorf("extractLevelMap: %w", err)
			}
		default:
			child, err := extractNonStruct(parent, newMapIndexValue, oldMapIndexValue, makeString(k), key)
			if err != nil {
				return fmt.Errorf("extractLevelMap: %w", err)
			}
			if child != nil {
				parent.AddChild(child, newValueField.Len())
			}
		}

		// the address cannot be set so setting it manually
		// check if type is a pointer and so we can get the address of the new value instead of the value itself
		// TODO: Probably need to make a copy of the value and then use, need test cases
		if oldType.Field(i).Type.Elem().Kind() == reflect.Ptr {
			oldValueField.SetMapIndex(k, newMapIndexValue.Addr())
		} else {
			oldValueField.SetMapIndex(k, newMapIndexValue)
		}
	}
	// find deletes
	for _, k := range oldValueField.MapKeys() {
		if !newValueField.MapIndex(k).IsValid() {
			deleteNonStruct(parent, makeString(k), key)
			oldValueField.SetMapIndex(k, reflect.Value{})
		}
	}
	return nil
}

func extractLevelPointer(parent *control.Diff, newValue, oldValue reflect.Value, i int, child *control.Diff) {
	if newValue.Field(i).IsNil() {
		if !oldValue.Field(i).IsNil() {
			child.Delete = true
			parent.AddChild(child, 10)
			if oldValue.CanSet() {
				oldValue.Set(newValue)
			}
		}
		return
	}
	if oldValue.Field(i).IsNil() {
		oldValue.Field(i).Set(reflect.New(newValue.Field(i).Elem().Type()))
	}
	var hasChildren bool
	err := extractChildren(parent, child, newValue.Field(i).Elem(), oldValue.Field(i).Elem(), &hasChildren)
	if err != nil {
		return
	}
}

func setValue(va reflect.Value, child *control.Diff) error {
	child.Value = &control.Object{}
	switch va.Kind() {
	case reflect.Invalid:
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value := va.Int()
		child.Value.Int64 = &value
	case reflect.Bool:
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value := va.Uint()
		child.Value.Uint64 = &value
	case reflect.Uintptr:
	case reflect.Float32:
		value := float32(va.Float())
		child.Value.Float32 = &value
	case reflect.Float64:
		value := va.Float()
		child.Value.Float64 = &value
	case reflect.String:
		value := va.String()
		child.Value.String_ = &value
	default:
		return fmt.Errorf("cannot setValue of type %s", va.Type().String())
	}
	return nil
}

func extractNonStruct(parent *control.Diff, newValue reflect.Value, oldValue reflect.Value, index any, key string) (*control.Diff, error) {
	if !equal(newValue, oldValue) {
		var indexObject control.Object
		switch v := index.(type) {
		case string:
			indexObject.String_ = &v
		case int:
			i := int64(v)
			indexObject.Int64 = &i
		case int32:
			i := int64(v)
			indexObject.Int64 = &i
		case int64:
			indexObject.Int64 = &v
		default:
			return nil, fmt.Errorf("extractNonStruct: %w", ErrUnsupportedType)
		}

		child := control.NewDiff(append(parent.Key, &control.Key{
			Key:   key,
			Index: &indexObject,
		}),
		)
		err := setValue(newValue, child)
		if err != nil {
			return nil, fmt.Errorf("extractNonStruct: %w", err)
		}
		if oldValue.CanSet() {
			oldValue.Set(newValue)
		}
		return child, nil
	}

	return nil, nil
}

func deleteNonStruct[i int | string](parent *control.Diff, index i, key string) {
	child := control.NewDelDiff(append(parent.Key, &control.Key{
		Key:   key,
		Index: control.NewObject(index),
	}),
	)
	parent.AddChild(child, 10)
}

func extractChildren(parent *control.Diff, child *control.Diff, newValue, oldValue reflect.Value, hasChildren *bool) error {
	err := extractLevel(child, newValue, oldValue)
	if err != nil {
		return fmt.Errorf("extractChildren: %w", err)
	}
	if !*hasChildren && child.Children != nil {
		parent.AddChild(child, 10)
		*hasChildren = true
	}
	return nil
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

func makeString(x reflect.Value) string {
	switch x.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(x.Int(), 10)
	case reflect.Bool:
		return strconv.FormatBool(x.Bool())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(x.Uint(), 10)
	case reflect.Uintptr:
		return fmt.Sprintf("%d", x.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", x.Float())
	case reflect.Complex64, reflect.Complex128:
		return fmt.Sprintf("%f", x.Complex())
	case reflect.Pointer:
		return makeString(x.Elem())
	case reflect.String:
		return x.String()
	default:
		panic("makeString: unsupported type " + x.Type().String())
	}
}
