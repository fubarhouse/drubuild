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
	"log"
	"net/mail"

	"os"
	"strings"

	"fmt"

	alias2 "github.com/fubarhouse/drubuild/util/alias"
	"github.com/fubarhouse/drubuild/util/user"
	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User management tasks",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		if !strings.Contains(pattern, "%v") {
			log.Fatalf("Specified pattern '%v' does not include alias modifier '%%v'.", pattern)
			os.Exit(1)
		}

		Email, err := mail.ParseAddress(user_email)
		if err != nil {
			log.Fatal("Input email address could not be parsed.")
			log.Fatal(err)
		}

		switch args[0] {
		case "block":
			for _, Alias := range strings.Split(aliases, ",") {
				Alias = strings.Trim(Alias, " ")
				Alias = strings.Replace(pattern, "%v", Alias, 1)
				log.Printf("Beginning to work with alias %v.", Alias)
				UserGroup := user.NewDrupalUserGroup()
				UserGroup.Populate(Alias)
				User := UserGroup.GetUser(user_name)
				if User.Name == user_name {
					if User.State == 0 {
						User.State = 1
						user.Block(alias, User.Name)
					} else {
						log.Printf("User '%v' is already blocked on %v\n", User.Name, Alias)
					}
				} else {
					log.Printf("User '%v' was not found on %v\n", User.Name, Alias)
				}
			}
			break;

		case "create":
			for _, Alias := range strings.Split(aliases, ",") {
				Alias = strings.Trim(Alias, " ")
				Alias = strings.Replace(pattern, "%v", Alias, 1)
				log.Printf("Beginning to work with alias %v.", Alias)
				UserGroup := user.NewDrupalUserGroup()
				UserGroup.Populate(Alias)
				User := UserGroup.GetUser(user_name)
				User.Alias = Alias
				User.Name = user_name
				User.Email = user_email
				user.Create(alias, user_name, user_email, user_password)
				if user_role != "" {
					if !user.HasRole(User.Alias, User.Name, user_role) {
						User.Roles = append(User.Roles, user_role)
					}
					user.RolesAdd(alias, user_name, User.Roles)
				}
			}
			break;

		case "delete":
			for _, Alias := range strings.Split(aliases, ",") {
				Alias = strings.Trim(Alias, " ")
				Alias = strings.Replace(pattern, "%v", Alias, 1)
				log.Printf("Beginning to work with alias %v.", Alias)
				UserGroup := user.NewDrupalUserGroup()
				UserGroup.Populate(Alias)
				User := UserGroup.GetUser(user_name)
				if User.Name == user_name {
					user.Delete(alias, name)
				} else {
					log.Printf("User '%v' was not found on %v", user_name, Alias)
				}
			}
			break;

		case "unblock":
			for _, Alias := range strings.Split(aliases, ",") {
				Alias = strings.Trim(Alias, " ")
				Alias = strings.Replace(pattern, "%v", Alias, 1)
				log.Printf("Beginning to work with alias %v.", Alias)
				UserGroup := user.NewDrupalUserGroup()
				UserGroup.Populate(Alias)
				User := UserGroup.GetUser(user_name)
				if User.Name == user_name {
					if User.State == 1 {
						User.State = 0
						user.Unblock(alias, User.Name)
					} else {
						log.Printf("User '%v' is already unblocked on %v", User.Name, Alias)
					}
				} else {
					log.Printf("User '%v' was not found on %v", user_name, Alias)
				}
			}
			break;

		case "verify":
			for _, ThisAlias := range strings.Split(aliases, ",") {
				ThisAlias = strings.Trim(ThisAlias, " ")
				ThisAlias = strings.Replace(pattern, "%v", ThisAlias, 1)
				log.Printf("Beginning to work with alias %v.", ThisAlias)
				UserGroup := user.NewDrupalUserGroup()
				Alias := alias2.NewAlias(ThisAlias, "", ThisAlias)
				if Alias.GetStatus() {
					UserGroup.Populate(Alias.GetName())
					User := UserGroup.GetUser(user_name)
					if User.Alias == "" {
						User.Alias = Alias.GetName()
					}
					if User.Name == "" {
						User.Name = user_name
					}
					User.Email = Email.Address

					if user_create {
						user.Create(alias, user_name, user_email, user_password)
					} else {
						log.Println("User creation not set, skipping.")
					}

					if user_block {
						User.State = 0
						user.Block(alias, User.Name)
					} else if user_unblock {
						User.State = 1
						user.Unblock(alias, User.Name)
					} else {
						log.Println("Block/Unblock action not set, skipping.")
					}

					if user_password != "" {
						user.SetPassword(alias, User.Name, user_password)
					} else {
						log.Println("Password not set, skipping.")
					}

					user.EmailChange(alias, user_name, user_email)

					if user_role != "" {
						User.Roles = append(User.Roles, user_role)
						user.RolesAdd(alias, user_name, User.Roles)
					} else {
						log.Println("Role not set, skipping.")
					}

				} else {
					log.Printf("Could not find alias %v\n", Alias.GetName())
				}
			}
			break;

		default:
			cmd.Usage()
			fmt.Println()
			log.Fatalln("no action specified.")
			os.Exit(1)
		}

	},
}

func init() {

	RootCmd.AddCommand(userCmd)

	userCmd.Flags().StringVarP(&user_name, "name", "n", "", "User name")
	userCmd.Flags().StringVarP(&user_email, "email", "e", "", "User email")
	userCmd.Flags().StringVarP(&user_password, "password", "s", "", "User password")
	userCmd.Flags().StringVarP(&user_role, "role", "r", "", "User role (case sensitive)")

	userCmd.Flags().StringVarP(&aliases, "aliases", "a", "", "Comma-separated list of drush aliases")
	userCmd.Flags().StringVarP(&pattern, "pattern", "p", "%v", "Pattern to match against drush aliases, where token is '%v'")

	userCmd.MarkFlagRequired("name")
	userCmd.MarkFlagRequired("aliases")

}
