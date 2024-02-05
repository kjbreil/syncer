package server

import (
	"context"
	"fmt"
	"github.com/kjbreil/syncer/combined"
	"github.com/kjbreil/syncer/control"
	"github.com/kjbreil/syncer/endpoint/settings"
	"google.golang.org/grpc"
	"net"
	"sync"
)

type Server struct {
	control.UnimplementedConfigServer
	grpcServer *grpc.Server
	errors     chan error

	combined *combined.Combined

	// extractor *extractor.Extractor
	// // server injector not used yet
	// injector *injector.Injector

	data   any
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
}

var (
	ErrServerExited   = fmt.Errorf("server exited")
	ErrServerListen   = fmt.Errorf("server could not start listening")
	ErrServerInjector = fmt.Errorf("server could not create injector")
)

func New(ctx context.Context, wg *sync.WaitGroup, data any, settings *settings.Settings, errors chan error) (*Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", settings.Port))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrServerListen, err)
	}
	var opts []grpc.ServerOption

	var s = &Server{
		grpcServer: grpc.NewServer(opts...),
		errors:     errors,
		// extractor:  ext,
		data: data,
		wg:   wg,
	}
	s.ctx, s.cancel = context.WithCancel(ctx)

	// s.injector, err = injector.New(data)
	s.combined, err = combined.New(data)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrServerInjector, err)
	}

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

func (s *Server) Pull(req *control.Request, srv control.Config_PullServer) error {
	switch req.GetType() {
	case control.Request_INIT:
		s.combined.Reset()
		fallthrough
	case control.Request_CHANGES:
		head := s.combined.Diff(s.data)
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
