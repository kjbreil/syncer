package injector

import (
	"reflect"

	"github.com/kjbreil/syncer/control"
)

func injectInterface(va reflect.Value, entry *control.Entry) error {
	// if it's a remove and the last KeyIndex then nil out the value
	if entry.GetRemove() && entry.IsLastKeyIndex() {
		va.Set(reflect.Zero(va.Type()))
		return nil
	}

	va = va.Elem()
	return add(va, entry)
}
