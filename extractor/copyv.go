package extractor

import (
	"reflect"
)

func copyValue(src reflect.Value) reflect.Value {
	dst := reflect.New(src.Type()).Elem()

	switch src.Kind() {
	case reflect.Pointer:
		if src.Elem().IsValid() {
			subV := copyValue(src.Elem())
			dst.Set(subV.Addr())
		}
	case reflect.Interface:
		if src.Elem().IsValid() {
			subV := copyValue(src.Elem())
			if dst.CanSet() {
				dst.Set(subV)
			}
		}
	case reflect.Struct:
		for i := 0; i < src.NumField(); i++ {
			dstFV := copyValue(src.Field(i))
			dst.Field(i).Set(dstFV)
		}
	case reflect.Map:
		keyType := src.Type().Key()
		valueType := src.Type().Elem()
		mapType := reflect.MapOf(keyType, valueType)
		if dst.CanSet() {
			dst.Set(reflect.MakeMapWithSize(mapType, 0))
		}
		for _, k := range src.MapKeys() {
			srcKV := src.MapIndex(k)
			dstKV := copyValue(srcKV)
			dst.SetMapIndex(k, dstKV)
		}
	case reflect.Slice:
		elemType := src.Type().Elem()
		sliceType := reflect.SliceOf(elemType)
		if dst.CanSet() {
			dst.Set(reflect.MakeSlice(sliceType, src.Len(), src.Len()))
		}
		for i := 0; i < src.Len(); i++ {
			subV := copyValue(src.Index(i))
			dst.Index(i).Set(subV)
		}
	default:
		if dst.CanSet() {
			dst.Set(src)
		}
	}

	return dst
}
