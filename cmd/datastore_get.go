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

// datastoreGetCmd represents the get command
var datastoreGetCmd = &cobra.Command{
	Use:          "get",
	Short:        "show datastore details",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()
		dataClient, err := createDataClient(ctx, addr)
		if err != nil {
			return err
		}
		req := &sdcpb.GetDataStoreRequest{
			DatastoreName: datastoreName,
		}
		// fmt.Println("request:")
		// fmt.Println(prototext.Format(req))
		rsp, err := dataClient.GetDataStore(ctx, req)
		if err != nil {
			return err
		}
		switch format {
		case "":
			fmt.Println(prototext.Format(rsp))
		case "json":
			b, err := json.MarshalIndent(rsp, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(b))
		}

		return nil
	},
}

func init() {
	datastoreCmd.AddCommand(datastoreGetCmd)
}
