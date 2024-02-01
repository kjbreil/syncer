package endpoint

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/kjbreil/syncer/endpoint/client"
	"github.com/kjbreil/syncer/endpoint/server"
	"log/slog"
	"math/big"
	"net"
	"os"
	"reflect"
	"sync"
	"time"
)

var (
	ErrNotPointer               = fmt.Errorf("data must be a pointer")
	ErrClientAlreadyConnected   = fmt.Errorf("client already connected")
	ErrClientServerNonAvailable = fmt.Errorf("client could not find available server to connect to")
)

// Endpoint contains both the server and the Client
// The clients first attempt to connect to external servers
// server then starts up
type Endpoint struct {
	port   int
	peers  []net.TCPAddr
	server *server.Server
	client *client.Client
	data   any
	Errors chan error
	logger *slog.Logger

	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
}

// New creates a new endpoint with the given port and peers
// Port is the port number of the server, all peer servers will listen on this port
func New(data any, port int, peers []net.TCPAddr) (*Endpoint, error) {
	if reflect.ValueOf(data).Kind() != reflect.Ptr {
		return nil, ErrNotPointer
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Endpoint{
		port:   port,
		peers:  peers,
		server: nil,
		client: nil,
		data:   data,
		ctx:    ctx,
		cancel: cancel,
		wg:     &sync.WaitGroup{},
		Errors: make(chan error, 100),
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}, nil
}

func (e *Endpoint) Run(onlyClient bool) {
	if e.Running() {
		return
	}
	// add two WG because there are two goroutines started in e.run
	e.wg.Add(2)
	go e.run(onlyClient)
	for !e.Running() {
		time.Sleep(100 * time.Millisecond)
	}
}

// IsServer returns true if the endpoint is running as a server
func (e *Endpoint) IsServer() bool {
	return e.server != nil
}

// Wait blocks until the endpoint is stopped
func (e *Endpoint) Wait() {
	e.wg.Wait()
}

func (e *Endpoint) run(onlyClient bool) {
	var err error

	checkPeersDuration := time.Minute
	checkPeersLast := time.Now()

	go func() {
		for {
			select {
			case <-e.ctx.Done():
				e.wg.Done()
				return
			case err := <-e.Errors:
				e.logger.Error(err.Error())
			}
		}
	}()

	for {
		if e.ctx.Err() != nil {
			e.wg.Done()
			return
		}
		if !e.Running() {
			err = e.tryPeers(false)
			if err == nil {
				e.logger.Info("Client Started")
			}
			if errors.Is(err, ErrClientServerNonAvailable) && !onlyClient {
				e.server, err = server.New(e.ctx, e.wg, e.data, e.port, e.Errors)
				if err == nil {
					e.logger.Info("Server Started")
					checkPeersLast = time.Now()
				}
			}
		}
		// check if the Client exists but the context is canceled
		if e.client != nil && !e.client.Running() {
			e.logger.Info("Client Stopped")
			e.client = nil
		}
		if e.server != nil && !e.server.Running() {
			e.logger.Info("Server Stopped")
			e.server = nil
		}
		if e.server != nil && time.Since(checkPeersLast) > checkPeersDuration {
			checkPeersLast = time.Now()
			_ = e.tryPeers(true)
		}

		// try and connect to peers using random milliseconds between 100 and 1000
		time.Sleep(time.Duration(randomInt(100, 1000)) * time.Millisecond)
	}
}

func (e *Endpoint) tryPeers(stop bool) error {
	var err error
	if e.client != nil {
		return ErrClientAlreadyConnected
	}
	for _, peer := range e.peers {
		e.client, err = client.New(e.ctx, e.wg, e.data, peer, e.Errors)
		if err == nil {
			if stop {
				e.client.ShutdownRemoteServer()
				continue
			}
			e.client.Init()
			return nil
		}
		// TODO: Check error for if there is an injector problem (return error) or not available (continue)
	}
	return ErrClientServerNonAvailable
}
func (e *Endpoint) Stop() {
	e.cancel()
}

// Running returns true if the endpoint is running
func (e *Endpoint) Running() bool {
	return e.server != nil || e.client != nil
}

func randomInt(l, h int) int {
	r, err := rand.Int(rand.Reader, big.NewInt(int64(h-l)))
	if err != nil {
		// if random generation fails return the middle of the low/high
		return h / l
	}

	return int(r.Int64()) + l
}
