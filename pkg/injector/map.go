package injector

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/kjbreil/syncer/pkg/control"
)

func injectMap(va reflect.Value, entry *control.Entry) error {
	if entry.GetCurrKey().HasNoIndex() {
		// no index on a map key and remove type make map nil
		if entry.GetRemove() {
			va.Set(reflect.New(va.Type()).Elem())
			return nil
		}
		return errors.New("map type without index")
	}

	// get the key and value types
	keyType := va.Type().Key()

	// if the map is empty, create a new map with the correct key and value types
	if va.Len() == 0 {
		// get the value type of the map element
		vt := va.Type().Elem()
		// create a map type with the correct key and value types
		mapType := reflect.MapOf(keyType, vt)
		// if the map can be set, create a new map with the correct size
		if va.CanSet() {
			va.Set(reflect.MakeMapWithSize(mapType, 0))
		}
	}

	mapKey, err := makeMapKey(keyType, entry)
	if err != nil {
		return err
	}

	if entry.GetRemove() && entry.IsLastKeyIndex() {
		va.SetMapIndex(mapKey, reflect.Value{})
		return nil
	}

	// create a variable to hold the indexed value
	mapValue, err := makeMapValue(va, entry, mapKey)
	if err != nil {
		return err
	}

	// check if the indexed key is zero
	if !mapKey.IsValid() {
		// return an error if the indexed key is zero
		return errors.New("mapKey is not valid")
	}

	// check if the indexed value is zero
	if !mapValue.IsValid() {
		// return an error if the indexed value is zero
		return errors.New("mapValue is not valid")
	}

	// set the indexed value in the map
	va.SetMapIndex(mapKey, mapValue)
	// return nil

	return nil
}

func makeMapValue(va reflect.Value, entry *control.Entry, mapKey reflect.Value) (reflect.Value, error) {
	mapValue := reflect.New(va.Type().Elem()).Elem()

	// get the current value if it exits in the map
	currValue := va.MapIndex(mapKey)
	// if we got a valid value then assign mapValue to the current value
	if currValue.IsValid() {
		mapValue.Set(currValue)
	}
	var err error
	switch mapValue.Kind() {
	case reflect.Struct:
		err = add(mapValue, entry)
	default:
		err = add(mapValue, entry.Advance())
	}
	if err != nil {
		return mapValue, err
	}

	return mapValue, nil
}

func makeMapKey(keyType reflect.Type, entry *control.Entry) (reflect.Value, error) {
	// create a variable to hold the indexed key
	var mapKey reflect.Value

	// based on the key type, set the indexed key
	switch keyType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		mapKey = reflect.ValueOf(int(entry.GetCurrentIndex().GetInt64()))
	case reflect.String:
		mapKey = reflect.ValueOf(entry.GetCurrentIndex().GetString_())
	case reflect.Bool:
		mapKey = reflect.ValueOf(entry.GetCurrentIndex().GetBool())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		mapKey = reflect.ValueOf(uint(entry.GetCurrentIndex().GetUint64()))
	case reflect.Float32:
		mapKey = reflect.ValueOf(float32(entry.GetCurrentIndex().GetFloat32()))
	case reflect.Float64:
		mapKey = reflect.ValueOf(entry.GetCurrentIndex().GetFloat64())
	default:
		return reflect.Value{}, fmt.Errorf("cannot create key of type %s", keyType.Kind())
	}

	// this handles if the key is a type definition, the keyType.Kind() will be the base type of the type definition
	// however the key needs to be set to the type defined
	newK := reflect.New(keyType).Elem()
	newK.Set(mapKey.Convert(keyType))

	mapKey = newK
	return mapKey, nil
}
