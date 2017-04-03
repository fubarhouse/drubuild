package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/golang-drush/make"
	"os"
)

func main() {

	var Address = flag.String("address", "http://localhost:8983", "http address of solr installation where solr version < 5.")
	var Name = flag.String("name", "", "Name of core to create")
	var Path = flag.String("path", "/var/solr", "Path to Solr data folder")
	var Resources = flag.String("resources", "", "Path to Solr resources for new cores")
	var DataPath = flag.String("data-path", "", "Path to data directory inside the Solr directory")

	flag.Parse()

	if *Name == "" {
		log.Infoln("Name input is empty")
	}

	if *Resources == "" {
		log.Infoln("Resources input is empty")
	}

	if *Name == "" || *Resources == "" {
		flag.Usage()
		os.Exit(1)
	}

	SolrCore := make.SolrCore{*Address, *Name, *Resources, *Path, *DataPath}
	log.Infoln("Starting Solr core installation task.")
	SolrCore.Install()
}
