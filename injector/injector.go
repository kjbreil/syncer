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

type injFn func(va reflect.Value, entry *control.Entry) error

var injFns map[reflect.Kind]injFn

func init() {
	injFns = map[reflect.Kind]injFn{
		reflect.Bool:       injectPrimitive,
		reflect.Int:        injectPrimitive,
		reflect.Int8:       injectPrimitive,
		reflect.Int16:      injectPrimitive,
		reflect.Int32:      injectPrimitive,
		reflect.Int64:      injectPrimitive,
		reflect.Uint:       injectPrimitive,
		reflect.Uint8:      injectPrimitive,
		reflect.Uint16:     injectPrimitive,
		reflect.Uint32:     injectPrimitive,
		reflect.Uint64:     injectPrimitive,
		reflect.Uintptr:    injectPrimitive,
		reflect.Float32:    injectPrimitive,
		reflect.Float64:    injectPrimitive,
		reflect.Complex64:  injectPrimitive,
		reflect.Complex128: injectPrimitive,
		reflect.String:     injectPrimitive,
		reflect.Array:      injectArray,
		reflect.Map:        injectMap,
		reflect.Ptr:        injectPointer,
		reflect.Interface:  injectInterface,
		reflect.Slice:      injectSlice,
		reflect.Struct:     injectStruct,
	}
}

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

	return add(v, entry)
	// return add(v, entry.Advance())
}

// Add adds a control entry to the data. Based on the data type either travels down the key's or sets the value
func add(v reflect.Value, entry *control.Entry) error {

	var err error
	if iFn, ok := injFns[v.Kind()]; ok {
		err = iFn(v, entry)
		if err != nil {
			return err
		}
	}

	return nil
}
