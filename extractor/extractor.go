package extractor

import (
	"fmt"
	"github.com/kjbreil/syncer/control"
	"reflect"
)

type Extractor struct {
	data any
}

type Differences struct {
	children []*Differences
	key      []*control.Key
	value    *control.Object
	delete   bool
}

func New(data any) *Extractor {

	t := reflect.TypeOf(data)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	dataStruct := reflect.New(t)
	aStruct := dataStruct.Interface()
	return &Extractor{
		data: aStruct,
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

func (d *Differences) Entries() []*control.Entry {
	return d.entries()
}
func (d *Differences) entries() []*control.Entry {
	var moulds []*control.Entry
	for _, c := range d.children {
		if len(c.children) > 0 {
			moulds = append(moulds, c.entries()...)
		} else {
			if c.delete {
				moulds = append(moulds, &control.Entry{
					Key:    c.key,
					Action: control.Entry_REMOVE,
				})
			} else {
				moulds = append(moulds, &control.Entry{
					Key:   c.key,
					Value: c.value,
				})
			}
		}
	}
	return moulds
}

func (ext *Extractor) Diff(data any) *Differences {
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
	head := &Differences{
		key: []*control.Key{
			{
				Key: headName,
			},
		},
		value: nil,
	}

	extractLevel(head, newValue, oldValue)

	return head
}

func extractLevel(parent *Differences, newValue reflect.Value, oldValue reflect.Value) {
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

		newValueFieldKind := newValue.Field(i).Kind()

		key := oldType.Field(i).Name
		child := &Differences{
			key: append(parent.key, &control.Key{
				Key: key,
			}),
		}
		hasChildren := false

		switch newValueFieldKind {
		case reflect.Pointer:
			extractLevelPointer(parent, newValue, oldValue, i, child, hasChildren)
		case reflect.Map:
			extractLevelMap(parent, newValue, oldValue, i, oldType, child, hasChildren, key)
		case reflect.Slice, reflect.Array:
			extractLevelSlice(parent, newValue, oldValue, i, key, child, hasChildren)
		case reflect.Struct:
			extractLevel(child, newValue.Field(i), oldValue.Field(i))
			if child.children != nil {
				parent.children = append(parent.children, child)
			}

			extractChildren(parent, child, newValue.Field(i), oldValue.Field(i), &hasChildren)
		default:
			if !equal(newValue.Field(i), oldValue.Field(i)) {
				child.value = &control.Object{}
				setValue(newValue.Field(i), child)

				parent.children = append(parent.children, child)
				if oldValue.Field(i).CanSet() {
					oldValue.Field(i).Set(newValue.Field(i))
				}
			}
		}
	}
	return
}

func extractLevelSlice(parent *Differences, newValue reflect.Value, oldValue reflect.Value, i int, key string, child *Differences, hasChildren bool) {
	shortest := min(newValue.Field(i).Len(), oldValue.Field(i).Len())

	for ii := 0; ii < shortest; ii++ {
		if equal(newValue.Field(i).Index(ii), oldValue.Field(i).Index(ii)) {
			continue
		}
		indexNewValue := newValue.Field(i).Index(ii)
		if indexNewValue.Type().Kind() == reflect.Ptr {
			extractNonStruct(parent, newValue.Field(i).Index(ii).Elem(), oldValue.Field(i).Index(ii).Elem(), ii, key)
		} else if indexNewValue.Type().Kind() != reflect.Struct {
			extractNonStruct(parent, newValue.Field(i).Index(ii), oldValue.Field(i).Index(ii), ii, key)
		} else {
			extractChildren(parent, child, newValue.Field(i).Index(ii), oldValue.Field(i).Index(ii), &hasChildren)
		}
	}
	// new value has more data than the olddata
	if newValue.Field(i).Len() > oldValue.Field(i).Len() {
		for ii := shortest; ii < newValue.Field(i).Len(); ii++ {
			// create a dataStruct of the type in the slice to append to the oldValue slice
			dataStruct := reflect.New(newValue.Field(i).Index(ii).Type()).Elem()
			// append that value to the oldValue slice
			oldValue.Field(i).Set(reflect.Append(oldValue.Field(i), dataStruct))
			// now extract
			indexNewValue := newValue.Field(i).Index(ii)
			if indexNewValue.Type().Kind() == reflect.Ptr {
				extractNonStruct(parent, newValue.Field(i).Index(ii).Elem(), oldValue.Field(i).Index(ii).Elem(), ii, key)
			} else if indexNewValue.Type().Kind() != reflect.Struct {
				extractNonStruct(parent, newValue.Field(i).Index(ii), oldValue.Field(i).Index(ii), ii, key)
			} else {
				extractChildren(parent, child, newValue.Field(i).Index(ii), oldValue.Field(i).Index(ii), &hasChildren)
			}
		}
	}
	// oldValue slice is longer than the newValue so items were deleted
	if oldValue.Field(i).Len() > newValue.Field(i).Len() {
		for ii := shortest; ii < oldValue.Field(i).Len(); ii++ {
			deleteNonStruct(parent, ii, key)
		}
	}

	reflect.Copy(oldValue.Field(i), newValue.Field(i))
}

func extractLevelMap(parent *Differences, newValue reflect.Value, oldValue reflect.Value, i int, oldType reflect.Type, child *Differences, hasChildren bool, key string) {
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

		if newValue.Field(i).MapIndex(k).Type().Kind() == reflect.Ptr {
			extractChildren(parent, child, newValue.Field(i).MapIndex(k).Elem(), oldValue.Field(i).MapIndex(k).Elem(), &hasChildren)
		} else if newValue.Field(i).MapIndex(k).Type().Kind() != reflect.Struct {
			extractNonStruct(parent, newValue.Field(i).MapIndex(k), oldValue.Field(i).MapIndex(k), makeString(k), key)
		} else {
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

func extractLevelPointer(parent *Differences, newValue reflect.Value, oldValue reflect.Value, i int, child *Differences, hasChildren bool) {
	if newValue.Field(i).IsNil() {
		if !oldValue.Field(i).IsNil() {
			child.delete = true
			parent.children = append(parent.children, child)
			if oldValue.CanSet() {
				oldValue.Set(newValue)
			}
		}
		return
	}
	if oldValue.Field(i).IsNil() {
		oldValue.Field(i).Set(reflect.New(newValue.Field(i).Elem().Type()))
	}
	extractChildren(parent, child, newValue.Field(i).Elem(), oldValue.Field(i).Elem(), &hasChildren)
}

func setValue(va reflect.Value, child *Differences) {
	switch va.Kind() {
	case reflect.Invalid:
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value := va.Int()
		child.value.Int64 = &value
	case reflect.Bool:
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value := va.Uint()
		child.value.Uint64 = &value
	case reflect.Uintptr:
	case reflect.Float32:
		value := float32(va.Float())
		child.value.Float32 = &value
	case reflect.Float64:
		value := va.Float()
		child.value.Float64 = &value
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
		child.value.String_ = &value
	case reflect.Struct:
	case reflect.UnsafePointer:
	}
}

func extractNonStruct(parent *Differences, newValue reflect.Value, oldValue reflect.Value, index interface{}, key string) {
	if !equal(newValue, oldValue) {
		key = fmt.Sprintf("%s[%v]", key, index)
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
		child := &Differences{
			key: append(parent.key, &control.Key{
				Key:   key,
				Index: &indexObject,
			}),
			value: &control.Object{},
		}
		setValue(newValue, child)

		parent.children = append(parent.children, child)
		if oldValue.CanSet() {
			oldValue.Set(newValue)
		}
	}

}

func deleteNonStruct[i int | string](parent *Differences, index i, key string) {
	key = fmt.Sprintf("%s[%v]", key, index)
	child := &Differences{
		key: append(parent.key, &control.Key{
			Key: key,
		}),
		delete: true,
	}
	parent.children = append(parent.children, child)
}

func extractChildren(parent *Differences, child *Differences, newValue reflect.Value, oldValue reflect.Value, hasChildren *bool) {
	extractLevel(child, newValue, oldValue)
	if !*hasChildren && child.children != nil {
		parent.children = append(parent.children, child)
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
