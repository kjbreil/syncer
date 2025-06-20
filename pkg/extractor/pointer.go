package extractor

import (
	"reflect"

	"github.com/kjbreil/syncer/pkg/control"
)

func extractPointer(newValue, oldValue reflect.Value, upperValue reflect.StructField, level int) (control.Entries, error) {
	if (!newValue.IsValid() || newValue.IsNil()) && (!oldValue.IsValid() || oldValue.IsNil()) {
		return nil, nil
	}

	if newValue.IsNil() && !oldValue.IsNil() {
		return control.Entries{control.NewRemoveEntry(level)}, nil
	}

	if !oldValue.IsValid() || oldValue.IsNil() {
		oldValue = reflect.New(newValue.Type()).Elem()
	}

	return extract(newValue.Elem(), oldValue.Elem(), upperValue, level, false)
}
