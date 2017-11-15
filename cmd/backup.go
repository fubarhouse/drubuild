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
	"os/exec"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Take a archive-dump snapshot of a local site",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		d, err := exec.LookPath("drush")
		if err != nil {
			log.Fatal("Drush was not found in your $PATH")
		}
		c := exec.Command(d, source, "archive-dump", "--destination=" + destination)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Run()
		c.Wait()
	},
}

func init() {
	RootCmd.AddCommand(backupCmd)
	// Flags
	backupCmd.Flags().StringVarP(&source, "source", "s", "", "Drush alias to use for operation")
	backupCmd.Flags().StringVarP(&destination, "destination", "d", "", "Path to Drush archive-dump destination")
	// Mark flags as required
	backupCmd.MarkFlagRequired("source")
	backupCmd.MarkFlagRequired("destination")
}
