// Copyright © 2017 Karl Hepworth <Karl.Hepworth@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"strings"

	"log"

		"github.com/spf13/cobra"
	"github.com/fubarhouse/drubuild/util/drush"
)

// runnerCmd represents the runner command
var runnerCmd = &cobra.Command{
	Use:   "runner",
	Short: "Runs a series of drush commands consecutively.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		for _, Alias := range strings.Split(aliases, ",") {
			log.Println("Actioning", Alias)
			Alias = strings.Replace(pattern, "%v", Alias, 1)
			Alias = strings.Trim(Alias, " ")
			for _, Command := range strings.Split(commands, ",") {
				Command = strings.Trim(Command, " ")
				drush.Run([]string{Alias, Command})
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(runnerCmd)
	runnerCmd.Flags().StringVarP(&aliases, "aliases", "a", "", "Comma-separated list of drush aliases")
	runnerCmd.Flags().StringVarP(&commands, "commands", "c", "", "Comma-separated list of commands to run")
	runnerCmd.Flags().StringVarP(&pattern, "pattern", "p", "%v", "Pattern to match against drush aliases, where token is '%v'")

	runnerCmd.MarkFlagRequired("aliases")
	runnerCmd.MarkFlagRequired("commands")
}
