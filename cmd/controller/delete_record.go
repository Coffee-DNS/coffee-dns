package main

import (
	"fmt"
	"os"

	"github.com/coffee-dns/coffee-dns/controller/client"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var deleteRecordCmd = &cobra.Command{
	Use:   "record",
	Short: "Delete DNS record",
	Run: func(cmd *cobra.Command, args []string) {
		if err := deleteRecord(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	deleteCmd.AddCommand(deleteRecordCmd)
	deleteRecordCmd.Flags().StringVarP(&recordKey, "name", "", "", "Record Name")
	deleteRecordCmd.MarkFlagRequired("name") // #nosec - Cobra flag error handling is not necessary
}

func deleteRecord() error {
	c, err := client.New(controllerEndpoint, controllerTLS)
	if err != nil {
		return err
	}

	if err := c.DeleteRecord(recordKey); err != nil {
		return errors.Wrap(err, "delete record failed")
	}

	fmt.Printf("Key %s deleted\n", recordKey)
	return nil
}
