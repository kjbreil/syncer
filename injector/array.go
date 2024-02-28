package injector

import (
	"github.com/kjbreil/syncer/control"
	"reflect"
)

func injectArray(va reflect.Value, entry *control.Entry) error {
	indexInt := int(entry.GetCurrentIndex().GetInt64())
	return add(va.Index(indexInt), entry.Advance())
}
