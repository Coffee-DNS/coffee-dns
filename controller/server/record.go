package server

import (
	"context"
	"fmt"
	"time"

	"github.com/coffee-dns/coffee-dns/controller/api"
)

func (c Controller) GetRecord(ctx context.Context, req *api.ControllerGetRecordReq) (*api.ControllerGetRecordResp, error) {
	c.log.Infof("get record request from %s: %s", reqAddress(ctx), req.RecordKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	resp, err := c.nsclient.GetRecord(ctx, req.RecordKey)
	if err != nil {
		c.log.Errorf("failed to get record %s: %s", req.RecordKey, err)
		return nil, fmt.Errorf("failed to get record %s, internal server error", req.RecordKey)
	}
	return &api.ControllerGetRecordResp{
		RecordKey:   resp.RecordKey,
		RecordValue: resp.RecordValue,
	}, nil
}

func (c Controller) CreateRecord(ctx context.Context, req *api.ControllerCreateRecordReq) (*api.ControllerCreateRecordResp, error) {
	c.log.Infof("create record request from %s: %s %s %s %d", reqAddress(ctx), req.RecordType, req.RecordKey, req.RecordValue, req.RecordTTL)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	rType := req.RecordType
	rKey := req.RecordKey
	rValue := req.RecordValue
	ttl := req.RecordTTL
	force := req.OverWrite

	_, err := c.nsclient.CreateRecord(ctx, rType, rKey, rValue, ttl, force)
	if err != nil {
		c.log.Errorf("failed to create record %s with value %s: %s", rKey, rValue, err)
		return nil, fmt.Errorf("failed to create record %s with value %s, internal server error", rKey, rValue)
	}

	x := api.ControllerCreateRecordResp{
		RecordUpdateURI: "record created, update uri not implemented yet",
	}
	return &x, nil
}

func (c Controller) DeleteRecord(ctx context.Context, req *api.ControllerDeleteRecordReq) (*api.ControllerDeleteRecordResp, error) {
	c.log.Infof("delete record request from %s: %s", reqAddress(ctx), req.RecordKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := c.nsclient.DeleteRecord(ctx, req.RecordKey)
	if err != nil {
		c.log.Errorf("failed to delete record %s: %s", req.RecordKey, err)
		return nil, fmt.Errorf("failed to delete record %s, internal server error", req.RecordKey)
	}
	return &api.ControllerDeleteRecordResp{}, nil
}
