package drush

import (
	"fmt"
	"os/exec"
)

// Multiple Drush objects as []Drush, known as lists here.

type DrushList struct {
	// Our structured data/object for Drush
	item []Drush
}

func NewDrushList() DrushList {
	// Creates a new container for []Drush objects
	return DrushList{}
}

func (drush *DrushList) RemoveIndex(indexes... int) {
	// Remove a single []Drush object to a []Drush slice based upon the index of the item.
	for _, index := range indexes {
		drush.item[index] = Drush{}
	}
}

func (drush *DrushList) RemoveCommand(item *Drush) {
	// Remove any []Drush object which is effectively identical to the provided to a []Drush slice.
	for index := range drush.item {
		if item.alias == drush.item[index].alias && item.command == drush.item[index].command && item.verbose == drush.item[index].verbose {
			drush.item[index] = Drush{}
		}
	}
}

func (drush *DrushList) Add(items... *Drush) {
	// Add any quantity of single []Drush objects to a []Drush slice.
	for _, item := range items {
		if item.alias != "" {
			// Rewrite this alias with the @ symbol as a prefix.
			item.alias = fmt.Sprintf("@%v", item.alias)
		}
		if item.verbose == true {
			// Rewrite this alias to include verbose when verbose is set to true.
			item.alias = fmt.Sprintf("%v --verbose", item.alias)
		}
		// Convert this item into the correct format.
		thisItem := Drush{item.alias, item.command, item.verbose}
		// Add this item to the pointer variable.
		drush.item = append(drush.item, thisItem)
	}
}

func (drush *DrushList) Output() (string, []error) {
	// Gets the output from a set of []Drush objects
	responsesArray, errorsArray := drush.Run()
	return string(responsesArray), errorsArray
}

func (drush *DrushList) Run() (string, []error) {
	// Runs a set of []Drush objects
	responses := ""
	errors := []error{}
	for index := range drush.item {
		if drush.item[index].alias != "" { drush.item[index].alias = fmt.Sprintf("@%v", drush.item[index].alias) }
		if drush.item[index].verbose == true { drush.item[index].alias = fmt.Sprintf("%v --verbose", drush.item[index].alias) }
		args := fmt.Sprintf("drush %v %v", drush.item[index].alias, drush.item[index].command)
		response, err := exec.Command("sh", "-c", args).Output()
		responses = fmt.Sprintf("%v\n%v", responses, string(response))
		errors = append(errors, err)
	}
	return responses, errors
}