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
	"os/exec"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Restore an archive-dump snapshot of a site to a local destination",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		d, err := exec.LookPath("drush")
		if err != nil {
			log.Fatal("Drush was not found in your $PATH")
		}

		db_user = viper.GetString("db_user")
		db_pass = viper.GetString("db_pass")
		db_host = viper.GetString("db_host")
		db_port = viper.GetInt("db_port")

		db := fmt.Sprintf("--db-url=mysql://%v:%v@%v:%v/%v", db_user, db_pass, db_host, db_port, name)

		if docroot != "" {
			if yes {
				c := exec.Command(d, destination, "archive-restore", source, "--destination="+docroot, "--overwrite", db)
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				c.Run()
				c.Wait()
			} else {
				c := exec.Command(d, destination, "archive-restore", source, "--destination="+docroot, db)
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				c.Run()
				c.Wait()
			}
		} else {
			if yes {
				c := exec.Command(d, destination, "archive-restore", source, "--overwrite", db)
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				c.Run()
				c.Wait()
			} else {
				c := exec.Command(d, destination, "archive-restore", source, db)
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				c.Run()
				c.Wait()
			}
		}
	},
}

func init() {

	RootCmd.AddCommand(importCmd)

	// Get the current working directory.
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Flags
	importCmd.Flags().StringVarP(&destination, "destination", "d", "", "Drush alias to use for operation")
	importCmd.Flags().StringVarP(&docroot, "docroot", "r", "", "Root of site (--destination flag passed to Drush)")
	importCmd.Flags().StringVarP(&source, "source", "s", dir, "Path to Drush archive-dump destination")
	importCmd.Flags().StringVarP(&name, "database-name", "b", "db_" + "", "The name of the destination database")
	importCmd.Flags().BoolVarP(&yes, "yes", "y", false, "Override the existing site available at 'docroot'?")
	// Mark flags as required
	importCmd.MarkFlagRequired("source")
	importCmd.MarkFlagRequired("destination")
}