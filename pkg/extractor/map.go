package extractor

import (
	"reflect"

	"github.com/kjbreil/syncer/pkg/control"
)

func extractMap(newValue, oldValue reflect.Value, upperType reflect.StructField, level int) (control.Entries, error) {
	level++

	if newValue.IsNil() && !oldValue.IsNil() {
		return control.Entries{control.NewRemoveEntry(level)}, nil
	}

	// if the oldValue is not valid then the map needs to be created
	if !oldValue.IsValid() || oldValue.IsNil() {
		// TODO: see if we can just copy the type of newValue
		keyType := upperType.Type.Key()
		valueType := upperType.Type.Elem()
		mapType := reflect.MapOf(keyType, valueType)
		oldValue = reflect.MakeMapWithSize(mapType, 0)
	}

	var entries control.Entries

	// look for values in oldValue not in newValue
	for _, k := range oldValue.MapKeys() {
		newMapIndexValue := newValue.MapIndex(k)

		if !newMapIndexValue.IsValid() {
			additions := control.Entries{control.NewRemoveEntry(level)}
			additions.AddIndex(k.Interface())
			entries = append(entries, additions...)
		}
	}

	for _, k := range newValue.MapKeys() {
		newMapIndexValue := newValue.MapIndex(k)
		oldMapIndexValue := oldValue.MapIndex(k)

		newMapIndexValue = reflect.Indirect(newMapIndexValue)
		oldMapIndexValue = reflect.Indirect(oldMapIndexValue)

		additions, err := extract(newMapIndexValue, oldMapIndexValue, upperType, level, false)
		if err != nil {
			return nil, err
		}
		additions.AddIndex(k.Interface())
		entries = append(entries, additions...)
	}

	return entries, nil
}
