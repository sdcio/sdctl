/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	sdcpb "github.com/sdcio/sdc-protos/sdcpb"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
)

var target string
var syncFile string

// var owner string
// var priority int32

// datastoreCreateCmd represents the create command
var datastoreCreateCmd = &cobra.Command{
	Use:          "create",
	Short:        "create datastore",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()
		dataClient, err := createDataClient(ctx, addr)
		if err != nil {
			return err
		}

		req := &sdcpb.CreateDataStoreRequest{
			DatastoreName: datastoreName,
		}

		var tg *sdcpb.Target
		if target != "" {
			f, err := os.ReadFile(target)
			if err != nil {
				return err
			}
			tg = &sdcpb.Target{}
			err = protojson.Unmarshal(f, tg)
			if err != nil {
				return err
			}
			req.Target = tg
		}

		if syncFile != "" {
			f, err := os.ReadFile(syncFile)
			if err != nil {
				return err
			}
			sync := &sdcpb.Sync{}
			err = protojson.Unmarshal(f, sync)
			if err != nil {
				return err
			}
			req.Sync = sync
		}

		req.Schema = &sdcpb.Schema{
			Name:    schemaName,
			Vendor:  schemaVendor,
			Version: schemaVersion,
		}

		fmt.Println("request:")
		fmt.Println(prototext.Format(req))
		rsp, err := dataClient.CreateDataStore(ctx, req)
		if err != nil {
			return err
		}
		fmt.Println("response:")
		fmt.Println(prototext.Format(rsp))
		return nil
	},
}

func init() {
	datastoreCmd.AddCommand(datastoreCreateCmd)

	datastoreCreateCmd.Flags().StringVarP(&target, "target", "", "", "target definition file")
	datastoreCreateCmd.Flags().StringVarP(&syncFile, "sync", "", "", "target sync definition file")
}
