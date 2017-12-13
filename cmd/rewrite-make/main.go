package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/drubuild/makeupdater"
)

func main() {
	var strMake = flag.String("make", "", "Absolute path to make file")
	flag.Parse()
	if *strMake != "" {
		Projects := makeupdater.GetProjectsFromMake(*strMake)
		makeupdater.GenerateMake(Projects, *strMake)
	} else {
		log.Infoln("Invalid input make, must be of type string and != ''")
	}
}
