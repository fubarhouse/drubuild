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
	"github.com/spf13/cobra"

	"github.com/fubarhouse/golang-drush/make"
	"github.com/spf13/viper"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Remove all traces of an installed site.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		db_user = viper.GetString("db_user")
		db_pass = viper.GetString("db_pass")
		db_host = viper.GetString("db_host")
		db_port = viper.GetInt("db_port")

		webserver = viper.GetString("webserver")

		alias_template = viper.GetString("alias_template")
		sites_php_template = viper.GetString("sites_php_template")
		virtualhost_path = viper.GetString("virtualhost_path")
		virtualhost_template = viper.GetString("virtualhost_template")

		if alias == "" {
			alias = domain
		}
		x := make.NewSite("none", name, destination, alias, "", domain, "", "")
		y := make.NewmakeDB(db_host, db_user, db_pass, db_port)
		x.DatabaseSet(y)
		x.ActionDestroy()
	},
}

func init() {
	RootCmd.AddCommand(destroyCmd)
	// Flags
	destroyCmd.Flags().StringVarP(&name, "name", "n", "", "The human-readable name for this site")
	destroyCmd.Flags().StringVarP(&alias, "alias", "a", "", "The drush alias for this site")
	destroyCmd.Flags().StringVarP(&destination, "destination", "p", "", "The path to where the site(s) exist.")
	destroyCmd.Flags().StringVarP(&domain, "domain", "d", "", "The domain this site is to use")
	// Mark as required
	destroyCmd.MarkFlagRequired("name")
	destroyCmd.MarkFlagRequired("destination")
	destroyCmd.MarkFlagRequired("domain")

	db_user = viper.GetString("db_user")
	db_pass = viper.GetString("db_pass")
	db_host = viper.GetString("db_host")
	db_port = viper.GetInt("db_port")
}
