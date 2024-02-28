package control

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewDiff(key *Key) *Diff {
	return &Diff{
		Key: key,
	}
}

func NewDelDiff(key *Key) *Diff {
	return &Diff{
		Key:    key,
		Delete: true,
	}
}

func (d *Diff) Timestamp() {
	if d == nil {
		return
	}
	d.Time = timestamppb.Now()
	for _, child := range d.GetChildren() {
		child.Timestamp()
	}
}

func (d *Diff) AddChild(child *Diff, length int) {
	if len(d.Children) == 0 {
		d.Children = make([]*Diff, 0, length)
	}
	d.Children = append(d.Children, child)
}

func (d *Diff) Entries() Entries {
	return d.entries([]*Key{d.GetKey()})
}

func (d *Diff) entries(keys []*Key) []*Entry {
	var entries Entries
	for _, c := range d.GetChildren() {
		if len(c.GetChildren()) > 0 {
			// make a new keys matching length of current keys and copy into
			// new keys
			newKeys := make([]*Key, len(keys), len(keys)+1)
			copy(newKeys, keys)
			newKeys = append(newKeys, c.GetKey())

			children := c.entries(newKeys)
			if len(children) > 0 {
				entries = append(entries, children...)
			}
		} else {
			if c.GetDelete() {
				entries = append(entries, &Entry{
					Key:    append(keys, c.GetKey()),
					Remove: true,
				})
			} else {
				k := c.GetKey()
				entries = append(entries, &Entry{
					Key:   append(keys, k),
					Value: c.GetValue(),
				})
			}
		}
	}
	return entries
}

func (d *Diff) SetValue(v reflect.Value) error {
	d.Value = &Object{}
	switch v.Kind() {
	case reflect.Invalid:
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value := v.Int()
		d.Value.Int64 = &value
	case reflect.Bool:
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value := v.Uint()
		d.Value.Uint64 = &value
	case reflect.Uintptr:
	case reflect.Float32:
		value := float32(v.Float())
		d.Value.Float32 = &value
	case reflect.Float64:
		value := v.Float()
		d.Value.Float64 = &value
	case reflect.String:
		value := v.String()
		d.Value.String_ = &value
	default:
		return fmt.Errorf("cannot setValue of type %s", v.Type().String())
	}
	return nil
}
