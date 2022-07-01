package main

import (
	"context"
	"fmt"
	"os"

	"github.com/coffee-dns/coffee-dns/nameserver/client"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var getRecordCmd = &cobra.Command{
	Use:   "record",
	Short: "Get DNS record",
	Run: func(cmd *cobra.Command, args []string) {
		if err := getRecord(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	getCmd.AddCommand(getRecordCmd)
	getRecordCmd.Flags().StringVarP(&recordKey, "name", "", "", "Record Name")
	getRecordCmd.MarkFlagRequired("name") // #nosec - Cobra flag error handling is not necessary
}

func getRecord() error {
	c, err := client.New(nameserverEndpoint, nameserverTLS)
	if err != nil {
		return err
	}

	value, err := c.GetRecord(context.Background(), recordKey)
	if err != nil {
		return errors.Wrap(err, "get record failed")
	}

	fmt.Printf("Key: %s, Value: %s\n", value.RecordKey, value.RecordValue)
	return nil
}
