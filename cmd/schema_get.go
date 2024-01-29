/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/iptecharch/schema-server/pkg/utils"
	sdcpb "github.com/iptecharch/sdc-protos/sdcpb"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/prototext"
)

var xpath string
var withDesc bool
var all bool

// schemaGetCmd represents the get command
var schemaGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get schema",
	// SilenceErrors: true,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		p, err := utils.ParsePath(xpath)
		if err != nil {
			return err
		}
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()
		schemaClient, err := createSchemaClient(ctx, addr)
		if err != nil {
			return err
		}
		req := &sdcpb.GetSchemaRequest{
			Path: p,
			Schema: &sdcpb.Schema{
				Name:    schemaName,
				Vendor:  schemaVendor,
				Version: schemaVersion,
			},
			WithDescription: withDesc,
		}
		fmt.Fprintln(os.Stderr, "request:")
		fmt.Println(prototext.Format(req))
		ctx, cancel2 := context.WithTimeout(cmd.Context(), timeout)
		defer cancel2()
		if all {
			return handleGetSchemaElems(ctx, schemaClient, req)
		}
		rsp, err := schemaClient.GetSchema(ctx, req)
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr, "response:")
		if format == "json" {
			b, err := json.MarshalIndent(rsp, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(b))
			return nil
		}
		fmt.Println(prototext.Format(rsp))
		return nil
	},
}

func init() {
	schemaCmd.AddCommand(schemaGetCmd)

	schemaGetCmd.PersistentFlags().StringVarP(&xpath, "path", "p", "", "xpath")
	schemaGetCmd.PersistentFlags().BoolVarP(&all, "all", "", false, "return all path elems schemas")
	schemaGetCmd.PersistentFlags().BoolVarP(&withDesc, "with-desc", "", false, "include YANG entries descriptions")
}

func handleGetSchemaElems(ctx context.Context, scc sdcpb.SchemaServerClient, req *sdcpb.GetSchemaRequest) error {
	stream, err := scc.GetSchemaElements(ctx, req)
	if err != nil {
		return err
	}
	for {
		rsp, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				return nil
			}
			return err
		}
		fmt.Println("response:")
		fmt.Println(prototext.Format(rsp))
	}
}
