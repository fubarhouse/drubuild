package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/golang-drush/makeupdater"
	"strings"
)

func main() {
	var strMake = flag.String("makes", "", "Comma-separated list of absolute paths to make files to update.")
	flag.Parse()

	// Remove double spaces.
	*strMake = strings.Replace(*strMake, "  ", " ", -1)

	if *strMake != "" {
		Makes := strings.Split(*strMake, ",")
		for _, Makefile := range Makes {
			Makefile = strings.Trim(Makefile, " ")
			makeupdater.UpdateMake(Makefile)
			makeupdater.FindDuplicatesInMake(Makefile)
		}
	} else {
		log.Infoln("Invalid make file input")
		flag.Usage()
	}
}
