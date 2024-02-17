package injector

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/kjbreil/syncer/control"
)

type Injector struct {
	data any
}

var (
	ErrNotPointer = errors.New("data is not a pointer")
)

// New creates a new injector with the given data.
// If the data is not a pointer, an error is returned.
func New(data any) (*Injector, error) {
	if reflect.ValueOf(data).Kind() != reflect.Ptr {
		return nil, ErrNotPointer
	}

	return &Injector{
		data: data,
	}, nil
}

// Add adds a control entry to the data.
func (inj *Injector) Add(ctrl *control.Entry) error {

	v := reflect.ValueOf(inj.data)

	// if it is a pointer follow to the real data
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	if t.Name() != ctrl.GetKey()[0].GetKey() {
		return fmt.Errorf("injector top level type mismatch %s  != %s", t.Name(), ctrl.GetKey()[0].GetKey())
	}

	return add(v, ctrl.Advance())
}

func add(v reflect.Value, ctrl *control.Entry) error {
	v = reflect.Indirect(v)

	if va := v.FieldByName(ctrl.GetKey()[0].GetKey()); va.IsValid() {
		if len(ctrl.GetKey()) > 1 {
			return add(va, ctrl.Advance())
		} else if va.CanSet() {
			return setValue(va, ctrl)
		}
	}

	return errors.New("key not found in data")
}

func setValue(va reflect.Value, ctrl *control.Entry) error {
	switch va.Kind() {
	case reflect.Slice:
		return setValueSlice(va, ctrl)
	case reflect.Map:
		return setValueMap(va, ctrl)
	case reflect.Struct:
	default:
		return ctrl.GetValue().SetValue(va)
	}
	return nil
}

func setValueMap(va reflect.Value, ctrl *control.Entry) error {
	if ctrl.GetKey()[0].GetIndex() == nil {
		return errors.New("map type without index")
	}

	keyType := va.Type().Key()
	valueType := va.Type().Elem().Kind()
	// convert the index into the right type
	var iKey reflect.Value
	var iValue reflect.Value

	switch keyType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		indexInt := int(ctrl.GetKey()[0].GetIndex()[0].GetInt64())
		iKey = reflect.ValueOf(indexInt)
	case reflect.String:
		iKey = reflect.ValueOf(ctrl.GetKey()[0].GetIndex()[0].GetString_())
	default:
		panic("I don't know what i'm doing here")
	}

	switch valueType {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		indexInt := int(ctrl.GetValue().GetInt64())
		iValue = reflect.ValueOf(indexInt)
	case reflect.String:
		iValue = reflect.ValueOf(ctrl.GetValue().GetString_())
	case reflect.Map:
		// TODO: Check if value exists and use that otherwise create new map
		if va.Len() == 0 {
			keyType = va.Type().Elem().Key()
			vt := va.Type().Elem().Elem()
			mapType := reflect.MapOf(keyType, vt)
			iValue = reflect.MakeMapWithSize(mapType, 0)
		} else {
			iValue = va.MapIndex(iKey)
		}
		if len(ctrl.GetKey()[0].GetIndex()) > 1 {
			ctrl.Key[0].Index = ctrl.GetKey()[0].GetIndex()[1:]
		}
		err := setValue(iValue, ctrl)
		if err != nil {
			return err
		}
	default:
		panic("I don't know what i'm doing here")
	}
	if iKey.IsZero() || iValue.IsZero() {
		return errors.New("keys or Value Value is zero")
	}
	if va.Len() == 0 {
		vt := va.Type().Elem()
		mapType := reflect.MapOf(keyType, vt)
		if va.CanSet() {
			va.Set(reflect.MakeMapWithSize(mapType, 0))
		}
	}

	va.SetMapIndex(iKey, iValue)
	return nil
}

func setValueSlice(va reflect.Value, ctrl *control.Entry) error {
	if len(ctrl.GetKey()[0].GetIndex()) == 0 {
		return errors.New("slice type without index")
	}
	// if ctrl is a delete entry then delete the current index+
	indexInt := int(ctrl.GetKey()[0].GetIndex()[0].GetInt64())
	if ctrl.GetRemove() {
		newSlice := reflect.MakeSlice(va.Type(), indexInt, indexInt)
		reflect.Copy(newSlice, va)
		va.Set(newSlice)
		return nil
	}

	// create a slice of the elements needed
	diff := indexInt + 1 - va.Len()
	if diff > 0 {
		newSlice := reflect.MakeSlice(va.Type(), diff, diff)
		va.Set(reflect.AppendSlice(va, newSlice))
	}

	// e.Key = e.GetKey()[1:]
	if len(ctrl.GetKey()[0].GetIndex()) > 1 {
		ctrl.Key[0].Index = ctrl.GetKey()[0].GetIndex()[1:]
	}
	return setValue(va.Index(indexInt), ctrl)
}
