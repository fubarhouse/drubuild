package command

import (
	"bufio"
	"fmt"
	log "github.com/Sirupsen/logrus"
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

// SetWorkingDir sets the specified working directory used with executed Drush commands.
func (drush *Command) SetWorkingDir(value string) {
	drush.workingdir = value
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

// SetVerbose changes the verbosity setting on executed Drush commands.
func (drush *Command) SetVerbose(value bool) {
	drush.verbose = value
}

// LiveOutput returns, and prints the live output of the executing program
// This will wait for completion before proceeding.
func (drush *Command) LiveOutput() error {
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
	comm = exec.Command("sh", "-c", pathDrush + " " + " " + args)
	comm.Dir = drush.GetWorkingDir()
	Pipe, _ := comm.StderrPipe()
	scanner := bufio.NewScanner(Pipe)
	go func() {
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "[error]") || strings.Contains(scanner.Text(), "[warning]") {
				log.Warnf("%s", scanner.Text())
			} else {
				log.Infof("%s", scanner.Text())
			}
		}
	}()
	err := comm.Start()
	err = comm.Wait()
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

// DrushDatabaseSync executes a database synchronisation task from a source to destination with the use of Drush.
func DrushDatabaseSync(srcAlias, destAlias string) {
	/* So our binary and this function combined support two-way traffic...  */
	x := NewDrushCommand()
	srcAlias = strings.Replace(srcAlias, "@", "", -1)
	destAlias = strings.Replace(destAlias, "@", "", -1)
	x.Set("", fmt.Sprintf("sql-sync @%v @%v -y", srcAlias, destAlias), true)
	_, err := x.Output()
	if err == nil {
		log.Infoln("Syncronised databases complete.")
	} else {
		log.Errorln("Could not syncronise databases.")
	}
}

// DrushFilesSync executes a file synchronisation task from a source to destination with the use of Drush.
func DrushFilesSync(srcAlias, destAlias string) {
	x := NewDrushCommand()
	srcAlias = strings.Replace(srcAlias, "@", "", -1)
	destAlias = strings.Replace(destAlias, "@", "", -1)
	x.Set("", fmt.Sprintf("--yes rsync --exclude-other-sites --exclude-conf @%v:%%files @%v:%%files", srcAlias, destAlias), true)
	_, err := x.Output()
	if err == nil {
		log.Infoln("Synced public file system.")
	} else {
		log.Warnln("Public file system has not been synced.")
	}
	x.Set("", fmt.Sprintf("--yes rsync --exclude-other-sites --exclude-conf @%v:%%private @%v:%%private", srcAlias, destAlias), true)
	_, err = x.Output()
	if err == nil {
		log.Infoln("Synced private file system.")
	} else {
		log.Warnln("Private file system has not been synced.")
	}
}

// DrushClearCache performs a cache clear task on an input site alias with the use of Drush.
func DrushClearCache(alias string) {
	drushCommand := NewDrushCommand()
	drushCommand.Set(alias, "cc all", false)
	_, err := drushCommand.Output()
	if err != nil {
		log.Warnln("Could not clear caches.", err)
	} else {
		log.Infoln("Caches cleared.")
	}
}

// DrushRebuildRegistry performs a registry rebuild task on an input site alias with the use of Drush.
func DrushRebuildRegistry(alias string) {
	drushCommand := NewDrushCommand()
	drushCommand.Set(alias, "rr", false)
	_, err := drushCommand.Output()
	if err != nil {
		log.Warnln("Could not rebuild registry.", err)
	} else {
		log.Infoln("Rebuilt registry.")
	}
}

// DrushUpdateDatabase performs a database update task on an input site alias with the use of Drush.
func DrushUpdateDatabase(alias string) {
	drushCommand := NewDrushCommand()
	drushCommand.Set(alias, "updb -y", false)
	_, err := drushCommand.Output()
	if err != nil {
		log.Warnln("Could not update database.", err)
	} else {
		log.Infoln("Updated database where possible.")
	}
}

// DrushDownloadToPath performs a database update task on an input site alias with the use of Drush.
func DrushDownloadToPath(path, project string, version int64) {
	majorversion := fmt.Sprintf("%v", version)
	drushCommand := NewDrushCommand()
	drushCommand.Set("", "pm-download --yes "+project+" --default-major="+majorversion+" --destination="+path, false)
	_, err := drushCommand.Output()
	if err != nil {
		log.Warnln("Could not download module ", project, err)
	} else {
		log.Infoln("Downloaded module", project)
	}
}

// DrushDownloadToAlias performs a database update task on an input site alias with the use of Drush.
func DrushDownloadToAlias(alias, project string) {
	drushCommand := NewDrushCommand()
	drushCommand.Set(alias, "pm-download --yes "+project, false)
	_, err := drushCommand.Output()
	if err != nil {
		log.Warnln("Could not download module ", project, err)
	} else {
		log.Infoln("Downloaded module", project)
	}
}

// DrushVariableSet Runs drush vset with a given variable name and value.
func DrushVariableSet(alias, variableName, variableValue string) {
	srcAlias := strings.Replace(alias, "@", "", -1)
	x := NewDrushCommand()
	x.Set(srcAlias, fmt.Sprintf("vset %v %v", variableName, variableValue), false)
	drushOut, err := x.Output()
	if err == nil {
		log.Infof("Successfully set %v to %v via Drush", variableName, variableValue)
	} else {
		log.Errorf("Could not set %v to %v via Drush: %v", variableName, variableValue, drushOut)
	}
}

// DrushVariableGet Runs drush vget with an exact variable name from the alias.
func DrushVariableGet(alias, variableName string) string {
	srcAlias := strings.Replace(alias, "@", "", -1)
	x := NewDrushCommand()
	x.Set(srcAlias, fmt.Sprintf("vget --exact %v", variableName), false)
	drushOut, err := x.Output()
	drushOutString := fmt.Sprintf("%s", drushOut)
	if strings.Contains(drushOutString, "No matching variable found") {
		log.Warnf("Variable %v was not found", variableName)
	} else if err == nil {
		log.Infof("Successfully retreived %v via Drush", variableName)
		drushOutString = strings.Replace(drushOutString, "[", "", -1)
		drushOutString = strings.Replace(drushOutString, "]", "", -1)
		drushOutString = strings.Replace(drushOutString, "\n", "", -1)
		return drushOutString
	} else {
		log.Errorf("Could not retreived %v via Drush: %v", variableName, drushOut)
	}
	return ""
}
