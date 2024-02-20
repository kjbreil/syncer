package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
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
	control.UnsafeConfigServer
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

func New(ctx context.Context, wg *sync.WaitGroup, data any, settings *settings.Settings, errors chan *slog.Record) (*Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", settings.Port))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrServerListen, err)
	}
	var opts []grpc.ServerOption

	var s = &Server{
		grpcServer: grpc.NewServer(opts...),
		logger:     slog.New(slogchannel.Option{Level: slog.LevelDebug, Channel: errors}.NewChannelHandler()),
		// extractor:  ext,
		data: data,
		wg:   wg,
	}
	s.ctx, s.cancel = context.WithCancel(ctx)

	maxRetryCount := 5
	retryCount := 0
	for {
		s.combined, err = combined.New(data)
		if err == nil {
			break
		}
		retryCount++
		if retryCount > maxRetryCount {
			return nil, fmt.Errorf("%w: %w", ErrServerInjector, err)
		}
		delay := time.Duration(math.Pow(2, float64(retryCount))) * time.Second
		time.Sleep(delay)
	}

	control.RegisterConfigServer(s.grpcServer, s)
	go func() {
		err := s.grpcServer.Serve(lis)
		if err != nil && err != grpc.ErrServerStopped {
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

func (s *Server) Pull(req *control.Request, srv control.Config_PullServer) error {
	switch req.GetType() {
	case control.Request_INIT:
		s.combined.Reset()
		fallthrough
	case control.Request_CHANGES:
		head, _ := s.combined.Diff(s.data)
		entries := head.Entries()
		for _, e := range entries {
			err := srv.Send(e)
			if err != nil {
				s.logger.Error(err.Error())
			}
		}
	}

	return nil
}

func (s *Server) Push(server control.Config_PushServer) error {
	// TODO implement me
	panic("implement me")
}

func (s *Server) PushPull(server control.Config_PushPullServer) error {
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
				head, err := s.combined.Diff(s.data)
				if err != nil {
					s.logger.Error(err.Error())
					mu.Unlock()
				}
				entries := head.Entries()
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
			_, _ = s.combined.Diff(s.data)
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
