// Package composer is a basic package to run composer tasks in a Drupal 8 docroot.
package composer

import (
	"errors"
		"io/ioutil"
	"github.com/fubarhouse/drubuild/util/command"
)

// Copy will copy a file to a destination.
func Copy(src, dest string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return errors.New("could not read " + src + ": " + err.Error())
	}
	err = ioutil.WriteFile(dest, data, 0644)
	if err != nil {
		return errors.New("could not write " + src + ": " + err.Error())
	}
	return nil
}

func Run(args []string) (string, error) {
	o, e := command.Run("composer", args)
	return o, e
}