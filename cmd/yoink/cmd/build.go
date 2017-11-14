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
	"github.com/mitchellh/go-homedir"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "The build process for Yoink",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		r := strings.Join([]string{home, "yoink"}, string(os.PathSeparator))

		db_user = viper.GetString("db_user")
		db_pass = viper.GetString("db_pass")
		db_host = viper.GetString("db_host")
		db_port = viper.GetInt("db_port")

		webserver = viper.GetString("webserver")

		alias_template = r + viper.GetString("alias_template")
		sites_php_template = r + viper.GetString("sites_php_template")
		virtualhost_path = r + viper.GetString("virtualhost_path")
		virtualhost_template = r + viper.GetString("virtualhost_template")

		if docroot == "" {
			log.Printf("docroot value is emptied, sub-folders will not be used.")
			timestamp = 0
		} else if timestamp == 0 {
			timestamp, _ = strconv.ParseInt(time.Now().Format("20060102150405"), 0, 0)
			log.Printf("Timestamp not specified, using %v", timestamp)
		}

		x := make.NewSite("", name, destination, alias, webserver, domain, virtualhost_path, virtualhost_template)
		x.AliasTemplate = alias_template
		x.Docroot = docroot
		if alias == "" {
			x.Alias = domain
		}
		log.Println(db_host, db_user, db_pass, db_port)
		y := make.NewmakeDB(db_host, db_user, db_pass, db_port)
		x.DatabaseSet(y)

		if workingCopy {
			x.WorkingCopy = true
		}

		if composer != "" {
			x.Composer = true
			composer2.InstallComposerCodebase(x.Name, x.TimeStampGet(), composer, x.Path, workingCopy)
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
		if timestamp != 0 {
			x.SymReinstall()
		}

		if !noInstall {
			x.ActionInstall()
		}

		if virtualhost_template != "" {
			if ok, err := os.Stat(virtualhost_template); err == nil {
				log.Infof("Found template %v for usage", ok.Name())
				x.Template = ok.Name()
			} else {
				log.Println("Could not find configured or default virtual host template.")
			}
		}

		if vhost {
			x.VhostInstall()
		}

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

		if drupal {
			x.AliasInstall(docroot)
		}
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
	buildCmd.Flags().StringVarP(&docroot, "docroot", "o", "docroot", "The folder to use for the built codebase.")
	buildCmd.Flags().BoolVarP(&drupal, "drupal", "r", true, "Mark the build process as a Drupal build.")
	buildCmd.Flags().BoolVarP(&noInstall, "no-install", "i", false, "Mark this build so that installation doesn't happen.")
	buildCmd.Flags().Int64VarP(&timestamp, "timestamp", "t", 0, "Optional timestamp in the format YYYYMMDDHHMMSS")
	buildCmd.Flags().BoolVarP(&vhost, "vhost", "v", true, "Include a virtual host as configured with this build.")
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
