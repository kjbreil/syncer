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
	for _, child := range d.GetChildren() {
		child.Timestamp()
	}
}

func (d *Diff) Entries() []*Entry {
	return d.entries()
}
func (d *Diff) entries() []*Entry {
	var molds []*Entry
	for _, c := range d.GetChildren() {
		if len(c.GetChildren()) > 0 {
			molds = append(molds, c.entries()...)
		} else {
			if c.GetDelete() {
				molds = append(molds, &Entry{
					Key:    c.GetKey(),
					Action: Entry_REMOVE,
				})
			} else {
				molds = append(molds, &Entry{
					Key:   c.GetKey(),
					Value: c.GetValue(),
				})
			}
		}
	}
	return molds
}
