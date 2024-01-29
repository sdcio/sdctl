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
	"github.com/spf13/cobra"

	"github.com/iptecharch/cache/pkg/client"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a cache instance",

	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client.New(cmd.Context(), &client.ClientConfig{
			Address:       addr,
			MaxReadStream: 1,
			Timeout:       timeout,
		})
		if err != nil {
			return err
		}
		err = c.Delete(cmd.Context(), cacheName)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	cacheCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVarP(&cacheName, "name", "n", "", "cache name")
}
