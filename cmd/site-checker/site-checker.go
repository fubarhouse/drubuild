/*
	Site-checker runs a slice of commands on a slice of Drush aliases matching a specified pattern.
*/
package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/golang-drush/command"
	"os"
	"strings"
)

func main() {
	var strAliases = flag.String("aliases", "", "comma-separated list of aliases for action")
	var strCommands = flag.String("commands", "rr,updb --yes,cc all", "comma-separated list of commands for action")
	var strPattern = flag.String("pattern", "%v", "A modifier which allows rewriting of aliases replacing '%v' in the pattern with the alias.")
	var boolVerbose = flag.Bool("verbose", false, "adds raw output to end of program.")
	flag.Parse()

	var FinalOutput []string

	if !strings.Contains(*strPattern, "%v") {
		log.Errorln("Specified pattern does not include alias modifier.")
	}

	if *strCommands == "" {
		log.Errorln("Commands are not specified.")
	}

	if *strAliases == "" {
		flag.Usage()
		log.Errorln("Aliases are not specified.")
	}

	if *strAliases != "" && *strCommands != "" {
		for _, Alias := range strings.Split(*strAliases, ",") {
			Alias = strings.Replace(*strPattern, "%v", Alias, 1)
			for _, Command := range strings.Split(*strCommands, ",") {
				FinalOutput = append(FinalOutput, fmt.Sprintf("\n\ndrush @%v %v\n", Alias, Command))
				DrushCommand := command.NewDrushCommand()
				DrushCommand.SetAlias(Alias)
				DrushCommand.SetCommand(Command)
				DrushCommandOut, DrushCommandError := DrushCommand.Output()
				if DrushCommandError != nil {
					log.Warnf("%v, %v, unsuccessful.", DrushCommand.GetAlias(), DrushCommand.GetCommand())
					StdOutLines := DrushCommandOut
					for _, StdOutLine := range StdOutLines {
						FinalOutput = append(FinalOutput, StdOutLine)
					}
				} else {
					log.Infof("%v, %v, successful.", DrushCommand.GetAlias(), DrushCommand.GetCommand())
					StdOutLines := DrushCommandOut
					for _, StdOutLine := range StdOutLines {
						FinalOutput = append(FinalOutput, StdOutLine)
					}
				}
			}
		}
		if *boolVerbose {
			for _, value := range FinalOutput {
				log.Println(value)
			}
		}
	} else {
		flag.Usage()
		os.Exit(1)
	}
}
