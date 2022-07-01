package server

import (
	"context"

	"google.golang.org/grpc/peer"
)

func reqAddress(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "unknown"
	}
	return p.Addr.String()
}
