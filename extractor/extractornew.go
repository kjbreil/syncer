package extractor

import (
	"fmt"
	"github.com/kjbreil/syncer/control"
	"reflect"
)

func (ext *Extractor) EntriesNew(data any) control.Entries {
	h, err := ext.DiffNew(data)
	if err != nil {
		panic(err)
	}
	return h.Entries()
}

func (ext *Extractor) DiffNew(data any) (*control.Diff, error) {
	// force single threaded access
	ext.mut.Lock()
	defer ext.mut.Unlock()

	newValue := reflect.ValueOf(data)
	if newValue.Kind() != reflect.Ptr {
		return nil, ErrNotPointer
	}

	// follow the pointer to the actual value
	newValue = reflect.Indirect(newValue)
	oldValue := reflect.Indirect(reflect.ValueOf(ext.data))

	head := control.NewDiff(&control.Key{Key: oldValue.Type().Name()})

	var err error
	head.Children, err = extractStruct(newValue, oldValue)
	if err != nil {
		return nil, err
	}

	head.Timestamp()

	return head, nil
}

func extractStruct(newValue, oldValue reflect.Value) ([]*control.Diff, error) {

	// TODO: Check if actually structs?

	// check if the oldValue is valid (exists) and create it if it does not
	if !oldValue.IsValid() {
		oldValue = reflect.New(newValue.Type()).Elem()
	}

	var children []*control.Diff

	// loop over the fields of the newValue finding the relevant matching field in the old value
	numFields := newValue.NumField()
	for i := 0; i < numFields; i++ {
		oldField := oldValue.FieldByName(newValue.Field(i).Type().Name())
		fmt.Println(oldField)
		newValueFieldKind := newValue.Field(i).Kind()

		switch newValueFieldKind {
		default:
			// if !equal(newValue.Field(i), oldValue.Field(i)) {
			// 	if len(children) == 0 {
			// 		children = make([]*control.Diff, 0, numFields)
			// 	}
			// 	key := oldType.Field(i).Name
			// 	child := control.NewDiff(append(parent.Key, &control.Key{
			// 		Key: key,
			// 	}),
			// 	)
			//
			// 	children = append(children, &control.Diff{
			//
			// 	})
			// 	child.Value = &control.Object{}
			// 	err := setValue(newValue.Field(i), child)
			// 	if err != nil {
			// 		return fmt.Errorf("extractLevel: %w", err)
			// 	}
			// 	parent.AddChild(child, numFields)
			// 	if oldField.CanSet() {
			// 		oldField.Set(newValue.Field(i))
			// 	}
			// }
		}
	}

	return children, nil
}
