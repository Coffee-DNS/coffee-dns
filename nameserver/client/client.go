package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/coffee-dns/coffee-dns/nameserver/api"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// NameServer is a Coffee DNS nameserver client
type NameServer struct {
	nameserver api.NameserverClient
}

// Status returns the status of the nameserver
func (c NameServer) Status(ctx context.Context) error {
	_, err := c.nameserver.Status(ctx, &api.NameserverHealthReq{})
	if err != nil {
		return errors.Wrap(err, "gRPC healthcheck")
	}
	return nil
}

// CreateRecord creates a DNS record
func (c NameServer) CreateRecord(ctx context.Context, rType, rKey, rValue string, ttl int32, force bool) (*api.NameserverCreateRecordResp, error) {
	req := api.NameserverCreateRecordReq{
		RecordType:  rType,
		RecordKey:   rKey,
		RecordValue: rValue,
		RecordTTL:   ttl,
		OverWrite:   force,
	}
	resp, err := c.nameserver.CreateRecord(ctx, &req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetRecord returns a DNS record
func (c NameServer) GetRecord(ctx context.Context, key string) (*api.NameserverGetRecordResp, error) {
	resp, err := c.nameserver.GetRecord(
		ctx,
		&api.NameserverGetRecordReq{
			RecordKey: key,
		},
	)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteRecord deletes a DNS record
func (c NameServer) DeleteRecord(ctx context.Context, key string) error {
	_, err := c.nameserver.DeleteRecord(
		ctx,
		&api.NameserverDeleteRecordReq{
			RecordKey: key,
		},
	)
	return err
}

// New returns a new Nameserver
func New(endpoint string, enableTLS bool) (NameServer, error) {
	//lint:ignore SA1019 https://github.com/Coffee-DNS/coffee-dns/issues/2
	secure := grpc.WithInsecure()
	if enableTLS {
		h2creds := credentials.NewTLS(&tls.Config{
			NextProtos: []string{"h2"},
			MinVersion: tls.VersionTLS12,
		})
		secure = grpc.WithTransportCredentials(h2creds)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	conn, err := grpc.DialContext(ctx, endpoint, secure, grpc.WithBlock())
	if err != nil {
		return NameServer{}, errors.Wrap(err, fmt.Sprintf("failed to connect to Coffee DNS Nameserver at %s", endpoint))
	}

	var n NameServer
	n.nameserver = api.NewNameserverClient(conn)

	return n, nil
}
