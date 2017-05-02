package alias

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"os/user"
	"strings"
)

// A struct for managing a single Drush Alias
type Alias struct {
	name string
	path string
	uri  string
}

// Instantiate an Alias struct
func NewAlias(name, path, alias string) *Alias {
	return &Alias{name, path, alias}
}

// Set the name field for an alias struct
func (Alias *Alias) SetName(value string) {
	Alias.name = value
}

// Get the name field for an alias struct
func (Alias *Alias) GetName() string {
	return Alias.name
}

// Set the uri field for an alias struct
func (Alias *Alias) SetUri(value string) {
	Alias.uri = value
}

// Get the uri field for an alias struct
func (Alias *Alias) GetUri() string {
	return Alias.uri
}

// Set the path field for an alias struct
func (Alias *Alias) SetPath(value string) {
	Alias.path = value
}

// Get the path field for an alias struct
func (Alias *Alias) GetPath() string {
	return Alias.path
}

// Install an alias from an alias struct
func (Alias *Alias) Install() {
	Root := Alias.GetPath()
	if strings.HasSuffix(Root, "latest") == true {
		Root = strings.TrimSuffix(Root, "latest")
	}
	if strings.HasSuffix(Root, ".") == true {
		Root = strings.TrimSuffix(Root, ".")
	}
	if strings.HasSuffix(Root, "_") == true {
		Root = strings.TrimSuffix(Root, "_")
	}
	Root = fmt.Sprintf("%v/%v.latest", Root, Alias.GetUri())

	data := map[string]string{
		"Name":   Alias.GetName(),
		"Root":   Root,
		"Alias":  Alias.GetUri(),
		"Domain": Alias.GetUri(),
	}
	usr, _ := user.Current()
	filedir := usr.HomeDir + "/.drush"
	filename := Alias.GetUri() + ".alias.drushrc.php"
	fullpath := filedir + "/" + filename

	buffer := []byte{60, 63, 112, 104, 112, 10, 36, 97, 108, 105, 97, 115, 101, 115, 91, 39, 65, 76, 73, 65, 83, 39, 93, 32, 61, 32, 97, 114, 114, 97, 121, 40, 10, 32, 32, 39, 114, 111, 111, 116, 39, 32, 61, 62, 32, 39, 82, 79, 79, 84, 39, 44, 10, 32, 32, 39, 117, 114, 105, 39, 32, 61, 62, 32, 39, 68, 79, 77, 65, 73, 78, 39, 44, 10, 32, 32, 39, 112, 97, 116, 104, 45, 97, 108, 105, 97, 115, 101, 115, 39, 32, 61, 62, 32, 97, 114, 114, 97, 121, 40, 10, 32, 32, 32, 32, 39, 37, 102, 105, 108, 101, 115, 39, 32, 61, 62, 32, 39, 115, 105, 116, 101, 115, 47, 78, 65, 77, 69, 47, 102, 105, 108, 101, 115, 39, 44, 10, 32, 32, 32, 32, 39, 37, 112, 114, 105, 118, 97, 116, 101, 39, 32, 61, 62, 32, 39, 115, 105, 116, 101, 115, 47, 78, 65, 77, 69, 47, 112, 114, 105, 118, 97, 116, 101, 39, 44, 10, 32, 32, 41, 44, 10, 41, 59, 10, 63, 62}
	tpl := fmt.Sprintf("%v", string(buffer[:]))
	tpl = strings.Replace(tpl, "NAME", data["Name"], -1)
	tpl = strings.Replace(tpl, "ROOT", data["Root"], -1)
	tpl = strings.Replace(tpl, "ALIAS", data["Alias"], -1)
	tpl = strings.Replace(tpl, "DOMAIN", data["Domain"], -1)

	_, statErr := os.Stat(fullpath)
	if statErr != nil {
		nf, err := os.Create(fullpath)
		if err != nil {
			log.Fatalln("Error creating file", err)
		}
		_, err = nf.WriteString(tpl)
		if err != nil {
			log.Warnln("Could not add alias", fullpath)
		} else {
			log.Infoln("Added alias", filename)
		}
		defer nf.Close()
	} else {
		log.Warnln("Alias already created")
	}
}

// Un-install an alias from an alias struct
func (Alias *Alias) Uninstall() {
	usr, _ := user.Current()
	filedir := usr.HomeDir + "/.drush"
	filename := Alias.GetUri() + ".alias.drushrc.php"
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

// Re-install an alias from an alias struct
func (Alias *Alias) Reinstall() {
	Alias.Uninstall()
	Alias.Install()

}

// Return the installation status of an alias struct
func (Alias *Alias) GetStatus() bool {
	_, err := os.Stat(getHome() + "/.drush/" + Alias.GetUri() + ".alias.drushrc.php")
	if err != nil {
		return false
	} else {
		return true
	}
}

// Print the installation status of an alias struct
func (Alias *Alias) PrintStatus() {
	_, err := os.Stat(getHome() + "/.drush/" + Alias.GetUri() + ".alias.drushrc.php")
	if err != nil {
		fmt.Println("false")
	} else {
		fmt.Println("true")
	}
}

// Local function to return the user home directory.
// Performs some validation in the process.
func getHome() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	return usr.HomeDir
}
