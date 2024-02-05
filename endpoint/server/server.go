package server

import (
	"context"
	"fmt"
	"github.com/kjbreil/syncer/control"
	"github.com/kjbreil/syncer/extractor"
	"google.golang.org/grpc"
	"net"
	"sync"
)

type Server struct {
	control.UnimplementedConfigServer
	grpcServer *grpc.Server
	errors     chan error
	extractor  *extractor.Extractor
	data       any
	ctx        context.Context
	cancel     context.CancelFunc
	wg         *sync.WaitGroup
}

var (
	ErrServerExited = fmt.Errorf("server exited")
	ErrServerListen = fmt.Errorf("server could not start listening")
)

func New(ctx context.Context, wg *sync.WaitGroup, data any, port int, errors chan error) (*Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrServerListen, err)
	}
	var opts []grpc.ServerOption

	ext := extractor.New(data)

	var s = &Server{
		grpcServer: grpc.NewServer(opts...),
		errors:     errors,
		extractor:  ext,
		data:       data,
		wg:         wg,
	}
	s.ctx, s.cancel = context.WithCancel(ctx)

	control.RegisterConfigServer(s.grpcServer, s)
	go func() {
		err := s.grpcServer.Serve(lis)
		if err != nil {
			s.errors <- fmt.Errorf("%w: %w", ErrServerExited, err)
		}
		s.errors <- ErrServerExited
		s.cancel()
	}()

	wg.Add(1)
	go func() {
		<-s.ctx.Done()
		s.grpcServer.Stop()
		wg.Done()
		return
	}()

	return s, nil
}

func (s *Server) Running() bool {
	return s.ctx.Err() == nil
}

func (s *Server) Update(req *control.Request, srv control.Config_UpdateServer) error {
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

func (s *Server) Control(_ context.Context, message *control.Message) (*control.Response, error) {
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