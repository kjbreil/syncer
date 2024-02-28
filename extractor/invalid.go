package extractor

import (
	"reflect"

	"github.com/kjbreil/syncer/control"
)

func extractInvalid(_, _ reflect.Value, _ reflect.StructField, _ int) (control.Entries, error) {
	panic("extractInvalid should not be called")
	return nil, nil
}
