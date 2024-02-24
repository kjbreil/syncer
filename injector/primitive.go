package injector

import (
	"fmt"
	"github.com/kjbreil/syncer/control"
	"reflect"
)

func injectPrimitive(va reflect.Value, entry *control.Entry) error {
	if va.CanSet() {
		return entry.GetValue().SetValue(va)
	} else {
		return fmt.Errorf("cannot set value for primitive %s", entry.GetCurrKey())
	}
}
