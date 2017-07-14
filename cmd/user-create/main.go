//
// user-create will create Drupal accounts on site alises based upon input values.
//
package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/golang-drush/user"
	"net/mail"
	"os"
	"strings"
)

func main() {
	var strAliases = flag.String("aliases", "", "comma-separated list of aliases for action")
	var strPattern = flag.String("pattern", "%v", "A modifier which allows rewriting of aliases replacing '%v' in the pattern with the alias.")
	var strUser = flag.String("user", "", "User name for validation, example 'Firstname Sirname'")
	var strRole = flag.String("role", "", "The role name to add to the user, if the user isn't a part of that role.")
	var strPassword = flag.String("password", "", "Password to reset to - what the password for each account should be.")
	var strEmail = flag.String("email", "", "Email to reset to - what the email for each account should be.")
	var strState = flag.Bool("active", true, "The active state of the user account for each account - should the account be active?")
	var boolCreate = flag.Bool("create-user", true, "If required, create the user if it doesn't exist on each alias.")
	flag.Parse()

	// Remove double spaces.
	*strAliases = strings.Replace(*strAliases, "  ", " ", -1)

	if !strings.Contains(*strPattern, "%v") {
		log.Errorln("Specified pattern does not include alias modifier.")
		flag.Usage()
		os.Exit(1)
	}

	if *strUser == "" {
		log.Errorln("User is not specified.")
		flag.Usage()
		os.Exit(1)
	}

	if *strAliases == "" {
		flag.Usage()
		log.Errorln("Aliases are not specified.")
		flag.Usage()
		os.Exit(1)
	}

	Email, err := mail.ParseAddress(*strEmail)
	if err != nil {
		log.Errorln("Input email address could not be parsed.")
		log.Fatal(err)
	}

	if *strAliases != "" && *strUser != "" {
		for _, Alias := range strings.Split(*strAliases, ",") {
			Alias = strings.Trim(Alias, " ")
			Alias = strings.Replace(*strPattern, "%v", Alias, 1)
			UserGroup := user.NewDrupalUserGroup()
			UserGroup.Populate(Alias)
			User := UserGroup.GetUser(*strUser)
			User.Email = Email.Address

			if *boolCreate {
				User.Create(*strPassword)
			}

			if *strState {
				User.State = 1
				User.StateChange()
			} else {
				User.State = 0
				User.StateChange()
			}

			if *strPassword != "" {
				User.SetPassword(*strPassword)
			}

			User.EmailChange()

			if !User.HasRole(*strRole) {
				User.Roles = append(User.Roles, *strRole)
			}

			User.RolesAdd()
		}
	} else {
		flag.Usage()
		os.Exit(1)
	}
}
