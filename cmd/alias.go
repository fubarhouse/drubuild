// Copyright © 2018 Karl Hepworth Karl.Hepworth@gmail.com
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
		log "github.com/Sirupsen/logrus"
	"os"
	"fmt"
	"os/user"
	"strings"
	"bufio"
)

// AliasTemplate is a basic Drush alias template
var AliasTemplate = `
<?php
	$aliases['{{ .Alias }}'] = array(
		'root' => '{{ .Root }}',
		'uri' => '{{ .Domain }}',
		'path-aliases' => array(
			'%files' => 'sites/{{ .Name }}/files',
			'%private' => 'sites/{{ .Name }}/private',
		),
	);
?>
`

func getDrushPath() string {
	usr, _ := user.Current()
	filedir := usr.HomeDir
	filename := domain + ".alias.drushrc.php"
	return strings.Join([]string{
	filedir,
	".drush",
	filename,
	}, string(os.PathSeparator))
}

// aliasCmd represents the alias command
var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "An isolated command for alias operations",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {

		switch args[0] {
		case "install":

			fullpath := getDrushPath()
			if _, err := os.Stat(fullpath); err == nil {
				log.Errorln("Alias file already exists, please uninstall it first.")
				os.Exit(1)
			}

			if file, err := os.Create(fullpath); err != nil {
				log.Errorln("Could not create file", fullpath, err)
				os.Exit(1)
			} else {
				if (name == "" || alias == "" || destination == "") {
					log.Errorln("Flags name, alias and destination are required for installation, run with --help for more information")
					os.Exit(1)
				}
				defer file.Close()
				w := bufio.NewWriter(file)
				AliasTemplate = strings.Replace(AliasTemplate, "{{ .Alias }}", alias, -1);
				AliasTemplate = strings.Replace(AliasTemplate, "{{ .Domain }}", domain, -1);
				AliasTemplate = strings.Replace(AliasTemplate, "{{ .Name }}", name, -1);
				AliasTemplate = strings.Replace(AliasTemplate, "{{ .Root }}", destination, -1);
				if _, err = fmt.Fprintf(w, AliasTemplate); err != nil {
					log.Errorln("Could not template file.", err)
				}
				log.Infoln("Drush alias successfully installed")
				w.Flush()
			}
			break

		case "uninstall":
			fullpath := getDrushPath()
			_, statErr := os.Stat(fullpath)
			if statErr == nil {
				err := os.Remove(fullpath)
				if err != nil {
					log.Warnln("Could not remove alias file", fullpath)
				} else {
					log.Infoln("Removed alias file", fullpath)
				}
			} else {
				log.Errorln("Alias file was not found, please install it first.")
				os.Exit(1)
			}
			break

		default:
			break;
		}
	},
}

func init() {
	RootCmd.AddCommand(aliasCmd)

	// Declare flags.
	aliasCmd.Flags().StringVarP(&name, "directory", "d", "", "The directory name under /sites which contains settings.php")
	aliasCmd.Flags().StringVarP(&alias, "alias", "a", "", "The drush alias for this site")
	aliasCmd.Flags().StringVarP(&destination, "root", "r", "", "The path to the root of the site")
	aliasCmd.Flags().StringVarP(&domain, "url", "u", "", "The domain of the site not including protocol or trailing slashes")

	// Declare required flags.
	aliasCmd.MarkFlagRequired("url")
}
