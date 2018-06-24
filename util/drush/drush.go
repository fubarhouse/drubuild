package drush

import "github.com/fubarhouse/drubuild/util/command"

// Run a drush command with the input arguments.
func Run(args []string) (string, error) {
	o, e := command.Run("drush", args)
	return o, e
}
