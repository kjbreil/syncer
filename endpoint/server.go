package endpoint

import (
	"context"
	"fmt"
	"github.com/kjbreil/syncer/control"
	"github.com/kjbreil/syncer/extractor"
	"google.golang.org/grpc"
	"net"
)

type server struct {
	control.UnimplementedConfigServer
	grpcServer *grpc.Server
	errors     chan error
	extractor  *extractor.Extractor
	data       any
	ctx        context.Context
	cancel     context.CancelFunc
}

var (
	ErrServerExited = fmt.Errorf("server exited")
	ErrServerListen = fmt.Errorf("server could not start listening")
)

func newServer(port int, data any, errors chan error) (*server, error) {

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrServerListen, err)
	}
	var opts []grpc.ServerOption

	ext := extractor.New(data)

	ctx, cancel := context.WithCancel(context.Background())
	var s = &server{
		grpcServer: grpc.NewServer(opts...),
		errors:     errors,
		extractor:  ext,
		data:       data,
		ctx:        ctx,
		cancel:     cancel,
	}

	control.RegisterConfigServer(s.grpcServer, s)
	go func() {
		err := s.grpcServer.Serve(lis)
		if err != nil {
			s.errors <- fmt.Errorf("%w: %w", ErrServerExited, err)
		}
		s.errors <- ErrServerExited
		cancel()
	}()

	return s, nil
}

func (s *server) Update(req *control.Request, srv control.Config_UpdateServer) error {
	switch req.GetType() {
	case control.Request_INIT:
		s.extractor.Reset()
		fallthrough
	case control.Request_CHANGES:
		head := s.extractor.Diff(s.data)
		entries := head.Entries()
		for _, e := range entries {
			err := srv.Send(e)
			if err != nil {
				s.errors <- err
			}
		}
	}

	return nil
}

func (s *server) Control(ctx context.Context, message *control.Message) (*control.Response, error) {
	switch message.Action {
	case control.Message_PING:
		return &control.Response{}, nil
	case control.Message_SHUTDOWN:
		s.cancel()
		return &control.Response{}, nil
	default:
		return &control.Response{}, nil
	}
}

func (s *server) stop() {
	s.grpcServer.Stop()
}
