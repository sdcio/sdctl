/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

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
	"fmt"
	"strings"

	sdcpb "github.com/iptecharch/sdc-protos/sdcpb"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/iptecharch/cache/pkg/cache"
	"github.com/iptecharch/cache/pkg/client"
	"github.com/iptecharch/cache/proto/cachepb"
)

var updatePaths []string
var deletePaths []string
var owner string
var priority int32

// modifyCmd represents the modify command
var modifyCmd = &cobra.Command{
	Use:   "modify",
	Short: "modify values in the cache",

	RunE: func(cmd *cobra.Command, _ []string) error {
		var store cache.Store
		switch storeName {
		case "config":
			store = cache.StoreConfig
		case "state":
			store = cache.StoreState
		case "intended":
			store = cache.StoreIntended
		case "intents":
			store = cache.StoreIntents
		default:
			return fmt.Errorf("unknown store name: %s", storeName)
		}

		c, err := client.New(cmd.Context(), &client.ClientConfig{
			Address:       addr,
			MaxReadStream: 1,
			Timeout:       timeout,
		})
		if err != nil {
			return err
		}

		dels := make([][]string, 0, len(deletePaths))

		for _, del := range deletePaths {
			dels = append(dels, strings.Split(del, ","))
		}

		upds := make([]*cachepb.Update, 0, len(updatePaths))
		for _, upd := range updatePaths {
			upd := strings.SplitN(upd, ":::", 3)
			if len(upd) != 3 {
				return fmt.Errorf("update %q is malformed", upd)
			}
			tv, err := toTypedValue(upd[1], upd[2])
			if err != nil {
				return err
			}
			b, err := proto.Marshal(tv)
			if err != nil {
				return err
			}
			upds = append(upds, &cachepb.Update{
				Path: strings.Split(upd[0], ","),
				Value: &anypb.Any{
					Value: b,
				},
			})
		}

		wo := &client.ClientOpts{
			Store:    store,
			Owner:    owner,
			Priority: priority,
		}
		err = c.Modify(cmd.Context(), cacheName, wo, dels, upds)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	cacheCmd.AddCommand(modifyCmd)
	modifyCmd.Flags().StringVarP(&cacheName, "name", "n", "", "cache name")
	modifyCmd.Flags().StringArrayVarP(&updatePaths, "update", "", []string{}, "path:::value to write")
	modifyCmd.Flags().StringArrayVarP(&deletePaths, "delete", "", []string{}, "paths to delete")
	modifyCmd.Flags().StringVarP(&storeName, "store", "s", "config", "cache store to modify")
	modifyCmd.Flags().StringVarP(&owner, "owner", "", "", "value owner for an intended store")
	modifyCmd.Flags().Int32VarP(&priority, "priority", "", 0, "owner priority for an intended store")
}

func toTypedValue(typ, val string) (*sdcpb.TypedValue, error) {
	// TODO: switch over typ
	return &sdcpb.TypedValue{
		Value: &sdcpb.TypedValue_StringVal{
			StringVal: val,
		},
	}, nil
}
