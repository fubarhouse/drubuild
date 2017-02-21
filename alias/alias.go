package alias

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
	"text/template"
)

type Alias struct {
	name string
	path string
	uri  string
}

func NewAlias(name, path, alias string) *Alias {
	return &Alias{name, path, alias}
}

func (Alias *Alias) SetName(value string) {
	Alias.name = value
}

func (Alias *Alias) GetName() string {
	return Alias.name
}

func (Alias *Alias) SetUri(value string) {
	Alias.uri = value
}

func (Alias *Alias) GetUri() string {
	return Alias.uri
}

func (Alias *Alias) SetPath(value string) {
	Alias.path = value
}

func (Alias *Alias) GetPath() string {
	return Alias.path
}

func (Alias *Alias) Install() {
	log.Println("Adding alias", Alias.uri)
	data := map[string]string{
		"Name":  Alias.GetName(),
		"Root":  strings.Replace(Alias.GetPath(), "_", ".", -1),
		"Alias": Alias.GetUri(),
	}
	usr, _ := user.Current()
	filename := usr.HomeDir + "/.drush/" + Alias.uri + ".alias.drushrc.php"
	tpl, err := template.ParseFiles("templates/alias-template.gotpl")
	if err != nil {
		log.Fatalln(err)
	}

	nf, err := os.Create(filename)
	if err != nil {
		log.Fatalln("error creating file", err)
	}
	defer nf.Close()

	err = tpl.Execute(nf, data)
	if err != nil {
		log.Fatalln(err)
	}
}

func (Alias *Alias) Uninstall() {
	log.Println("Removing alias", Alias.uri)
	os.Remove(getHome() + "/.drush/" + Alias.GetUri() + ".alias.drushrc.php")

}

func (Alias *Alias) Reinstall() {
	Alias.Uninstall()
	Alias.Install()

}

func (Alias *Alias) GetStatus() bool {
	_, err := os.Stat(getHome() + "/.drush/" + Alias.GetUri() + ".alias.drushrc.php")
	if err != nil {
		return false
	} else {
		return true
	}
}

func (Alias *Alias) PrintStatus() {
	_, err := os.Stat(getHome() + "/.drush/" + Alias.GetUri() + ".alias.drushrc.php")
	if err != nil {
		fmt.Println("false")
	} else {
		fmt.Println("true")
	}
}

func getHome() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}
