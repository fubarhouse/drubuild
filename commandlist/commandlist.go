package commandlist

import (
	"fmt"
	"github.com/fubarhouse/golang-drush/command"
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

// RemoveIndex removes a single []Command object to a []Command slice based upon the index of the item.
func (drush *CommandList) RemoveIndex(indexes ...int) {
	for _, index := range indexes {
		drush.item[index] = &command.Command{}
	}
}

// RemoveCommand removes any []Command object which is effectively identical to the provided to a []Command slice.
func (drush *CommandList) RemoveCommand(item *command.Command) {
	for index := range drush.item {
		if item.GetAlias() == drush.item[index].GetAlias() && item.GetCommand() == drush.item[index].GetCommand() && item.GetVerbose() == drush.item[index].GetVerbose() {
			drush.item[index] = &command.Command{}
		}
	}
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
