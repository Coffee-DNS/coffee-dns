package main

import (
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete resources",
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
