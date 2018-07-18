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
	"fmt"
		"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/fubarhouse/drubuild/util/drush"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Restore an archive-dump snapshot of a site to a local destination",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		db_user = viper.GetString("db_user")
		db_pass = viper.GetString("db_pass")
		db_host = viper.GetString("db_host")
		db_port = viper.GetInt("db_port")

		db := fmt.Sprintf("--db-url=mysql://%v:%v@%v:%v/%v", db_user, db_pass, db_host, db_port, name)

		if docroot != "" {
			if yes {
				drush.Run([]string{destination, "archive-restore", source, "--destination="+docroot, "--overwrite", db})
			} else {
				drush.Run([]string{destination, "archive-restore", source, "--destination="+docroot, db})
			}
		} else {
			if yes {
				drush.Run([]string{destination, "archive-restore", source, "--overwrite", db})
			} else {
				drush.Run([]string{destination, "archive-restore", source, db})
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