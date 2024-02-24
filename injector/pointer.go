package injector

import (
	"github.com/kjbreil/syncer/control"
	"reflect"
)

func injectPointer(va reflect.Value, entry *control.Entry) error {
	// if its a remove and the last KeyIndex then nil out the value
	if entry.GetRemove() && entry.IsLastKeyIndex() {
		va.Set(reflect.Zero(va.Type()))
		return nil
	}
	// make the value if it is nil
	if va.IsNil() {
		newVa := reflect.New(va.Type().Elem())
		va.Set(newVa)
	}

	return add(va.Elem(), entry)
}
