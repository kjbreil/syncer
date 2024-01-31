package injector

import (
	"fmt"
	"reflect"
	"strconv"
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

func (inj *Injector) Add(inja *Injectable) error {
	return add(inj.data, inja)
}
func add(data any, inja *Injectable) error {
	v := reflect.ValueOf(data)

	// if its a pointer follow to the real data
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	if t.Name() != inja.key[0].Key {
		return fmt.Errorf("type mismatch %s!= %s", t.Name(), inja.key[0].Key)
	}
	inja.Advance()

	for i := 0; i < v.NumField(); i++ {
		if t.Field(i).Name == inja.key[0].Key {
			va := v.Field(i)
			if len(inja.key) > 1 {
				return add(va.Interface(), inja.Advance())
			} else {
				if va.CanSet() {
					return setValue(va, inja)
				}
			}
		}
	}
	return fmt.Errorf("key not found in data")
}

func setValue(va reflect.Value, inja *Injectable) error {
	switch va.Kind() {
	case reflect.String:
		va.SetString(inja.value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.Atoi(inja.value)
		if err != nil {
			return err
		}
		va.SetInt(int64(i))
	case reflect.Slice:
		if inja.key[0].Index == nil {
			return fmt.Errorf("slice type without index")
		}
		indexInt := int(inja.key[0].Index.GetInt64())

		fmt.Println(va.Type())
		// create a slice of the elements needed
		diff := indexInt + 1 - va.Len()
		if diff > 0 {
			newSlice := reflect.MakeSlice(va.Type(), diff, diff)
			va.Set(reflect.AppendSlice(va, newSlice))
		}

		return setValue(va.Index(indexInt), inja)
	case reflect.Map:
		if inja.key[0].Index == nil {
			return fmt.Errorf("map type without index")
		}

		keyType := va.Type().Key()
		valueType := va.Type().Elem().Kind()
		// convert the index into the right type
		var iKey reflect.Value
		var iValue reflect.Value

		switch keyType.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			indexInt := int(inja.key[0].Index.GetInt64())
			iKey = reflect.ValueOf(indexInt)
		case reflect.String:
			iKey = reflect.ValueOf(inja.key[0].Index.GetString_()).Elem()
		default:
			panic("I don't know what i'm doing here")
		}

		switch valueType {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			indexInt, err := strconv.Atoi(inja.value)
			if err != nil {
				return fmt.Errorf("could not convert %s to int", inja.value)
			}
			iValue = reflect.ValueOf(indexInt)
		case reflect.String:
			iValue = reflect.ValueOf(inja.value)
		default:
			panic("I don't know what i'm doing here")
		}
		zeroValue := reflect.Value{}
		if iKey == zeroValue || iValue == zeroValue {
			return fmt.Errorf("keys or value Value is zero")
		}
		va.SetMapIndex(iKey, iValue)
		return nil
		// indexInt, err := strconv.Atoi(*inja.key[0].index)
		// if err!= nil {
		//     return fmt.Errorf("could not convert %s to int", *inja.key[0].index)
		// }
		// return setValue(va.MapIndex(reflect.ValueOf(indexInt)), inja)
	case reflect.Struct:
	default:
		panic("setValue used on unknown type")
	}
	return nil
}
