package command

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
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

func DrushRebuildRegistry(alias string) {
	drushCommand := NewDrushCommand()
	drushCommand.Set(alias, "rr", false)
	_, err := drushCommand.Output()
	if err != nil {
		log.Warnln("Could not rebuild registry...", err)
	} else {
		log.Infoln("Rebuilt registry.")
	}
}
