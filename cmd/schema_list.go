/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/prototext"

	sdcpb "github.com/sdcio/sdc-protos/sdcpb"
)

// schemaListCmd represents the list command
var schemaListCmd = &cobra.Command{
	Use:          "list",
	Short:        "list schemas",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()
		schemaClient, err := createSchemaClient(ctx, addr)
		if err != nil {
			return err
		}
		req := &sdcpb.ListSchemaRequest{}
		fmt.Println("request:")
		fmt.Println(prototext.Format(req))
		ctx, cancel2 := context.WithTimeout(cmd.Context(), timeout)
		defer cancel2()
		schemaList, err := schemaClient.ListSchema(ctx, req)
		if err != nil {
			return err
		}
		fmt.Println("response:")
		switch format {
		case "proto":
			fmt.Println(prototext.Format(schemaList))
		case "json":
			b, err := json.MarshalIndent(schemaList.GetSchema(), "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(b))
		}

		return nil
	},
}

func init() {
	schemaCmd.AddCommand(schemaListCmd)
}
