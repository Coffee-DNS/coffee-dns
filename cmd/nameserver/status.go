package main

import (
	"context"
	"fmt"
	"os"

	"github.com/coffee-dns/coffee-dns/nameserver/client"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of the Coffee DNS Nameserver",
	Run: func(cmd *cobra.Command, args []string) {
		if err := status(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	rootCmd.PersistentFlags().StringVar(&nameserverEndpoint, "endpoint", "localhost:53", "Coffee DNS Controller gRCP endpoint")
	rootCmd.PersistentFlags().BoolVar(&nameserverTLS, "tls", false, "Enable or disable TLS while communicating with the controller endpoint")
}

func status() error {
	c, err := client.New(nameserverEndpoint, nameserverTLS)
	if err != nil {
		return err
	}

	if err := c.Status(context.Background()); err != nil {
		return err
	}

	fmt.Println("Coffee DNS Nameserver is healthy")
	return nil
}
