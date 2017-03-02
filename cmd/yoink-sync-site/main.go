package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/golang-drush/command"
	"os"
)

func main() {

	var SourceAlias = flag.String("source-alias", "", "Alias of target site")
	var DestAlias = flag.String("dest-alias", "", "Alias of source site")
	var SyncDB = flag.Bool("db", false, "Mark database for syncronization")
	var SyncFiles = flag.Bool("files", false, "Mark files for syncronization")

	// Usage:
	// -local-alias="mysite.dev" \
	// -remote-alias="mysite.dev" \

	flag.Parse()

	if *SourceAlias == "" || *DestAlias == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *SyncDB == true {
		log.Infoln("Database was marked for syncing, working now...")
		command.DrushDatabaseSync(*SourceAlias, *DestAlias)
	}
	if *SyncFiles == true {
		log.Infoln("Files were marked for syncing, working now...")
		command.DrushFilesSync(*SourceAlias, *DestAlias)
	}
	if *SyncDB == true || *SyncFiles == true {
		log.Infoln("Attempting to rebuild registries...")
		command.DrushRebuildRegistry(*DestAlias)
	}
}
