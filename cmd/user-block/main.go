//
// user-block will block all Drupal accounts on site alises based upon input values.
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
	var strUser = flag.String("user", "", "User name for blocking, example 'Firstname Sirname'")
	flag.Parse()

	// Remove double spaces.
	*strUser = strings.Replace(*strUser, "  ", " ", -1)

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
				if User.State == 1 {
					User.State = 0
					User.StateChange()
				} else {
					log.Infof("User '%v' is already blocked on %v", User.Name, Alias)
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
