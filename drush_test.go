package drush

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	//"fmt"
)

func thisAlias() string {
	return "none"
}

func thisCommand() string {
	return ""
}

func TestDrush(t *testing.T) {

	alias := thisAlias()
	command := thisCommand()
	verbose := false
	instance := NewDrush(alias, command, verbose)
	runcomm, runerr := instance.Output()

	Convey("Testing Single Drush command capability", t, func() {
		So(instance.alias, ShouldEqual, "@none")
		So(instance.command, ShouldEqual, command)
		So(instance.verbose, ShouldEqual, false)
		So(runerr, ShouldBeNil)
		So(runcomm, ShouldNotBeNil)
	})

	Convey("Testing Multiple Drush command capability", t, func() {

		// Create a new list object
		drushList := NewDrushList()

		// Create commands to add into list
		command1 := NewDrush("none", "", false)
		command2 := NewDrush("none", "", false)
		command3 := NewDrush("none", "", false)
		command4 := NewDrush("none", "", false)

		// Verify contents of each command
		So(command1, ShouldNotBeNil)
		So(command2, ShouldNotBeNil)
		So(command3, ShouldNotBeNil)
		So(command4, ShouldNotBeNil)

		// Ensure list is empty
		So(drushList, ShouldNotBeNil)

		// Add items into list
		drushList.Add(command1, command2, command3, command4)

		// Verify commands have entered array
		So(drushList.item[0].verbose, ShouldEqual, false)
		So(drushList.item[1].verbose, ShouldEqual, false)
		So(drushList.item[2].verbose, ShouldEqual, false)
		So(drushList.item[3].verbose, ShouldEqual, false)

		// Remove item 1
		drushList.RemoveIndex(1, 2)

		// Execute Drush commands
		outputArray, errorArray := drushList.Output()

		// Ensure the error object of each item is as expected.
		So(errorArray[0], ShouldNotBeEmpty)
		So(errorArray[1], ShouldBeEmpty)
		So(errorArray[2], ShouldBeEmpty)
		So(errorArray[3], ShouldNotBeEmpty)

		// Verify the command object of each item is as expected.
		So(outputArray[0], ShouldNotBeNil)
		So(outputArray[1], ShouldNotBeNil)
		So(outputArray[2], ShouldNotBeNil)
		So(outputArray[3], ShouldNotBeNil)
	})
}