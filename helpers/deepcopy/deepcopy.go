package deepcopy

import (
	"reflect"
)

type copier func(src reflect.Value) (dst reflect.Value)

var copiers map[reflect.Kind]copier

func init() {
	copiers = map[reflect.Kind]copier{
		reflect.Bool:       deepCopyPrimitive,
		reflect.Int:        deepCopyPrimitive,
		reflect.Int8:       deepCopyPrimitive,
		reflect.Int16:      deepCopyPrimitive,
		reflect.Int32:      deepCopyPrimitive,
		reflect.Int64:      deepCopyPrimitive,
		reflect.Uint:       deepCopyPrimitive,
		reflect.Uint8:      deepCopyPrimitive,
		reflect.Uint16:     deepCopyPrimitive,
		reflect.Uint32:     deepCopyPrimitive,
		reflect.Uint64:     deepCopyPrimitive,
		reflect.Uintptr:    deepCopyPrimitive,
		reflect.Float32:    deepCopyPrimitive,
		reflect.Float64:    deepCopyPrimitive,
		reflect.Complex64:  deepCopyPrimitive,
		reflect.Complex128: deepCopyPrimitive,
		reflect.String:     deepCopyPrimitive,
		reflect.Array:      deepCopyArray,
		reflect.Map:        deepCopyMap,
		reflect.Ptr:        deepCopyPointer,
		reflect.Interface:  deepCopyInterface,
		reflect.Slice:      deepCopySlice,
		reflect.Struct:     deepCopyStruct,
	}
}

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
	return deepCopy(src)
}

func deepCopy(src reflect.Value) reflect.Value {
	dst, valid := newDst(src)
	if !valid {
		return dst
	}

	if c, ok := copiers[src.Kind()]; ok {
		dst.Set(c(src))
	}

	// Return the destination value.
	return dst
}

func deepCopyStruct(src reflect.Value) reflect.Value {
	dst, valid := newDst(src)
	if !valid {
		return dst
	}

	for i := 0; i < src.NumField(); i++ {
		if !src.Field(i).CanInterface() || !dst.Field(i).CanSet() {
			continue
		}
		dst.Field(i).Set(deepCopy(src.Field(i)))
	}
	return dst
}

func deepCopyInterface(src reflect.Value) reflect.Value {
	dst, valid := newDst(src)
	if !valid {
		return dst
	}

	if src.Elem().IsValid() {
		dst.Set(deepCopy(src.Elem()))
	}
	return dst
}

func deepCopyPointer(src reflect.Value) reflect.Value {
	dst, valid := newDst(src)
	if !valid {
		return dst
	}

	if src.Elem().IsValid() {
		dst.Set(deepCopy(src.Elem()).Addr())
	}
	return dst
}

func deepCopyMap(src reflect.Value) reflect.Value {
	dst, valid := newDst(src)
	if !valid {
		return dst
	}

	dst.Set(reflect.MakeMapWithSize(src.Type(), src.Len()))
	for _, k := range src.MapKeys() {
		dst.SetMapIndex(k, deepCopy(src.MapIndex(k)))
	}
	return dst
}

func deepCopySlice(src reflect.Value) reflect.Value {
	dst, valid := newDst(src)
	if !valid {
		return dst
	}

	elemType := src.Type().Elem()
	sliceType := reflect.SliceOf(elemType)
	dst.Set(reflect.MakeSlice(sliceType, 0, src.Len()))
	for i := 0; i < src.Len(); i++ {
		dst = reflect.Append(dst, deepCopy(src.Index(i)))
	}
	return dst
}
func deepCopyArray(src reflect.Value) reflect.Value {
	dst, valid := newDst(src)
	if !valid {
		return dst
	}
	elemType := src.Type().Elem()
	dst.Set(reflect.New(reflect.ArrayOf(src.Len(), elemType)))
	for i := 0; i < src.Len(); i++ {
		dst.Index(i).Set(deepCopy(src.Index(i)))
	}
	return dst
}

func deepCopyPrimitive(src reflect.Value) reflect.Value {
	dst, valid := newDst(src)
	if !valid {
		return dst
	}
	dst.Set(src)
	return dst
}

func newDst(src reflect.Value) (reflect.Value, bool) {
	// Create a new destination value that is a copy of the source value.
	dst := reflect.Indirect(reflect.New(src.Type()))

	return dst, dst.CanSet() || src.CanInterface()
}
