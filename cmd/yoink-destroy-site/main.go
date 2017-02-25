package main

import (
	"flag"
	"github.com/fubarhouse/golang-drush/make"
	"os"
)

func main() {

	var Path = flag.String("path", "", "Path to site")
	var Site = flag.String("site", "", "Shortname of site")
	var Domain = flag.String("domain", "", "Domain of site")
	var Alias = flag.String("alias", "", "Alias of site")

	// Usage:
	// -path="/path/to/site" \
	// -site="mysite" \
	// -domain="mysite.dev" \
	// -alias="mysite.dev" \

	flag.Parse()

	if string(*Site) == "" || string(*Alias) == "" || string(*Path) == "" || string(*Domain) == "" {
		flag.Usage()
		os.Exit(1)
	}

	x := make.NewSite("none", string(*Site), string(*Path), string(*Alias), "nginx", string(*Domain), "/etc/nginx/sites-enabled")
	y := make.NewmakeDB("127.0.0.1", "root", "root", 3306)
	x.DatabaseSet(y)
	x.ActionDestroy()
}
