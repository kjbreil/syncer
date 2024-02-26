package extractor

import (
	"fmt"
	"github.com/kjbreil/syncer/control"
	"reflect"
)

func extractInterface(newValue, oldValue reflect.Value, upperType reflect.StructField, level int) (control.Entries, error) {

	nvK, ovK := newValue.Kind().String(), oldValue.Kind().String()
	fmt.Println(nvK, ovK)

	return extract(newValue.Elem(), oldValue.Elem(), upperType, level)
}
