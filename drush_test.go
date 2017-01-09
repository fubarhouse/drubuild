package drush

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
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
	instance := NewDrush(alias, command)
	runcomm, runerr := instance.Output()

	Convey("Testing NewDrush()", t, func() {
		So(instance.alias, ShouldEqual, "none")
		So(instance.command, ShouldEqual, command)
	})
	Convey("Testing drush execution", t, func() {
		So(runerr, ShouldBeNil)
		So(runcomm, ShouldNotBeNil)
	})
}