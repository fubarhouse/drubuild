package user

import (
	"fmt"
	"strings"

					"github.com/fubarhouse/drubuild/util/drush"
)

// DrupalUser represents fields from Drupals user table, as well as roles.
type DrupalUser struct {
	Alias string
	UID   int
	Name  string
	Email string
	State int
	Roles []string
}

// SetRoles will allocate a valid and accurate value to the Roles field in a given DrupalUser object.
func SetRole(alias, name string) {
	var RolesCommand = fmt.Sprintf("user-information '%v' --fields=roles | cut -d: -f2", name)
	cmdRolesOut, _ := drush.Run([]string{alias, RolesCommand})
	Roles := []string{}
	for _, Role := range strings.Split(cmdRolesOut, "\n") {
		Role = strings.TrimSpace(Role)
		if Role != "" {
			Roles = append(Roles, Role)
		}
	}
}

// Delete will delete a user from a Drupal site, but only if it exists.
func Delete(alias, name string) {
	var Command = fmt.Sprintf("user-cancel --yes '%v'", name)
	drush.Run([]string{alias, Command})
}

// Create will create a user from a Drupal site, but only if does not exist.
func Create(alias, name, email, password string) {
	var Command = fmt.Sprintf("user-create '%v' --mail='%v' --password='%v'", name, email, password)
	drush.Run([]string{alias, Command})
}

// Unblock will change the status of the user to the value specified in *DrupalUser.State
// There is a built-in verification process here, so a separate verification method is not required.
func Unblock(alias, name string) {
	var Command = fmt.Sprintf("%v '%v'", "user-unblock", name)
	drush.Run([]string{alias, Command})
}

// Block will change the status of the user to the value specified in *DrupalUser.State
// There is a built-in verification process here, so a separate verification method is not required.
func Block(alias, name string) {
	var Command = fmt.Sprintf("%v '%v'", "user-block", name)
	drush.Run([]string{alias, Command})
}

// SetPassword will set the password of a user.
// Action will be performed, as there is no password validation available.
func SetPassword(alias, name, password string) {
	var Command = fmt.Sprintf("user-password \"%v\" --password=\"%v\"", name, password)
	drush.Run([]string{alias, Command})
}

// EmailChange will change the email of the target if the email address
// does not match the email address in the DrupalUser object.
func EmailChange(alias, name, email string) {
	UserGroup := NewDrupalUserGroup()
	UserGroup.Populate(alias)
	User := UserGroup.GetUser(name)
	if User.Email != email && UserGroup.FindUser(name) {
		var Command = "sqlq \"UPDATE users SET init='" + User.Email + "', mail='" + email + "' WHERE name='" + name + "';\""
		drush.Run([]string{alias, Command})

	} else if User.Email == email {
		fmt.Println("Email address already matches, not changing.")
	}
}

// HasRole will determine if the user has a given String in the list of roles, which will return as a Boolean.
func HasRole(alias, name, role string) bool {
	// TODO: Rewrite.
	return false
}

// RolesAdd will add all associated roles to the target user,
// when not present in the DrupalUser object.
//
// TODO: This code calls duplicates where user already has a role.
func RolesAdd(alias, name string, roles []string) {
	UserGroup := NewDrupalUserGroup()
	UserGroup.Populate(alias)
	SetRole(alias, name)

	for _, Role := range roles {
		if !HasRole(alias, name, Role) {
			fmt.Println(!HasRole(alias, name, Role))

			var Command = fmt.Sprintf("user-add-role --name='%v' '%v'", name, Role)
			drush.Run([]string{alias, Command})
		} else {
			fmt.Printf("User already has role '%v'\n", Role)
		}
	}
}

// RolesRemove will remove all associated roles to the target user,
// when present in the DrupalUser object.
func RolesRemove(alias, name string, roles []string) {
	// if not "authenticated user" {
	// if user has role, and the role needs to be removed, remove the role. {
	UserGroup := NewDrupalUserGroup()
	UserGroup.Populate(alias)
	SetRole(alias, name)
	for _, Role := range roles {
		if Role != "authenticated user" {
			if HasRole(alias, name, Role) {
				var Command = fmt.Sprintf("user-add-role --name='%v' '%v'", name, Role)
				drush.Run([]string{alias, Command})
			}
		}
	}
}
