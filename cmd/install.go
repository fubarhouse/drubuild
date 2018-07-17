// Copyright © 2018 Karl Hepworth <karl.hepworth@gmail.com>
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
	"bufio"
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/fubarhouse/drubuild/util/drush"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a Drupal website",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		var sites_php_template_data = `
<?php

/**
 * @file
 * Configuration file for Drupal's multi-site directory aliasing feature.
 */

 $sites['default'] = '{{ .Name }}';
 $sites['{{ .Domain }}'] = '{{ .Name }}';

?>
`

		// Specify directories:
		if source != "" {
			sitedir := strings.Join([]string{destination, "sites", source}, string(os.PathSeparator))
			log.Infoln("Creating subdirectory", sitedir)
			dirErr := os.MkdirAll(sitedir, 0755)
			if dirErr != nil {
				log.Errorln("Unable to create directory", sitedir, dirErr)
			} else {
				log.Infoln("Created directory", sitedir)
			}
		}

		// sites.php
		{
			sitesdir := strings.Join([]string{destination, "sites"}, string(os.PathSeparator))
			filename := strings.Join([]string{sitesdir, "sites.php"}, string(os.PathSeparator))
			log.Infoln("Creating file", filename)

			if file, err := os.Create(filename); err != nil {
				log.Errorln("Could not create file", filename, err)
				os.Exit(1)
			} else {
				defer file.Close()
				w := bufio.NewWriter(file)
				sites_php_template_data = strings.Replace(sites_php_template_data, "{{ .Domain }}", domain, -1)
				sites_php_template_data = strings.Replace(sites_php_template_data, "{{ .Name }}", source, -1)
				if _, err = fmt.Fprintf(w, sites_php_template_data); err != nil {
					log.Errorln("Could not template file.", err)
				}
				log.Infoln("Drush alias successfully installed")
				w.Flush()
			}
		}

		// drush site-install:
		{
			drush.Run([]string{"site-install", "--root=" + destination, "--yes", "--sites-subdir=" + name, fmt.Sprintf("--db-url=mysql://%v:%v@%v:%v/%v", db_user, db_pass, db_host, db_port, name)})
		}
	},
}

func init() {
	RootCmd.AddCommand(installCmd)

	// Get $PWD
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Parameters/Flags
	installCmd.Flags().StringVarP(&name, "name", "n", "default", "The database name")
	installCmd.Flags().StringVarP(&source, "sites-subdir", "f", "default", "The sites directory")
	installCmd.Flags().StringVarP(&name, "domain", "d", "default", "The sites directory")
	installCmd.Flags().StringVarP(&db_host, "host", "t", "", "The database host")
	installCmd.Flags().StringVarP(&db_pass, "password", "s", "", "The database password")
	installCmd.Flags().StringVarP(&db_user, "user", "u", "", "The database user name")
	installCmd.Flags().StringVarP(&destination, "path", "p", pwd, "The directory of the Drupal codebase")
	installCmd.Flags().IntVarP(&db_port, "port", "o", 3306, "The database port")

	// All inputs are mandatory because we use a multisite setup.
	installCmd.MarkFlagRequired("name")
	installCmd.MarkFlagRequired("sites-subdir")
	installCmd.MarkFlagRequired("domain")
	installCmd.MarkFlagRequired("host")
	installCmd.MarkFlagRequired("password")
	installCmd.MarkFlagRequired("user")
	installCmd.MarkFlagRequired("path")
	installCmd.MarkFlagRequired("port")
}
