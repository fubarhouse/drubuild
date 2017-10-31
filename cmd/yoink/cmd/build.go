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
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"

	composer2 "github.com/fubarhouse/golang-drush/composer"
	"github.com/fubarhouse/golang-drush/make"
	"github.com/spf13/viper"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "The build process for Yoink",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if timestamp == 0 {
			timestamp, _ = strconv.ParseInt(time.Now().Format("20060102150405"), 0, 0)
			log.Printf("Timestamp not specified, using %v", timestamp)
		}

		x := make.NewSite("", name, destination, alias, webserver, domain, virtualhost_path, virtualhost_template)
		x.AliasTemplate = alias_template
		if alias == "" {
			x.Alias = domain
		}
		y := make.NewmakeDB(db_host, db_user, db_pass, db_port)
		x.DatabaseSet(y)

		if workingCopy {
			x.WorkingCopy = true
		}

		if composer != "" {
			x.Composer = true
			composer2.InstallComposerCodebase(x.Name, x.TimeStampGet(), composer, x.Path)
		} else if makes != "" {
			x.Make = makes
			MakefilesFormatted := strings.Replace(makes, " ", "", -1)
			MakeFiles := strings.Split(MakefilesFormatted, ",")

			if rewriteSource != "" && rewriteDestination != "" {
				x.MakeFileRewriteSource = rewriteSource
				x.MakeFileRewriteDestination = rewriteDestination
			}
			x.ActionRebuildCodebase(MakeFiles)
		} else {
			cmd.Usage()
			fmt.Println()
			log.Fatalln("makes and/or composer values were not specified")
			os.Exit(1)
		}

		x.InstallSiteRef(sites_php_template)
		x.SymReinstall()
		x.ActionInstall()

		if virtualhost_template != "" {
			if ok, err := os.Stat(virtualhost_template); err == nil {
				log.Infof("Found template %v for usage", ok.Name())
				x.Template = ok.Name()
			} else {
				log.Println("Could not find configured or default virtual host template.")
			}
		}
		x.VhostInstall()

		if alias_template != "" {
			if ok, err := os.Stat(alias_template); err != nil {
				log.Infof("Found template %v", ok.Name())
				x.AliasTemplate = ok.Name()
			} else {
				t := fmt.Sprintf("%v/src/github.com/fubarhouse/golang-drush/cmd/yoink/templates/alias.gotpl", os.Getenv("GOPATH"))
				log.Infof("Could not find template %v, using %v", ok.Name(), t)
			}
		}
		if alias == "" {
			x.Alias = domain
		}
		x.AliasInstall()
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)
	// Required flags
	buildCmd.Flags().StringVarP(&name, "name", "n", "", "The human-readable name for this site")
	buildCmd.Flags().StringVarP(&alias, "alias", "a", "", "The drush alias for this site")
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
	buildCmd.MarkFlagRequired("destination")
	buildCmd.MarkFlagRequired("domain")

	// Set configurables to defaults.
	viper.SetDefault("db_user", "root")
	viper.SetDefault("db_pass", "root")
	viper.SetDefault("db_host", "127.0.0.1")
	viper.SetDefault("db_port", 3306)
	viper.SetDefault("webserver", "nginx")
	viper.SetDefault("alias_template", "")
	viper.SetDefault("sites_php_template", "")
	viper.SetDefault("virtualhost_path", "/etc/nginx/sites-available")
	viper.SetDefault("virtualhost_template", "")

	// Database
	db_user = viper.GetString("db_user")
	db_pass = viper.GetString("db_pass")
	db_host = viper.GetString("db_host")
	db_port = viper.GetInt("db_port")

	// Sites.php template
	sites_php_template = viper.GetString("sites_php_template")

	// Alias template
	alias_template = viper.GetString("alias_template")

	// Virtual host
	webserver = viper.GetString("webserver")
	virtualhost_path = viper.GetString("virtualhost_path")
	virtualhost_template = viper.GetString("virtualhost_template")
}
