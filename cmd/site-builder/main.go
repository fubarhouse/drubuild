package main

import (
	"flag"
	"fmt"
	"github.com/fubarhouse/golang-drush/aliases"
	"github.com/fubarhouse/golang-drush/make"
	"os"
	"strings"
)

func main() {

	var Path = flag.String("path", "", "Path to site")
	var Site = flag.String("site", "", "Shortname of site")
	var Environment = flag.String("environment", "", "Name of environment (acts as Drush alias filter initially)")
	var Filter = flag.String("filter", "", "Filter the available drush aliases with this sting to find a concise list of sites")
	var Makes = flag.String("makes", "", "Comma-separated list of make files to use")
	var Action = flag.String("action", "", "action to perform (build|destroy)")
	var BuildID = flag.String("build", "", "optional timestamp of site")

	// Usage:
	// -path="/path/to/site" -site="mysite" -filter="mysite.dev" -makes="/path/to/make1.make, /path/to/make2.make" -action="build"

	flag.Parse()

	if string(*Site) == "" || string(*Action) == "" {
		flag.Usage()
		os.Exit(1)
	}

	Aliases := aliases.NewAliasList()
	Aliases.Generate(string(*Filter))
	Aliases.Filter(string(*Environment))

	Sites := []string{}
	for _, name := range Aliases.GetNames() {
		name = strings.Replace(name, "."+string(*Environment), "", -1)
		Sites = append(Sites, name)
	}

	x := make.NewSite(string(*Makes), string(*Site), string(*Path), string(*Site), "nginx", "/etc/nginx/sites-enabled")
	y := make.NewmakeDB("127.0.0.1", "root", "root", 3306)
	x.DatabaseSet(y)
	if string(*BuildID) == "" {
		x.TimeStampReset()
	} else {
		x.TimeStampSet(string(*BuildID))
	}

	MakefilesFormatted := strings.Replace(x.Make, " ", "", -1)
	MakeFiles := strings.Split(MakefilesFormatted, ",")

	if string(*Action) == "build" {
		x.ActionRebuildCodebase(MakeFiles)
		x.InstallSiteRef()
		x.ActionInstall()
		x.VhostInstall()
		x.AliasInstall()
		x.ActionDatabaseSyncLocal(fmt.Sprintf("@%v", string(*Site))) // Needs work!
		x.RebuildRegistry()
		x.SymReinstall(x.TimeStampGet())
		x.RestartWebServer()
	} else if string(*Action) == "destroy" {
		x.ActionDestroy()
	}
}
