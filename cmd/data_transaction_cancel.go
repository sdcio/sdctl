/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	sdcpb "github.com/sdcio/sdc-protos/sdcpb"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/prototext"
)

// dataSetIntentCmd represents the set-intent command
var dataTransactionCancelCmd = &cobra.Command{
	Use:          "cancel",
	Short:        "Cancel a transaction",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, _ []string) error {

		if transactionId == "" {
			return fmt.Errorf("")
		}

		req := &sdcpb.TransactionCancelRequest{
			TransactionId: transactionId,
			DatastoreName: datastoreName,
		}

		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()
		dataClient, err := createDataClient(ctx, addr)
		if err != nil {
			return err
		}
		fmt.Println("request:")
		fmt.Println(prototext.Format(req))
		rsp, err := dataClient.TransactionCancel(ctx, req)
		if err != nil {
			return err
		}
		fmt.Println("response:")
		fmt.Println(prototext.Format(rsp))
		return nil
	},
}

func init() {
	dataTransactionCmd.AddCommand(dataTransactionCancelCmd)
}
