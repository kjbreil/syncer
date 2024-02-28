package injector

import (
	"errors"
	"github.com/kjbreil/syncer/control"
	"reflect"
)

func injectSlice(va reflect.Value, entry *control.Entry) error {
	// no index, either an error or full remove the slice
	if entry.GetCurrKey().HasNoIndex() {
		// no index on a map key and remove type make map nil
		if entry.GetRemove() {
			va.Set(reflect.New(va.Type()).Elem())
			return nil
		}
		return errors.New("slice type without index")
	}
	// get the int representing the current index
	indexInt := int(entry.GetCurrentIndex().GetInt64())
	// if entry is a delete entry then delete the current index+
	if entry.GetRemove() {
		newSlice := reflect.MakeSlice(va.Type(), indexInt, indexInt)
		reflect.Copy(newSlice, va)
		va.Set(newSlice)
		return nil
	}

	// if the current slice is not large enough for the new entry create a new slice
	// with the current entries and nil
	diff := indexInt + 1 - va.Len()
	if diff > 0 {
		newSlice := reflect.MakeSlice(va.Type(), diff, diff)
		va.Set(reflect.AppendSlice(va, newSlice))
	}

	return add(va.Index(indexInt), entry.Advance())
}
