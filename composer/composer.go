// Package composer is a basic package to run composer tasks in a Drupal 8 docroot.
package composer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/drubuild/makeupdater"
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

// GetProjects will return all projects in a make file with the format of a DrupalProject.
// This function will slowly take over a lot of similar functionality.
func GetProjects(fullpath string) []DrupalProject {
	Projects := []DrupalProject{}
	for _, Project := range makeupdater.GetProjectsFromMake(fullpath) {
		catCmd := fmt.Sprintf("cat %v | grep \"projects\\[%v\\]\"", fullpath, Project)
		y, e := exec.Command("sh", "-c", catCmd).CombinedOutput()
		if e != nil {
			log.Warnf("Could not execute `%v`\n", catCmd)
		}
		project := DrupalProject{}
		for _, Line := range strings.Split(string(y), "\n") {
			if strings.Contains(Line, "projects["+Project+"][version] = ") {
				Version := strings.Split(Line, "=")
				Version[1] = strings.Trim(Version[1], " ")
				Version[1] = strings.Replace(Version[1], "\"", "", -1)
				project.Version = Version[1]
			} else if strings.Contains(Line, "projects["+Project+"][subdir] = ") {
				Subdir := strings.Split(Line, "=")
				Subdir[1] = strings.Trim(Subdir[1], " ")
				Subdir[1] = strings.Replace(Subdir[1], "\"", "", -1)
				project.Subdir = Subdir[1]
			} else if strings.Contains(Line, "projects["+Project+"][patch] = ") {
				Patch := strings.Split(Line, "=")
				Patch[1] = strings.Trim(Patch[1], " ")
				Patch[1] = strings.Replace(Patch[1], "\"", "", -1)
				project.Patch = Patch[1]
			}
		}
		if Project != "drupal" {
			project.Project = Project
			Projects = append(Projects, project)
		}
	}
	return Projects
}

// InstallProjects will install drupal composer projects via composer.
// It will do this for every project in the make file, or as a []DrupalProject input.
func InstallProjects(Projects []DrupalProject, Path string) {
	for _, Project := range Projects {
		log.Infof("Processing drupal/%v:%v:", Project.Project, Project.Version)
		ProjectString := fmt.Sprintf("drupal/%v:%v", Project.Project, Project.Version)
		cpCmd := exec.Command("sh", "-c", "cd "+Path+" && composer require "+ProjectString)
		_, cpErr := cpCmd.CombinedOutput()
		if cpErr != nil {
			log.Errorln("Could not complete:", cpErr)
		}
	}
}

// FindComposerJSONFiles will find all the composer.json files inside custom modules/themes
// in a Drupal 8 (or even 7) file system. It will return the paths  in a []string
// which can be iterated over for processing.
func FindComposerJSONFiles(Path string) []string {
	fileList := []string{}
	fmt.Println(len(fileList))
	e := filepath.Walk(Path, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})

	if e != nil {
		log.Warnf("Could not scan for composer.json files under %v: %v\n", Path, e)
	}

	results := []string{}
	for _, file := range fileList {
		if strings.Contains(file, "/custom/") && strings.HasSuffix(file, "composer.json") {
			results = append(results, file)
		}
	}

	return results
}

// InstallComposerJSONFiles will accept a []string of paths
// and run a composer install over each of the files found.
func InstallComposerJSONFiles(Paths []string) {
	for _, v := range Paths {
		v = strings.Replace(v, "composer.json", "", -1)
		cpCmd := exec.Command("composer", "install", "--prefer-dist", "--working-dir="+v)
		cpOut, cpErr := cpCmd.CombinedOutput()
		if cpErr != nil {
			log.Errorln("Could not complete:", string(cpOut), cpErr)
		} else {
			log.Println(string(cpOut))
		}
	}
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
