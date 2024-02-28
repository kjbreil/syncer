package extractor

import (
	"github.com/kjbreil/syncer/control"
	"reflect"
)

func extractInterface(newValue, oldValue reflect.Value, upperType reflect.StructField, level int) (control.Entries, error) {
	// the base type of the interface is invalid on both new and old values, effectively equal
	if !newValue.Elem().IsValid() && !oldValue.Elem().IsValid() {
		return nil, nil
	}

	if !newValue.Elem().IsValid() && oldValue.Elem().IsValid() {
		return control.Entries{control.NewRemoveEntry(level)}, nil
	}

	return extract(newValue.Elem(), oldValue.Elem(), upperType, level, false)
}
