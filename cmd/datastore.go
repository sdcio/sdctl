// Copyright 2024 Nokia
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"time"

	sdcpb "github.com/sdcio/sdc-protos/sdcpb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var datastoreName string

// datastoreCmd represents the datastore command
var datastoreCmd = &cobra.Command{
	Use:   "datastore",
	Short: "manipulate datastores",
}

func init() {
	rootCmd.AddCommand(datastoreCmd)
	datastoreCmd.PersistentFlags().StringVarP(&datastoreName, "ds", "", "", "datastore (main) name")
}

func createDataClient(ctx context.Context, addr string) (sdcpb.DataServerClient, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	cc, err := grpc.DialContext(ctx, addr, grpcClientDialOpts()...)
	if err != nil {
		return nil, err
	}
	return sdcpb.NewDataServerClient(cc), nil
}
