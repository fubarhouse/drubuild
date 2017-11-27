// Package composer is a basic package to run composer tasks in a Drupal 8 docroot.
package composer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// DrupalProject is a type to provide both name and verison of a given Drupal project.
type DrupalProject struct {
	Project string
	Version string
	Patch   string
	Subdir  string
}

// GetPath will get the relative file path of a given project.
// This assumes that the composer.json file is setup for a custom
// package, where the path is allocated to a custom path for a given
// package. The string {$name} will be replaced with the input project name.
// The input project name should reflect the dependency declaration in
// composer.json. An example of this is drupal/views
func GetPath(fullpath, project string) (string, error) {
	var projectType string
	var projectPath string
	name := strings.Split(project, "/")[1]
	// Find composer.
	composer, e := exec.LookPath("composer")
	if e != nil {
		return "", e
	}
	// Get project information.
	c := exec.Command(composer, "show", project)
	c.Dir = fullpath
	o, oe := c.Output()
	if e != nil {
		return "", oe
	}
	// Get the package type from the output of the above command.
	for _, pr := range strings.Split(string(o), "\n") {
		if strings.Contains(pr, "type") {
			s := strings.Split(pr, ":")[1]
			projectType = strings.Trim(s, " ")
		}
	}
	// Get the information required from composer.json.
	p := strings.Join([]string{fullpath, string(os.PathSeparator), "composer.json"}, string(os.PathSeparator))
	data, de := ioutil.ReadFile(p)
	if de != nil {
		return "", oe
	}
	// Begin processing the file looking for the project path.
	allData := strings.Split(string(data), "\n")
	for _, d := range allData {
		cd := fmt.Sprintf("\"%v\":", projectType)
		if strings.Contains(d, cd) {
			// We found the project type declaration, pull out and process the piece required.
			projectPath = strings.Split(d, ":")[1]
			projectPath = strings.Trim(projectPath, " ")
			projectPath = strings.Replace(projectPath, string(os.PathSeparator)+"{$name}", "", -1)
			projectPath = strings.Replace(projectPath, "\"", "", -1)
			projectPath = strings.Replace(projectPath, ",", "", -1)
			projectPath = strings.TrimRight(projectPath, string(os.PathSeparator))
			projectPath = strings.Join([]string{projectPath, name}, string(os.PathSeparator))
		}
	}

	// Return the end result.
	return projectPath, nil
}

// copy will copy a file to a destination.
func copy(src, dest string) error {
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

// InstallComposerCodebase will accept a []string of paths
// and run a composer install over each of the files found.
func InstallComposerCodebase(Name, Timestamp string, ComposerFile, Destination string, workingCopy bool) {
	Destination += "/" + Name + Timestamp
	// Identify if copying the file is required.
	ComposerPath := strings.Replace(ComposerFile, "/composer.json", "", -1)
	ComposerPath = strings.TrimRight(ComposerPath, "/")
	ComposerDestination := strings.TrimRight(Destination, "/") + "/" + Name + Timestamp

	if _, err := os.Stat(Destination); err != nil {
		ok := os.MkdirAll(Destination, 0700)
		if ok != nil {
			log.Fatalf("could not create directory %v: %v", Destination, ok.Error())
		}
	}

	if !strings.HasSuffix(ComposerDestination, ComposerPath) {
		log.Infof("composer.json not found, copying from %v", ComposerFile)
		e := copy(ComposerFile, Destination+"/composer.json")
		if e != nil {
			log.Warnln(e)
		} else {
			log.Printf("Copied %v to %v\n", ComposerFile, Destination+"/composer.json")
		}
	} else {
		log.Infof("%v/composer.json was found, not copying", ComposerDestination)
	}
	var c string
	if workingCopy {
		c = "composer install --prefer-source"
	} else {
		c = "composer install --prefer-dist"
	}
	cpCmd := exec.Command("sh", "-c", c)
	cpCmd.Dir = Destination
	cpCmd.Stdout = os.Stdout
	cpCmd.Stderr = os.Stderr
	cpCmd.Run()
	cpCmd.Wait()
}
