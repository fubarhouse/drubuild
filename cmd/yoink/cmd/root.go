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

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	// timestamp is an int64 which is representational of a date format in the
	// format of YYYYMMDDHHMMSS. An example of this is
	// 19700101000000.
	// if the timestamp is not set, this will change the file-system, drush
	// aliases and virtual hosts which are created or modified.
	//
	// these formats are used depending on this timestamp, and if composer
	// is being used, or drush make is used:
	//  - drush, timestamp
	//     /path/to/sites/mysite/mysite.timestamp/
	//  - drush, no timestamp
	//     /path/to/sites/mysite/mysite/
	//  - composer, timestamp
	//     /path/to/sites/mysite/mysite.timestamp/docroot
	//  - composer, no timestamp
	//     /path/to/sites/mysite/mysite/docroot
	//
	timestamp int64

	// name is the human-readable name for the target of this application.
	name string

	// source is the source alias or path to a source file to be used.
	// the specific action will be determined based on the command.
	source string

	// destination is the destination alias or path for the desired action.
	// it will be determined based upon the command in use.
	destination string

	// alias is the destination drush alias this site should be using.
	// in many places this will default to the domain name if not specified.
	alias string

	// domain is the destination domain to be used when setting up a new site
	domain string

	// makes is a comma-separated list of legacy make files to be used.
	// it will be automatically superseded by the use of the composer flag
	// however there is a lot of available deprecated functionality here.
	// Deprecated: use composer instead.
	makes string

	// when working with make files, you can tell the system to rewrite
	// a given module branch to change via a unique string inside the make
	// file(s). this represents the source of that change, what string is to be
	// replaced in the generated make file.
	//
	// Deprecated: use composer instead.
	rewriteSource string

	// when working with make files, you can tell the system to rewrite
	// a given module branch to change via a unique string inside the make
	// file(s). this represents the destination result of that change, what the
	// string is to be replaced to be in the generated make file.
	//
	// Deprecated: use composer instead.
	rewriteDestination string

	// composer represents the path to the composer file to be used.
	// it also represents a source file, in the event a composer.json file
	// does not exist at the destination path.
	// this flag will also supersede the necessity and the functionality
	// associated with legacy make files.
	composer string

	// webserver indicates which template should be used when using the
	// init command, which will generate a new template for you to use in the
	// destination folder, which will be at the location of the destination flag
	// or the current directory.
	//
	// This flag only accepts the values "apache", "httpd" and "nginx".
	webserver string

	// working-copy identifies if the build should leave the .git file-system
	// in-tact during the build. this would be useful when you are expecting
	// to send a file system to production, or for local development.
	// a working-copy is the local file-system including any development
	// file system data associated with each project/module.
	workingCopy bool

	// db_user is the string which represents the configured user account.
	// this user account should have permission to create databases, and
	// this user can be configured at $HOME/golang-drush.yml, and defaults
	// to 'root'.
	db_user string

	// db_pass is an unprotected string which represents the configured user
	// password. this user account should have permission to create
	// databases, and this password can be configured at
	// $HOME/golang-drush.yml, and defaults to 'root'.
	db_pass string

	// db_host is the string which represents the configured database host
	// this host path can be configured at $HOME/golang-drush.yml, and
	// defaults to '127.0.0.1'.
	db_host string

	// db_port is an integer which represents the configured database port
	// this port path can be configured at $HOME/golang-drush.yml, and
	// defaults to 3306.
	db_port int
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "yoink",
	Short: "A Drupal build system.",
	Long:  ``,
}

// Execute is the main function for the root command.
func Execute() {

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name "golang-drush" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName("golang-drush")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
