package deepcopy

import (
	"reflect"
)

// DeepCopy copies the value of a reflect.Value.
//
// The copied value is a deep copy, meaning that all nested values (e.g. pointers,
// slices, and maps) are copied as well.
//
// If the destination value is not addressable, a new addressable value is created.
//
// If the source value is a pointer, the pointer is dereferenced before copying. If
// the source value is an interface, its underlying value is copied.
//
// If the source value is a struct, all fields are copied recursively. If the source
// value is a map, all keys and values are copied recursively. If the source value is
// a slice, all elements are copied recursively.
//
// If the source value is not a supported type, it is copied directly.
//
// The returned value is the destination value.
func DeepCopy(src reflect.Value) reflect.Value {
	// Create a new destination value that is a copy of the source value.
	dst := reflect.Indirect(reflect.New(src.Type()))

	// Check if the destination value is addressable. If not return
	if !dst.CanSet() || !src.CanInterface() {
		return dst
	}

	switch src.Kind() {
	// If the source value is a pointer, dereference it and copy the pointed-to value.
	case reflect.Pointer:
		if src.Elem().IsValid() {
			subV := DeepCopy(src.Elem())
			dst.Set(subV.Addr())
		}
	// If the source value is an interface, copy the underlying value.
	case reflect.Interface:
		if src.Elem().IsValid() {
			subV := DeepCopy(src.Elem())
			dst.Set(subV)
		}
	// If the source value is a struct, copy all fields recursively.
	case reflect.Struct:
		for i := 0; i < src.NumField(); i++ {
			if !dst.Field(i).CanInterface() || !dst.Field(i).CanInterface() {
				continue
			}
			dstFV := DeepCopy(src.Field(i))
			dst.Field(i).Set(dstFV)
		}
	// If the source value is a map, copy all keys and values recursively.
	case reflect.Map:
		return deepCopyMap(src, dst)
	// If the source value is a slice, copy all elements recursively.
	case reflect.Slice:
		return deepCopySlice(src, dst)
	// For all other value kinds, copy the value directly.
	default:
		dst.Set(src)
	}

	// Return the destination value.
	return dst
}

func deepCopyMap(src reflect.Value, dst reflect.Value) reflect.Value {
	dst.Set(reflect.MakeMapWithSize(src.Type(), src.Len()))
	for _, k := range src.MapKeys() {
		srcKV := src.MapIndex(k)
		dstKV := DeepCopy(srcKV)
		dst.SetMapIndex(k, dstKV)
	}
	return dst
}

func deepCopySlice(src reflect.Value, dst reflect.Value) reflect.Value {
	elemType := src.Type().Elem()
	sliceType := reflect.SliceOf(elemType)
	dst.Set(reflect.MakeSlice(sliceType, 0, src.Len()))
	for i := 0; i < src.Len(); i++ {
		subV := DeepCopy(src.Index(i))
		dst = reflect.Append(dst, subV)
	}
	return dst
}
