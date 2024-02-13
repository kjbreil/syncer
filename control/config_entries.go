package control

import "fmt"

type Entries []*Entry

func (ent *Entries) Struct() string {
	var s string
	for _, e := range *ent {
		s += e.Struct()
	}

	return s
}

func (ent Entries) Equals(other Entries) bool {
	if len(ent) != len(other) {
		return false
	}

	for i, e := range ent {
		if !e.Equals(other[i]) {
			return false
		}

	}

	return true
}

func (ent Entries) Diff(other Entries) *Entries {
	var diff Entries
	for i, e := range ent {
		if !e.Equals(other[i]) {
			diff = append(diff, e)
		}
	}

	return &diff
}

func (e *Entry) Struct() string {
	var s string
	s += "{\n\tKey: []*control.Key{\n"
	for _, k := range e.GetKey() {
		s += fmt.Sprintf("\t\t{\n\t\t\tKey: \"%s\",\n\t\t\tIndex: %s,\n\t\t},\n", k.GetKey(), k.GetIndex().Struct())
	}
	s += "\t},\n"
	if e.GetRemove() {
		s += "\tRemove: true,\n"
	} else {
		s += fmt.Sprintf("\tValue: %s,\n", e.GetValue().Struct())
	}

	s += "},\n"

	return s
}

func (o *Object) Struct() string {
	var s string
	s += "&control.Object{"
	if o != nil {
		switch {
		case o.String_ != nil:
			s += fmt.Sprintf("String_: control.MakePtr(\"%s\")", o.GetString_())
		case o.Int64 != nil:
			s += fmt.Sprintf("Int64: control.MakePtr(int64(%d))", o.GetInt64())
		case o.Uint64 != nil:
			s += fmt.Sprintf("Uint64: control.MakePtr(uint64(%d)),\n", o.GetUint64())
		case o.Float32 != nil:
			s += fmt.Sprintf("Float32: control.MakePtr(%f)", o.GetFloat32())
		case o.Float64 != nil:
			s += fmt.Sprintf("Float64: control.MakePtr(%f)", o.GetFloat64())
		case o.Bool != nil:
			s += fmt.Sprintf("Bool: control.MakePtr(%t)", o.GetBool())
		case o.Bytes != nil:
			s += fmt.Sprintf("Bytes: []byte{0x%02x}", o.GetBytes())
		}
	}
	s += "}"
	return s
}

func (e *Entry) Equals(other *Entry) bool {
	if e == nil && other == nil {
		return true
	}

	if e == nil || other == nil {
		return false
	}

	if len(e.GetKey()) != len(other.GetKey()) {
		return false
	}

	for i, k := range e.GetKey() {
		if k.GetKey() != other.GetKey()[i].GetKey() {
			return false
		}
	}

	if e.GetValue() == nil && other.GetValue() == nil {
		return true
	}

	if e.GetValue() == nil || other.GetValue() == nil {
		return false
	}

	return e.GetValue().Equals(other.GetValue())
}

func (o *Object) Equals(other *Object) bool {
	if o == nil && other == nil {
		return true
	}

	if o == nil || other == nil {
		return false
	}

	if o.GetString_() != other.GetString_() {
		return false
	}

	if o.GetInt64() != other.GetInt64() {
		return false
	}

	if o.GetUint64() != other.GetUint64() {
		return false
	}

	if o.GetFloat32() != other.GetFloat32() {
		return false
	}

	if o.GetFloat64() != other.GetFloat64() {
		return false
	}

	if o.GetBool() != other.GetBool() {
		return false
	}

	if o.GetBytes() != nil {
		panic("need to make bytes comparison")
	}

	return true
}

func MakePtr[V any](v V) *V {
	return &v
}
