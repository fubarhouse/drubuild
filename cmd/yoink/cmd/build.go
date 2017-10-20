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
	"strings"
	"os"

	"github.com/spf13/cobra"

	"github.com/fubarhouse/golang-drush/make"
	composer2 "github.com/fubarhouse/golang-drush/composer"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "The build process for Yoink",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		var quitOut bool
		var quitMessages []string
		if makes == "" && composer == "" {
			quitMessages = append(quitMessages, "makes or composer values were not specified")
			quitOut = true
		}

		if quitOut == true {
			cmd.Usage()
			fmt.Println()
			for _, v := range quitMessages {
				fmt.Println(v)
			}
			os.Exit(1)
		}

		x := make.NewSite("", name, destination, alias, "", domain, "", "")
		y := make.NewmakeDB("127.0.0.1", "root", "root", 3306)
		x.DatabaseSet(y)
		if timestamp == 0 {
			x.TimeStampReset()
		} else {
			x.TimeStampSet(string(timestamp))
		}
		if workingCopy {
			x.WorkingCopy = true
		}

		if composer != "" {
			x.Composer = true
			composer2.InstallComposerCodebase(x.Name, x.TimeStampGet(), composer, x.Path)
		} else {
			x.Make = makes
			MakefilesFormatted := strings.Replace(makes, " ", "", -1)
			MakeFiles := strings.Split(MakefilesFormatted, ",")

			if rewriteSource != "" && rewriteDestination != "" {
				x.MakeFileRewriteSource = rewriteSource
				x.MakeFileRewriteDestination = rewriteDestination
			}
			x.ActionRebuildCodebase(MakeFiles)
		}
		x.InstallSiteRef()
		x.ActionInstall()
		x.SymReinstall()
		//x.VhostInstall()
		//x.AliasInstall()
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)
	// Required flags
	buildCmd.Flags().StringVarP(&name, "name", "n", "", "The human-readable name for this site")
	//buildCmd.Flags().StringVarP(&alias, "alias", "a", "", "The drush alias for this site")
	buildCmd.Flags().StringVarP(&destination, "destination", "p", "", "The path to where the site(s) exist.")
	buildCmd.Flags().StringVarP(&domain, "domain", "d", "", "The domain this site is to use")
	// Very important but not completely needed > 0 is needed though.
	buildCmd.Flags().StringVarP(&makes, "makes", "m", "", "A comma-separated list of make files for use")
	buildCmd.Flags().StringVarP(&composer, "composer", "c", "", "Path to the composer.json file.")
	// Optional flags
	buildCmd.Flags().Int64VarP(&timestamp, "timestamp", "t", 0, "Optional timestamp in the format YYYYMMDDHHMMSS")
	// Deprecated flags
	buildCmd.Flags().StringVarP(&rewriteSource, "rewrite-source", "x", "", "The rewrite string source")
	buildCmd.Flags().StringVarP(&rewriteDestination, "rewrite-destination", "y", "", "The rewrite string destination")
	buildCmd.Flags().BoolVarP(&workingCopy, "working-copy", "w", false, "Mark as a working-copy during the build process.")
	// Hide deprecated fields.
	buildCmd.Flags().MarkHidden("rewrite-source")
	buildCmd.Flags().MarkHidden("rewrite-destination")
	// Mark required flags.
	buildCmd.MarkFlagRequired("name")
	buildCmd.MarkFlagRequired("alias")
	buildCmd.MarkFlagRequired("destination")
	buildCmd.MarkFlagRequired("domain")
}
