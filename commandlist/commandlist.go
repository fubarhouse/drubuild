package commandlist

import (
	"fmt"
	"github.com/fubarhouse/drubuild/command"
	"os/exec"
)

// CommandList supports multiple Command objects as []Command, known as lists here.
type CommandList struct {
	// Our structured data/object for Command
	item []*command.Command
}

// NewDrushCommandList creates a new container for []Command objects
func NewDrushCommandList() CommandList {
	return CommandList{}
}

// Add adds any quantity of single []Command objects to a []Command slice.
func (drush *CommandList) Add(items ...*command.Command) {
	for index, item := range items {
		if item.GetAlias() != "" {
			// Rewrite this alias with the @ symbol as a prefix.
			item.SetAlias(fmt.Sprintf("@%v", item.GetAlias()))
		}
		if item.GetVerbose() == true {
			// Rewrite this alias to include verbose when verbose is set to true.
			drush.item[index].SetAlias(fmt.Sprintf("%v --verbose", drush.item[index].GetAlias()))
		}
		// Add this item to the pointer variable.
		drush.item = append(drush.item, items[index])
	}
}

// Output will return the output from the command
func (drush *CommandList) Output() (string, []error) {
	// Gets the output from a set of []Command objects
	responsesArray, errorsArray := drush.Run()
	return string(responsesArray), errorsArray
}

// Run runs a set of []Command objects
func (drush *CommandList) Run() (string, []error) {
	responses := ""
	errors := []error{}
	for index := range drush.item {
		if drush.item[index].GetAlias() != "" {
			drush.item[index].SetAlias(fmt.Sprintf("%v", drush.item[index].GetAlias()))
		}
		if drush.item[index].GetVerbose() == true {
			drush.item[index].SetAlias(fmt.Sprintf("%v --verbose", drush.item[index].GetAlias()))
		}
		args := fmt.Sprintf("drush %v %v", drush.item[index].GetAlias(), drush.item[index].GetCommand())
		response, err := exec.Command("sh", "-c", args).Output()
		responses = fmt.Sprintf("%v\n%v", responses, string(response))
		errors = append(errors, err)
	}
	return responses, errors
}
