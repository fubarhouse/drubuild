package user

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/golang-drush/command"
	"strings"
)

// DrupalUser represents fields from Drupals user table, as well as roles.
type DrupalUser struct {
	Alias 	string
	UID 	int
	Name	string
	Email	string
	State	int
	Roles	[]string
}

// NewDrupalUser generates a new DrupalUser object.
func NewDrupalUser() DrupalUser {
	return DrupalUser{}
}

// SetRoles will allocate a valid and accurate value to the Roles field in a given DrupalUser object.
func (DrupalUser *DrupalUser) SetRoles() {
	var RolesCommand= fmt.Sprintf("user-information '%v' --fields=roles | cut -d: -f2", DrupalUser.Name)
	cmd := command.NewDrushCommand()
	cmd.Set(DrupalUser.Alias, RolesCommand, false)
	cmdRolesOut, cmdRolesErr := cmd.CombinedOutput()
	if cmdRolesErr != nil {
		log.Errorln("Could not execute Drush user-information:", cmdRolesErr.Error())
	}
	Roles := []string{}
	for _, Role := range strings.Split(string(cmdRolesOut), "\n") {
		Role = strings.TrimSpace(Role)
		if Role != "" {
			Roles = append(Roles, Role)
		}
	}
	DrupalUser.Roles = Roles
}

// Delete will delete a user from a Drupal site, but only if it exists.
func (DrupalUser *DrupalUser) Delete() {
	UserGroup := NewDrupalUserGroup()
	UserGroup.Populate(DrupalUser.Alias)
	if UserGroup.FindUser(DrupalUser.Name) {
		var Command = fmt.Sprintf("user-cancel --yes '%v'", DrupalUser.Name)
		cmd := command.NewDrushCommand()
		cmd.Set(DrupalUser.Alias, Command, false)
		_, cmdErr := cmd.CombinedOutput()
		if cmdErr != nil {
			log.Warnf("Could not remove user %v on site %v: %v", DrupalUser.Name, DrupalUser.Alias, cmdErr.Error())
		} else {
			log.Infof("Removed user %v on site %v.", DrupalUser.Name, DrupalUser.Alias)
		}
	}
}

// Create will create a user from a Drupal site, but only if does not exist.
func (DrupalUser *DrupalUser) Create(Password string) {
	UserGroup := NewDrupalUserGroup()
	UserGroup.Populate(DrupalUser.Alias)
	if !UserGroup.FindUser(DrupalUser.Name) {
		var Command = fmt.Sprintf("user-create '%v' --mail='%v' --password='%v'", DrupalUser.Name, DrupalUser.Email, Password)
		cmd := command.NewDrushCommand()
		cmd.Set(DrupalUser.Alias, Command, false)
		_, cmdErr := cmd.CombinedOutput()
		if cmdErr != nil {
			log.Warnf("Could not create user %v on site %v: %v", DrupalUser.Name, DrupalUser.Alias, cmdErr.Error())
		} else {
			log.Infof("Created user %v on site %v.", DrupalUser.Name, DrupalUser.Alias)
		}
	}
}

// StateChange will change the status of the user to the value specified in *DrupalUser.State
// There is a built-in verification process here, so a separate verification method is not required.
func (DrupalUser *DrupalUser) StateChange() {
	// Get the absolutely correct User object.
	UserGroup := NewDrupalUserGroup()
	UserGroup.Populate(DrupalUser.Alias)
	User := UserGroup.GetUser(DrupalUser.Name)

	if User.State != DrupalUser.State {
		State := "user-block"
		if User.State == 0 {
			State = "user-unblock"
		}
		cmd := command.NewDrushCommand()
		var Command= fmt.Sprintf("%v '%v'", State, DrupalUser.Name)
		cmd.Set(DrupalUser.Alias, Command, false)
		_, cmdErr := cmd.CombinedOutput()
		if cmdErr != nil {
			log.Warnf("Could not perform action %v for user %v on site %v: %v", State, DrupalUser.Name, DrupalUser.Alias, cmdErr.Error())
		} else {
			log.Infof("Performed action %v for user %v on site %v", State, DrupalUser.Name, DrupalUser.Alias)
		}
	}
}

