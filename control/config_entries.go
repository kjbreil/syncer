package control

import (
	"strings"
)

// Entries is a slice of Entry structs
type Entries []*Entry

// AddKey adds a new key to the first key in each Entry
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

// AddIndex adds a new index to the first key in each Entry
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

// Struct returns a string representation of the Entries struct
func (ent *Entries) Struct() string {
	var builder strings.Builder
	for _, e := range *ent {
		builder.WriteString(e.Struct())
	}

	return builder.String()
}

// Equals returns true if the Entries are equal, false otherwise
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

// Diff returns a slice of Entries that are different between the two slices
func (ent Entries) Diff(other Entries) *Entries {
	var diff Entries
	for i, e := range ent {
		if !e.Equals(other[i]) {
			diff = append(diff, e)
		}
	}

	return &diff
}
