package command

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Command is our structured data/object for Command
type Command struct {
	alias      string
	command    string
	verbose    bool
	workingdir string
}

const pathDrush = "/usr/local/bin/drush"

// NewDrushCommand creates a new container for []Command objects
func NewDrushCommand() *Command {
	return &Command{}
}

// Set changes all given values for a Drush command object.
func (drush *Command) Set(alias string, command string, verbose bool) {
	drush.alias = alias
	drush.command = command
	drush.verbose = verbose
	drush.workingdir = "."
}

// GetWorkingDir returns the specified working directory used with executed Drush commands.
func (drush *Command) GetWorkingDir() string {
	return drush.workingdir
}

// GetAlias returns the alias used to executed Drush commands.
func (drush *Command) GetAlias() string {
	return drush.alias
}

// SetAlias changes the alias used to executed Drush commands.
func (drush *Command) SetAlias(value string) {
	drush.alias = value
}

// GetCommand returns the command string on executed Drush commands.
func (drush *Command) GetCommand() string {
	return drush.command
}

// SetCommand changes the command string on executed Drush commands.
func (drush *Command) SetCommand(value string) {
	drush.command = value
}

// GetVerbose returns the verbosity setting on executed Drush commands.
func (drush *Command) GetVerbose() bool {
	return drush.verbose
}

func (drush *Command) RawOutput() error {
	if strings.Contains(drush.alias, "@") == true {
		drush.alias = strings.Replace(drush.alias, "@", "", -1)
	}
	if drush.alias != "" {
		drush.alias = fmt.Sprintf("@%v", drush.alias)
	}
	if drush.verbose == true {
		drush.command = fmt.Sprintf("%v --verbose", drush.command)
	}
	args := fmt.Sprintf("%v %v", drush.alias, drush.command)

	comm := new(exec.Cmd)
	comm = exec.Command("sh", "-c", pathDrush+" "+" "+args)
	comm.Dir = drush.GetWorkingDir()
	comm.Stderr = os.Stderr
	comm.Stdout = os.Stdout
	err := comm.Start()
	comm.Wait()
	return err
}

// Output gets the output from a single Command object, does not support []Command items.
func (drush *Command) Output() ([]string, error) {
	comm, err := drush.Run()
	response := filepath.SplitList(string(comm))
	return response, err
}

// CombinedOutput will return the CombinedOutput of a command.
func (drush *Command) CombinedOutput() ([]byte, error) {
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
	if drush.GetWorkingDir() != "." {
		comm, err := exec.Command("sh", "-c", "cd "+drush.workingdir+" && "+pathDrush+" "+args).CombinedOutput()
		return comm, err
	}
	comm, err := exec.Command("sh", "-c", pathDrush+" "+args).CombinedOutput()
	return comm, err
}

// Run runs an individual Command object, does not support []Command items.
func (drush *Command) Run() ([]byte, error) {
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
	if drush.GetWorkingDir() != "." {
		comm, err := exec.Command("sh", "-c", "cd "+drush.workingdir+" && "+pathDrush+" "+args).CombinedOutput()
		return comm, err
	}
	comm, err := exec.Command("sh", "-c", pathDrush+" "+args).CombinedOutput()
	return comm, err
}
