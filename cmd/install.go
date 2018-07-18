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
				log.Infoln("sites.php successfully installed")
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
	installCmd.Flags().StringVarP(&name, "name", "n", "", "The database name")
	installCmd.Flags().StringVarP(&source, "sites-subdir", "f", "default", "The sites directory")
	installCmd.Flags().StringVarP(&domain, "domain", "d", "", "The sites directory")
	installCmd.Flags().StringVarP(&db_host, "host", "t", "127.0.0.1", "The database host")
	installCmd.Flags().StringVarP(&db_pass, "password", "s", "", "The database password")
	installCmd.Flags().StringVarP(&db_user, "user", "u", "", "The database user name")
	installCmd.Flags().StringVarP(&destination, "path", "p", pwd, "The directory of the Drupal codebase")
	installCmd.Flags().IntVarP(&db_port, "port", "o", 3306, "The database port")

	// All inputs are mandatory because we use a multisite setup.
	installCmd.MarkFlagRequired("name")
	installCmd.MarkFlagRequired("domain")
	installCmd.MarkFlagRequired("password")
	installCmd.MarkFlagRequired("user")
}
