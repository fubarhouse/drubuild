package user

import (
	"fmt"
	"testing"
)

func TestUserListInstansiation(t *testing.T) {
	Group := NewDrupalUserGroup()
	Type := fmt.Sprintf("%T", Group)
	fmt.Println(Type)
	if Type != "user.DrupalUserList" {
		t.Errorf("Type name could not be tested to expectation, expected 'user.DrupalUserList' but got '%v'", Type)
	}
}

func TestUserCreation(t *testing.T) {
	Group := NewDrupalUserGroup()
	NewUser := NewDrupalUser()
	NewUser.Alias = "Test"
	NewUser.Email = "test@test.com"
	NewUser.Name = "Test"
	NewUser.Roles = []string{"Anonymous User"}
	NewUser.State = 1
	NewUser.UID = 666
	Group = append(Group, NewUser)
	if len(Group) != 1 {
		t.Error("Group length was not as expected")
	}
}

func TestUserSearchName(t *testing.T) {
	Group := NewDrupalUserGroup()
	NewUser := NewDrupalUser()
	NewUser.Alias = "Test"
	NewUser.Email = "test@test.com"
	NewUser.Name = "Test"
	NewUser.Roles = []string{"Anonymous User"}
	NewUser.State = 1
	NewUser.UID = 666
	Group = append(Group, NewUser)
	if !Group.FindUser(NewUser.Name) {
		t.Error("Test user was not found in the object.")
	}
}

func TestUserSearchEmail(t *testing.T) {
	Group := NewDrupalUserGroup()
	NewUser := NewDrupalUser()
	NewUser.Alias = "Test"
	NewUser.Email = "test@test.com"
	NewUser.Name = "Test"
	NewUser.Roles = []string{"Anonymous User"}
	NewUser.State = 1
	NewUser.UID = 666
	Group = append(Group, NewUser)
	if !Group.FindUser(NewUser.Email) {
		t.Error("Test user was not found in the object.")
	}
}

func TestUserGet(t *testing.T) {
	Group := NewDrupalUserGroup()
	NewUser := NewDrupalUser()
	NewUser.Alias = "Test"
	NewUser.Email = "test@test.com"
	NewUser.Name = "Test"
	NewUser.Roles = []string{"Anonymous User"}
	NewUser.State = 1
	NewUser.UID = 666
	Group = append(Group, NewUser)
	TestUser := Group.GetUser(NewUser.Name)
	if TestUser.Name != NewUser.Name {
		t.Error("Test user was not returned from getter.")
	}
}
