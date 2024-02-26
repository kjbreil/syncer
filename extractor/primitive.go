package extractor

import (
	"github.com/kjbreil/syncer/control"
	"github.com/kjbreil/syncer/helpers/equal"
	"reflect"
)

func extractPrimitive(newValue, oldValue reflect.Value, upperType reflect.StructField, level int) (control.Entries, error) {
	if !equal.Equal(newValue, oldValue) {
		entry := control.NewEntry(level, reflect.Indirect(newValue).Interface())
		return control.Entries{entry}, nil
	}
	return nil, nil
}
