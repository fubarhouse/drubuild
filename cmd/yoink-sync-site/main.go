package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/golang-drush/make"
	"os"
)

func main() {

	var LocalAlias = flag.String("local-alias", "", "Alias of target site")
	var RemoteAlias = flag.String("remote-alias", "", "Alias of source site")
	var SyncDB = flag.Bool("db", false, "Mark database for syncronization")
	var SyncFiles = flag.Bool("files", false, "Mark files for syncronization")

	// Usage:
	// -local-alias="mysite.dev" \
	// -remote-alias="mysite.dev" \

	flag.Parse()

	if *LocalAlias == "" || *RemoteAlias == "" {
		flag.Usage()
		os.Exit(1)
	}

	x := make.Site{}
	x.Alias = *LocalAlias
	if *SyncDB == true {
		log.Infoln("Database was marked for syncing, working now...")
		x.ActionDatabaseSyncLocal(*RemoteAlias)
	}
	if *SyncFiles == true {
		log.Infoln("Files were marked for syncing, working now...")
		x.ActionFilesSyncLocal(*RemoteAlias)
	}
	if *SyncDB == true || *SyncFiles == true {
		log.Infoln("Attempting to rebuild registries...")
		x.RebuildRegistry()
	}
}
