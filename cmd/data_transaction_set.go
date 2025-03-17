/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"time"

	sdcpb "github.com/sdcio/sdc-protos/sdcpb"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/prototext"
)

var (
	dryRun             bool
	intents            []string
	replaceIntent      string
	transactionTimeout time.Duration
)

// dataSetIntentCmd represents the set-intent command
var dataTransactionSetCmd = &cobra.Command{
	Use:          "set",
	Short:        "Start a transaction",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, _ []string) error {

		transactionTimeoutInt32 := int32(transactionTimeout.Seconds())

		req := &sdcpb.TransactionSetRequest{
			TransactionId: transactionId,
			DatastoreName: datastoreName,
			DryRun:        dryRun,
			Intents:       []*sdcpb.TransactionIntent{},
			Timeout:       &transactionTimeoutInt32,
		}

		for _, transactionIntentFile := range intents {
			ti, err := LoadTransactionIntent(transactionIntentFile)
			if err != nil {
				return err
			}
			sdcpbTi, err := ti.ToSdcpbTransactionIntent()
			if err != nil {
				return err
			}
			req.Intents = append(req.Intents, sdcpbTi)
		}

		if replaceIntent != "" {
			ti, err := LoadTransactionIntent(replaceIntent)
			if err != nil {
				return err
			}
			sdcpbTi, err := ti.ToSdcpbTransactionIntent()
			if err != nil {
				return err
			}
			req.ReplaceIntent = sdcpbTi
		}

		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()
		dataClient, err := createDataClient(ctx, addr)
		if err != nil {
			return err
		}
		fmt.Println("request:")
		fmt.Println(prototext.Format(req))
		rsp, err := dataClient.TransactionSet(ctx, req)
		if err != nil {
			return err
		}

		fmt.Println("response:")
		fmt.Println(prototext.Format(rsp))

		errs := rsp.GetErrors()
		if len(errs) > 0 {
			return fmt.Errorf("validation errors present")
		}

		return nil
	},
}

func init() {
	dataTransactionCmd.AddCommand(dataTransactionSetCmd)
	dataTransactionSetCmd.Flags().StringSliceVarP(&intents, "intent", "i", nil, "path to intent file")
	dataTransactionSetCmd.Flags().StringVar(&replaceIntent, "replace-intent", "", "path to replace-intent file")
	dataTransactionSetCmd.Flags().DurationVar(&transactionTimeout, "transaction-timeout", 60*time.Second, "timeout before the transaction is auto-rolbacked")
	dataTransactionSetCmd.MarkFlagFilename("intent", "json")
	dataTransactionSetCmd.MarkFlagFilename("replace-intent", "json")
}
