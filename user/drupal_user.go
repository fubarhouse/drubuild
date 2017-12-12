package user

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/drubuild/command"
	"os"
	"os/exec"
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

// NewDrupalUser generates a new DrupalUser object.
func NewDrupalUser() DrupalUser {
	return DrupalUser{}
}

// SetRoles will allocate a valid and accurate value to the Roles field in a given DrupalUser object.
func (DrupalUser *DrupalUser) SetRoles() {
	var RolesCommand = fmt.Sprintf("user-information '%v' --fields=roles | cut -d: -f2", DrupalUser.Name)
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
	d, e := exec.LookPath("drush")
	if e != nil {
		panic(e)
	}
	var Command = fmt.Sprintf("user-cancel --yes '%v'", DrupalUser.Name)
	cmd := exec.Command(d, DrupalUser.Alias, Command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	cmd.Wait()
}

// Create will create a user from a Drupal site, but only if does not exist.
func (DrupalUser *DrupalUser) Create(Password string) {
	d, e := exec.LookPath("drush")
	if e != nil {
		panic(e)
	}
	var Command = fmt.Sprintf("user-create '%v' --mail='%v' --password='%v'", DrupalUser.Name, DrupalUser.Email, Password)
	cmd := exec.Command(d, DrupalUser.Alias, Command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	cmd.Wait()
}

// Unblock will change the status of the user to the value specified in *DrupalUser.State
// There is a built-in verification process here, so a separate verification method is not required.
func (DrupalUser *DrupalUser) Unblock() {
	d, e := exec.LookPath("drush")
	if e != nil {
		panic(e)
	}
	var Command = fmt.Sprintf("%v '%v'", "user-unblock", DrupalUser.Name)
	cmd := exec.Command(d, DrupalUser.Alias, Command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	cmd.Wait()
}

// Block will change the status of the user to the value specified in *DrupalUser.State
// There is a built-in verification process here, so a separate verification method is not required.
func (DrupalUser *DrupalUser) Block() {
	d, e := exec.LookPath("drush")
	if e != nil {
		panic(e)
	}
	var Command = fmt.Sprintf("%v '%v'", "user-block", DrupalUser.Name)
	cmd := exec.Command(d, DrupalUser.Alias, Command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	cmd.Wait()
}

// SetPassword will set the password of a user.
// Action will be performed, as there is no password validation available.
func (DrupalUser *DrupalUser) SetPassword(Password string) {
	d, e := exec.LookPath("drush")
	if e != nil {
		panic(e)
	}
	var Command = fmt.Sprintf("user-password \"%v\" --password=\"%v\"", DrupalUser.Name, Password)
	cmd := exec.Command(d, DrupalUser.Alias, Command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	cmd.Wait()
}

// EmailChange will change the email of the target if the email address
// does not match the email address in the DrupalUser object.
func (DrupalUser *DrupalUser) EmailChange() {
	UserGroup := NewDrupalUserGroup()
	UserGroup.Populate(DrupalUser.Alias)
	User := UserGroup.GetUser(DrupalUser.Name)
	if User.Email != DrupalUser.Email && UserGroup.FindUser(DrupalUser.Name) {
		d, e := exec.LookPath("drush")
		if e != nil {
			panic(e)
		}
		var Command = "sqlq \"UPDATE users SET init='" + User.Email + "', mail='" + DrupalUser.Email + "' WHERE name='" + DrupalUser.Name + "';\""
		cmd := exec.Command(d, DrupalUser.Alias, Command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		cmd.Wait()

		_, cmdErr := cmd.CombinedOutput()
		log.Println(cmdErr)
		if cmdErr != nil {
			log.Warnf("User email address was not changed")
			log.Println(Command)
		} else {
			fmt.Println("User email address was changed")
		}
	} else if User.Email == DrupalUser.Email {
		fmt.Println("Email address already matches, not changing.")
	}
}

// HasRole will determine if the user has a given String in the list of roles, which will return as a Boolean.
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
//
// TODO: This code calls duplicates where user already has a role.
func (DrupalUser *DrupalUser) RolesAdd() {
	UserGroup := NewDrupalUserGroup()
	UserGroup.Populate(DrupalUser.Alias)
	User := UserGroup.GetUser(DrupalUser.Name)
	User.SetRoles()

	for _, Role := range DrupalUser.Roles {
		if !User.HasRole(Role) {
			fmt.Println(!User.HasRole(Role))

			d, e := exec.LookPath("drush")
			if e != nil {
				panic(e)
			}
			var Command = fmt.Sprintf("user-add-role --name='%v' '%v'", DrupalUser.Name, Role)
			cmd := exec.Command(d, DrupalUser.Alias, Command)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
			cmd.Wait()
		} else {
			fmt.Printf("User already has role '%v'\n", Role)
		}
	}
}

// RolesRemove will remove all associated roles to the target user,
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
				d, e := exec.LookPath("drush")
				if e != nil {
					panic(e)
				}
				var Command = fmt.Sprintf("user-add-role --name='%v' '%v'", DrupalUser.Name, Role)
				cmd := exec.Command(d, DrupalUser.Alias, Command)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
				cmd.Wait()
			}
		}
	}
}
