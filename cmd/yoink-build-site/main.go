package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
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
	var CustomTemplate = flag.String("template", "", "Absolute path to a custom template, which falls back to a given default.")
	var FilepathPublic = flag.String("public-files", "files", "Path under site directory to create public files directory.")
	var FilepathPrivate = flag.String("private-files", "files/private", "Path under site directory to create private files directory.")
	var FilepathTemporary = flag.String("temp-files", "files/private/temp", "Path under site directory to create temporary files directory.")
	var RewriteStringSource = flag.String("rewrite-source", "", "A string of text to replace in the make file before building.")
	var RewriteStringDestination = flag.String("rewrite-dest", "", "A string of text to replace the rewrite-source value with before building.")
	var WorkingCopy = flag.Bool("working-copy", false, "Apply --working-copy to to drush during any make processes.")

	// Templates will replace the following strings in the provided template with the values inputted to the program.
	// Failing to provide this file path, it will fall-back to a template stored inside the go code.
	// "Name": Name of the site
	// "Domain": Configured domain of the site

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

	x := make.NewSite(*Makes, *Site, *Path, *Alias, *WebserverName, *Domain, *VHostDir, *CustomTemplate)
	y := make.NewmakeDB("127.0.0.1", "root", "root", 3306)
	x.DatabaseSet(y)
	if *BuildID == "" {
		x.TimeStampReset()
	} else {
		x.TimeStampSet(*BuildID)
	}
	if *WorkingCopy {
		x.WorkingCopy = true
	}

	MakefilesFormatted := strings.Replace(*Makes, " ", "", -1)
	MakeFiles := strings.Split(MakefilesFormatted, ",")

	if *RewriteStringSource != "" && *RewriteStringDestination != "" {
		x.MakeFileRewriteSource = *RewriteStringSource
		x.MakeFileRewriteDestination = *RewriteStringDestination
	}
	x.ActionRebuildCodebase(MakeFiles)
	x.InstallSiteRef()
	x.InstallFileSystem(*FilepathPublic)
	x.InstallFileSystem(*FilepathPrivate)
	x.InstallFileSystem(*FilepathTemporary)
	x.ActionInstall()
	x.SymReinstall()
	x.VhostInstall()
	x.AliasInstall()
	x.StopWebServer()
	x.StartWebServer()
}
