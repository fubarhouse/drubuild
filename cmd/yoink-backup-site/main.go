package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/drubuild/make"
	"os"
)

func main() {

	var Alias = flag.String("alias", "", "Alias of site")
	var Destination = flag.String("destination", "/tmp/drush-ard.tar", "Destination file for backup.")

	flag.Parse()

	if *Alias == "" {
		log.Infoln("Alias input is empty")
		flag.Usage()
		os.Exit(1)
	}

	x := make.NewSite("", "", "", *Alias, "", "", "", "")
	x.ActionBackup(*Destination)
}
