package command

import (
	"testing"
)

func TestCreateNewCommand(t *testing.T) {
	// Test the creation of a drush command
	x := NewDrushCommand()
	x.Set("","cc drush", false)
	if x.alias != "" && x.command != "cc drush" && x.verbose != false {
		t.Error("Test failed")
	}
}

func TestCreateNewCommandExecution(t *testing.T) {
	// Test the execution of a drush command
	x := NewDrushCommand()
	x.Set("","cc drush", false)
	_, cmdErr := x.Output()
	if cmdErr != nil {
		t.Error("Test failed")
	}
}

