package server

import (
	"context"
	"net"
	"os"

	"github.com/coffee-dns/coffee-dns/nameserver/api"
	"github.com/coffee-dns/coffee-dns/nameserver/persist"
	"google.golang.org/grpc"

	"github.com/coffee-dns/coffee-dns/internal/log"
	"github.com/miekg/dns"
)

// Server is a Coffee DNS nameserver
type Server struct {
	*dns.Server
	APIConf      APIConfig
	ResolverConf ResolverConfig
	Persister    persist.Persist
	Logger       *log.Logger

	api.UnimplementedNameserverServer
}

// APIConfig is an API configuration for a Coffee DNS Server
type APIConfig struct {
	Address       string
	Port          int
	listenAddress string
}

// ResolverConfig is a resolver configuration for a Coffee DNS Server
type ResolverConfig struct {
	Address       string
	Port          int
	listenAddress string
}

// Start starts the nameserver
func (s *Server) Start() error {
	go s.startResolver()
	return s.startAPI()
}

func (s *Server) startResolver() {
	s.Logger.Infof("Starting resolver on %s", s.Server.Addr)
	if err := s.Server.ListenAndServe(); err != nil {
		s.Logger.Errorf("failed to set udp listener %s", err.Error())
	}

	// TODO: pass context or channel. If resolver fails, we need to kill
	// the pod, which is not possible because this function is run as a go routine
	os.Exit(1)
}

func (s *Server) startAPI() error {
	lis, err := net.Listen("tcp", s.APIConf.listenAddress)
	if err != nil {
		s.Logger.Errorf("failed to start api: %s", err)
	}

	var ops []grpc.ServerOption
	grpcServer := grpc.NewServer(ops...)

	api.RegisterNameserverServer(grpcServer, s)

	s.Logger.Infof("starting coffee dns nameserver grpc interface on %s", s.APIConf.listenAddress)

	return grpcServer.Serve(lis)
}

// ServeDNS implements github.com/miekg/dns Handler interface
func (s *Server) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	if len(r.Question) < 1 {
		s.Logger.Tracef("dns request had less than one question")
		return
	}

	msg := dns.Msg{}
	msg.SetReply(r)
	switch t := r.Question[0].Qtype; t {
	case dns.TypeA:
		msg.Authoritative = true
		domain := msg.Question[0].Name

		req := api.NameserverGetRecordReq{
			RecordKey: domain,
		}
		r, err := s.GetRecord(context.Background(), &req)
		if err != nil {
			s.Logger.Tracef("got error looking up domain %s: %s", domain, err)
			return
		}
		address := r.RecordValue

		s.Logger.Tracef("responding with address %s for domain %s", address, domain)
		msg.Answer = append(msg.Answer, &dns.A{
			Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
			A:   net.ParseIP(address),
		})

	}
	if err := w.WriteMsg(&msg); err != nil {
		s.Logger.Errorf("response failed: %s", err)
	}
}
