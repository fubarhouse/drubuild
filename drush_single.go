package drush

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

type Drush struct {
	// Our structured data/object for Drush
	alias string
	command string
	verbose bool
}

// Individual item functions & methods.

func NewDrush(alias string, command string, verbose bool) *Drush {
	// Create a new Drush object with values.
	drushOpts := new(Drush)
	drushOpts.alias = alias
	drushOpts.command = command
	drushOpts.verbose = verbose
	return drushOpts
}

func (drush *Drush) Output() ([]string, error) {
	// Gets the output from a single Drush object, does not support []Drush items.
	comm, err := drush.Run()
	response := filepath.SplitList(string(comm))
	return response, err
}

func (drush *Drush) Run() ([]byte, error) {
	// Run an individual Drush object, does not support []Drush items.
	if drush.alias != "" { drush.alias = fmt.Sprintf("@%v", drush.alias) }
	if drush.verbose == true { drush.alias = fmt.Sprintf("%v --verbose", drush.alias) }
	args := fmt.Sprintf("drush %v %v", drush.alias, drush.command)
	comm, err := exec.Command("sh", "-c", args).Output()
	return comm, err
}