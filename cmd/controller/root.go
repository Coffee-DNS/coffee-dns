package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	defaultControllerEndpoint = "controller.coffeedns.net:443"
	defaultControllerTLS      = true
)

// variables set by command flags
var (
	controllerEndpoint string
	controllerTLS      bool

	recordType     string
	recordKey      string
	recordValue    string
	recordTTL      string // TODO: this should be a custom duration type
	recordOverwite bool
)

var rootCmd = &cobra.Command{
	Use:   "coffee",
	Short: "Manage your Coffee DNS Account",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&controllerEndpoint, "endpoint", defaultControllerEndpoint, "Coffee DNS Controller gRCP endpoint")
	rootCmd.PersistentFlags().BoolVar(&controllerTLS, "tls", defaultControllerTLS, "Enable or disable TLS while communicating with the controller endpoint")
}

func initConfig() {

}

func main() {
	Execute()
}
