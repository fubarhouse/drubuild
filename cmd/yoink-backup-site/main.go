package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/golang-drush/make"
	"os"
)

// Deprecated: use Yoink instead.
func main() {

	var Alias = flag.String("alias", "", "Alias of site")
	var Destination = flag.String("destination", "/tmp/drush-ard.tar", "Destination file for backup.")

	flag.Parse()

	if *Alias == "" {
		log.Infoln("Alias input is empty")
		flag.Usage()
		os.Exit(1)
	}

	log.Warnln("This binary has been deprecated in favor of `yoink`.")
	x := make.NewSite("", "", "", *Alias, "", "", "", "")
	x.ActionBackup(*Destination)
}
