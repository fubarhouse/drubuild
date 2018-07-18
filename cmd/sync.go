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
	"log"
	"os"
		"fmt"

	"github.com/spf13/cobra"
	"github.com/fubarhouse/drubuild/util/drush"
)

// syncCmd represents the backup command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Execute drush sql-sync or rsync between two drush aliases",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if syncFiles {
			{
				if name == "" {
					log.Fatal("rsync target was not provided, please specify target with --target.")
				}
				fsPu := fmt.Sprintf("%v:%%%v", source, name)
				fdPu := fmt.Sprintf("%v:%%%v", destination, name)
				if yes {
					drush.Run([]string{"--yes", "rsync", fsPu, fdPu, "--exclude-other-sites", "--exclude-conf"})
				} else {
					drush.Run([]string{"rsync", fsPu, fdPu, "--exclude-other-sites", "--exclude-conf"})
				}
			}
		}
		if syncDatabase {
			if yes {
				drush.Run([]string{"--yes", "sql-sync", source, destination})
			} else {
				drush.Run([]string{"sql-sync", source, destination})
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(syncCmd)

	// Get the current working directory.
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	syncCmd.Flags().StringVarP(&source, "source", "s", "", "Drush alias to use as source")
	syncCmd.Flags().StringVarP(&destination, "destination", "d", dir, "Drush alias to use as destination")
	syncCmd.Flags().StringVarP(&name, "target", "t", "", "The name of the path alias in the drush alias. ie files, public, private, temp")
	syncCmd.Flags().BoolVarP(&syncDatabase, "database", "b", false, "Flag database for sync action.")
	syncCmd.Flags().BoolVarP(&syncFiles, "files", "f", false, "Flag files for sync action.")
	syncCmd.Flags().BoolVarP(&yes, "yes", "y", false, "Use command with --yes")

	syncCmd.MarkFlagRequired("source")
	syncCmd.MarkFlagRequired("destination")
}
