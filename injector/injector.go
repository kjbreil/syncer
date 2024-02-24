package injector

import (
	"errors"
	"fmt"
	"github.com/kjbreil/syncer/control"
	"reflect"
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

// AddAll adds multiple entries to the data.
func (inj *Injector) AddAll(entries control.Entries) error {
	for _, e := range entries {
		err := inj.Add(e)
		if err != nil {
			return err
		}
	}
	return nil
}

// Add adds a control entry to the data.
func (inj *Injector) Add(entry *control.Entry) error {

	v := reflect.ValueOf(inj.data)

	// if it is a pointer follow to the real data
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	if t.Name() != entry.GetKey()[entry.GetKeyI()].GetKey() {
		return fmt.Errorf("injector top level type mismatch %s  != %s", t.Name(), entry.GetKey()[entry.GetKeyI()].GetKey())
	}

	return add(v, entry.Advance())
}

// Add adds a control entry to the data. Based on the data type either travels down the key's or sets the value
func add(v reflect.Value, entry *control.Entry) error {
	v = reflect.Indirect(v)

	if va := v.FieldByName(entry.GetKey()[entry.GetKeyI()].GetKey()); va.IsValid() {
		va = reflect.Indirect(va)
		switch va.Kind() {
		case reflect.Slice:
			return setValueSlice(va, entry)
		case reflect.Map:
			return setValueMap(va, entry)
		case reflect.Interface:
			va = va.Elem()
			fallthrough
		case reflect.Struct:
			return add(va, entry.Advance())
		case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			if va.CanSet() {
				return setValue(va, entry)
			} else {
				return fmt.Errorf("cannot set value in add for %s", entry.GetKey()[entry.GetKeyI()].GetKey())
			}
		default:
			return fmt.Errorf("cannot add value for %s type %s", entry.GetKey()[entry.GetKeyI()].GetKey(), va.Kind())
		}
	}
	// return errors.New("key not found in data")

	return errors.New("injector add reached end when it should not have")
}

func setValue(va reflect.Value, entry *control.Entry) error {
	switch va.Kind() {
	case reflect.Slice:
		return setValueSlice(va, entry)
	case reflect.Map:
		return setValueMap(va, entry)
	case reflect.Interface:
		va = va.Elem()
		fallthrough
	case reflect.Struct:
		return add(va, entry.Advance())
	case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return entry.GetValue().SetValue(va)
	default:
		return fmt.Errorf("cannot set value in setValue for %s of type %s", entry.GetKey()[entry.GetKeyI()].GetKey(), va.Kind())
	}
}

// setValueMap sets a value in a map based on the control entry
func setValueMap(va reflect.Value, entry *control.Entry) error {
	// check if the key is indexed
	if entry.GetCurrIndexObjects() == nil {
		// no index on a map key and remove type make map nil
		if entry.Remove {
			va.Set(reflect.New(va.Type()).Elem())
			return nil
		}
		// return an error if the key is not indexed
		return errors.New("map type without index")
	}

	// get the key and value types
	keyType := va.Type().Key()
	valueType := va.Type().Elem().Kind()
	// if valueType is pointer we need to get the kind of the actual valueType
	if valueType == reflect.Ptr {
		valueType = va.Type().Elem().Elem().Kind()
	}

	// create a variable to hold the indexed key
	var iKey reflect.Value

	// based on the key type, set the indexed key
	switch keyType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		iKey = reflect.ValueOf(int(entry.GetCurrIndexObjects()[0].GetInt64()))
	case reflect.String:
		// get the index as a string
		iKey = reflect.ValueOf(entry.GetCurrIndexObjects()[0].GetString_())
		// TODO: Handle below casese
	case reflect.Bool:
		iKey = reflect.ValueOf(entry.GetCurrIndexObjects()[0].GetBool())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		iKey = reflect.ValueOf(int(entry.GetCurrIndexObjects()[0].GetUint64()))
	case reflect.Float32:
		iKey = reflect.ValueOf(int(entry.GetCurrIndexObjects()[0].GetFloat32()))
	case reflect.Float64:
		iKey = reflect.ValueOf(int(entry.GetCurrIndexObjects()[0].GetFloat64()))
	// case reflect.Pointer:
	// case reflect.Interface:
	// case reflect.Struct:
	default:
		return fmt.Errorf("cannot create key of type %s", keyType.Kind())
	}

	newK := reflect.New(keyType).Elem()
	newK.Set(iKey.Convert(keyType))

	iKey = newK

	if entry.Last() && entry.Remove {
		va.SetMapIndex(iKey, reflect.Value{})
		return nil
	}

	// create a variable to hold the indexed value
	var iValue reflect.Value

	// based on the value type, set the indexed value
	switch valueType {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// get the index as an int
		indexInt := int(entry.GetValue().GetInt64())
		// set the indexed value to the int
		iValue = reflect.ValueOf(indexInt)
	case reflect.String:
		// get the index as a string
		iValue = reflect.ValueOf(entry.GetValue().GetString_())
	case reflect.Map:
		// if the map is empty, create a new map with the correct key and value types
		if va.Len() == 0 {
			// get the key and value types of the map element
			keyType := va.Type().Elem().Key()
			vt := va.Type().Elem().Elem()
			// create a map type with the correct key and value types
			mapType := reflect.MapOf(keyType, vt)
			// create a new map with the correct size
			iValue = reflect.MakeMapWithSize(mapType, 0)
		} else {
			// get the value at the indexed key
			iValue = va.MapIndex(iKey)
		}
		// if the key is indexed, advance the control entry
		if len(entry.GetCurrIndexObjects()) > 1 {
			entry.Key[entry.GetKeyI()].Index = entry.GetCurrIndexObjects()[1:]
		}
		// set the indexed value based on the advanced control entry
		err := setValue(iValue, entry)
		// return an error if one occurs
		if err != nil {
			return err
		}

	case reflect.Struct, reflect.Interface:
		// generate a value to be used for the new value at the key
		// map values obtained by reflect cannot be set since they are not addressable so we need to get the current
		// value and set a new value to the current value then modify said value and then assign it to the map
		iValue = reflect.New(va.Type().Elem()).Elem()

		// get the current value if it exits in the map
		currValue := va.MapIndex(iKey)
		// copy the current value to the iValue so we can modify
		if currValue.IsValid() {
			iValue.Set(currValue)
		}
		if iValue.Kind() == reflect.Interface {
			if !currValue.IsValid() {
				panic("cannot create a interface object with the entry given")
			}
			iValue = iValue.Elem()
		}

		// add the indexed value based on the advanced control entry
		err := add(iValue, entry.Advance())
		// return an error if one occurs
		if err != nil {
			return err
		}
	default:
		// panic if the value type is not supported
		panic("I don't know what I'm doing here")
	}

	// check if the indexed key is zero
	if !iKey.IsValid() {
		// return an error if the indexed key is zero
		return errors.New("iKey is not valid")
	}

	// check if the indexed value is zero
	if !iValue.IsValid() {
		// return an error if the indexed value is zero
		return errors.New("iValue is not valid")
	}
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
	// set the indexed value in the map
	va.SetMapIndex(iKey, iValue)
	// return nil
	return nil
}

func setValueSlice(va reflect.Value, entry *control.Entry) error {
	if len(entry.GetCurrIndexObjects()) == 0 {
		// no index on a map key and remove type make map nil
		if entry.Remove {
			va.Set(reflect.New(va.Type()).Elem())
			return nil
		}
		return errors.New("slice type without index")
	}
	// if entry is a delete entry then delete the current index+
	indexInt := int(entry.GetCurrIndexObjects()[0].GetInt64())
	if entry.GetRemove() {
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
	if len(entry.GetCurrIndexObjects()) > 1 {
		entry.AdvanceCurrKeyIndex()
	}
	return setValue(va.Index(indexInt), entry)
}
