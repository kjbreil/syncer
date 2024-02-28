package injector

import (
	"reflect"

	"github.com/kjbreil/syncer/control"
)

func injectArray(va reflect.Value, entry *control.Entry) error {
	indexInt := int(entry.GetCurrentIndex().GetInt64())
	return add(va.Index(indexInt), entry.Advance())
}
