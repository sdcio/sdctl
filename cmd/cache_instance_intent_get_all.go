// Copyright 2024 Nokia
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"strings"

	"github.com/sdcio/cache/pkg/client"
	"github.com/spf13/cobra"
)

var intent_get_all = &cobra.Command{
	Use:   "intents-getall",
	Short: "getall intent",

	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client.New(cmd.Context(), &client.ClientConfig{
			Address:       addr,
			MaxReadStream: 1,
			Timeout:       timeout,
		})
		if err != nil {
			return err
		}

		data, err := c.InstanceIntentsGetAll(cmd.Context(), cacheName, exceptNames)
		if err != nil {
			return err
		}
		result := make([]string, 0, len(data))
		for _, x := range data {
			result = append(result, x.String())
		}
		fmt.Println(strings.Join(result, "----------\n"))
		return nil
	},
}

func init() {
	intent_get_all.Flags().StringVarP(&cacheName, "name", "n", "", "cache name")
	intent_get_all.Flags().StringArrayVarP(&exceptNames, "except", "e", []string{}, "except names, intents that should be skipped")
	cacheCmd.AddCommand(intent_get_all)
}
