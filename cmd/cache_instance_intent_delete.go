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
	"github.com/sdcio/cache/pkg/client"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var intent_delete = &cobra.Command{
	Use:   "intent-delete",
	Short: "delete intent",

	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client.New(cmd.Context(), &client.ClientConfig{
			Address:       addr,
			MaxReadStream: 1,
			Timeout:       timeout,
		})
		if err != nil {
			return err
		}

		err = c.InstanceIntentDelete(cmd.Context(), cacheName, intentName)
		return err
	},
}

func init() {
	intent_delete.Flags().StringVarP(&cacheName, "name", "n", "", "cache name")
	intent_delete.Flags().StringVarP(&intentName, "intent", "i", "", "intent name")
	cacheCmd.AddCommand(intent_delete)
}
