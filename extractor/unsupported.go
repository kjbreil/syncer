package extractor

import (
	"github.com/kjbreil/syncer/control"
	"reflect"
)

func extractUnsupported(newValue reflect.Value, oldValue reflect.Value, upperType reflect.StructField, level int) (control.Entries, error) {
	return nil, nil
}
