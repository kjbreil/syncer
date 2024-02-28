package extractor

import (
	"github.com/kjbreil/syncer/control"
	"github.com/kjbreil/syncer/helpers/equal"
	"reflect"
)

func extractArray(newValue, oldValue reflect.Value, upperType reflect.StructField, level int) (control.Entries, error) {

	var entries control.Entries
	level++
	for i := 0; i < newValue.Len(); i++ {
		newIndexValue, oldIndexValue := newValue.Index(i), oldValue.Index(i)

		if equal.Equal(newIndexValue, oldIndexValue) {
			continue
		}
		additions, err := extract(newIndexValue, oldIndexValue, upperType, level, false)
		if err != nil {
			return nil, err
		}
		additions.AddIndex(i)
		entries = append(entries, additions...)
	}

	return entries, nil
}
