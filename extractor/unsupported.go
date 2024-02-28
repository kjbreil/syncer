package extractor

import (
	"reflect"

	"github.com/kjbreil/syncer/control"
)

func extractUnsupported(_, _ reflect.Value, _ reflect.StructField, _ int) (control.Entries, error) {
	return nil, nil
}
