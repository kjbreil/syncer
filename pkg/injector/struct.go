package injector

import (
	"fmt"
	"reflect"

	"github.com/kjbreil/syncer/pkg/control"
)

func injectStruct(va reflect.Value, entry *control.Entry) error {
	entry.Advance()
	if structFieldValue := va.FieldByName(entry.GetCurrKeyString()); structFieldValue.IsValid() {
		return add(structFieldValue, entry)
	}
	return fmt.Errorf("field %s not found in struct", entry.GetCurrKeyString())
}
