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
	switch vv := v.(type) {
	case string:
		return &Object{String_: &vv}
	case *string:
		return &Object{String_: v.(*string)}
	case int:
		vvv := int64(vv)
		return &Object{Int64: &vvv}
	case int8:
		vvv := int64(vv)
		return &Object{Int64: &vvv}
	case int16:
		vvv := int64(vv)
		return &Object{Int64: &vvv}
	case int32:
		vvv := int64(vv)
		return &Object{Int64: &vvv}
	case int64:
		return &Object{Int64: &vv}
	case uint:
		vvv := uint64(vv)
		return &Object{Uint64: &vvv}
	case uint8:
		vvv := uint64(vv)
		return &Object{Uint64: &vvv}
	case uint16:
		vvv := uint64(vv)
		return &Object{Uint64: &vvv}
	case uint32:
		vvv := uint64(vv)
		return &Object{Uint64: &vvv}
	case uint64:
		return &Object{Uint64: &vv}
	case float32:
		return &Object{Float32: &vv}
	case float64:
		return &Object{Float64: &vv}
	case bool:
		return &Object{Bool: &vv}
	case []byte:
		return &Object{Bytes: vv}
	case Object:
		return &vv
	case *Object:
		return vv
	}

	return nil
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
