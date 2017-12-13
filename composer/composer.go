// composer is a basic package to run composer tasks in a Drupal 8 docroot.
package composer

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/drubuild/makeupdater"
	"os/exec"
	"strings"
	"path/filepath"
	"os"
)

// DrupalProject is a type to provide both name and verison of a given Drupal project.
type DrupalProject struct {
	Project string
	Version string
	Patch   string
	Subdir  string
}

// GetProjects will return all projects in a make file with the format of a DrupalProject.
// This function will slowly take over a lot of similar functionality.
func GetProjects(fullpath string) []DrupalProject {
	Projects := []DrupalProject{}
	for _, Project := range makeupdater.GetProjectsFromMake(fullpath) {
		catCmd := fmt.Sprintf("cat %v | grep \"projects\\[%v\\]\"", fullpath, Project)
		y, _ := exec.Command("sh", "-c", catCmd).CombinedOutput()
		DrupalProject := DrupalProject{}
		for _, Line := range strings.Split(string(y), "\n") {
			if strings.Contains(Line, "projects["+Project+"][version] = ") {
				Version := strings.Split(Line, "=")
				Version[1] = strings.Trim(Version[1], " ")
				Version[1] = strings.Replace(Version[1], "\"", "", -1)
				DrupalProject.Version = Version[1]
			} else if strings.Contains(Line, "projects["+Project+"][subdir] = ") {
				Subdir := strings.Split(Line, "=")
				Subdir[1] = strings.Trim(Subdir[1], " ")
				Subdir[1] = strings.Replace(Subdir[1], "\"", "", -1)
				DrupalProject.Subdir = Subdir[1]
			} else if strings.Contains(Line, "projects["+Project+"][patch] = ") {
				Patch := strings.Split(Line, "=")
				Patch[1] = strings.Trim(Patch[1], " ")
				Patch[1] = strings.Replace(Patch[1], "\"", "", -1)
				DrupalProject.Patch = Patch[1]
			}
		}
		if Project != "drupal" {
			DrupalProject.Project = Project
			Projects = append(Projects, DrupalProject)
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
	filepath.Walk(Path, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})

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
		cpCmd := exec.Command("composer", "install", "--prefer-dist", "--working-dir=" + v)
		cpOut, cpErr := cpCmd.CombinedOutput()
		if cpErr != nil {
			log.Errorln("Could not complete:", string(cpOut), cpErr)
		} else {
			log.Println(string(cpOut))
		}
	}
}