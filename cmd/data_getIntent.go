/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	sdcpb "github.com/sdcio/sdc-protos/sdcpb"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/prototext"
)

// dataGetIntentCmd represents the get-intent command
var dataGetIntentCmd = &cobra.Command{
	Use:          "get-intent",
	Short:        "get intent data",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		req := &sdcpb.GetIntentRequest{
			DatastoreName: datastoreName,
			Intent:        intentName,
			Format:        sdcpb.Format_Intent_Format_JSON,
		}
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()
		dataClient, err := createDataClient(ctx, addr)
		if err != nil {
			return err
		}
		fmt.Println("request:")
		fmt.Println(prototext.Format(req))
		rsp, err := dataClient.GetIntent(ctx, req)
		if err != nil {
			return err
		}

		switch {
		case rsp.GetBlob() != nil:
			fmt.Println("Blob Data:")
			var v any
			json.Unmarshal(rsp.GetBlob(), &v)
			prettyByte, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(prettyByte))
		case rsp.GetProto() != nil:
			fmt.Println("response:")
			fmt.Println(prototext.Format(rsp))
		}

		return nil
	},
}

func init() {
	dataCmd.AddCommand(dataGetIntentCmd)
	dataGetIntentCmd.Flags().StringVarP(&intentName, "intent", "", "", "intent name")
}
