package server

import (
	"context"
	"fmt"
	"time"

	"github.com/coffee-dns/coffee-dns/controller/api"
)

// Status returns the status of the controller
func (c Controller) Status(ctx context.Context, _ *api.ControllerHealthReq) (*api.ControllerHealthResp, error) {
	c.log.Trace("healthcheck from ", reqAddress(ctx))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := c.nsclient.Status(ctx)
	if err != nil {
		// avoid publishing the exact error because it could contain
		// private information such as internal ip addresses or kubernetes
		// service names
		return &api.ControllerHealthResp{}, fmt.Errorf("failed to connect to Coffee DNS resolver service")
	}
	return &api.ControllerHealthResp{}, nil
}