// SetPassword will set the password of a user.
// Action will be performed, as there is no password validation available.
func (DrupalUser *DrupalUser) SetPassword(Password string) {
	var Command = fmt.Sprintf("user-password \"%v\" --password=\"%v\"", DrupalUser.Name, Password)
	cmd := command.NewDrushCommand()
	cmd.Set(DrupalUser.Alias, Command, false)
	_, cmdErr := cmd.CombinedOutput()
	if cmdErr != nil {
		log.Warnf("Could not complete password change for user %v on site %v: %v", DrupalUser.Name, DrupalUser.Alias, cmdErr.Error())
	} else {
		log.Infof("Password for user %v on site %v has been changed.", DrupalUser.Name, DrupalUser.Alias)
	}
}

// EmailChange will change the email of the target if the email address
// does not match the email address in the DrupalUser object.
func (DrupalUser *DrupalUser) EmailChange() {
	UserGroup := NewDrupalUserGroup()
	UserGroup.Populate(DrupalUser.Alias)
	User := UserGroup.GetUser(DrupalUser.Name)
	if User.Email != DrupalUser.Email && UserGroup.FindUser(DrupalUser.Name) {
		var Command= "sqlq \"UPDATE users SET init='" + User.Email + "', mail='" + DrupalUser.Email + "' WHERE name='" + DrupalUser.Name + "';\""
		cmd := command.NewDrushCommand()
		cmd.Set(DrupalUser.Alias, Command, false)
		_, cmdErr := cmd.CombinedOutput()
		if cmdErr != nil {
			log.Warnf("Could not change email for user %v on site %v from %v to %v: %v", DrupalUser.Name, DrupalUser.Alias, User.Email, DrupalUser.Email, cmdErr.Error())
		} else {
			log.Infof("Changed email for user %v on site %v from %v to %v, clear caches if results are unexpected.", DrupalUser.Name, DrupalUser.Alias, User.Email, DrupalUser.Email)
		}
	}
}

func (DrupalUser *DrupalUser) HasRole(Role string) bool {
	for _, value := range DrupalUser.Roles {
		if value == Role {
			return true
		}
	}
	return false
}

// RolesAdd will add all associated roles to the target user,
// when not present in the DrupalUser object.
func (DrupalUser *DrupalUser) RolesAdd() {
	UserGroup := NewDrupalUserGroup()
	UserGroup.Populate(DrupalUser.Alias)
	User := UserGroup.GetUser(DrupalUser.Name)
	User.SetRoles()
	for _, Role := range DrupalUser.Roles {
		if Role != "authenticated user" {
			if !User.HasRole(Role) {
				var Command = fmt.Sprintf("user-add-role --name='%v' '%v'", DrupalUser.Name, Role)
				cmd := command.NewDrushCommand()
				cmd.Set(DrupalUser.Alias, Command, false)
				_, cmdErr := cmd.CombinedOutput()
				if cmdErr != nil {
					log.Warnf("Could not add role %v to use %v on site %v: %v", Role, DrupalUser.Name, DrupalUser.Alias, cmdErr.Error())
				} else {
					log.Infof("Added user %v to role %v on site %v.", DrupalUser.Name, Role, DrupalUser.Alias)
				}
			}
		}
	}
}

// RolesAdd will remove all associated roles to the target user,
// when present in the DrupalUser object.
func (DrupalUser *DrupalUser) RolesRemove() {
	// if not "authenticated user" {
	// if user has role, and the role needs to be removed, remove the role. {
	UserGroup := NewDrupalUserGroup()
	UserGroup.Populate(DrupalUser.Alias)
	User := UserGroup.GetUser(DrupalUser.Name)
	User.SetRoles()
	for _, Role := range DrupalUser.Roles {
		if Role != "authenticated user" {
			if User.HasRole(Role) {
				var Command = fmt.Sprintf("user-remove-role --name='%v' '%v'", DrupalUser.Name, Role)
				cmd := command.NewDrushCommand()
				cmd.Set(DrupalUser.Alias, Command, false)
				_, cmdErr := cmd.CombinedOutput()
				if cmdErr != nil {
					log.Warnf("Could not remove role %v on user %v on site %v: %v", Role, DrupalUser.Name, DrupalUser.Alias, cmdErr.Error())
				} else {
					log.Infof("Removed user %v from role %v on site %v.", DrupalUser.Name, Role, DrupalUser.Alias)
				}
			}
		}
	}
}