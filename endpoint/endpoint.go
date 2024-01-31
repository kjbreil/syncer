package endpoint

import (
	"fmt"
	"github.com/kjbreil/syncer/control"
	"github.com/kjbreil/syncer/extractor"
	"github.com/kjbreil/syncer/injector"
	"google.golang.org/grpc"
	"log"
	"net"
)

// Endpoint contains both the server and the client
// The clients first attempt to connect to external servers
// server then starts up
type Endpoint struct {
	port      int
	peers     map[string]*client
	server    *server
	client    *client
	data      any
	extractor *extractor.Extractor
	injector  *injector.Injector
}

// New creates a new endpoint with the given port and peers
// Port is the port number of the server, all peer servers will listen on this port
func New(data any, port int, peers []net.IP) (*Endpoint, error) {
	inj, err := injector.New(data)
	if err != nil {
		return nil, err
	}

	peersMap := make(map[string]*client)
	for _, peer := range peers {
		peersMap[string(peer.To16())] = nil
	}

	return &Endpoint{
		port:      port,
		peers:     peersMap,
		server:    newServer(),
		client:    newClient(),
		data:      data,
		extractor: extractor.New(data),
		injector:  inj,
	}, nil
}

func (e *Endpoint) Run(onlyClient bool) error {
	var err error
	if !onlyClient {
		err = e.runServer()
		if err != nil {
			return err
		}
	}
	for p, c := range e.peers {
		if c != nil {
			continue
		}
		e.peers[p] = e.runClient(net.IP(p), e.data)
		e.peers[p].init()
	}
	return err
}

func (e *Endpoint) runServer() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", e.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	control.RegisterConfigServer(grpcServer, e.server)
	err = grpcServer.Serve(lis)

	return err
}

func (e *Endpoint) runClient(peer net.IP, data any) *client {
	var opts []grpc.DialOption

	addr := fmt.Sprintf("%s:%d", peer.String(), e.port)

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil
	}
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)
	c := control.NewConfigClient(conn)

	inj, err := injector.New(data)
	if err != nil {
		panic(err)
	}

	return &client{
		connected: true,
		c:         c,
		inj:       inj,
	}
}
