package drush

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

type Drush struct {
	alias, command string
}

func NewDrush(a string, b string) Drush {
	return Drush{a, b}
}

func (drush *Drush) Run() ([]byte, error) {
	args := fmt.Sprintf("drush @%v %v", drush.alias, drush.command)
	comm, err := exec.Command("sh", "-c", args).Output()
	return comm, err
}

func (drush *Drush) Output() ([]string, error) {
	comm, err := drush.Run()
	response := filepath.SplitList(string(comm))
	return response, err
}