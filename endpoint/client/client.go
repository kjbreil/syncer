package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/kjbreil/syncer/combined"
	"github.com/kjbreil/syncer/control"
	"github.com/kjbreil/syncer/endpoint/settings"
	slogchannel "github.com/samber/slog-channel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
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
	data any

	logger *slog.Logger
}

func New(ctx context.Context, wg *sync.WaitGroup, data any, peer net.TCPAddr, errs chan *slog.Record, settings *settings.Settings) (*Client, error) {
	var err error

	c := &Client{
		peer:     peer,
		logger:   slog.New(slogchannel.Option{Level: slog.LevelDebug, Channel: errs}.NewChannelHandler()),
		settings: settings,
		data:     data,
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
		defer wg.Done()
		for {
			select {
			case <-time.After(time.Second * 5):
				_, err = c.c.Control(c.ctx, &control.Message{Action: control.Message_PING})
				if err != nil {
					c.logger.Error(err.Error())
					c.cancel()
				}
			case <-c.ctx.Done():

				err = c.conn.Close()
				if err != nil {
					c.logger.Error(err.Error())
					return
				}
				return
			}
		}
	}()

	if settings.AutoUpdate {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.PushPull()
		}()
	}

	_, err = c.c.Control(c.ctx, &control.Message{Action: control.Message_PING})
	if err != nil {
		return nil, c.closeWithError(fmt.Errorf("%w: %w", ErrClientNotAvailable, err))
	}

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
		c.logger.Error(fmt.Errorf("Client.Init(): %w", err).Error())
		return
	}
	c.processUpdate(update)
}

// ShutdownRemoteServer requests to shut down the server
func (c *Client) ShutdownRemoteServer() {
	_, err := c.c.Control(c.ctx, &control.Message{Action: control.Message_SHUTDOWN})
	if err != nil {
		c.logger.Error(fmt.Errorf("Client.ShutdownRemoteServer(): %w", err).Error())
		return
	}
}

func (c *Client) PushPull() {
	client, err := c.c.PushPull(c.ctx)
	if err != nil {
		c.logger.Error(fmt.Errorf("Client.PushPull(): %w", err).Error())
		return
	}
	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	checkInterval := time.Millisecond * 100

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer c.cancel()
		for {
			select {
			case <-time.After(checkInterval):
				mu.Lock()
				head, err := c.combined.Diff(c.data)
				if err != nil {
					c.logger.Error(err.Error())
					mu.Unlock()
				}
				entries := head.Entries()
				for _, e := range entries {
					err := client.Send(e)
					if err != nil {
						c.logger.Error(err.Error())
						mu.Unlock()
						return
					}
				}
				mu.Unlock()
			case <-c.ctx.Done():
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer c.cancel()
		for {
			e, err := client.Recv()
			if errors.Is(err, io.EOF) {
				return
			}
			if stat, ok := status.FromError(err); ok {
				switch stat.Code() {
				case codes.OK:
				case codes.Canceled:
					return
				case codes.Unavailable:
					c.logger.Error("Client.PushPull() GRPC Server became unavailable:")
					return
				default:
					c.logger.Error(fmt.Sprintf("Client.PushPull() GRPC error: %s", stat.String()))
					return
				}
			}
			if err != nil {
				c.logger.Error(fmt.Errorf("Client.PushPull(): %w", err).Error())
				return
			}
			mu.Lock()
			err = c.combined.Add(e)
			_, _ = c.combined.Diff(c.data)
			mu.Unlock()
			if err != nil {
				c.logger.Error(fmt.Errorf("Client.PushPull(): %w", err).Error())
				return
			}
		}
	}()
	c.logger.Info("Client.PushPull() started")

	wg.Wait()
	c.logger.Info("Client.PushPull() stopped")
}

func (c *Client) Changes() {
	update, err := c.c.Pull(c.ctx, &control.Request{Type: control.Request_CHANGES})
	if err != nil {
		c.logger.Error(fmt.Errorf("Client.changes(): %w", err).Error())
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
			c.logger.Error(fmt.Errorf("Client.processUpdate(): %w", err).Error())
			c.cancel()
			return
		}
		err = c.combined.Add(cfg)
		if err != nil {
			c.logger.Error(err.Error())
		}
	}
}

func (c *Client) closeWithError(err error) error {
	c.cancel()
	return err
}
