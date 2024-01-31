package control

import (
	"fmt"
	"reflect"
)

func (e *Entry) Advance() *Entry {
	e.Key = e.Key[1:]
	return e
}

func NewObject(v any) *Object {
	switch v.(type) {
	case string:
		vv := v.(string)
		return &Object{String_: &vv}
	case *string:
		return &Object{String_: v.(*string)}
	case int:
		vv := int64(v.(int))
		return &Object{Int64: &vv}
	case int32:
		vv := int64(v.(int32))
		return &Object{Int64: &vv}
	case int64:
		vv := v.(int64)
		return &Object{Int64: &vv}
	case uint:
		vv := uint64(v.(uint))
		return &Object{Uint64: &vv}
	case uint32:
		vv := uint64(v.(uint32))
		return &Object{Uint64: &vv}
	case uint64:
		vv := v.(uint64)
		return &Object{Uint64: &vv}
	case float32:
		vv := v.(float32)
		return &Object{Float32: &vv}
	case float64:
		vv := v.(float64)
		return &Object{Float64: &vv}
	case bool:
		vv := v.(bool)
		return &Object{Bool: &vv}
	case []byte:
		return &Object{Bytes: v.([]byte)}
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
