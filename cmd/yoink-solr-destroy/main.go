package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/golang-drush/make"
	"os"
)

func main() {

	// TODO Finish this for Solr 3,4,5,6.

	var Address = flag.String("address", "http://localhost:8983", "http address of solr installation where solr version < 5.")
	var Name = flag.String("name", "", "Name of core to create")
	var Path = flag.String("path", "/var/solr", "Path to Solr data folder")

	flag.Parse()

	if *Name == "" {
		log.Infoln("Name input is empty")
		flag.Usage()
		os.Exit(1)
	}

	SolrCore := make.SolrCore{*Address, *Name, "", *Path, false}
	log.Infoln("Starting Solr core uninstallation task.")
	SolrCore.Uninstall()
}
