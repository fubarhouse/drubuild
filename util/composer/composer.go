// Package composer is a basic package to run composer tasks in a Drupal 8 docroot.
package composer

import (
	"errors"
	"io/ioutil"
	"github.com/fubarhouse/drubuild/util/command"
	log "github.com/Sirupsen/logrus"
	"fmt"
	"os"
	"strings"
)

// Copy will copy a file to a destination.
func Copy(src, dest string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return errors.New("could not read " + src + ": " + err.Error())
	} else {
		log.Infof("Successfully read data from %v", src)
	}
	dest = strings.Join([]string{dest, string(os.PathSeparator), "composer.json"}, string(os.PathSeparator))
	err = ioutil.WriteFile(dest, data, 0644)
	if err != nil {
		return errors.New("could not write " + dest + ": " + err.Error())
	} else {
		log.Infof("Successfully wrote data to %v", dest)
	}
	return nil
}

func Run(args []string) (string, error) {
	fmt.Println(args)
	o, e := command.Run("composer", args)
	return o, e
}