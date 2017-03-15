package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/golang-drush/command"
	"github.com/fubarhouse/golang-drush/make"
	"os"
	"strings"
)

func main() {

	var Path = flag.String("path", "", "Path to site")
	var Site = flag.String("site", "", "Shortname of site")
	var Domain = flag.String("domain", "", "Domain of site")
	var Alias = flag.String("alias", "", "Alias of site")
	var Makes = flag.String("makes", "", "Comma-separated list of make files to use")
	var BuildID = flag.String("build", "", "optional timestamp of site")
	var VHostDir = flag.String("vhost-dir", "/etc/nginx/sites-enabled", "Directory containing virtual host file(s)")
	var WebserverName = flag.String("webserver-name", "nginx", "The name of the web service on the server.")

	// Usage:
	// -path="/path/to/site" \
	// -site="mysite" \
	// -domain="mysite.dev" \
	// -alias="mysite.dev" \
	// -makes="/path/to/make1.make, /path/to/make2.make" \

	flag.Parse()

	if *Site == "" {
		log.Infoln("Site input is empty")
	}
	if *Alias == "" {
		log.Infoln("Alias input is empty")
	}
	if *Path == "" {
		log.Infoln("Path input is empty")
	}
	if *Domain == "" {
		log.Infoln("Domain input is empty")
	}
	if *Makes == "" {
		log.Infoln("Makes input is empty")
	}

	if *Site == "" || *Alias == "" || *Makes == "" || *Path == "" || *Domain == "" {
		flag.Usage()
		os.Exit(1)
	}

	x := make.NewSite(*Makes, *Site, *Path, *Alias, *WebserverName, *Domain, *VHostDir)
	y := make.NewmakeDB("127.0.0.1", "root", "root", 3306)
	x.DatabaseSet(y)
	if *BuildID == "" {
		x.TimeStampReset()
	} else {
		x.TimeStampSet(*BuildID)
	}

	MakefilesFormatted := strings.Replace(*Makes, " ", "", -1)
	MakeFiles := strings.Split(MakefilesFormatted, ",")

	x.ActionRebuildCodebase(MakeFiles)
	x.InstallSiteRef()
	x.InstallPrivateFileSystem()
	x.ActionInstall()
	x.SymReinstall(x.TimeStampGet())
	x.VhostInstall()
	x.AliasInstall()
	command.DrushUpdateDatabase(x.Alias)
	command.DrushRebuildRegistry(x.Alias)
	x.RestartWebServer()
}
