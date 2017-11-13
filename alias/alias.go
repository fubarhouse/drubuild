package alias

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"text/template"

	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

// Alias is a struct for managing a single Drush Alias
type Alias struct {
	Docroot	 string
	name     string
	path     string
	uri      string
	template string
}

func (a *Alias) Template() string {
	return a.template
}

func (a *Alias) SetTemplate(template string) {
	a.template = template
}

// NewAlias instantiates an Alias struct
func NewAlias(name, path, alias string) *Alias {
	alias = strings.Replace(alias, "@", "", -1)
	Command := exec.Command("drush", "sa", "@"+alias)
	CommandOut, _ := Command.CombinedOutput()
	if strings.Contains(string(CommandOut), "Could not find the alias") {
		log.Warnln(string(CommandOut))
		return &Alias{}
	} else {
		return &Alias{"", name, path, alias, ""}
	}
}

// SetName sets the name field for an alias struct
func (Alias *Alias) SetName(value string) {
	Alias.name = value
}

// GetName gets the name field for an alias struct
func (Alias *Alias) GetName() string {
	return Alias.name
}

// SetURI sets the uri field for an alias struct
func (Alias *Alias) SetURI(value string) {
	Alias.uri = value
}

// GetURI gets the uri field for an alias struct
func (Alias *Alias) GetURI() string {
	return Alias.uri
}

// SetPath sets the path field for an alias struct
func (Alias *Alias) SetPath(value string) {
	Alias.path = value
}

// GetPath gets the path field for an alias struct
func (Alias *Alias) GetPath() string {
	return Alias.path
}

// Install an alias from an alias struct
func (Alias *Alias) Install() {
	Root := fmt.Sprintf("%v/%v.latest/%v/", Alias.GetPath(), Alias.GetURI(), Alias.Docroot)
	data := map[string]string{
		"Name":   Alias.GetName(),
		"Root":   Root,
		"Alias":  Alias.GetURI(),
		"Domain": Alias.GetURI(),
	}
	usr, _ := user.Current()
	filedir := usr.HomeDir + "/.drush"
	filename := Alias.GetURI() + ".alias.drushrc.php"
	fullpath := filedir + "/" + filename

	defaultTemplate := ""
	if Alias.template == "" {
		defaultTemplate = fmt.Sprintf("%v/src/github.com/fubarhouse/golang-drush/cmd/yoink/templates/alias.gotpl", os.Getenv("GOPATH"))
	} else {
		defaultTemplate = Alias.template
	}

	t := template.New("alias")
	if _, err := os.Stat(defaultTemplate); err != nil {
		log.Warnln("default drush alias template could not be found, source files do not exist.")
	} else {
		log.Infof("Found template %v for usage", defaultTemplate)
		defaultData, _ := ioutil.ReadFile(defaultTemplate)
		t.Parse(string(defaultData))
	}

	os.Remove(fullpath)
	file, _ := os.Create(fullpath)
	tplErr := t.Execute(file, data)

	if tplErr == nil {
		log.Infof("Successfully templated alias to file %v", fullpath)
	} else {
		log.Warnf("Error templating alias to file %v", fullpath)
	}
}

// Uninstall un-installs an alias from an alias struct
func (Alias *Alias) Uninstall() {
	usr, _ := user.Current()
	filedir := usr.HomeDir + "/.drush"
	filename := Alias.GetURI() + ".alias.drushrc.php"
	fullpath := filedir + "/" + filename
	_, statErr := os.Stat(fullpath)
	if statErr == nil {
		err := os.Remove(fullpath)
		if err != nil {
			log.Warnln("Could not remove alias file", fullpath)
		} else {
			log.Infoln("Removed alias file", fullpath)
		}
	} else {
		log.Warnln("Alias file was not found.", fullpath)
	}

}

// Reinstall re-installs an alias from an alias struct
func (Alias *Alias) Reinstall() {
	Alias.Uninstall()
	Alias.Install()

}

// GetStatus returns the installation status of an alias struct
func (Alias *Alias) GetStatus() bool {
	_, err := os.Stat(getHome() + "/.drush/" + Alias.GetURI() + ".alias.drushrc.php")
	if err != nil {
		Command := exec.Command("drush", "sa")
		CommandOut, _ := Command.CombinedOutput()
		if strings.Contains(string(CommandOut), Alias.GetName()) {
			return true
		} else {
			return false
		}
	}
	return true
}

// PrintStatus prints the installation status of an alias struct
func (Alias *Alias) PrintStatus() {
	_, err := os.Stat(getHome() + "/.drush/" + Alias.GetURI() + ".alias.drushrc.php")
	if err != nil {
		fmt.Println("false")
	} else {
		fmt.Println("true")
	}
}

// getHome returns the user home directory.
// Performs some validation in the process.
func getHome() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	return usr.HomeDir
}
