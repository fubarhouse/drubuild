// Copyright © 2017 Karl Hepworth karl.hepworth@gmail.com
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

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
