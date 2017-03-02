package command

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type Command struct {
	// Our structured data/object for Command
	alias   string
	command string
	verbose bool
}

const PATH_DRUSH = "/usr/local/bin/drush"

func NewDrushCommand() *Command {
	// Creates a new container for []Command objects
	return &Command{}
}

func (drush *Command) Set(alias string, command string, verbose bool) {
	drush.alias = alias
	drush.command = command
	drush.verbose = verbose
}

func (drush *Command) GetAlias() string {
	return drush.alias
}

func (drush *Command) SetAlias(value string) {
	drush.alias = value
}

func (drush *Command) GetCommand() string {
	return drush.command
}

func (drush *Command) SetCommand(value string) {
	drush.command = value
}

func (drush *Command) GetVerbose() bool {
	return drush.verbose
}

func (drush *Command) SetVerbose(value bool) {
	drush.verbose = value
}

func (drush *Command) Output() ([]string, error) {
	// Gets the output from a single Command object, does not support []Command items.
	comm, err := drush.Run()
	response := filepath.SplitList(string(comm))
	return response, err
}

func (drush *Command) Run() ([]byte, error) {
	// Run an individual Command object, does not support []Command items.
	if strings.Contains(drush.alias, "@") == true {
		drush.alias = strings.Replace(drush.alias, "@", "", -1)
	}
	if drush.alias != "" {
		drush.alias = fmt.Sprintf("@%v", drush.alias)
	}
	if drush.verbose == true {
		drush.alias = fmt.Sprintf("%v --verbose", drush.alias)
	}
	args := fmt.Sprintf("%v %v", drush.alias, drush.command)
	comm, err := exec.Command("sh", "-c", PATH_DRUSH+" "+args).Output()
	return comm, err
}
