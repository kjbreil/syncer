package extractor

import (
	"fmt"
	"github.com/kjbreil/syncer/control"
	"reflect"
)

type Extractor struct {
	data    any
	history []*control.Diff
}

func New(data any) *Extractor {

	t := reflect.TypeOf(data)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	dataStruct := reflect.New(t)
	aStruct := dataStruct.Interface()
	return &Extractor{
		data:    aStruct,
		history: make([]*control.Diff, 0, 2),
	}
}

func (ext *Extractor) Data() any {
	return ext.data
}

func (ext *Extractor) Reset() {
	t := reflect.TypeOf(ext.data)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	dataStruct := reflect.New(t)
	aStruct := dataStruct.Interface()

	ext.data = aStruct
}

func (ext *Extractor) Diff(data any) *control.Diff {
	newValue := reflect.ValueOf(data)
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

	headName := oldType.Name()
	head := control.NewDiff([]*control.Key{
		{
			Key: headName,
		},
	},
	)

	extractLevel(head, newValue, oldValue)

	head.Timestamp()

	ext.addHistory(head)

	return head
}

func (ext *Extractor) addHistory(head *control.Diff) {
	// if length of history equal to capacity drop first item and move everything down one
	// if len(ext.history) == cap(ext.history) {
	//     ext.history = ext.history[1:]
	// }
	// ext.history = append(ext.history, head)

	if len(head.Children) == 0 {
		return
	}
	fmt.Println(len(ext.history), cap(ext.history))
	if len(ext.history) == cap(ext.history) {
		for i := 0; i < len(ext.history)-1; i++ {
			ext.history[i] = ext.history[i+1]
		}
		ext.history[len(ext.history)-1] = head
	} else {
		ext.history = append(ext.history, head)
	}

}

func extractLevel(parent *control.Diff, newValue reflect.Value, oldValue reflect.Value) {
	newType := newValue.Type()

	zv := reflect.Value{}
	if oldValue == zv {
		dataStruct := reflect.New(newValue.Type()).Elem()
		oldValue = dataStruct
	}
	oldType := oldValue.Type()

	numFields := newValue.NumField()
	for i := 0; i < numFields; i++ {
		if newType.Field(i).Name != oldType.Field(i).Name {
			panic("this shouldn't happen")
		}

		etag := newType.Field(i).Tag.Get("extractor")
		if etag == "-" {
			continue
		}

		newValueFieldKind := newValue.Field(i).Kind()

		key := oldType.Field(i).Name
		child := control.NewDiff(append(parent.Key, &control.Key{
			Key: key,
		}),
		)
		hasChildren := false

		switch newValueFieldKind {
		case reflect.Pointer:
			extractLevelPointer(parent, newValue, oldValue, i, child)
		case reflect.Map:
			extractLevelMap(parent, newValue, oldValue, i, oldType, child, key)
		case reflect.Slice, reflect.Array:
			extractLevelSlice(parent, newValue, oldValue, i, key, child)
		case reflect.Struct:
			extractLevel(child, newValue.Field(i), oldValue.Field(i))
			if child.Children != nil {
				parent.Children = append(parent.Children, child)
			}
			extractChildren(parent, child, newValue.Field(i), oldValue.Field(i), &hasChildren)
		default:
			if !equal(newValue.Field(i), oldValue.Field(i)) {
				child.Value = &control.Object{}
				setValue(newValue.Field(i), child)

				parent.Children = append(parent.Children, child)
				if oldValue.Field(i).CanSet() {
					oldValue.Field(i).Set(newValue.Field(i))
				}
			}
		}
	}
}

func extractLevelSlice(parent *control.Diff, newValue reflect.Value, oldValue reflect.Value, i int, key string, child *control.Diff) {
	newFieldValue, oldFieldValue := newValue.Field(i), oldValue.Field(i)
	shortest := min(newFieldValue.Len(), oldFieldValue.Len())
	var hasChildren bool
	for ii := 0; ii < shortest; ii++ {
		newIndexValue, oldIndexValue := newFieldValue.Index(ii), oldFieldValue.Index(ii)

		if equal(newIndexValue, oldIndexValue) {
			continue
		}
		indexNewValue := newIndexValue
		if indexNewValue.Type().Kind() == reflect.Ptr {
			extractNonStruct(parent, newIndexValue.Elem(), oldIndexValue.Elem(), ii, key)
		} else if indexNewValue.Type().Kind() != reflect.Struct {
			extractNonStruct(parent, newIndexValue, oldIndexValue, ii, key)
		} else {
			var hasChildren bool
			extractChildren(parent, child, newIndexValue, oldIndexValue, &hasChildren)
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
			// now extract
			if newIndexValue.Type().Kind() == reflect.Ptr {
				extractNonStruct(parent, newIndexValue.Elem(), oldIndexValue.Elem(), ii, key)
			} else if newIndexValue.Type().Kind() != reflect.Struct {
				extractNonStruct(parent, newIndexValue, oldIndexValue, ii, key)
			} else {
				extractChildren(parent, child, newIndexValue, oldIndexValue, &hasChildren)
			}
		}
	}
	// oldValue slice is longer than the newValue so items were deleted
	if oldFieldValue.Len() > newFieldValue.Len() {
		for ii := shortest; ii < oldFieldValue.Len(); ii++ {
			deleteNonStruct(parent, ii, key)
		}
	}

	reflect.Copy(oldFieldValue, newFieldValue)
}

