package makeupdater

// Note this package is exclusively compatible with Drupal 7 make files.

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// replaceTextInFile will replace a string of test in a file.
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

// removeChar will remove particular characters from a string.
func removeChar(input string, chars ...string) string {
	for _, value := range chars {
		input = strings.Replace(input, value, "", -1)
	}
	return input
}

// inArray will return the quanity of specific input values in the input slice.
func inArray(input []string, subject string) int {
	counter := 0
	for _, value := range input {
		if value == subject {
			counter++
		}
	}
	return counter
}

// UpdateMake will update the version numbers in a specified make file
func UpdateMake(fullpath string) {
	_, err := os.Stat(fullpath)
	if err != nil {
		panic(err)
	}
	projects := GetProjectsFromMake(fullpath)
	count := 0
	for _, project := range projects {
		if project != "" {
			catCmd := "cat " + fullpath + " | grep \"projects\\[" + project + "\\]\" | grep version | cut -d '=' -f2"
			z, _ := exec.Command("sh", "-c", catCmd).Output()
			for _, stream := range strings.Split(string(z), "\n") {
				stream = strings.Replace(stream, " ", "", -1)
				stream = strings.Replace(stream, "\"", "", -1)
				if stream != "" {
					x, _ := exec.Command("sh", "-c", "drush pm-releases --pipe "+project+" | grep Recommended | cut -d',' -f2").Output()
					versionNew := removeChar(string(x), " ", "7.x-", "\"", "\n", "[", "]")
					if !strings.Contains(stream, versionNew) {
						fmt.Printf("Replacing %v v%v with v%v\n", project, stream, versionNew)
						replaceTextInFile(fullpath, fmt.Sprintf("projects[%v][version] = \"%v\"\n", project, stream), fmt.Sprintf("projects[%v][version] = \"%v\"\n", project, versionNew))
						count++
					}
				}
			}
		}
	}
	if count == 0 {
		fmt.Printf("%v is already up to date.", fullpath)
	}
}

// FindDuplicatesInMake will find and report Duplicate projects in Drupal make files.
// It will not return a value.
func FindDuplicatesInMake(makefile string) {
	projects := GetProjectsFromMake(makefile)
	// Run a short report containing information on all duplicates.
	for _, project := range projects {
		projectCounter := 0
		if project != "" {
			catCmd := "cat " + makefile + " | grep \"projects\\[" + project + "\\]\" | grep version | cut -d '=' -f2"
			z, _ := exec.Command("sh", "-c", catCmd).Output()
			for _, stream := range strings.Split(string(z), "\n") {
				if stream != "" {
					projectCounter ++
				}
			}
			if projectCounter > 1 {
				fmt.Printf("Found %v instances of project %v\n", projectCounter, project)
			}
		}
	}
}

// GetProjectsFromMake returns a list of projects from a given make file
func GetProjectsFromMake(fullpath string) []string {
	Projects := []string{}
	catCmd := "cat " + fullpath + " | grep projects | cut -d'[' -f2 | cut -d']' -f1 | uniq | sort"
	y, _ := exec.Command("sh", "-c", catCmd).Output()
	rawProjects := strings.Split(string(y), "\n")
	for _, project := range rawProjects {
		project = strings.Replace(project, " ", "", -1)
		if project != "" && project != "projects" {
			if inArray(Projects, project) == 0 {
				Projects = append(Projects, project)
			}
		}
	}
	return Projects
}

// GenerateMake takes a []string of projects and writes out a make file
// Modules are added with the latest recommended version.
func GenerateMake(Projects []string, File string) {
	headerLines := []string{}
	headerLines = append(headerLines, "; Generated by make-updater")
	headerLines = append(headerLines, "; Script created by Fubarhouse")
	headerLines = append(headerLines, "; Toolkit available at github.com/fubarhouse/golang-drush/...")
	headerLines = append(headerLines, "core = 7.x")
	headerLines = append(headerLines, "api = 2")
	headerLines = append(headerLines, "")

	// Rewrite core, if core is in the original Projects list.

	for _, Project := range Projects {
		coreAppended := 0
		if Project == "drupal" {
			if coreAppended == 0 {
				headerLines = append(headerLines, "; core")
				x, _ := exec.Command("sh", "-c", "drush pm-releases --pipe drupal | grep Recommended | cut -d',' -f2").Output()
				ProjectVersion := removeChar(string(x), " ", "7.x-", "\"", "\n", "[", "]")
				headerLines = append(headerLines, "projects[drupal][type] = \"core\"")
				headerLines = append(headerLines, fmt.Sprintf("projects[drupal][version] = \"%v\"", ProjectVersion))
				headerLines = append(headerLines, "projects[drupal][download][type] = \"get\"")
				headerLines = append(headerLines, fmt.Sprintf("projects[drupal][download][url] = \"https://ftp.drupal.org/files/projects/drupal-%v.tar.gz\"", ProjectVersion))
				headerLines = append(headerLines, "")
				coreAppended++
			}
		}
	}

	// Rewrite contrib
	headerLines = append(headerLines, "; modules")
	headerLines = append(headerLines, "defaults[projects][subdir] = contrib")
	headerLines = append(headerLines, "")

	for _, Project := range Projects {

		if Project != "drupal" {
			x, y := exec.Command("sh", "-c", "drush pm-releases --pipe "+Project+" | grep Recommended | cut -d',' -f2").Output()
			if y == nil {
				ProjectVersion := removeChar(string(x), " ", "7.x-", "\"", "\n", "[", "]")
				ProjectType := "contrib"
				if ProjectVersion == "" {
					ProjectType = "custom"
				}
				headerLines = append(headerLines, fmt.Sprintf("projects[%v][version] = \"%v\"", Project, ProjectVersion))
				headerLines = append(headerLines, fmt.Sprintf("projects[%v][type] = \"module\"", Project))
				headerLines = append(headerLines, fmt.Sprintf("projects[%v][subdir] = \"%v\"", Project, ProjectType))
				headerLines = append(headerLines, fmt.Sprint(""))
			}
		}
	}

	// Print to path File

	newFile, _ := os.Create(File)
	for _, line := range headerLines {
		fmt.Fprintln(newFile, line)
	}
	newFile.Sync()
	defer newFile.Close()

}
