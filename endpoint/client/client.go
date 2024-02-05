package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/kjbreil/syncer/control"
	"github.com/kjbreil/syncer/extractor"
	"github.com/kjbreil/syncer/injector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"net"
	"sync"
	"time"
)

var (
	ErrClientNotAvailable = fmt.Errorf("could not dial Client")
	ErrClientInjector     = fmt.Errorf("client could not create injector")
)

type Client struct {
	c    control.ConfigClient
	conn *grpc.ClientConn
	peer net.TCPAddr

	ctx       context.Context
	cancel    context.CancelFunc
	injector  *injector.Injector
	extractor *extractor.Extractor
	errors    chan error
}

func New(ctx context.Context, wg *sync.WaitGroup, data any, peer net.TCPAddr, errs chan error) (*Client, error) {
	var err error

	c := &Client{
		peer:   peer,
		errors: errs,
	}

	c.ctx, c.cancel = context.WithCancel(ctx)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	addr := fmt.Sprintf("%s:%d", peer.IP.String(), peer.Port)

	c.conn, err = grpc.Dial(addr, opts...)
	if err != nil {
		c.cancel()
		return nil, ErrClientNotAvailable
	}

	c.c = control.NewConfigClient(c.conn)

	wg.Add(1)
	go func() {
		for {
			select {
			case <-time.After(time.Second):
				_, err = c.c.Control(c.ctx, &control.Message{Action: control.Message_PING})
				if err != nil {
					c.errors <- err
					c.cancel()
				}
			case <-c.ctx.Done():
				wg.Done()
				err = c.conn.Close()
				if err != nil {
					c.errors <- err
					return
				}
				return
			}
		}
	}()

	_, err = c.c.Control(c.ctx, &control.Message{Action: control.Message_PING})
	if err != nil {
		return nil, c.closeWithError(fmt.Errorf("%w: %w", ErrClientNotAvailable, err))
	}

	c.extractor = extractor.New(data)
	c.injector, err = injector.New(data)
	if err != nil {
		return nil, c.closeWithError(fmt.Errorf("%w: %w", ErrClientInjector, err))
	}

	return c, nil
}

func (c *Client) Running() bool {
	return c.ctx.Err() == nil
}

// Init requests to init data from the server
func (c *Client) Init() {
	update, err := c.c.Update(c.ctx, &control.Request{Type: control.Request_INIT})
	if err != nil {
		c.errors <- fmt.Errorf("Client.Init(): %w", err)
		return
	}
	c.processUpdate(update)
}

// ShutdownRemoteServer requests to shut down the server
func (c *Client) ShutdownRemoteServer() {
	_, err := c.c.Control(c.ctx, &control.Message{Action: control.Message_SHUTDOWN})
	if err != nil {
		c.errors <- fmt.Errorf("Client.ShutdownRemoteServer(): %w", err)
		return
	}
}

func (c *Client) Changes() {
	update, err := c.c.Update(c.ctx, &control.Request{Type: control.Request_INIT})
	if err != nil {
		c.errors <- fmt.Errorf("Client.changes(): %w", err)
		return
	}
	c.processUpdate(update)
}

func (c *Client) processUpdate(update control.Config_UpdateClient) {
	for {
		cfg, err := update.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			if err != nil {
				c.errors <- fmt.Errorf("Client.processUpdate(): %w", err)
			}
			c.cancel()
			return
		}
		err = c.injector.Add(cfg)
		if err != nil {
			c.errors <- err
		}
	}
}

func (c *Client) closeWithError(err error) error {
	c.cancel()
	return err
}
