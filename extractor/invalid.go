package extractor

import (
	"github.com/kjbreil/syncer/control"
	"reflect"
)

func extractInvalid(newValue reflect.Value, oldValue reflect.Value, upperType reflect.StructField, level int) (control.Entries, error) {
	panic("extractInvalid should not be called")
	return nil, nil
}
