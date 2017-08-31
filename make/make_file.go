package make

import (
	"bytes"
	//"github.com/ghodss/yaml"
	"io/ioutil"
	"os/exec"
	"reflect"
	"strings"
	"unicode"
)

type Makefile interface {
	ParseJSON() ([]byte, error) // Parse JSON data for processing
	ParseYML() ([]byte, error)  // Parse a Drupal 8 make file for processing
	ParseINF() ([]byte, error)  // Parse a Drupal 7 make file for processing
}

type Make struct {
	Path string
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

// hasPrefix will return true if the first non-whitespace bytes in buf is prefix.
func hasPrefix(buf []byte, prefix []byte) bool {
	trim := bytes.TrimLeftFunc(buf, unicode.IsSpace)
	return bytes.HasPrefix(trim, prefix)
}

// hasJSONPrefix returns true if the provided buffer
// appears to start with a JSON open brace.
func hasJSONPrefix(buf []byte) bool {
	return hasPrefix(buf, []byte("{"))
}

// ParseJSON will parse a []string into JSON.
// Implementation is not yet complete.
func (Make *Make) ParseJSON() ([]string, error) {
	return []string{}, nil
}

// modifyProject is a mechanic which adds/modifies a key/value
// in a parsed INF thing, which is yet to become a reality as
// it currently is under active development.
func modifyProject(name, key, value string, Projects []DrupalProject) {
	for _, Project := range Projects {
		if Project.Name == reflect.Indirect(reflect.ValueOf(DrupalProject{})).Type().Field(0).Name {
		}
	}
}

// ParseINF parses Drupal 7 make files into JSON.
// TODO: This needs urgent attention before it can be used..
func (Make *Make) parseINF() ([]string, error) {
	Projects := []string{}
	catCmd := "cat " + Make.Path
	Lines, _ := exec.Command("sh", "-c", catCmd).Output()

	for _, Line := range strings.Split(string(Lines), "\n") {
		if strings.HasPrefix(Line, "project") {
			// CLean up projects
			Line = strings.Replace(Line, "projects", "", 1)
			Line = strings.Replace(Line, "[", " ", -1)
			Line = strings.Replace(Line, "]", " ", -1)
			Line = strings.Replace(Line, "=", " ", -1)
			Line = strings.Replace(Line, "\"", " ", -1)
			Line = strings.Replace(Line, "  ", " ", -1)
			Line = strings.Replace(Line, "  ", " ", -1)
			Line = strings.Replace(Line, "  ", " ", -1)
			//modifyProject("drupal", "subdir", "contrib", string())
		} else if strings.HasPrefix(Line, "libraries") {
			Line = strings.Replace(Line, "libraries", "", 1)
			Line = strings.Replace(Line, "[", " ", -1)
			Line = strings.Replace(Line, "]", " ", -1)
			Line = strings.Replace(Line, "=", " ", -1)
			Line = strings.Replace(Line, "\"", " ", -1)
			Line = strings.Replace(Line, "  ", " ", -1)
			Line = strings.Replace(Line, "  ", " ", -1)
		}
	}

	//rawProjects := strings.Split(string(y), "\n")
	//for _, project := range rawProjects {
	//	project = strings.Replace(project, " ", "", -1)
	//	if project != "" && project != "projects" {
	//		if inArray(Projects, project) == 0 {
	//			Projects = append(Projects, project)
	//		}
	//	}
	//}

	return Projects, nil
}

// ParseYML parses Drupal 8 make files into JSON.
func (Make *Make) ParseYML() ([]byte, error) {
	data, err := ioutil.ReadFile(Make.Path)
	if err != nil {
		panic(err)
	}
	return data, err
	//parse_data, parse_error := yaml.YAMLToJSON(data)
	//return parse_data, parse_error
}
