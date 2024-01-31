package endpoint

import (
	"context"
	"github.com/kjbreil/syncer/control"
	"github.com/kjbreil/syncer/injector"
)

type client struct {
	c         control.ConfigClient
	connected bool
	ctx       context.Context
	inj       *injector.Injector
}

func (c *client) init() {
	update, err := c.c.Update(c.ctx, &control.Request{Type: control.Request_INIT})
	if err != nil {
		return
	}
	for {
		cfg, err := update.Recv()
		if err != nil {
			continue
		}
		c.inj.Add(cfg)
	}
}

func newClient() *client {
	ctx, _ := context.WithCancel(context.Background())
	return &client{
		ctx: ctx,
	}
}
