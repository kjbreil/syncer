package control

import (
	"strings"
)

type Entries []*Entry

func (ent *Entries) AddKey(key string) {
	for _, e := range *ent {
		e.Key = append([]*Key{&Key{
			Key: key,
		}}, e.Key...)
	}
}

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
