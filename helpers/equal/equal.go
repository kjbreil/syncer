package equal

import (
	"fmt"
	"reflect"
)

// Equal returns true if the two values are equal, false otherwise.
// Differs from reflect.Value.Equal in that it follows and compares the value behind pointers
// and will compare any type of int or uint or float against itself.
// floats do suffer from float math and generally a float32 does not match a float64
func Equal(n, o reflect.Value) bool {
	if !sameKind(n, o) {
		nVi, oVi := n.Interface(), o.Interface()
		fmt.Println(nVi, oVi)
		return false
	}
	// if both are invalid then they are Equal
	switch n.Kind() {
	case reflect.Pointer:
		return Equal(n.Elem(), o.Elem())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return n.Int() == o.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return n.Uint() == o.Uint()
	case reflect.String:
		return n.String() == o.String()
	case reflect.Bool:
		return n.Bool() == o.Bool()
	case reflect.Float32, reflect.Float64:
		return n.Float() == o.Float()
	case reflect.Complex64, reflect.Complex128:
		return n.Complex() == o.Complex()
	case reflect.Slice, reflect.Array:
		if n.Len() != o.Len() {
			return false
		}
		for i := 0; i < n.Len(); i++ {
			if !Equal(n.Index(i), o.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Struct:
		for i := 0; i < n.NumField(); i++ {
			if !Equal(n.Field(i), o.Field(i)) {
				return false
			}
		}
		return true
	case reflect.Interface:
		return Equal(n.Elem(), o.Elem())
	case reflect.Chan:
		if n.Type().ChanDir() != o.Type().ChanDir() {
			return false
		}
		if n.Type().Elem() != o.Type().Elem() {
			return false
		}
		return true
	case reflect.Func:
		nT, oT := n.Type(), o.Type()
		if nT != oT {
			return false
		}
		if nT.NumIn() != oT.NumIn() || nT.NumOut() != oT.NumOut() || nT.IsVariadic() != oT.IsVariadic() {
			return false
		}
		for i := 0; i < nT.NumIn(); i++ {
			if nT.In(i) != oT.In(i) {
				return false
			}
		}
		for i := 0; i < nT.NumOut(); i++ {
			if nT.Out(i) != oT.Out(i) {
				return false
			}
		}
		return true
	case reflect.Invalid:
		return o.Kind() == reflect.Invalid
	case reflect.Map:
		if n.Len() != o.Len() {
			return false
		}
		for _, k := range n.MapKeys() {
			if !Equal(n.MapIndex(k), o.MapIndex(k)) {
				return false
			}
		}

		return true
	case reflect.UnsafePointer:
		return n.UnsafePointer() == o.UnsafePointer()
	}
	return false
}

func Any[T any](n, o T) bool {
	return Equal(reflect.ValueOf(n), reflect.ValueOf(o))
}

func sameKind(n, o reflect.Value) bool {

	switch n.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return o.Kind() == reflect.Int || o.Kind() == reflect.Int8 || o.Kind() == reflect.Int16 || o.Kind() == reflect.Int32 || o.Kind() == reflect.Int64
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return o.Kind() == reflect.Uint || o.Kind() == reflect.Uint8 || o.Kind() == reflect.Uint16 || o.Kind() == reflect.Uint32 || o.Kind() == reflect.Uint64
	case reflect.Float32, reflect.Float64:
		return o.Kind() == reflect.Float32 || o.Kind() == reflect.Float64
	case reflect.Complex64, reflect.Complex128:
		return o.Kind() == reflect.Complex64 || o.Kind() == reflect.Complex128
	case reflect.Struct:
		if !o.IsValid() || !n.IsValid() {
			return false
		}
		return o.Type() == n.Type()
	case reflect.Interface:
		return sameKind(n.Elem(), o.Elem())
	default:
		return o.Kind() == n.Kind()
	}
}
