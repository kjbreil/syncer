package control

import (
	"strings"
)

type Entries []*Entry

func (ent *Entries) AddKey(key string) {
	for _, e := range *ent {
		if len(e.Key) > 0 && e.Key[0].GetKey() == "" {
			e.Key[0].Key = key
		} else {
			e.Key = append([]*Key{{
				Key: key,
			}}, e.Key...)
		}
	}
}

func (ent *Entries) AddIndex(index any) {
	for _, e := range *ent {
		if len(e.Key) > 0 && e.Key[0].GetKey() == "" {
			e.Key[0].Index = NewObjects(NewObject(index), e.Key[0].Index...)
		} else {
			e.Key = append([]*Key{{
				Index: NewObjects(NewObject(index)),
			}}, e.Key...)
		}
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
