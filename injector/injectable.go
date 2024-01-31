package injector

import "github.com/kjbreil/syncer/control"

type Injectable struct {
	key    []*control.Key
	value  string
	action action
}

type action int

const (
	Add action = iota
	Remove
)

func (inja *Injectable) Advance() *Injectable {
	inja.key = inja.key[1:]
	return inja
}
