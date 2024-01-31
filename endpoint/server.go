package endpoint

import (
	"context"
	"github.com/kjbreil/syncer/control"
)

type server struct {
	control.UnimplementedConfigServer
}

func newServer() *server {
	return &server{}
}

func (s server) Update(req *control.Request, srv control.Config_UpdateServer) error {
	switch req.GetType() {
	case control.Request_CHANGES:
	case control.Request_INIT:
	}

	return nil
}

func (s server) Control(ctx context.Context, message *control.Message) (*control.Response, error) {
	// TODO implement me
	panic("implement me")
}
