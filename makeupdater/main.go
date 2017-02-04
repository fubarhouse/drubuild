package makeupdater

// Note this package is exclusively compatible with Drupal 7 make files.

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

func replaceTextInFile(fullPath string, oldString string, newString string) {
	read, err := ioutil.ReadFile(fullPath)
	if err != nil {
		panic(err)
	}
	newContents := strings.Replace(string(read), oldString, newString, -1)
	err = ioutil.WriteFile(fullPath, []byte(newContents), 0)
	if err != nil {
		panic(err)
	}
}

func UpdateMake(fullpath string) []string {
	fmt.Println("Processing " + fullpath + "...")
	affectedProjects := []string{}
	catCmd := "cat " + fullpath + " | grep projects | cut -d'[' -f2 | cut -d']' -f1 | uniq | sort"
	y, _ := exec.Command("sh", "-c", catCmd).Output()
	projects := strings.Split(string(y), "\n")
	for _, project := range projects {
		if project != "" {
			catCmd = "cat " + fullpath + " | grep \"projects\\[" + project + "\\]\" | grep version | cut -d '\"' -f2"
			z, _ := exec.Command("sh", "-c", catCmd).Output()
			versionOld := strings.Replace(string(z), " ", "", -1)
			versionOld = strings.Replace(string(z), "\"", "", -1)
			versionOld = strings.Replace(string(z), "\n", "", -1)

			x, _ := exec.Command("sh", "-c", "drush pm-releases --pipe "+project+" | grep Recommended | cut -d',' -f2").Output()
			versionNew := strings.Replace(string(x), "\n", "", -1)
			versionNew = strings.Replace(versionNew, "7.x-", "", -1)
			versionNew = strings.Replace(versionNew, "\n", "", -1)
			versionNew = strings.Replace(versionNew, " ", "", -1)
			versionNew = strings.Replace(versionNew, "[", "", -1)
			versionNew = strings.Replace(versionNew, "]", "", -1)
			if versionOld != versionNew && versionNew != "" {
				fmt.Printf("Replacing %v v%v with v%v\n", project, versionOld, versionNew)
				affectedProjects = append(affectedProjects, project)
				replaceTextInFile(fullpath, fmt.Sprintf("projects[%v][version] = \"%v\"", project, versionOld), fmt.Sprintf("projects[%v][version] = \"%v\"", project, versionNew))
			}
		}
	}
	return affectedProjects
}
