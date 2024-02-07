package combined

import (
	"github.com/kjbreil/syncer/control"
	"github.com/kjbreil/syncer/extractor"
	"github.com/kjbreil/syncer/injector"
)

type Combined struct {
	extractor *extractor.Extractor
	injector  *injector.Injector
}

func New(data any) (*Combined, error) {
	var err error
	c := Combined{}
	c.extractor = extractor.New(data)
	c.injector, err = injector.New(data)

	return &c, err
}

func (c *Combined) Add(cfg *control.Entry) error {
	return c.injector.Add(cfg)
}

func (c *Combined) Reset() {
	c.extractor.Reset()
}

func (c *Combined) Diff(data any) (*control.Diff, error) {
	return c.extractor.Diff(data)
}
