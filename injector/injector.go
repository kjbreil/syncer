package injector

import (
	"fmt"
	"github.com/kjbreil/syncer/control"
	"reflect"
)

type Injector struct {
	data any
}

var (
	ErrNotPointer = fmt.Errorf("data is not a pointer")
)

func New(data any) (*Injector, error) {
	if reflect.ValueOf(data).Kind() != reflect.Ptr {
		return nil, ErrNotPointer
	}

	return &Injector{
		data: data,
	}, nil
}

func (inj *Injector) Add(ctrl *control.Entry) error {
	return add(inj.data, ctrl)
}

func (inj *Injector) Data() any {
	return inj.data
}
func add(data any, ctrl *control.Entry) error {
	v := reflect.ValueOf(data)

	// if its a pointer follow to the real data
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	if t.Name() != ctrl.Key[0].Key {
		return fmt.Errorf("type mismatch %s!= %s", t.Name(), ctrl.Key[0].Key)
	}
	ctrl.Advance()

	for i := 0; i < v.NumField(); i++ {
		if t.Field(i).Name == ctrl.Key[0].Key {
			va := v.Field(i)
			if len(ctrl.Key) > 1 {
				return add(va.Interface(), ctrl.Advance())
			} else if va.CanSet() {
				return setValue(va, ctrl)
			}
		}
	}
	return fmt.Errorf("key not found in data")
}

func setValue(va reflect.Value, ctrl *control.Entry) error {
	switch va.Kind() {
	case reflect.Slice:
		return setValueSlice(va, ctrl)
	case reflect.Map:
		return setValueMap(va, ctrl)
	case reflect.Struct:
	default:
		return ctrl.Value.SetValue(va)
	}
	return nil
}

func setValueMap(va reflect.Value, ctrl *control.Entry) error {
	if ctrl.Key[0].Index == nil {
		return fmt.Errorf("map type without index")
	}

	keyType := va.Type().Key()
	valueType := va.Type().Elem().Kind()
	// convert the index into the right type
	var iKey reflect.Value
	var iValue reflect.Value

	switch keyType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		indexInt := int(ctrl.Key[0].Index.GetInt64())
		iKey = reflect.ValueOf(indexInt)
	case reflect.String:
		iKey = reflect.ValueOf(ctrl.Key[0].Index.GetString_()).Elem()
	default:
		panic("I don't know what i'm doing here")
	}

	switch valueType {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		indexInt := int(ctrl.Value.GetInt64())
		iValue = reflect.ValueOf(indexInt)
	case reflect.String:
		iValue = reflect.ValueOf(ctrl.Value.GetString_())
	default:
		panic("I don't know what i'm doing here")
	}
	zeroValue := reflect.Value{}
	if iKey == zeroValue || iValue == zeroValue {
		return fmt.Errorf("keys or Value Value is zero")
	}
	va.SetMapIndex(iKey, iValue)
	return nil
}

func setValueSlice(va reflect.Value, ctrl *control.Entry) error {
	if ctrl.Key[0].Index == nil {
		return fmt.Errorf("slice type without index")
	}
	indexInt := int(ctrl.Key[0].Index.GetInt64())
	// create a slice of the elements needed
	diff := indexInt + 1 - va.Len()
	if diff > 0 {
		newSlice := reflect.MakeSlice(va.Type(), diff, diff)
		va.Set(reflect.AppendSlice(va, newSlice))
	}

	return setValue(va.Index(indexInt), ctrl)
}
