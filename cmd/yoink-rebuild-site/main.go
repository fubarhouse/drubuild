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
	var BuildID = flag.String("build", "", "optional timestamp of site")

	// Usage:
	// -path="/path/to/site" \
	// -makes="/path/to/make1.make, /path/to/make2.make" \
	// -buildid="20170101"

	flag.Parse()

	if string(*Makes) == "" || string(*Path) == "" {
		flag.Usage()
		os.Exit(1)
	}

	x := make.Site{}
	x.Make = *Makes
	MakefilesFormatted := strings.Replace(x.Make, " ", "", -1)

	x.Name = ""
	x.Path = *Path
	x.Make = MakefilesFormatted

	if string(*BuildID) == "" {
		x.TimeStampSet("")
	} else {
		x.TimeStampSet(string(*BuildID))
	}

	MakeFiles := strings.Split(MakefilesFormatted, ",")
	x.ActionRebuildCodebase(MakeFiles)
}
