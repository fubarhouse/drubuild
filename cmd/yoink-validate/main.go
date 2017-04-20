package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/golang-drush/alias"
	"github.com/fubarhouse/golang-drush/make"
	"github.com/fubarhouse/golang-drush/vhost"
	"os"
	"os/user"
)

func main() {

	var WebserverDir = flag.String("vhost-dir", "/etc/nginx/sites-enabled", "Directory containing virtual host file(s)")
	var Webserver = flag.String("webserver-name", "nginx", "The name of the web service on the server.")

	flag.Parse()

	log.Println("Instanciating Alias")
	Alias := alias.NewAlias("temporaryAlias", "/tmp", "temporaryAlias")
	log.Println("Checking folder for Alias")
	usr, _ := user.Current()
	filedir := usr.HomeDir + "/.drush"
	_, statErr := os.Stat(filedir)
	if statErr != nil {
		log.Println("Could not find", filedir)
	} else {
		log.Println("Found", filedir)
	}
	log.Println("Installing Alias")
	Alias.Install()
	log.Println("Uninstalling Alias")
	Alias.Uninstall()

	log.Println("Instanciating Vhost")
	VirtualHost := vhost.NewVirtualHost("temporaryVhost", "/tmp", *Webserver, "temporary.vhost", *WebserverDir)
	log.Println("Checking folder for Vhost")
	_, statErr = os.Stat(*WebserverDir)
	if statErr != nil {
		log.Println("Could not find", *WebserverDir)
	} else {
		log.Println("Found", *WebserverDir)
	}
	log.Println("Installing Vhost")
	VirtualHost.Install()
	log.Println("Uninstalling Vhost")
	VirtualHost.Uninstall()

	//log.Println("Instanciating Solr core")
	//SolrCore := make.SolrCore{"http://localhost:8983", "blah", "/acquia/scripts/conf", "/var/solr"}
	//log.Println("Installing Solr core")
	//SolrCore.Install()
	//log.Println("Uninstalling Solr core")
	//SolrCore.Uninstall()
}
