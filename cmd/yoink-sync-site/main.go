package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/golang-drush/command"
	"os"
	"strings"
)

func main() {

	var SourceAlias = flag.String("source-alias", "", "Alias of source site")
	var DestAlias = flag.String("dest-alias", "", "Alias of destination site")
	var Forbidden = flag.String("forbid", "", "For automation/security purposes, do not allow destination aliases to contain this string.")
	var SyncDB = flag.Bool("db", false, "Mark database for synchronization")
	var SyncFiles = flag.Bool("files", false, "Mark files for synchronization")
	var FilepathVerification = flag.Bool("verify-files", false, "Boolean which tells yoink-sync-site to run drush vsets for file system variables, as a verification step.")
	var FilepathPublic = flag.String("public-files", "files", "Path under site directory to create public files directory.")
	var FilepathPrivate = flag.String("private-files", "files/private", "Path under site directory to create private files directory.")
	var FilepathTemporary = flag.String("temp-files", "files/private/temp", "Path under site directory to create temporary files directory.")

	// Usage:
	// -local-alias="mysite.dev" \
	// -remote-alias="mysite.dev" \
	// -db \
	// -files

	flag.Parse()

	if *SourceAlias == "" {
		log.Infoln("Source input is empty")
	}
	if *DestAlias == "" {
		log.Infoln("Destination input is empty")
	}
	if !*SyncDB {
		log.Infoln("Database flag is switched off")
	}
	if !*SyncFiles {
		log.Infoln("Files flag is switched off")
	}

	if *SourceAlias == "" || *DestAlias == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *SyncDB {
		log.Infoln("Database was marked for syncing, working now...")
		if !strings.Contains(*DestAlias, *Forbidden) {
			command.DrushDatabaseSync(*SourceAlias, *DestAlias)
		} else {
			log.Errorln("Destination alias contains a forbidden string")
		}
	}
	if *SyncFiles {
		log.Infoln("Files were marked for syncing, working now...")
		if !strings.Contains(*DestAlias, *Forbidden) {
			command.DrushFilesSync(*SourceAlias, *DestAlias)
		} else {
			log.Errorln("Destination alias contains a forbidden string")
		}
	}
	if *FilepathVerification {
		command.DrushVariableSet(*DestAlias, "file_public_path", *FilepathPublic)
		command.DrushVariableSet(*DestAlias, "file_private_path", *FilepathPrivate)
		command.DrushVariableSet(*DestAlias, "file_temporary_path", *FilepathTemporary)
	}
	if *SyncDB || *SyncFiles {
		log.Infoln("Attempting to rebuild registries...")
		command.DrushRebuildRegistry(*DestAlias)
	}
}
