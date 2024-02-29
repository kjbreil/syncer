package server

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
	"google.golang.org/grpc/status"
)

type Server struct {
	control.UnsafeControlServer
	grpcServer *grpc.Server

	logger   *slog.Logger
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
	ErrServerExited   = errors.New("server exited")
	ErrServerListen   = errors.New("server could not start listening")
	ErrServerInjector = errors.New("server could not create injector")
)

func New(ctx context.Context, wg *sync.WaitGroup, data any, stngs *settings.Settings, errChan chan *slog.Record) (*Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", stngs.Port))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrServerListen, err)
	}
	var opts []grpc.ServerOption

	s := &Server{
		grpcServer: grpc.NewServer(opts...),
		logger:     slog.New(slogchannel.Option{Level: slog.LevelDebug, Channel: errChan}.NewChannelHandler()),
		// extractor:  ext,
		data: data,
		wg:   wg,
	}
	s.ctx, s.cancel = context.WithCancel(ctx)

	s.combined, err = combined.New(s.ctx, data)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrServerInjector, err)
	}

	control.RegisterControlServer(s.grpcServer, s)
	go func() {
		err := s.grpcServer.Serve(lis)
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			s.logger.Error(fmt.Errorf("%w: %w", ErrServerExited, err).Error())
		}
		s.logger.Error(ErrServerExited.Error())

		s.cancel()
	}()

	wg.Add(1)
	go func() {
		<-s.ctx.Done()
		s.grpcServer.Stop()
		wg.Done()
	}()

	return s, nil
}

func (s *Server) Running() bool {
	return s.ctx.Err() == nil
}

func (s *Server) AddExtHandler(ext func() error) {
	s.combined.ExtractorChanges(ext)
}

func (s *Server) AddInjHandler(inj func() error) {
	s.combined.InjectorChanges(inj)
}

func (s *Server) Control(_ context.Context, message *control.Message) (*control.Response, error) {
	switch message.GetAction() {
	case control.Message_PING:
		return &control.Response{}, nil
	case control.Message_SHUTDOWN:
		s.cancel()
		return &control.Response{}, nil
	default:
		return &control.Response{}, nil
	}
}

func (s *Server) Pull(req *control.Request, srv control.Control_PullServer) error {
	switch req.GetType() {
	case control.Request_INIT:
		s.combined.Reset()
		fallthrough
	case control.Request_CHANGES:
		entries, _ := s.combined.Entries(s.data)
		for _, e := range entries {
			err := srv.Send(e)
			if err != nil {
				s.logger.Error(err.Error())
			}
		}
	}

	return nil
}

func (s *Server) Push(server control.Control_PushServer) error {
	// TODO implement me
	panic("implement me")
}

func (s *Server) PushPull(server control.Control_PushPullServer) error {
	ctx, cancel := context.WithCancel(s.ctx)

	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	checkInterval := time.Second

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		for {
			select {
			case <-time.After(checkInterval):
				mu.Lock()
				entries, err := s.combined.Entries(s.data)
				if err != nil {
					s.logger.Error(err.Error())
					mu.Unlock()
				}
				for _, e := range entries {
					err = server.Send(e)
					if err != nil {
						s.logger.Error(err.Error())
						mu.Unlock()
						return
					}
				}
				mu.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		for {
			e, err := server.Recv()
			if errors.Is(err, io.EOF) {
				return
			}
			if stat, ok := status.FromError(err); ok {
				switch stat.Code() {
				case codes.OK:
				case codes.Canceled:
					return
				default:
					s.logger.Error("Server.PushPull() GRPC error: %s" + stat.String())
					return
				}
			}

			if err != nil {
				s.logger.Error(fmt.Errorf("Server.PushPull(): %w", err).Error())
				return
			}
			mu.Lock()
			err = s.combined.Add(e)
			_, _ = s.combined.Entries(s.data)
			mu.Unlock()
			if err != nil {
				s.logger.Error(fmt.Errorf("Server.PushPull(): %w", err).Error())
				return
			}
		}
	}()

	s.logger.Info("Server.PushPull() started")
	wg.Wait()
	s.logger.Info("Server.PushPull() stopped")

	return nil
}
