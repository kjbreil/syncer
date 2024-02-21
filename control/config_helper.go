package control

import (
	"fmt"
	"reflect"
)

func (e *Entry) Advance() *Entry {
	// e.Key = e.GetKey()[1:]
	e.KeyI++
	return e
}

func NewObjects(v any, oo ...*Object) []*Object {
	objects := make([]*Object, len(oo)+1)
	objects[0] = NewObject(v)
	for i, o := range oo {
		objects[i+1] = o
	}

	return objects
}

//nolint:funlen
func NewObject(v any) *Object {
	va := reflect.Indirect(reflect.ValueOf(v))
	if !va.IsValid() {
		return nil
	}

	// TODO: Should panic if is not a valid type for NewObject
	var o Object
	switch va.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		o.Int64 = MakePtr(va.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		o.Uint64 = MakePtr(va.Uint())
	case reflect.Float32:
		o.Float32 = MakePtr(float32(va.Float()))
	case reflect.Float64:
		o.Float64 = MakePtr(va.Float())
	case reflect.String:
		o.String_ = MakePtr(va.String())
	case reflect.Bool:
		o.Bool = MakePtr(va.Bool())
	case reflect.Slice:
		if va.Type().Elem().Kind() == reflect.Uint8 {
			o.Bytes = va.Bytes()
		}
	}

	switch vv := v.(type) {
	case Object:
		return &vv
	case *Object:
		return vv
	}

	return &o
}

func (o *Object) SetValue(va reflect.Value) error {
	switch va.Kind() {
	case reflect.String:
		va.SetString(o.GetString_())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i := int(o.GetInt64())
		va.SetInt(int64(i))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i := uint(o.GetUint64())
		va.SetUint(uint64(i))
	case reflect.Float32:
		f := o.GetFloat32()
		va.SetFloat(float64(f))
	case reflect.Float64:
		f := o.GetFloat64()
		va.SetFloat(f)
	case reflect.Bool:
		b := o.GetBool()
		va.SetBool(b)
	default:
		return fmt.Errorf("SetValue used on unknown kind: %s", va.Kind())
	}
	return nil
}
