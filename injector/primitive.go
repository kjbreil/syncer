package injector

import (
	"fmt"
	"reflect"

	"github.com/kjbreil/syncer/control"
)

func injectPrimitive(va reflect.Value, entry *control.Entry) error {
	if va.CanSet() {
		return entry.GetValue().SetValue(va)
	}
	return fmt.Errorf("cannot set value for primitive %s", entry.GetCurrKeyString())
}
