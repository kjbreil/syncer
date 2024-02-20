package combined

import (
	"errors"
	"fmt"
	"github.com/kjbreil/syncer/control"
	"github.com/kjbreil/syncer/extractor"
	"github.com/kjbreil/syncer/injector"
)

// Combined represents a combined configuration of an extractor and an injector.
type Combined struct {
	// extractor is the configuration of the extractor.
	extractor *extractor.Extractor
	// injector is the configuration of the injector.
	injector *injector.Injector
}

// New creates a new Combined instance.
func New(data any) (*Combined, error) {
	if data == nil {
		return nil, errors.New("data is nil")
	}
	var err error
	c := Combined{}
	c.extractor, err = extractor.New(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create extractor: %w", err)
	}
	c.injector, err = injector.New(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create injector: %w", err)
	}
	return &c, nil
}

// Add adds a new entry to the control file.
func (c *Combined) Add(cfg *control.Entry) error {
	return c.injector.Add(cfg)
}

// Reset resets the Combined instance.
func (c *Combined) Reset() {
	c.extractor.Reset()
}

// Diff returns the difference between the current configuration and the given data.
func (c *Combined) Diff(data any) (*control.Diff, error) {
	return c.extractor.Diff(data)
}
