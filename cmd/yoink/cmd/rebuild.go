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
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	composer2 "github.com/fubarhouse/golang-drush/composer"
	"github.com/fubarhouse/golang-drush/make"
)

// rebuildCmd represents the rebuild command
var rebuildCmd = &cobra.Command{
	Use:   "rebuild",
	Short: "Rebuild a site",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		x := make.Site{}
		x.TimeStampSet("")
		x.Name = name
		x.Path = destination
		x.WorkingCopy = workingCopy
		x.Composer = false

		if composer != "" {
			x.Composer = true
			composer2.InstallComposerCodebase(x.Name, x.TimeStampGet(), composer, x.Path, x.WorkingCopy)
		} else {

			if makes == "" {
				fmt.Println("Error: Required flag(s) \"makes\" have/has not been set")
				cmd.Usage()
				fmt.Println("\nRequired flag(s) \"makes\" have/has not been set")
				os.Exit(1)
			}

			x.Make = makes
			MakefilesFormatted := strings.Replace(makes, " ", "", -1)
			MakeFiles := strings.Split(MakefilesFormatted, ",")
			if rewriteSource != "" && rewriteDestination != "" {
				x.MakeFileRewriteSource = rewriteSource
				x.MakeFileRewriteDestination = rewriteDestination
			}
			x.ActionRebuildCodebase(MakeFiles)
		}
	},
}

func init() {
	RootCmd.AddCommand(rebuildCmd)
	rebuildCmd.Flags().StringVarP(&name, "name", "n", "", "The human-readable name for this site")
	rebuildCmd.Flags().StringVarP(&destination, "destination", "p", "", "The path to where the site(s) exist.")
	rebuildCmd.Flags().StringVarP(&composer, "composer", "c", "", "Path to the composer.json file.")
	rebuildCmd.Flags().BoolVarP(&workingCopy, "working-copy", "w", false, "Mark as a working-copy during the build process.")
	// Deprecated flags
	rebuildCmd.Flags().StringVarP(&makes, "makes", "m", "", "A comma-separated list of make files for use")
	rebuildCmd.Flags().StringVarP(&rewriteSource, "rewrite-source", "x", "", "The rewrite string source")
	rebuildCmd.Flags().StringVarP(&rewriteDestination, "rewrite-destination", "y", "", "The rewrite string destination")
	// Markers
	rebuildCmd.Flags().MarkHidden("rewrite-source")
	rebuildCmd.Flags().MarkHidden("rewrite-destination")
	rebuildCmd.MarkFlagRequired("name")
	rebuildCmd.MarkFlagRequired("destination")
}
