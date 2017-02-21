package main

import (
	"flag"
	"github.com/fubarhouse/golang-drush/makeupdater"
)

func main() {
	// TODO: add error handler for input make file
	// TODO: remove the fact the program always outputs "already updated"
	var strMake = flag.String("make", "", "Absolute path to make file")
	flag.Parse()
	if *strMake != "" {
		makeupdater.UpdateMake(*strMake)
	} else {
		flag.Usage()
	}
}
