package alias

import (
	"os"
	"os/user"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/drubuild/util/drush"
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
	c, _ := drush.Run([]string{"drush", "sa", "@"+alias})
	if strings.Contains(c, "Could not find the alias") {
		log.Warnln(c)
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

// GetStatus returns the installation status of an alias struct
func (Alias *Alias) GetStatus() bool {
	_, err := os.Stat(getHome() + "/.drush/" + Alias.GetURI() + ".alias.drushrc.php")
	if err != nil {
		CommandOut, _ := drush.Run([]string{"sa"})
		if strings.Contains(string(CommandOut), Alias.GetName()) {
			return true
		} else {
			return false
		}
	}
	return true
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
