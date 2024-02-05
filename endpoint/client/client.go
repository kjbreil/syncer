package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/kjbreil/syncer/combined"
	"github.com/kjbreil/syncer/control"
	"github.com/kjbreil/syncer/endpoint/settings"
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

	ctx    context.Context
	cancel context.CancelFunc

	settings *settings.Settings

	combined *combined.Combined
	// injector *injector.Injector
	// client extractor not used yet
	// extractor *extractor.Extractor

	errors chan error
}

func New(ctx context.Context, wg *sync.WaitGroup, data any, peer net.TCPAddr, errs chan error, settings *settings.Settings) (*Client, error) {
	var err error

	c := &Client{
		peer:     peer,
		errors:   errs,
		settings: settings,
	}

	c.ctx, c.cancel = context.WithCancel(ctx)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithBlock())

	addr := fmt.Sprintf("%s:%d", peer.IP.String(), peer.Port)

	dialCtx, cancel := context.WithTimeout(c.ctx, time.Second)
	defer cancel()
	c.conn, err = grpc.DialContext(dialCtx, addr, opts...)
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

	if settings.AutoUpdate {
		wg.Add(1)
		go func() {
			for {
				select {
				case <-time.After(time.Second):
					c.Changes()
				case <-c.ctx.Done():
					wg.Done()
					return
				}
			}
		}()
	}

	_, err = c.c.Control(c.ctx, &control.Message{Action: control.Message_PING})
	if err != nil {
		return nil, c.closeWithError(fmt.Errorf("%w: %w", ErrClientNotAvailable, err))
	}

	// c.extractor = extractor.New(data)
	// c.injector, err = injector.New(data)
	c.combined, err = combined.New(data)

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
	update, err := c.c.Pull(c.ctx, &control.Request{Type: control.Request_INIT})
	if err != nil {
		c.errors <- fmt.Errorf("Client.Init(): %w", err)
		return
	}
	c.processUpdate(update)

	// s, err := c.c.Push(c.ctx)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// for {
	// 	cfg, err := s.Send()
	// 	if errors.Is(err, io.EOF) {
	// 		break
	// 	}
	// 	if err != nil {
	// 		if err != nil {
	// 			c.errors <- fmt.Errorf("Client.processUpdate(): %w", err)
	// 		}
	// 		c.cancel()
	// 		return
	// 	}
	// 	err = c.injector.Add(cfg)
	// 	if err != nil {
	// 		c.errors <- err
	// 	}
	// }

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
	update, err := c.c.Pull(c.ctx, &control.Request{Type: control.Request_CHANGES})
	if err != nil {
		c.errors <- fmt.Errorf("Client.changes(): %w", err)
		return
	}
	c.processUpdate(update)
}

func (c *Client) processUpdate(update control.Config_PullClient) {
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
		err = c.combined.Add(cfg)
		if err != nil {
			c.errors <- err
		}
	}
}

func (c *Client) closeWithError(err error) error {
	c.cancel()
	return err
}
