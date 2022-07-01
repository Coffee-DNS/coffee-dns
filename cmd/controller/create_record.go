package main

import (
	"fmt"
	"os"

	"github.com/coffee-dns/coffee-dns/controller/client"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var createRecordCmd = &cobra.Command{
	Use:   "record",
	Short: "Create DNS record",
	Run: func(cmd *cobra.Command, args []string) {
		if err := createRecord(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	createCmd.AddCommand(createRecordCmd)

	createRecordCmd.Flags().StringVarP(&recordType, "type", "", "", "Record type")
	createRecordCmd.Flags().StringVarP(&recordKey, "name", "", "", "Record Name")
	createRecordCmd.Flags().StringVarP(&recordValue, "value", "", "", "Record Value")
	createRecordCmd.Flags().StringVarP(&recordTTL, "ttl", "", "12h", "Record TTL")
	createRecordCmd.Flags().BoolVarP(&recordOverwite, "force", "", false, "Replace record if it already exists")

	createRecordCmd.MarkFlagRequired("type")  // #nosec - Cobra flag error handling is not necessary
	createRecordCmd.MarkFlagRequired("name")  // #nosec - Cobra flag error handling is not necessary
	createRecordCmd.MarkFlagRequired("value") // #nosec - Cobra flag error handling is not necessary
}

func createRecord() error {
	c, err := client.New(controllerEndpoint, controllerTLS)
	if err != nil {
		return err
	}

	uri, err := c.CreateRecord(
		recordType,
		recordKey,
		recordValue,
		3000, // TODO: Actually use recordTTL after converting to int32
		recordOverwite,
	)
	if err != nil {
		return errors.Wrap(err, "create record failed")
	}

	fmt.Printf("Record created. Update URI: %s\n", uri)
	return nil
}
