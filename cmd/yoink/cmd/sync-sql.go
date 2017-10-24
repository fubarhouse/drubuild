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
	"fmt"
)

// syncCmd represents the backup command
var syncCmd = &cobra.Command{
	Use:   "sync-sql",
	Short: "Execute drush sql-sync between two drush aliases" +
		"Note: Drush does not allow remote-remote syncing.",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if source == "" || destination == "" {
			cmd.Usage()
			fmt.Println("\nsource and/or destination are not set")
			os.Exit(1)
		}
		d, err := exec.LookPath("drush")
		if err != nil {
			log.Fatal("Drush was not found in your $PATH")
		}
		c := exec.Command(d, "sql-sync", source, destination, "--yes")
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Run()
		c.Wait()
	},
}

func init() {
	RootCmd.AddCommand(syncCmd)
	syncCmd.Flags().StringVarP(&source, "source", "s", "", "Drush alias to use as source")
	syncCmd.Flags().StringVarP(&destination, "destination", "d", "", "Drush alias to use as destination")
}
