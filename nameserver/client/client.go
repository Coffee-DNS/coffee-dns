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

type NameServer struct {
	nameserver api.NameserverClient
}

func (c NameServer) Status(ctx context.Context) error {
	_, err := c.nameserver.Status(ctx, &api.NameserverHealthReq{})
	if err != nil {
		return errors.Wrap(err, "gRPC healthcheck")
	}
	return nil
}

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

func (c NameServer) DeleteRecord(ctx context.Context, key string) error {
	_, err := c.nameserver.DeleteRecord(
		ctx,
		&api.NameserverDeleteRecordReq{
			RecordKey: key,
		},
	)
	return err
}

func New(endpoint string, enableTls bool) (NameServer, error) {
	secure := grpc.WithInsecure()
	if enableTls {
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
