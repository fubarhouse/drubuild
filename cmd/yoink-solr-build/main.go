package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/golang-drush/make"
	"os"
)

func main() {

	// TODO Finish this for Solr 3,4,5,6.

	var Address = flag.String("address", "localhost:8983", "http address of solr installation where solr version < 5.")
	var Name = flag.String("name", "", "Name of core to create")
	var Path = flag.String("path", "/var/solr/data", "Path so Solr data folder")
	var Resources = flag.String("resources", "/opt/solr/example/files/conf", "Path to Solr resources for new cores")

	// Path will need to be addressed based upon the major solr version
	// Version 4: /var/solr
	// Version 5: /var/solr/data

	// Usage:

	flag.Parse()

	if *Name == "" {
		log.Infoln("Name input is empty")
	}

	SolrCore := make.SolrCore{*Address, *Name, *Resources, *Path}
	log.Infoln("Starting Solr core installation task.")
	SolrCore.Install()
}
