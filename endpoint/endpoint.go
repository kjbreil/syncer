package endpoint

import (
	"context"
	"errors"
	"fmt"
	"github.com/kjbreil/syncer/control"
	"github.com/kjbreil/syncer/injector"
	"io"
	"log"
	"log/slog"
	"math/rand"
	"net"
	"os"
	"time"
)

// Endpoint contains both the server and the client
// The clients first attempt to connect to external servers
// server then starts up
type Endpoint struct {
	port     int
	peers    []net.TCPAddr
	server   *server
	client   *client
	data     any
	injector *injector.Injector
	Errors   chan error
	started  bool
	logger   *slog.Logger
	ctx      context.Context
	cancel   context.CancelFunc
}

// New creates a new endpoint with the given port and peers
// Port is the port number of the server, all peer servers will listen on this port
func New(data any, port int, peers []net.TCPAddr) (*Endpoint, error) {
	// TODO: Validate that data is a pointer
	inj, err := injector.New(data)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Endpoint{
		port:     port,
		peers:    peers,
		server:   nil,
		client:   nil,
		data:     data,
		ctx:      ctx,
		cancel:   cancel,
		injector: inj,
		Errors:   make(chan error, 100),
		logger:   slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}, nil
}

func (e *Endpoint) Run(onlyClient bool) {
	go e.run(onlyClient)
	for !e.started {
		time.Sleep(100 * time.Millisecond)
	}
}

func (e *Endpoint) IsServer() bool {
	return e.server != nil
}

func (e *Endpoint) run(onlyClient bool) {
	var err error

	checkPeersDuration := time.Minute
	checkPeersLast := time.Now()

	for {
		if e.ctx.Err() != nil {
			return
		}
		if !e.started {
			err = e.tryPeers(false)
			if err == nil {
				e.logger.Info("Client Started")
				e.started = true
			}
			if errors.Is(err, ErrClientServerNonAvailable) && !onlyClient {
				e.server, err = newServer(e.port, e.data, e.Errors)
				if err == nil {
					e.logger.Info("Server Started")
					e.started = true
					checkPeersLast = time.Now()
				}
			}
		}
		// check if the client exists but the context is canceled
		if e.client != nil && e.client.ctx.Err() != nil {
			e.logger.Info("Client Stopped")
			e.client = nil
		}
		if e.server != nil && e.server.ctx.Err() != nil {
			e.logger.Info("Server Stopped")
			e.server = nil
		}
		if e.server == nil && e.client == nil {
			e.started = false
		}
		if e.server != nil && time.Now().Sub(checkPeersLast) > checkPeersDuration {
			checkPeersLast = time.Now()
			_ = e.tryPeers(true)
		}

		// rand.Seed(time.Now().UnixNano())
		r := rand.Intn(900) + 100
		// try and connect to peers
		time.Sleep(time.Duration(r) * time.Millisecond)
	}

}

func (e *Endpoint) getChanges() error {
	update, err := e.client.c.Update(e.client.ctx, &control.Request{Type: control.Request_INIT})
	if err != nil {
		return err
	}
	for {
		in, err := update.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("client.RouteChat failed: %v", err)
		}
		log.Printf("Got message Value %s", in.Value)

	}
	return nil
}

var (
	ErrClientAlreadyConnected   = fmt.Errorf("client already connected")
	ErrClientServerNonAvailable = fmt.Errorf("client could not find available server to connect to")
)

func (e *Endpoint) tryPeers(stop bool) error {
	var err error
	if e.client != nil {
		return ErrClientAlreadyConnected
	}
	for _, p := range e.peers {
		e.client, err = newClient(p, e.data)
		if err == nil {
			if stop {
				e.client.shutdownRemoteServer()
				continue
			} else {
				e.client.init()
			}
			return nil
		}
		// TODO: Check error for if there is an injector problem (return error) or not available (continue)
	}
	return ErrClientServerNonAvailable
}

func (e *Endpoint) Stop() {
	e.cancel()
	if e.client != nil {
		e.client.stop()
	}
	if e.server != nil {
		e.server.stop()
	}
}

var (
	ErrClientNotAvailable = fmt.Errorf("could not dial client")
	ErrClientInjector     = fmt.Errorf("injector could not be created")
)
