package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/coffee-dns/coffee-dns/controller/api"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Controller struct {
	controller api.ControllerClient
}

func (c Controller) Status() error {
	_, err := c.controller.Status(context.Background(), &api.ControllerHealthReq{})
	if err != nil {
		return errors.Wrap(err, "gRPC healthcheck")
	}
	return nil
}

func (c Controller) CreateRecord(rType, rKey, rValue string, ttl int32, force bool) (string, error) {
	req := api.ControllerCreateRecordReq{
		RecordType:  rType,
		RecordKey:   rKey,
		RecordValue: rValue,
		RecordTTL:   ttl,
		OverWrite:   force,
	}
	resp, err := c.controller.CreateRecord(context.Background(), &req)
	if err != nil {
		return "", err
	}
	return resp.RecordUpdateURI, nil
}

func (c Controller) GetRecord(key string) (string, error) {
	resp, err := c.controller.GetRecord(
		context.Background(),
		&api.ControllerGetRecordReq{
			RecordKey: key,
		},
	)
	if err != nil {
		return "", err
	}
	return resp.RecordValue, nil
}

func (c Controller) DeleteRecord(key string) error {
	_, err := c.controller.DeleteRecord(
		context.Background(),
		&api.ControllerDeleteRecordReq{
			RecordKey: key,
		},
	)
	return err
}

func New(endpoint string, enableTls bool) (Controller, error) {
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
		return Controller{}, errors.Wrap(err, fmt.Sprintf("failed to connect to Coffee DNS Controller at %s", endpoint))
	}

	var c Controller
	c.controller = api.NewControllerClient(conn)

	return c, nil
}
