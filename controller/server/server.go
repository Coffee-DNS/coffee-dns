package server

import (
	"fmt"
	"net"

	"github.com/coffee-dns/coffee-dns/controller/api"
	"github.com/coffee-dns/coffee-dns/internal/log"

	"github.com/coffee-dns/coffee-dns/nameserver/client"
	"google.golang.org/grpc"
)

// Controller is a Coffee DNS controlplane
type Controller struct {
	port struct {
		grpc uint
	}

	nsclient client.NameServer

	log *log.Logger
	api.UnimplementedControllerServer
}

// Start starts the controlplane
func (c Controller) Start() error {
	nameserverAddress := "nameserver:5555"
	conn, err := client.New(nameserverAddress, false)
	if err != nil {
		return fmt.Errorf("failed to connect to nameserver service at endpoint %s: %s", nameserverAddress, err)
	}
	c.nsclient = conn
	c.log.Infof("connected to resolver service grcp interface at %s", nameserverAddress)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", c.port.grpc))
	if err != nil {
		c.log.Fatal(err, 1)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	api.RegisterControllerServer(grpcServer, c)
	c.log.Infof("starting coffee dns controller grpc interface on port %d", c.port.grpc)
	return grpcServer.Serve(lis)
}

// New returns a new Coffee DNS Controller
func New(grpcPort uint, logger *log.Logger) (Controller, error) {
	var c Controller
	c.port.grpc = grpcPort
	c.log = logger
	return c, nil
}
