package main

import (
	"flag"
	"github.com/fubarhouse/golang-drush/make"
	"os"
	"strings"
)

func main() {

	var Path = flag.String("path", "", "Path to site")
	var Makes = flag.String("makes", "", "Comma-separated list of make files to use")

	// Usage:
	// -path="/path/to/site" \
	// -site="mysite" \
	// -domain="mysite.dev" \
	// -alias="mysite.dev" \
	// -makes="/path/to/make1.make, /path/to/make2.make" \

	flag.Parse()

	if *Makes == "" || *Path == "" {
		flag.Usage()
		os.Exit(1)
	}

	x := make.Site{}
	x.TimeStampSet("")
	x.Name = ""
	x.Path = *Path

	MakefilesFormatted := strings.Replace(*Makes, " ", "", -1)
	MakeFiles := strings.Split(MakefilesFormatted, ",")

	x.ActionRebuildCodebase(MakeFiles)
}
