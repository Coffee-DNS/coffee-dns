package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/coffee-dns/coffee-dns/discovery/node"
	"github.com/coffee-dns/coffee-dns/internal/log"
	"github.com/gin-gonic/gin"
)

var logLevel = "info"

const logLevelEnv = "COFFEE_LOG_LEVEL"

type server struct {
	Address           string
	Port              int
	DetectionInterval int
	Logger            *log.Logger

	lock sync.Mutex

	nodes []k8snode
}

// https://github.com/yyyar/gobetween/wiki/Discovery#json
type k8snode struct {
	Host     string `json:"host,omitempty"`
	Port     string `json:"port,omitempty"`
	Weight   int    `json:"weight,omitempty"`
	Priority int    `json:"priority,omitempty"`
	Sni      string `json:"sni,omitempty"`
}

func init() {
	switch l := os.Getenv(logLevelEnv); l {
	case "trace", "info", "warn", "error":
		logLevel = l
	}
}

func main() {
	args := os.Args

	if len(args) != 2 || args[1] != "server" {
		fmt.Println("expected one argument 'server'")
		os.Exit(1)
	}

	s := server{
		Address:           "0.0.0.0",
		Port:              8080,
		DetectionInterval: 60, // seconds

		Logger: log.NewJSONLogger(logLevel),
	}
	s.Logger.SetOutput(os.Stdout)

	go s.detect()
	s.server()
}

func (s *server) server() {
	addr := s.ListenAddress()
	s.Logger.Info("starting http server on %s", addr)
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, s.Nodes())
	})
	if err := router.Run(addr); err != nil {
		s.Logger.Errorf("server exited with error: %s", err)
	}
	s.Logger.Infof("server stopped")
}

func (s *server) detect() {
	s.Logger.Infof("starting node detection service with interval %d", s.DetectionInterval)

	port := "30053"
	weight := 1
	priority := 1

	for {
		s.lock.Lock()

		s.Logger.Infof("aquired lock, beginning node detection")

		nodes, err := node.Nodes()
		if err != nil {
			s.Logger.Errorf("failed to list nodes: %s", err)
			s.lock.Unlock()
			continue
		}

		// Reset the node list
		s.nodes = []k8snode{}

		for _, n := range nodes {
			if n.Spec.Unschedulable {
				s.Logger.Tracef("skipping unschedulable node: %s", n.Name)
				continue
			}

			for _, address := range n.Status.Addresses {
				if address.Type != "InternalIP" {
					s.Logger.Tracef("skipping non InternalIP address type '%s' for node %s", address.Type, n.Name)
					continue
				}

				node := k8snode{
					Host:     address.Address,
					Port:     port,
					Weight:   weight,
					Priority: priority,
				}
				s.Logger.Tracef("adding node %s with address %s", n.Name, node.Host)
				s.nodes = append(s.nodes, node)
			}
		}

		s.lock.Unlock()

		s.Logger.Infof("node detection finished, lock released")

		time.Sleep(time.Second * time.Duration(s.DetectionInterval))
	}
}

func (s *server) Nodes() []k8snode {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.nodes
}

func (s *server) ListenAddress() string {
	return fmt.Sprintf("%s:%d", s.Address, s.Port)
}
