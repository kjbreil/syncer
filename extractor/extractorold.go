package extractor

import (
	"fmt"
	"github.com/kjbreil/syncer/control"
	"reflect"
	"strconv"
)

func (ext *Extractor) addHistory(head *control.Diff) {
	// if length of history equal to capacity drop first item and move everything down one
	// if len(ext.history) == cap(ext.history) {
	//     ext.history = ext.history[1:]
	// }
	// ext.history = append(ext.history, head)

	if len(head.GetChildren()) == 0 {
		return
	}
	if len(ext.history) == cap(ext.history) {
		for i := 0; i < len(ext.history)-1; i++ {
			ext.history[i] = ext.history[i+1]
		}
		ext.history[len(ext.history)-1] = head
	} else {
		ext.history = append(ext.history, head)
	}
}

func setValue(va reflect.Value, child *control.Diff) error {
	child.Value = &control.Object{}
	switch va.Kind() {
	case reflect.Invalid:
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value := va.Int()
		child.Value.Int64 = &value
	case reflect.Bool:
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value := va.Uint()
		child.Value.Uint64 = &value
	case reflect.Uintptr:
	case reflect.Float32:
		value := float32(va.Float())
		child.Value.Float32 = &value
	case reflect.Float64:
		value := va.Float()
		child.Value.Float64 = &value
	case reflect.String:
		value := va.String()
		child.Value.String_ = &value
	default:
		return fmt.Errorf("cannot setValue of type %s", va.Type().String())
	}
	return nil
}

func equal(n, o reflect.Value) bool {
	if n.Kind() != o.Kind() {
		return false
	}
	switch n.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return n.Int() == o.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return n.Uint() == o.Uint()
	case reflect.String:
		newS := n.String()
		oldS := o.String()
		return newS == oldS
		// return n.String() == o.String()
	case reflect.Bool:
		return n.Bool() == o.Bool()
	case reflect.Float32, reflect.Float64:
		return n.Float() == o.Float()
	case reflect.Complex64, reflect.Complex128:
		return n.Complex() == o.Complex()
	default:
		return false
	}
}

func makeString(x reflect.Value) string {
	switch x.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(x.Int(), 10)
	case reflect.Bool:
		return strconv.FormatBool(x.Bool())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(x.Uint(), 10)
	case reflect.Uintptr:
		return fmt.Sprintf("%d", x.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", x.Float())
	case reflect.Complex64, reflect.Complex128:
		return fmt.Sprintf("%f", x.Complex())
	case reflect.Pointer:
		return makeString(x.Elem())
	case reflect.String:
		return x.String()
	default:
		panic("makeString: unsupported type " + x.Type().String())
	}
}
