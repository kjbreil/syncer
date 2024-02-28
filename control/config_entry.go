package control

import (
	"fmt"
	"strings"
)

func NewEntry(level int, value any) *Entry {
	return &Entry{
		Key:   make([]*Key, 0, level),
		Value: NewObject(value),
	}
}
func NewBlankEntry(level int) *Entry {
	return &Entry{
		Key: make([]*Key, 0, level),
	}
}
func NewRemoveEntry(level int) *Entry {
	return &Entry{
		Key:    make([]*Key, 0, level),
		Remove: true,
	}
}

// Advance either advances the index or the key
func (e *Entry) Advance() *Entry {
	if len(e.GetCurrIndexObjects()) > 1 {
		e.Key[e.GetKeyI()].Index = e.GetCurrIndexObjects()[1:]
	} else if int(e.KeyI) < len(e.GetKey())-1 {
		e.KeyI++
	}
	return e
}

func (e *Entry) IsLastKeyIndex() bool {
	return len(e.GetKey()) == 0 || (int(e.KeyI) == len(e.GetKey())-1 && e.IsLastIndex())
}

func (e *Entry) GetCurrIndexObjects() Objects {
	if e == nil {
		return nil
	}
	if len(e.GetKey()) == 0 {
		return nil
	}
	return e.GetKey()[e.GetKeyI()].GetIndex()
}

func (e *Entry) GetCurrKey() string {
	if e == nil {
		return ""
	}
	if len(e.GetKey()) == 0 {
		return ""
	}
	return e.GetKey()[e.GetKeyI()].GetKey()
}
func (e *Entry) GetCurrIndex() *Object {
	if len(e.GetCurrIndexObjects()) == 0 {
		return nil
	}
	return e.GetCurrIndexObjects()[0]
}

func (e *Entry) IsLastIndex() bool {
	return e.GetCurrIndex() == nil || len(e.GetCurrIndexObjects()) == 1
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

func (e *Entry) Struct() string {
	var sb strings.Builder
	sb.WriteString("{\n\tKey: []*control.Key{\n")
	for _, k := range e.GetKey() {
		key := k.GetKey()
		indexObject := Objects(k.GetIndex())
		sb.WriteString("\t\t{\n")
		sb.WriteString(fmt.Sprintf("\t\t\tKey: \"%s\",\n", key))
		if len(indexObject) > 0 {
			sb.WriteString(fmt.Sprintf("\t\t\tIndex: %s,\n", indexObject.Struct()))
		}
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
