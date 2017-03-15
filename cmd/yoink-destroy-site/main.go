package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/golang-drush/make"
	"os"
)

func main() {

	var Path = flag.String("path", "", "Path to site")
	var Site = flag.String("site", "", "Shortname of site")
	var Domain = flag.String("domain", "", "Domain of site")
	var Alias = flag.String("alias", "", "Alias of site")
	var VHostDir = flag.String("vhost-dir", "/etc/nginx/sites-enabled", "Directory containing virtual host file(s)")
	var WebserverName = flag.String("webserver-name", "nginx", "The name of the web service on the server.")

	// Usage:
	// -path="/path/to/site" \
	// -site="mysite" \
	// -domain="mysite.dev" \
	// -alias="mysite.dev" \

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

	if *Site == "" || *Alias == "" || *Path == "" || *Domain == "" {
		flag.Usage()
		os.Exit(1)
	}

	x := make.NewSite("none", *Site, *Path, *Alias, *WebserverName, *Domain, *VHostDir)
	y := make.NewmakeDB("127.0.0.1", "root", "root", 3306)
	x.DatabaseSet(y)
	x.ActionDestroy()
}
