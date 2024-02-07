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

// New creates a new Combined instance.
func New(data any) (*Combined, error) {
	// c is the new Combined instance.
	var err error
	c := Combined{}

	// c.extractor is initialized to a new extractor.Extractor instance.
	c.extractor = extractor.New(data)

	// c.injector is initialized to a new injector.Injector instance.
	c.injector, err = injector.New(data)

	// Return the new Combined instance and any error.
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
