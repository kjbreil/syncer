package endpoint

import (
	"context"
	"crypto/rand"
	"errors"
	"log/slog"
	"math/big"
	"net"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/kjbreil/syncer/endpoint/client"
	"github.com/kjbreil/syncer/endpoint/server"
	settings2 "github.com/kjbreil/syncer/endpoint/settings"
)

var (
	ErrNotPointer               = errors.New("data must be a pointer")
	ErrClientAlreadyConnected   = errors.New("client already connected")
	ErrClientServerNonAvailable = errors.New("client could not find available server to connect to")
)

// Endpoint contains both the server and the Client
// The clients first attempt to connect to external servers
// server then starts up.
type Endpoint struct {
	// port   int `extractor:"-"`
	// peers  []net.TCPAddr
	settings *settings2.Settings
	localIP  []net.IP
	server   *server.Server    `extractor:"-"`
	client   *client.Client    `extractor:"-"`
	data     any               `extractor:"-"`
	Errors   chan *slog.Record `extractor:"-"`
	logger   *slog.Logger      `extractor:"-"`

	ctx    context.Context    `extractor:"-"`
	cancel context.CancelFunc `extractor:"-"`
	wg     *sync.WaitGroup    `extractor:"-"`
}

// New creates a new Endpoint with the given data and settings.
// The data must be a pointer to a struct.
// If the data is not a pointer, an error is returned.
func New(data any, stngs *settings2.Settings) (*Endpoint, error) {
	if reflect.ValueOf(data).Kind() != reflect.Ptr {
		return nil, ErrNotPointer
	}

	ctx, cancel := context.WithCancel(context.Background())

	ep := &Endpoint{
		settings: stngs,
		server:   nil,
		client:   nil,
		data:     data,
		ctx:      ctx,
		cancel:   cancel,
		wg:       &sync.WaitGroup{},
		Errors:   make(chan *slog.Record, 100),
		logger:   slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	return ep, nil
}

// Run starts the endpoint. If onlyClient is true, the endpoint will only run as a client.
// If the endpoint is already running, this function does nothing.
func (e *Endpoint) Run(onlyClient bool) {
	if e.Running() {
		return
	}
	e.ctx, e.cancel = context.WithCancel(context.Background())
	// add two WG because there are two goroutines started in e.run
	e.wg.Add(2)
	go e.run(onlyClient)
	// for !e.Running() {
	// 	time.Sleep(100 * time.Millisecond)
	// }
}

// IsServer returns true if the endpoint is running as a server.
func (e *Endpoint) IsServer() bool {
	return e.server != nil
}

// Wait blocks until the endpoint is stopped.
func (e *Endpoint) Wait() {
	e.wg.Wait()
}

// SetLogger sets the logger for the endpoint.
func (e *Endpoint) SetLogger(handler slog.Handler) {
	e.logger = slog.New(handler)
}

func (e *Endpoint) run(onlyClient bool) {
	e.localIP = nil

	var err error
	e.ctx, e.cancel = context.WithCancel(context.Background())
	checkPeersDuration := time.Minute
	checkPeersLast := time.Now()

	go func() {
		for {
			select {
			case <-e.ctx.Done():
				e.wg.Done()
				return
			case err := <-e.Errors:
				e.logger.Handler().Handle(e.ctx, *err)
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
				e.server, err = server.New(e.ctx, e.wg, e.data, e.settings, e.Errors)

				if err == nil {
					e.logger.Info("Server Started")
					checkPeersLast = time.Now()

					ifaces, err := net.Interfaces()
					if err != nil {
						continue
					}
					e.localIP = nil
					// handle err
					for _, i := range ifaces {
						addrs, err := i.Addrs()
						if err != nil {
							continue
						}
						for _, addr := range addrs {
							var ip net.IP
							switch v := addr.(type) {
							case *net.IPNet:
								ip = v.IP
							case *net.IPAddr:
								ip = v.IP
							}
							e.localIP = append(e.localIP, ip)
							// process IP address
						}
					}
				}
			}
		}
		// check if the Client exists but the context is canceled
		if e.client != nil && !e.client.Running() {
			e.logger.Info("Client Stopped")
			e.localIP = nil
			e.client = nil
		}
		if e.server != nil && !e.server.Running() {
			e.logger.Info("Server Stopped")
			e.localIP = nil
			e.server = nil
		}
		if e.server != nil && time.Since(checkPeersLast) > checkPeersDuration {
			e.logger.Info("testing for other servers")
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
peerLoop:
	for _, peer := range e.settings.Peers {
		for _, ip := range e.localIP {
			if ip.Equal(peer.IP) {
				continue peerLoop
			}
		}
		e.client, err = client.New(e.ctx, e.wg, e.data, peer, e.Errors, e.settings)
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
	e.logger.Info("stopping endpoint")
	e.cancel()
	e.logger.Info("endpoint stopped")
	e.wg.Wait()
	e.client = nil
	e.server = nil
}

// Running returns true if the endpoint is running.
func (e *Endpoint) Running() bool {
	return e.server != nil || e.client != nil
}

func (e *Endpoint) ClientUpdate() {
	if e.client != nil {
		e.client.Changes()
	}
}

func randomInt(l, h int) int {
	r, err := rand.Int(rand.Reader, big.NewInt(int64(h-l)))
	if err != nil {
		// if random generation fails return the middle of the low/high
		return h / l
	}

	return int(r.Int64()) + l
}
