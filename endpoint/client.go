package endpoint

import (
	"context"
	"errors"
	"fmt"
	"github.com/kjbreil/syncer/control"
	"github.com/kjbreil/syncer/injector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"net"
	"time"
)

type client struct {
	c    control.ConfigClient
	conn *grpc.ClientConn
	peer net.TCPAddr

	ctx    context.Context
	cancel context.CancelFunc
	inj    *injector.Injector
}

func (c *client) init() {
	update, err := c.c.Update(c.ctx, &control.Request{Type: control.Request_INIT})
	if err != nil {
		return
	}
	c.processUpdate(update)
}
func (c *client) shutdownRemoteServer() {
	_, err := c.c.Control(c.ctx, &control.Message{Action: control.Message_SHUTDOWN})
	if err != nil {
		return
	}
}

func (c *client) changes() {
	update, err := c.c.Update(c.ctx, &control.Request{Type: control.Request_INIT})
	if err != nil {
		return
	}
	c.processUpdate(update)
}

func (c *client) processUpdate(update control.Config_UpdateClient) {
	for {
		cfg, err := update.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatalf("client.init failed: %v", err)
		}
		c.inj.Add(cfg)
	}
}

func newClient(peer net.TCPAddr, data any) (*client, error) {
	var err error

	ctx, cancel := context.WithCancel(context.Background())

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	addr := fmt.Sprintf("%s:%d", peer.IP.String(), peer.Port)

	c := &client{
		peer:   peer,
		ctx:    ctx,
		cancel: cancel,
	}

	c.conn, err = grpc.Dial(addr, opts...)
	if err != nil {
		cancel()
		return nil, ErrClientNotAvailable
	}

	c.c = control.NewConfigClient(c.conn)
	go func() {
		for {
			select {
			case <-time.After(time.Second):
				_, err = c.c.Control(ctx, &control.Message{Action: control.Message_PING})
				if err != nil {
					cancel()
				}
			case <-ctx.Done():
				c.conn.Close()
				return
			}
		}
	}()

	_, err = c.c.Control(ctx, &control.Message{Action: control.Message_PING})
	if err != nil {
		return nil, c.closeWithError(ErrClientNotAvailable)
	}

	c.inj, err = injector.New(data)
	if err != nil {
		return nil, c.closeWithError(ErrClientNotAvailable)
	}

	return c, nil
}

func (c *client) closeWithError(err error) error {
	c.cancel()
	return err
}

func (c *client) stop() {
	c.cancel()
}
