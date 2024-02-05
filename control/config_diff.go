package control

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewDiff(key []*Key) *Diff {
	return &Diff{
		Key: key,
	}
}

func NewDelDiff(key []*Key) *Diff {
	return &Diff{
		Key:    key,
		Delete: true,
	}
}

func (d *Diff) Timestamp() {
	d.Time = timestamppb.Now()
	for _, child := range d.Children {
		child.Timestamp()
	}
}

func (d *Diff) Entries() []*Entry {
	return d.entries()
}
func (d *Diff) entries() []*Entry {
	var moulds []*Entry
	for _, c := range d.Children {
		if len(c.Children) > 0 {
			moulds = append(moulds, c.entries()...)
		} else {
			if c.Delete {
				moulds = append(moulds, &Entry{
					Key:    c.Key,
					Action: Entry_REMOVE,
				})
			} else {
				moulds = append(moulds, &Entry{
					Key:   c.Key,
					Value: c.Value,
				})
			}
		}
	}
	return moulds
}
