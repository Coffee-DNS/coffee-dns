package main

import (
	"fmt"
	"os"

	"github.com/coffee-dns/coffee-dns/internal/log"
	"github.com/spf13/cobra"
)

// variables set by command flags
var (
	nameserverEndpoint string
	nameserverTLS      bool

	recordType     string
	recordKey      string
	recordValue    string
	recordTTL      string // TODO: this should be a custom duration type
	recordOverwite bool
)

var rootCmd = &cobra.Command{
	Use:   "nameserver",
	Short: "Coffee DNS Nameserver",
}

var logger *log.Logger

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	logger = log.NewJSONLogger("trace")
	logger.SetOutput(os.Stdout)
}

func main() {
	Execute()
}
