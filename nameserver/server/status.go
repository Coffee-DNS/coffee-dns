package server

import (
	"context"

	"github.com/coffee-dns/coffee-dns/nameserver/api"
)

func (s Server) Status(ctx context.Context, req *api.NameserverHealthReq) (*api.NameserverHealthResp, error) {
	// TODO: Actually return a status check that reflects the services health
	s.Logger.Trace("healthcheck from ", reqAddress(ctx))
	return &api.NameserverHealthResp{}, nil
}
