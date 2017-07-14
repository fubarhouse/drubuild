//
// user-unblock will unblock all Drupal accounts on site alises based upon input values.
//
package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/golang-drush/user"
	"os"
	"strings"
)

func main() {
	var strAliases = flag.String("aliases", "", "comma-separated list of aliases for action")
	var strPattern = flag.String("pattern", "%v", "A modifier which allows rewriting of aliases replacing '%v' in the pattern with the alias.")
	var strUser = flag.String("user", "", "User name for unblocking, example 'Firstname Sirname'")
	flag.Parse()

	// Remove double spaces.
	*strAliases = strings.Replace(*strAliases, "  ", " ", -1)

	if !strings.Contains(*strPattern, "%v") {
		log.Errorln("Specified pattern does not include alias modifier.")
		flag.Usage()
		os.Exit(1)
	}

	if *strUser == "" {
		log.Errorln("User name is required and not specified.")
		flag.Usage()
		os.Exit(1)
	}

	if *strAliases == "" {
		flag.Usage()
		log.Errorln("Aliases are not specified.")
		flag.Usage()
		os.Exit(1)
	}

	if *strAliases != "" && *strUser != "" {
		for _, Alias := range strings.Split(*strAliases, ",") {
			Alias = strings.Trim(Alias, " ")
			Alias = strings.Replace(*strPattern, "%v", Alias, 1)
			UserGroup := user.NewDrupalUserGroup()
			UserGroup.Populate(Alias)
			User := UserGroup.GetUser(*strUser)
			if User.Name == *strUser {
				if User.State == 0 {
					User.State = 1
					User.StateChange()
				} else {
					log.Infof("User '%v' is already unblocked on %v", User.Name, Alias)
				}
			} else {
				log.Infof("User '%v' was not found on %v", User.Name, Alias)
			}
		}
	} else {
		flag.Usage()
		os.Exit(1)
	}
}
