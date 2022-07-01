package main

import (
	"fmt"
	"os"

	"github.com/coffee-dns/coffee-dns/controller/common"
	"github.com/coffee-dns/coffee-dns/controller/server"
	"github.com/coffee-dns/coffee-dns/internal/log"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start Coffee DNS Controller",
	Run: func(cmd *cobra.Command, args []string) {
		if err := startServer(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func startServer() error {
	var (
		grpcServer common.Port
		err        error
	)

	grpcServer, err = common.GetPortFromEnv(common.EnvGRPCPort)
	if err != nil {
		return err
	}

	s, err := server.New(uint(grpcServer), log.NewJSONLogger(os.Getenv(common.EnvLogLevel)))
	if err != nil {
		return err
	}

	return s.Start()
}
