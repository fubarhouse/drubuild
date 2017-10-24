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
	"os"

	"fmt"

	"github.com/spf13/cobra"
)

// syncCmd represents the backup command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise a set of templates in the provided destination path",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if destination == "" {
			cmd.Usage()
			fmt.Println("\ndestination is not set")
			os.Exit(1)
		}
		//fmt.Sprint(templateAlias, templateSitesPhp, templateVhostApache, templateVhostHttpd, templateVhostNginx)
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&webserver, "webserver", "w", "", "Name of webserver. Supports apache, httpd, nginx.")
	initCmd.Flags().StringVarP(&destination, "destination", "d", "", "Destination path to where the templates will be installed.")
}
