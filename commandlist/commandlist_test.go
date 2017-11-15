package commandlist

import (
	"fmt"
	"github.com/fubarhouse/drubuild/command"
	"log"
	"testing"
)

func TestCreateNewCommand(t *testing.T) {
	// Test a drush command creation for the next test
	y := command.NewDrushCommand()
	y.Set("", "cc drush", false)

	if y.GetCommand() != "cc drush" {
		t.Error("Test failed")
	}
}

func TestCreateNewCommandList(t *testing.T) {
	// Test the creation of a drush command object
	y := NewDrushCommandList()
	x := command.NewDrushCommand()
	x.Set("", "cc drush", false)
	y.Add(x)
}

func TestCreateNewCommandListExecution(t *testing.T) {
	// Test the creation of package object
	y := NewDrushCommandList()
	x := command.NewDrushCommand()
	x.Set("", "cc drush", false)
	y.Add(x)
	_, drushError := y.Output()
	if fmt.Sprint(drushError) != "[<nil>]" {
		log.Println(drushError)
		t.Error("Test failed")
	}
}
