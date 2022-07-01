package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/coffee-dns/coffee-dns/nameserver/api"
	"github.com/coffee-dns/coffee-dns/nameserver/record"
)

func (s *Server) GetRecord(ctx context.Context, req *api.NameserverGetRecordReq) (*api.NameserverGetRecordResp, error) {
	s.Logger.Infof(
		"get record request from %s: %s",
		reqAddress(ctx),
		req.RecordKey,
	)

	if !strings.HasSuffix(req.RecordKey, ".") {
		req.RecordKey += "."
	}

	v, err := s.Persister.Get(req.RecordKey)
	if err != nil {
		s.Logger.Errorf("failed to get record %s: %s", req.RecordKey, err)
		// TODO: we should check to see if record does not exist or if there was an
		// internal server error
		return nil, fmt.Errorf("record does not exist")
	}

	return &api.NameserverGetRecordResp{
		RecordKey:   req.RecordKey,
		RecordValue: v,
	}, nil
}

func (s *Server) CreateRecord(ctx context.Context, req *api.NameserverCreateRecordReq) (*api.NameserverCreateRecordResp, error) {
	s.Logger.Infof(
		"create record request from %s: %s %s %s %d",
		reqAddress(ctx),
		req.RecordType,
		req.RecordKey,
		req.RecordValue,
		req.RecordTTL,
	)

	if !strings.HasSuffix(req.RecordKey, ".") {
		req.RecordKey += "."
	}

	r := record.Record{
		Hostname: req.RecordKey,
		Value:    req.RecordValue,
		Type:     "A", // We only support ipv4 A records
	}

	if err := s.Persister.Set(r); err != nil {
		s.Logger.Errorf("failed to set record %s with value %s: %s", req.RecordKey, req.RecordValue, err)
		return nil, fmt.Errorf("failed to set record %s with value %s, internal server error", req.RecordKey, req.RecordValue)
	}

	s.Logger.Infof("record %s created with value %s", req.RecordKey, req.RecordValue)

	return &api.NameserverCreateRecordResp{}, nil
}

func (s *Server) DeleteRecord(ctx context.Context, req *api.NameserverDeleteRecordReq) (*api.NameserverDeleteRecordResp, error) {
	s.Logger.Infof(
		"delete record request from %s: %s",
		reqAddress(ctx),
		req.RecordKey,
	)

	if !strings.HasSuffix(req.RecordKey, ".") {
		req.RecordKey += "."
	}

	if err := s.Persister.Delete(req.RecordKey); err != nil {
		s.Logger.Errorf("failed to delete record %s: %s", req.RecordKey, err)
		return nil, fmt.Errorf("failed to delete record %s, internal server error", req.RecordKey)
	}

	s.Logger.Infof("record %s deleted", req.RecordKey)
	return &api.NameserverDeleteRecordResp{}, nil
}
