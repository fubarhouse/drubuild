// composer is a basic package to run composer tasks in a Drupal 8 docroot.
package composer

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/golang-drush/makeupdater"
	"os/exec"
	"strings"
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
