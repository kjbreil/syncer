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
	return add(inj.data, ctrl)
}

func add(data any, ctrl *control.Entry) error {
	v := reflect.ValueOf(data)

	// if it is a pointer follow to the real data
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	if t.Name() != ctrl.GetKey()[0].GetKey() {
		return fmt.Errorf("type mismatch %s!= %s", t.Name(), ctrl.GetKey()[0].GetKey())
	}
	ctrl.Advance()

	for i := 0; i < v.NumField(); i++ {
		if t.Field(i).Name == ctrl.GetKey()[0].GetKey() {
			va := v.Field(i)
			if len(ctrl.GetKey()) > 1 {
				return add(va.Interface(), ctrl.Advance())
			} else if va.CanSet() {
				return setValue(va, ctrl)
			}
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
		indexInt := int(ctrl.GetKey()[0].GetIndex().GetInt64())
		iKey = reflect.ValueOf(indexInt)
	case reflect.String:
		iKey = reflect.ValueOf(ctrl.GetKey()[0].GetIndex().GetString_()).Elem()
	default:
		panic("I don't know what i'm doing here")
	}

	switch valueType {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		indexInt := int(ctrl.GetValue().GetInt64())
		iValue = reflect.ValueOf(indexInt)
	case reflect.String:
		iValue = reflect.ValueOf(ctrl.GetValue().GetString_())
	default:
		panic("I don't know what i'm doing here")
	}
	if iKey.IsZero() || iValue.IsZero() {
		return errors.New("keys or Value Value is zero")
	}
	va.SetMapIndex(iKey, iValue)
	return nil
}

func setValueSlice(va reflect.Value, ctrl *control.Entry) error {
	if ctrl.GetKey()[0].GetIndex() == nil {
		return errors.New("slice type without index")
	}
	// if ctrl is a delete entry then delete the current index+
	indexInt := int(ctrl.GetKey()[0].GetIndex().GetInt64())
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

	return setValue(va.Index(indexInt), ctrl)
}
