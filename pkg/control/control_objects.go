package control

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"strings"
)

type Objects []*Object

func (o Objects) Equals(other Objects) bool {
	if len(o) != len(other) {
		return false
	}
	for i, io := range o {
		if !io.Equals(other[i]) {
			return false
		}
	}
	return true
}

func (o Objects) Struct() string {
	var sb strings.Builder

	sb.WriteString("control.NewObjects(")
	for i, o := range o {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(o.Struct())
	}
	sb.WriteString(")")

	return sb.String()
}

func (o *Object) Equals(other *Object) bool {
	if o == nil && other == nil {
		return true
	}

	if o == nil || other == nil {
		return false
	}

	if o.GetString_() != other.GetString_() {
		return false
	}

	if o.GetInt64() != other.GetInt64() {
		return false
	}

	if o.GetUint64() != other.GetUint64() {
		return false
	}

	if o.GetFloat32() != other.GetFloat32() {
		return false
	}

	if o.GetFloat64() != other.GetFloat64() {
		return false
	}

	if o.GetBool() != other.GetBool() {
		return false
	}

	if !bytes.Equal(o.GetBytes(), other.GetBytes()) {
		return false
	}

	return true
}

func (o *Object) Struct() string {
	if o == nil {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("control.NewObject(")
	switch {
	case o.String_ != nil:
		sb.WriteString(fmt.Sprintf("control.MakePtr(\"%s\")", o.GetString_()))
	case o.Int64 != nil:
		sb.WriteString(fmt.Sprintf("control.MakePtr(int64(%d))", o.GetInt64()))
	case o.Uint64 != nil:
		sb.WriteString(fmt.Sprintf("control.MakePtr(uint64(%d))", o.GetUint64()))
	case o.Float32 != nil:
		sb.WriteString(fmt.Sprintf("control.MakePtr(float32(%f))", o.GetFloat32()))
	case o.Float64 != nil:
		sb.WriteString(fmt.Sprintf("control.MakePtr(float64(%f))", o.GetFloat64()))
	case o.Bool != nil:
		sb.WriteString(fmt.Sprintf("control.MakePtr(%t)", o.GetBool()))
	case o.Bytes != nil:
		sb.WriteString(fmt.Sprintf("[]byte(%v)", o.GetBytes()))
	}
	sb.WriteString(")")
	return sb.String()
}

func NewObjects(v any, oo ...*Object) []*Object {
	objects := make([]*Object, len(oo)+1)
	objects[0] = NewObject(v)
	for i, o := range oo {
		objects[i+1] = o
	}

	return objects
}

func NewObject(v any) *Object {
	va := reflect.Indirect(reflect.ValueOf(v))
	if !va.IsValid() {
		return nil
	}

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
	case reflect.Complex64:
		c := complex64(va.Complex())
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint32(buf[0:4], math.Float32bits(real(c)))
		binary.LittleEndian.PutUint32(buf[4:8], math.Float32bits(imag(c)))
		o.Bytes = buf
	case reflect.Complex128:
		c := va.Complex()
		buf := make([]byte, 16)
		binary.LittleEndian.PutUint64(buf[0:8], math.Float64bits(real(c)))
		binary.LittleEndian.PutUint64(buf[8:16], math.Float64bits(imag(c)))
		o.Bytes = buf
	case reflect.Slice:
		if va.Type().Elem().Kind() == reflect.Uint8 {
			o.Bytes = va.Bytes()
		}
	default:
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
		va.SetInt(o.GetInt64())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		va.SetUint(o.GetUint64())
	case reflect.Float32:
		f := o.GetFloat32()
		va.SetFloat(float64(f))
	case reflect.Float64:
		f := o.GetFloat64()
		va.SetFloat(f)
	case reflect.Bool:
		b := o.GetBool()
		va.SetBool(b)
	case reflect.Complex64:
		b := o.GetBytes()
		if len(b) == 8 {
			r := math.Float32frombits(binary.LittleEndian.Uint32(b[0:4]))
			i := math.Float32frombits(binary.LittleEndian.Uint32(b[4:8]))
			va.SetComplex(complex(float64(r), float64(i)))
		}
	case reflect.Complex128:
		b := o.GetBytes()
		if len(b) == 16 {
			r := math.Float64frombits(binary.LittleEndian.Uint64(b[0:8]))
			i := math.Float64frombits(binary.LittleEndian.Uint64(b[8:16]))
			va.SetComplex(complex(r, i))
		}
	default:
		return fmt.Errorf("SetValue used on unknown kind: %q", va.Kind())
	}
	return nil
}
