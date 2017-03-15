package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/golang-drush/make"
	"os"
)

func main() {

	var Address = flag.String("address", "localhost:8983", "http address of solr installation where solr version < 5.")
	var Name = flag.String("name", "", "Name of core to create")
	var Path = flag.String("path", "/var/solr/data", "Path so Solr data folder")
	var Resources = flag.String("resources", "/opt/solr/example/files/conf", "Path to Solr resources for new cores")
	var Create = flag.Bool("create", false, "Add create flag to mark creation")
	var Delete = flag.Bool("delete", false, "Add delete flag to maek deletion")

	// Path will need to be addressed based upon the major solr version
	// Version 4: /var/solr
	// Version 5: /var/solr/data

	// Usage:
	// -local-alias="mysite.dev" \

	flag.Parse()

	if !*Create {
		log.Infoln("Create flag has not been switched on")
	}
	if !*Delete {
		log.Infoln("Delete flag has not been switched on")
	}
	if *Name == "" {
		log.Infoln("Name input is empty")
	}

	if !*Create || !*Delete {
		if *Name == "" {
			flag.Usage()
			os.Exit(1)
		}
	}

	SolrCore := make.SolrCore{*Address, *Name, *Resources, *Path}

	if *Delete == true {
		log.Infoln("Starting Solr core uninstallation task.")
		SolrCore.Uninstall()
	}
	if *Create == true {
		log.Infoln("Starting Solr core installation task.")
		SolrCore.Install()
	}
}
