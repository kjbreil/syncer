package injector

import (
	"fmt"
	"github.com/kjbreil/syncer/control"
	"reflect"
)

func injectStruct(va reflect.Value, entry *control.Entry) error {
	entry.Advance()
	if structFieldValue := va.FieldByName(entry.GetCurrKeyString()); structFieldValue.IsValid() {
		return add(structFieldValue, entry)
	}
	return fmt.Errorf("field %s not found in struct", entry.GetCurrKeyString())
}
