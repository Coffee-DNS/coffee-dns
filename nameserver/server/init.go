package server

import (
	"fmt"
	"net"

	"github.com/miekg/dns"
)

func (s *Server) Init() error {
	if err := s.initAPI(); err != nil {
		return fmt.Errorf("failed to init grpc interface: %s", err)
	}

	if err := s.initResolver(); err != nil {
		return fmt.Errorf("failed to init resolver: %s", err)
	}

	return nil
}

func (s *Server) initAPI() error {
	if s.APIConf.Address == "" {
		return fmt.Errorf("invalid grpc address")
	}

	if s.APIConf.Port == 0 {
		return fmt.Errorf("invalid grpc port")
	}

	s.APIConf.listenAddress = fmt.Sprintf("%s:%d", s.APIConf.Address, s.APIConf.Port)

	return nil
}

func (s *Server) initResolver() error {
	s.ResolverConf.listenAddress = fmt.Sprintf("%s:%d", s.ResolverConf.Address, s.ResolverConf.Port)

	if _, err := net.ResolveUDPAddr("udp", s.ResolverConf.listenAddress); err != nil {
		return fmt.Errorf("failed to resolve listen_address: %s", err)

	}
	s.Server = &dns.Server{
		Addr: s.ResolverConf.listenAddress,
		Net:  "udp",
	}

	s.Server.Handler = s

	return nil
}
