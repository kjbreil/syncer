package control

import "strings"

type Objects []*Object

func (o Objects) Equals(other Objects) bool {
	if len(o) != len(other) {
		return false
	}
	for i, io := range o {
		if !io.Equals(other[i]) {
			return false
		}
	}
	return true
}

func (o Objects) Struct() string {
	var sb strings.Builder

	sb.WriteString("control.NewObjects(")
	for i, o := range o {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(o.Struct())
	}
	sb.WriteString(")")

	return sb.String()
}