func extractLevelMap(parent *control.Diff, newValue reflect.Value, oldValue reflect.Value, i int, oldType reflect.Type, child *control.Diff, key string) {
	// Make the map for the oldValue if it doesn't exist
	if oldValue.Field(i).Len() == 0 {
		keyType := oldType.Field(i).Type.Key()
		valueType := oldType.Field(i).Type.Elem()
		mapType := reflect.MapOf(keyType, valueType)
		if oldValue.Field(i).CanSet() {
			oldValue.Field(i).Set(reflect.MakeMapWithSize(mapType, 0))
		}
	}
	for _, k := range newValue.Field(i).MapKeys() {
		// append that value to the oldValue slice
		zeroValue := reflect.Value{}
		if oldValue.Field(i).MapIndex(k) == zeroValue {
			// create a dataStruct of the type in the slice to append to the oldValue slice
			dataStruct := reflect.New(newValue.Field(i).MapIndex(k).Type()).Elem()
			oldValue.Field(i).SetMapIndex(k, dataStruct)
		}

		var hasChildren bool
		switch newValue.Field(i).MapIndex(k).Type().Kind() {
		case reflect.Ptr:
			extractChildren(parent, child, newValue.Field(i).MapIndex(k).Elem(), oldValue.Field(i).MapIndex(k).Elem(), &hasChildren)
		case reflect.Struct:
			extractNonStruct(parent, newValue.Field(i).MapIndex(k), oldValue.Field(i).MapIndex(k), makeString(k), key)
		default:
			extractChildren(parent, child, newValue.Field(i).MapIndex(k), oldValue.Field(i).MapIndex(k), &hasChildren)
		}

		// the address cannot be set so setting it manually
		oldValue.Field(i).SetMapIndex(k, newValue.Field(i).MapIndex(k))
	}
	// find deletes
	for _, k := range oldValue.Field(i).MapKeys() {

		zeroValue := reflect.Value{}
		if newValue.Field(i).MapIndex(k) == zeroValue {
			deleteNonStruct(parent, makeString(k), key)
			oldValue.Field(i).SetMapIndex(k, reflect.Value{})
		}
	}
}

func extractLevelPointer(parent *control.Diff, newValue reflect.Value, oldValue reflect.Value, i int, child *control.Diff) {
	if newValue.Field(i).IsNil() {
		if !oldValue.Field(i).IsNil() {
			child.Delete = true
			parent.Children = append(parent.Children, child)
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
	extractChildren(parent, child, newValue.Field(i).Elem(), oldValue.Field(i).Elem(), &hasChildren)
}

func setValue(va reflect.Value, child *control.Diff) {
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
	case reflect.Complex64, reflect.Complex128:
	case reflect.Chan:
	case reflect.Func:
	case reflect.Interface:
	case reflect.Map:
		panic("cannot make string of slice, should explode out")
	case reflect.Pointer:
	case reflect.Slice, reflect.Array:
		panic("cannot make string of slice, should explode out")
	case reflect.String:
		value := va.String()
		child.Value.String_ = &value
	case reflect.Struct:
	case reflect.UnsafePointer:
	}
}

func extractNonStruct(parent *control.Diff, newValue reflect.Value, oldValue reflect.Value, index interface{}, key string) {
	if !equal(newValue, oldValue) {
		var indexObject control.Object
		switch index.(type) {
		case string:
			s := index.(string)
			indexObject.String_ = &s
		case int:
			i := int64(index.(int))
			indexObject.Int64 = &i
		case int32:
			i := int64(index.(int32))
			indexObject.Int64 = &i
		case int64:
			i := index.(int64)
			indexObject.Int64 = &i
		default:
			panic("extractNonStruct: unsupported type")
		}

		child := control.NewDiff(append(parent.Key, &control.Key{
			Key:   key,
			Index: &indexObject,
		}),
		)
		setValue(newValue, child)

		parent.Children = append(parent.Children, child)
		if oldValue.CanSet() {
			oldValue.Set(newValue)
		}
	}

}

func deleteNonStruct[i int | string](parent *control.Diff, index i, key string) {
	child := control.NewDelDiff(append(parent.Key, &control.Key{
		Key:   key,
		Index: control.NewObject(index),
	}),
	)
	parent.Children = append(parent.Children, child)
}

func extractChildren(parent *control.Diff, child *control.Diff, newValue reflect.Value, oldValue reflect.Value, hasChildren *bool) {
	extractLevel(child, newValue, oldValue)
	if !*hasChildren && child.Children != nil {
		parent.Children = append(parent.Children, child)
		*hasChildren = true
	}
}

func equal(n reflect.Value, o reflect.Value) bool {

	if n.Kind() != o.Kind() {
		return false
	}
	switch n.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return n.Int() == o.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return n.Uint() == o.Uint()
	case reflect.String:
		return n.String() == o.String()
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
		return fmt.Sprintf("%d", x.Int())
	case reflect.Bool:
		return fmt.Sprintf("%t", x.Bool())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", x.Uint())
	case reflect.Uintptr:
		return fmt.Sprintf("%d", x.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", x.Float())
	case reflect.Complex64, reflect.Complex128:
		return fmt.Sprintf("%f", x.Complex())
	case reflect.Chan:
	case reflect.Func:
	case reflect.Interface:
	case reflect.Map:
		panic("cannot make string of slice, should explode out")
	case reflect.Pointer:
		return makeString(x.Elem())
	case reflect.Slice, reflect.Array:
		panic("cannot make string of slice, should explode out")
	case reflect.String:
		return x.String()
	case reflect.Struct:
	case reflect.UnsafePointer:
	}
	return ""
}
