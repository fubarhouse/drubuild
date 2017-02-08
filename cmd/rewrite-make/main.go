package main

import (
	"flag"
	"github.com/fubarhouse/golang-drush/makeupdater"
)

func main() {
	var strMake = flag.String("make", "", "Absolute path to make file")
	flag.Parse()
	if *strMake != "" {
		Projects := makeupdater.GetProjectsFromMake(*strMake)
		makeupdater.GenerateMake(Projects, *strMake)
	}
}
