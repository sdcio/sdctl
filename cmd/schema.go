/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"

	sdcpb "github.com/iptecharch/sdc-protos/sdcpb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var schemaName string
var schemaVendor string
var schemaVersion string

// schemaCmd represents the schema command
var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "query/change schema(s)",
}

func init() {
	rootCmd.AddCommand(schemaCmd)
	schemaCmd.PersistentFlags().StringVar(&schemaName, "name", "", "schema name")
	schemaCmd.PersistentFlags().StringVar(&schemaVendor, "vendor", "", "schema vendor")
	schemaCmd.PersistentFlags().StringVar(&schemaVersion, "version", "", "schema version")

}

func grpcClientDialOpts() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxRcvMsg)),
	}
}
func createSchemaClient(ctx context.Context, addr string) (sdcpb.SchemaServerClient, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	cc, err := grpc.DialContext(ctx, addr, grpcClientDialOpts()...)
	if err != nil {
		return nil, err
	}
	return sdcpb.NewSchemaServerClient(cc), nil
}
