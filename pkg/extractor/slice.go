package extractor

import (
	"reflect"

	"github.com/kjbreil/syncer/pkg/control"
	"github.com/kjbreil/syncer/pkg/equal"
)

func extractSlice(newValue, oldValue reflect.Value, upperType reflect.StructField, level int) (control.Entries, error) {
	if newValue.IsNil() && !oldValue.IsNil() {
		return control.Entries{control.NewRemoveEntry(level)}, nil
	}

	// make the old slice match the new slice
	// oldValue is shorter, add the extra entries and just run compare
	if oldValue.Len() < newValue.Len() {
		// make a new slice for oldValue of capacity the newValue Slice
		newOldSlice := reflect.MakeSlice(newValue.Type(), newValue.Len(), newValue.Cap())
		// copy the values from the oldSlice into the newOldSlice
		reflect.Copy(newOldSlice, oldValue)
		// set the oldSlice to the newOldSlice
		oldValue = newOldSlice
		// if value is a pointer loop over and create a zero value entry for each element
	}

	var entries control.Entries
	level++

	// newValue is shorter, set a delete starting at the index of the difference
	if newValue.Len() < oldValue.Len() {
		additions := control.Entries{control.NewRemoveEntry(level)}
		additions.AddIndex(newValue.Len())
		entries = append(entries, additions...)
	}

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
