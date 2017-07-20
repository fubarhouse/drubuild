package user

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/golang-drush/command"
	"strconv"
	"strings"
)

// DrupalUserList is a custom type for a slice for Drupal Users in the form of a DrupalUser struct
type DrupalUserList []DrupalUser

// NewDrupalUserGroup generates a new DrupalUserList object.
func NewDrupalUserGroup() DrupalUserList {
	return DrupalUserList{}
}

// FindUser will return a boolean if the query sting is found inside
// the DrupalUser objects of a DrupalUserList of a DrupalUserList object.
func (DrupalUserList *DrupalUserList) FindUser(query string) bool {
	for _, DrupalUser := range *DrupalUserList {
		// Search by User Name
		if DrupalUser.Name == query {
			return true
		}
		// Search by Email
		if DrupalUser.Email == query {
			return true
		}
		// Search by UID
		if fmt.Sprint(DrupalUser.UID) == query {
			return true
		}
	}
	return false
}

// GetUser returns a full user object from a NewDrupalUserGroup object including the Roles field filled in.
func (DrupalUserList *DrupalUserList) GetUser(query string) DrupalUser {
	for _, User := range *DrupalUserList {
		// Search by User Name
		if User.Name == query {
			User.SetRoles()
			return User
		}
	}
	return DrupalUser{}
}

// Populate will populate a DrupalUserList object with the Users from a given alias.
// Existing users in the DrupalUserList object will not be overridden.
func (DrupalUserList *DrupalUserList) Populate(Alias string) {
	DrupalUsers := []DrupalUser{}
	var Command = fmt.Sprint("sqlq \"SELECT uid,name,mail,status FROM users;\"")
	cmd := command.NewDrushCommand()
	cmd.Set(Alias, Command, false)
	cmdOut, cmdErr := cmd.CombinedOutput()
	if cmdErr != nil {
		log.Warnln("Could not execute Drush sql-query:", cmdErr.Error())
	}
	for _, UserID := range strings.Split(string(cmdOut), "\n") {
		UserInfo := strings.Split(UserID, "\t")
		if UserInfo[0] != "" && UserInfo[1] != "" {
			UserState := 0
			if UserInfo[3] == "1" {
				UserState = 1
			}
			UserID, _ := strconv.Atoi(UserInfo[0])
			DrupalUser := DrupalUser{
				Alias, UserID, UserInfo[1], UserInfo[2], UserState, []string{},
			}
			DrupalUsers = append(DrupalUsers, DrupalUser)
		}
	}
	// Ensure previously inputted values do not get overridden.
	for _, DrupalUser := range DrupalUsers {
		*DrupalUserList = append(*DrupalUserList, DrupalUser)
	}
}
