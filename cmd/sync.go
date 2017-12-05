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
	"log"
	"os"
	"os/exec"

	"fmt"

	"github.com/spf13/cobra"
)

// syncCmd represents the backup command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Execute drush sql-sync or rsync between two drush aliases",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		d, err := exec.LookPath("drush")
		if err != nil {
			log.Fatal("Drush was not found in your $PATH")
		}
		if syncFiles {
			{
				fsPu := fmt.Sprintf("@%v:%%public", source)
				fdPu := fmt.Sprintf("@%v:%%public", destination)
				c := exec.Command(d, "rsync", fsPu, fdPu, "--yes", "--exclude-other-sites", "--exclude-conf")
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				c.Run()
				c.Wait()
			}
			{
				fsPr := fmt.Sprintf("@%v:%%private", source)
				fdPr := fmt.Sprintf("@%v:%%private", destination)
				c := exec.Command(d, "rsync", fsPr, fdPr, "--yes", "--exclude-other-sites", "--exclude-conf")
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				c.Run()
				c.Wait()
			}
		}
		if syncDatabase {
			c := exec.Command(d, "sql-sync", source, destination, "--yes")
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Run()
			c.Wait()
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
	syncCmd.Flags().BoolVarP(&syncDatabase, "database", "b", false, "Flag database for sync action.")
	syncCmd.Flags().BoolVarP(&syncFiles, "files", "f", false, "Flag files for sync action.")

	syncCmd.MarkFlagRequired("source")
	syncCmd.MarkFlagRequired("destination")
}
