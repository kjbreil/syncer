package extractor

import (
	"fmt"
	"reflect"

	"github.com/kjbreil/syncer/control"
)

func extractStruct(newValue, oldValue reflect.Value, _ reflect.StructField, level int) (control.Entries, error) {
	// TODO: This should check if oldValue is valid and return a delete if it is
	if !newValue.IsValid() {
		return nil, fmt.Errorf("extractStruct: newValue is not valid")
	}
	// check if the oldValue is valid (exists) and create it if it does not
	if !oldValue.IsValid() {
		oldValue = reflect.New(newValue.Type()).Elem()
	}

	var entries control.Entries

	for i := 0; i < newValue.NumField(); i++ {
		// skip if extractor tag is set
		etag := newValue.Type().Field(i).Tag.Get("extractor")
		if etag == "-" {
			continue
		}
		if !newValue.Field(i).CanInterface() {
			continue
		}
		level++
		fieldEntry, err := extract(newValue.Field(i), oldValue.Field(i), newValue.Type().Field(i), level, true)
		if err != nil {
			return nil, err
		}
		if fieldEntry != nil {
			entries = append(entries, fieldEntry...)
		}
	}

	return entries, nil
}
