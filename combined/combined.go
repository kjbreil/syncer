package combined

import (
	"context"
	"errors"
	"fmt"
	"github.com/kjbreil/syncer/control"
	"github.com/kjbreil/syncer/extractor"
	"github.com/kjbreil/syncer/injector"
	"time"
)

// Combined represents a combined configuration of an extractor and an injector.
type Combined struct {
	// extractor is the configuration of the extractor.
	extractor *extractor.Extractor
	// injector is the configuration of the injector.
	injector *injector.Injector

	extractorChanges func() error
	extractorChgChan chan struct{}
	injectorChanges  func() error
	injectorChgChan  chan struct{}

	Debounce time.Duration

	ctx    context.Context
	cancel context.CancelFunc
}

var (
	ErrExtractorChangeFn = errors.New("failed to run extractor changes function")
	ErrInjectorChangeFn  = errors.New("failed to run injector changes function")
)

// New creates a new Combined instance.
func New(ctx context.Context, data any) (*Combined, error) {
	if ctx == nil {
		return nil, errors.New("context is nil")
	}
	if data == nil {
		return nil, errors.New("data is nil")
	}
	var err error
	c := Combined{}
	c.ctx, c.cancel = context.WithCancel(ctx)

	c.extractor, err = extractor.New(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create extractor: %w", err)
	}

	c.injector, err = injector.New(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create injector: %w", err)
	}

	c.extractorChgChan = make(chan struct{}, 1)
	c.injectorChgChan = make(chan struct{}, 1)

	c.Debounce = time.Second * 2

	go c.injectorChangeDebounce()
	go c.extractorChangeDebounce()

	return &c, nil
}

func (c *Combined) injectorChangeDebounce() {
	timer := time.NewTimer(24 * 365 * time.Hour)
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-timer.C:
			// TODO: make logger in combined
			if c.injectorChanges != nil {
				_ = c.injectorChanges()
				// if err!= nil {
				//     c.logger.Error(err.Error())
				// }
			}
		case <-c.injectorChgChan:
			timer = time.NewTimer(c.Debounce)
		}
	}
}

func (c *Combined) extractorChangeDebounce() {
	timer := time.NewTimer(24 * 365 * time.Hour)
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-timer.C:
			// TODO: make logger in combined
			if c.extractorChanges != nil {
				_ = c.extractorChanges()
				// if err!= nil {
				//     c.logger.Error(err.Error())
				// }
			}
		case <-c.extractorChgChan:
			timer = time.NewTimer(c.Debounce)
		}
	}
}

// ExtractorChanges sets the function to be executed when the extractor configuration changes.
func (c *Combined) ExtractorChanges(fn func() error) {
	c.extractorChanges = fn
}

// InjectorChanges sets the function to be executed when the injector configuration changes.
func (c *Combined) InjectorChanges(fn func() error) {
	c.injectorChanges = fn
}

// Add adds a new entry to the control file.
func (c *Combined) Add(cfg *control.Entry) error {
	c.injectorChgChan <- struct{}{}
	return c.injector.Add(cfg)
}

// Reset resets the Combined instance.
func (c *Combined) Reset() {
	c.extractor.Reset()
}

// Diff returns the difference between the current configuration and the given data.
func (c *Combined) Diff(data any) (*control.Diff, error) {
	head, err := c.extractor.Diff(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create diff in extractor: %w", err)
	}
	if head == nil {
		return nil, nil
	}
	c.extractorChgChan <- struct{}{}
	return head, nil
}

// Diff returns the difference between the current configuration and the given data.
func (c *Combined) Entries(data any) (control.Entries, error) {
	head, err := c.extractor.Diff(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create diff in extractor: %w", err)
	}
	if head == nil {
		return nil, nil
	}
	entries := head.Entries()
	c.extractorChgChan <- struct{}{}
	return entries, nil
}

// Close stops the Combined instance and closes all open resources.
func (c *Combined) Close() error {
	c.cancel()
	close(c.extractorChgChan)
	close(c.injectorChgChan)
	return nil
}
