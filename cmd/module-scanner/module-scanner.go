package main

// This command will accept any amount of Drush aliases and Drupal module names in
// a comma separated format (ie "p1,p2,p3") and find out if the input aliases are
// using the input modules, and it will return a count of the total of which are
// enabled.

import (
	"flag"
	"fmt"
	"github.com/fubarhouse/golang-drush/alias"
	"github.com/fubarhouse/golang-drush/aliases"
	"github.com/fubarhouse/golang-drush/command"
	"strings"
)

func main() {
	var strAliases = flag.String("aliases", "", "alias1,alias2,alias3")
	var strModules = flag.String("modules", "", "views,features,admin_menu")
	var boolVerbose = flag.Bool("verbose", false, "false")
	flag.Parse()
	if *strAliases != "" && *strModules != "" {
		aliasList := aliases.NewAliasList()
		aliases := strings.Split(*strAliases, ",")
		modules := strings.Split(*strModules, ",")
		for _, value := range aliases {
			thisAliasA := strings.Replace(value, "@", "", -1)
			thisAliasA = fmt.Sprintf("@%v", thisAliasA)
			thisAlias := alias.NewAlias("", "", thisAliasA)
			aliasList.Add(thisAlias)
		}
		for _, module := range modules {
			count := 0
			thisModule := module
			for _, value := range aliasList.GetAliasNames() {
				cmd := command.NewDrushCommand()
				cmd.SetAlias(value)
				cmd.SetCommand("pm-info " + thisModule + " --fields=status")
				output, outputErr := cmd.Run()
				if outputErr != nil {
					fmt.Printf("Error: (%v) %v\n", cmd.GetAlias(), outputErr)
				}
				if strings.Contains(string(output), "enabled") {
					if *boolVerbose {
						fmt.Printf("Found module %v on site %v\n", thisModule, cmd.GetAlias())
					}
					count++
				} else {
					if *boolVerbose {
						fmt.Printf("Did not find module %v on site %v\n", thisModule, cmd.GetAlias())
					}
				}
			}
			fmt.Printf("Out of the %v tested sites, %v have the module %v installed.\n", aliasList.Count(), count, thisModule)
		}
	} else {
		flag.Usage()
	}
}
