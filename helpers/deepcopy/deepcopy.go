package deepcopy

import (
	"reflect"
)

type copyFn func(dst, src reflect.Value)

var copyFns map[reflect.Kind]copyFn

func init() {
	copyFns = map[reflect.Kind]copyFn{
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

// DeepCopy copies the value of a reflect.Value returning a new reflect.Value
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

// Any takes a value of any type and returns it as the same type after a deep copy.
// This function is useful when you need to work with a value of an unknown type
// and want to ensure that the value is copied before making any changes to it.
func Any[T any](src T) T {
	srcV := reflect.ValueOf(src)
	dst := deepCopy(srcV)
	return dst.Interface().(T)
}

func deepCopy(src reflect.Value) reflect.Value {
	dst := reflect.Indirect(reflect.New(src.Type()))

	if c, ok := copyFns[src.Kind()]; ok {
		c(dst, src)
	}

	// Return the destination value.
	return dst
}

func deepCopyStruct(dst, src reflect.Value) {
	for i := 0; i < src.NumField(); i++ {
		if !src.Field(i).CanInterface() || !dst.Field(i).CanSet() {
			continue
		}
		dst.Field(i).Set(deepCopy(src.Field(i)))
	}
	return
}

func deepCopyInterface(dst, src reflect.Value) {
	if src.Elem().IsValid() {
		dst.Set(deepCopy(src.Elem()))
	}
}

func deepCopyPointer(dst, src reflect.Value) {
	if src.Elem().IsValid() {
		dst.Set(deepCopy(src.Elem()).Addr())
	}
}

func deepCopyMap(dst, src reflect.Value) {
	if src.IsNil() {
		return
	}
	dst.Set(reflect.MakeMapWithSize(src.Type(), src.Len()))
	for _, k := range src.MapKeys() {
		dst.SetMapIndex(k, deepCopy(src.MapIndex(k)))
	}
}

func deepCopySlice(dst, src reflect.Value) {
	if src.IsNil() {
		return
	}
	elemType := src.Type().Elem()
	sliceType := reflect.SliceOf(elemType)
	dst.Set(reflect.MakeSlice(sliceType, src.Len(), src.Cap()))
	for i := 0; i < src.Len(); i++ {
		dst.Index(i).Set(deepCopy(src.Index(i)))
	}
}
func deepCopyArray(dst, src reflect.Value) {
	elemType := src.Type()
	dst.Set(reflect.New(reflect.ArrayOf(src.Len(), elemType.Elem())).Elem())
	for i := 0; i < src.Len(); i++ {
		dst.Index(i).Set(deepCopy(src.Index(i)))
	}
}

func deepCopyPrimitive(dst, src reflect.Value) {
	dst.Set(src)
}
