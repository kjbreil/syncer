package extractor

import (
	"reflect"

	"github.com/kjbreil/syncer/control"
	"github.com/kjbreil/syncer/helpers/equal"
)

func extractPrimitive(newValue, oldValue reflect.Value, _ reflect.StructField, level int) (control.Entries, error) {
	if !equal.Equal(newValue, oldValue) {
		entry := control.NewEntry(level, reflect.Indirect(newValue).Interface())
		return control.Entries{entry}, nil
	}
	return nil, nil
}
