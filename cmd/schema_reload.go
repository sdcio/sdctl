/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	sdcpb "github.com/iptecharch/sdc-protos/sdcpb"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/prototext"
)

// schemaReloadCmd represents the reload command
var schemaReloadCmd = &cobra.Command{
	Use:          "reload",
	Short:        "reload schema",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()
		schemaClient, err := createSchemaClient(ctx, addr)
		if err != nil {
			return err
		}
		req := &sdcpb.ReloadSchemaRequest{
			Schema: &sdcpb.Schema{
				Name:    schemaName,
				Vendor:  schemaVendor,
				Version: schemaVersion,
			},
		}
		fmt.Println("request:")
		fmt.Println(prototext.Format(req))
		ctx, cancel2 := context.WithTimeout(cmd.Context(), timeout)
		defer cancel2()
		rsp, err := schemaClient.ReloadSchema(ctx, req)
		if err != nil {
			return err
		}
		fmt.Println("response:")
		fmt.Println(prototext.Format(rsp))
		return nil
	},
}

func init() {
	schemaCmd.AddCommand(schemaReloadCmd)
}
