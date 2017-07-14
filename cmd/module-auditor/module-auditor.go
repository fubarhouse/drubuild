package main

// This command will accept any amount of Drush aliases and Drupal module names in
// a comma separated format (ie "p1,p2,p3") and find out if the input aliases are
// using the input modules from the specified input make files.

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/golang-drush/alias"
	"github.com/fubarhouse/golang-drush/aliases"
	"github.com/fubarhouse/golang-drush/command"
	"os/exec"
	"strings"
)

func processModules(modules []string, aliasList aliases.AliasList, boolVerbose bool) {
	for _, module := range modules {
		count := 0
		thisModule := strings.Replace(module, " ", "", -1)
		for _, value := range aliasList.GetAliasNames() {
			cmd := command.NewDrushCommand()
			cmd.SetAlias(value)
			cmd.SetCommand("pm-info " + thisModule + " --fields=status")
			output, outputErr := cmd.Run()
			if outputErr != nil {
				fmt.Printf("Error: (%v) %v\n", cmd.GetAlias(), outputErr)
			}
			if strings.Contains(string(output), "enabled") {
				if boolVerbose {
					log.Printf("%v installed on %v\n", thisModule, cmd.GetAlias())
				}
				count++
			} else if strings.Contains(string(output), "was not found") {
				cmdQ := command.NewDrushCommand()
				cmdQ.SetAlias(value)
				cmdQ.SetCommand("sql-query \"SELECT name from system where name = " + thisModule + "\"")
				outputQ, _ := cmdQ.Run()
				if strings.Contains(string(outputQ), thisModule) == true {
					if boolVerbose {
						log.Printf("Error: %v is installed and missing on %v", thisModule, cmd.GetAlias())
					} else {
						log.Printf("Error: %v is installed and missing on %v", thisModule, cmd.GetAlias())
					}
				} else {
					if boolVerbose {
						log.Printf("%v is not installed and missing from %v", thisModule, cmd.GetAlias())
					}
				}
			} else {
				if boolVerbose {
					log.Printf("%v not installed on %v\n", thisModule, cmd.GetAlias())
				}
			}
		}
		if boolVerbose {
			log.Printf("%v/%v: %v\n", count, aliasList.Count(), thisModule)
		} else {
			if count >= 0 {
				log.Printf("%v/%v: %v\n", count, aliasList.Count(), thisModule)
			}
		}
	}
}

func main() {
	var strAliases = flag.String("aliases", "", "alias1,alias2,alias3")
	var strModules = flag.String("modules", "", "views,features,admin_menu")
	var strMakefiles = flag.String("makes", "", "/path/to/make.make,/path/to/make-other.make")
	var strPattern = flag.String("pattern", "%v", "A pattern to cross-reference the list of aliases, where %v in this string represents the alias.")
	var boolVerbose = flag.Bool("verbose", false, "false")
	flag.Parse()

	// Trim each comma-separated entry.
	*strAliases = strings.Replace(*strAliases, "  ", " ",-1)
	*strAliases = strings.Replace(*strAliases, ", ", ",",-1)
	*strAliases = strings.Replace(*strAliases, " ,", ",",-1)

	*strModules = strings.Replace(*strModules, "  ", " ",-1)
	*strModules = strings.Replace(*strModules, ", ", ",",-1)
	*strModules = strings.Replace(*strModules, " ,", ",",-1)

	*strMakefiles = strings.Replace(*strMakefiles, "  ", " ",-1)
	*strMakefiles = strings.Replace(*strMakefiles, ", ", ",",-1)
	*strMakefiles = strings.Replace(*strMakefiles, " ,", ",",-1)


	var getModulesFromMake = false
	var projects []string
	var MakeProjects []string

	if *strPattern != "" || !strings.Contains(*strPattern, "%v") {
		log.Errorln("Invalid pattern, must contain '%v'.")
	}

	if *strMakefiles != "" {
		MakefileNames := strings.Split(*strMakefiles, ",")
		for _, Makefile := range MakefileNames {
			catCmd := "cat " + Makefile + " | grep projects | cut -d'[' -f2 | cut -d']' -f1 | uniq | sort"
			y, _ := exec.Command("sh", "-c", catCmd).Output()
			projects = strings.Split(string(y), "\n")
			for _, Project := range projects {
				MakeProjects = append(MakeProjects, Project)
			}
		}
		if len(MakeProjects) != 0 {
			getModulesFromMake = true
		}
	}

	if (*strAliases != "" && *strModules != "") || (*strAliases != "" && getModulesFromMake == true) {
		aliasList := aliases.NewAliasList()
		aliases := strings.Split(*strAliases, ",")
		modules := strings.Split(*strModules, ",")
		if len(MakeProjects) != 0 {
			modules = MakeProjects
		}
		for _, value := range aliases {
			thisAliasA := strings.Replace(value, "@", "", -1)
			thisAliasA = strings.Replace(value, " ", "", -1)
			thisAliasA = fmt.Sprintf("@%v", thisAliasA)
			thisAliasA = strings.Replace(*strPattern, "%v", thisAliasA, -1)
			thisAlias := alias.NewAlias("", "", thisAliasA)
			aliasList.Add(thisAlias)
		}
		processModules(modules, *aliasList, *boolVerbose)
	} else {
		flag.Usage()
	}
}
