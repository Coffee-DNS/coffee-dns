package main

import (
	"fmt"
	"os"

	"github.com/coffee-dns/coffee-dns/controller/client"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of the Coffee DNS Controller",
	Run: func(cmd *cobra.Command, args []string) {
		if err := status(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func status() error {
	c, err := client.New(controllerEndpoint, controllerTLS)
	if err != nil {
		return err
	}

	if err := c.Status(); err != nil {
		return err
	}

	fmt.Println("Coffee DNS Controller is healthy")
	return nil
}
