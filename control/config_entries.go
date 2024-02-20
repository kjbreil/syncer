package control

import (
	"fmt"
	"strings"
)

type Entries []*Entry

func (ent *Entries) Struct() string {
	var builder strings.Builder
	for _, e := range *ent {
		builder.WriteString(e.Struct())
	}

	return builder.String()
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
	var sb strings.Builder
	sb.WriteString("{\n\tKey: []*control.Key{\n")
	for _, k := range e.GetKey() {
		key := k.GetKey()
		index := Objects(k.GetIndex()).Struct()
		sb.WriteString("\t\t{\n")
		sb.WriteString(fmt.Sprintf("\t\t\tKey: \"%s\",\n", key))
		sb.WriteString(fmt.Sprintf("\t\t\tIndex: %s,\n", index))
		sb.WriteString("\t\t},\n")
	}
	sb.WriteString("\t},\n")
	var valueString string
	if e.GetRemove() {
		valueString = "\tRemove: true,\n"
	} else {
		valueString = fmt.Sprintf("\tValue: %s,\n", e.GetValue().Struct())
	}
	sb.WriteString(valueString)
	sb.WriteString("},\n")
	return sb.String()
}

func (o *Object) Struct() string {
	if o == nil {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("&control.Object{")
	if o.String_ != nil {
		sb.WriteString(fmt.Sprintf("String_: control.MakePtr(\"%s\")", o.GetString_()))
	} else if o.Int64 != nil {
		sb.WriteString(fmt.Sprintf("Int64: control.MakePtr(strconv.FormatInt(%d, 10))", o.GetInt64()))
	} else if o.Uint64 != nil {
		sb.WriteString(fmt.Sprintf("Uint64: control.MakePtr(strconv.FormatUint(%d, 10))", o.GetUint64()))
	} else if o.Float32 != nil {
		sb.WriteString(fmt.Sprintf("Float32: control.MakePtr(strconv.FormatFloat(float64(%f), 'f', -1, 32))", o.GetFloat32()))
	} else if o.Float64 != nil {
		sb.WriteString(fmt.Sprintf("Float64: control.MakePtr(strconv.FormatFloat(%f, 'f', -1, 64))", o.GetFloat64()))
	} else if o.Bool != nil {
		sb.WriteString(fmt.Sprintf("Bool: control.MakePtr(strconv.FormatBool(%t))", o.GetBool()))
	} else if o.Bytes != nil {
		sb.WriteString(fmt.Sprintf("Bytes: []byte(hex.EncodeToString(%v))", o.GetBytes()))
	}
	sb.WriteString("}")
	return sb.String()
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
		if !Objects(k.GetIndex()).Equals(other.GetKey()[i].GetIndex()) {
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
