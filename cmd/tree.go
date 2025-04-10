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
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var flat bool
var details bool

// treeCmd represents the tree command
var treeCmd = &cobra.Command{
	Use:   "tree",
	Short: "print the commands tree",
	RunE: func(cmd *cobra.Command, args []string) error {
		if flat {
			treeFlat(rootCmd, "")
			return nil
		}
		tree(rootCmd, "")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(treeCmd)

	treeCmd.Flags().BoolVar(&flat, "flat", false, "print flat commands tree")
	treeCmd.Flags().BoolVar(&details, "details", false, "print commands flags")
}

func tree(c *cobra.Command, indent string) error {
	fmt.Printf("%s", c.Use)
	if !c.HasSubCommands() {
		if c.HasLocalFlags() && details {
			sections := make([]string, 0)
			c.LocalFlags().VisitAll(func(flag *pflag.Flag) {
				flagSection := ""
				if flag.Shorthand != "" && flag.ShorthandDeprecated == "" {
					flagSection = fmt.Sprintf("[-%s | --%s]", flag.Shorthand, flag.Name)
				} else {
					flagSection = fmt.Sprintf("[--%s]", flag.Name)
				}
				sections = append(sections, flagSection)
			})
			fmt.Printf(" %s", strings.Join(sections, " "))
		}
	}
	fmt.Printf("\n")
	subCmds := c.Commands()
	numSubCommands := len(subCmds)
	for i, subC := range subCmds {
		add := " │   "
		if i == numSubCommands-1 {
			fmt.Printf("%s └─── ", indent)
			add = "     "
		} else {
			fmt.Printf("%s ├─── ", indent)
		}

		err := tree(subC, indent+add)
		if err != nil {
			return err
		}
	}
	return nil
}

func treeFlat(c *cobra.Command, prefix string) {
	prefix += " " + c.Use
	fmt.Println(prefix)
	for _, subC := range c.Commands() {
		treeFlat(subC, prefix)
	}
}
