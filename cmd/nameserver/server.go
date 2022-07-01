package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/coffee-dns/coffee-dns/nameserver/persist"
	"github.com/coffee-dns/coffee-dns/nameserver/server"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start Coffee DNS Nameserver",
	Run: func(cmd *cobra.Command, args []string) {
		if err := startServer(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

var (
	grpcAddress string
	grpcPort    int

	resolverAddress string
	resolverPort    int
)

const (
	envCloudDatastoreProjectID = "CLOUD_DATASTORE_PROJECT_ID"
	envCloudDatastoreKind      = "CLOUD_DATASTORE_KIND"
)

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.PersistentFlags().StringVar(&grpcAddress, "grpc-address", "0.0.0.0", "Nameserver address to bind to")
	serverCmd.PersistentFlags().IntVar(&grpcPort, "grpc-port", 5555, "Nameserver port")
	serverCmd.PersistentFlags().StringVar(&resolverAddress, "resolver-address", "0.0.0.0", "Nameserver address to bind to")
	serverCmd.PersistentFlags().IntVar(&resolverPort, "resolver-port", 53, "Nameserver port")
}

func startServer() error {
	// TODO: attempt to detect project id if env not set
	projectID := os.Getenv(envCloudDatastoreProjectID)
	if projectID == "" {
		return fmt.Errorf("failed to detect google project using environment variable %s", envCloudDatastoreProjectID)
	}
	kind := os.Getenv(envCloudDatastoreKind)
	if kind == "" {
		return fmt.Errorf("failed to detect google cloud datastore kind using environment variable %s", envCloudDatastoreKind)
	}
	var persister persist.Persist = &persist.Datastore{
		ProjectID: projectID,
		Kind:      kind,
	}
	if err := persister.Init(); err != nil {
		return fmt.Errorf("failed to init persistent backend: %s", err)
	}

	p := os.Getenv("NAMESERVER_PORT")
	if p != "" {
		port, err := strconv.Atoi(p)
		if err == nil {
			logger.Infof("Setting resolver port from environment: %d", port)
			resolverPort = port
		}
	}

	ns := server.Server{
		APIConf: server.APIConfig{
			Address: grpcAddress,
			Port:    grpcPort,
		},
		ResolverConf: server.ResolverConfig{
			Address: resolverAddress,
			Port:    resolverPort,
		},
		Persister: persister,
		Logger:    logger,
	}

	if err := ns.Init(); err != nil {
		return fmt.Errorf("failed to init server: %s", err)
	}

	return ns.Start()
}
