package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/golang-drush/command"
	"strings"
)

func main() {
	var strAliases = flag.String("aliases", "", "comma-separated list of aliases for action")
	var strCommands = flag.String("commands", "rr,updb --yes,cc all", "comma-separated list of commands for action")
	var boolVerbose = flag.Bool("verbose", false, "adds raw output to end of program.")
	flag.Parse()

	var FinalOutput []string

	if *strAliases != "" {
		for _, Alias := range strings.Split(*strAliases, ",") {

			for _, Command := range strings.Split(*strCommands, ",") {

				FinalOutput = append(FinalOutput, fmt.Sprintf("\n\ndrush @%v %v\n", Alias, Command))

				DrushCommand := command.NewDrushCommand()
				DrushCommand.SetAlias(Alias)
				DrushCommand.SetCommand(Command)
				DrushCommandOut, DrushCommandError := DrushCommand.Output()

				if DrushCommandError != nil {
					log.Warnf("%v, %v, unsuccessful.", DrushCommand.GetAlias(), DrushCommand.GetCommand())
					StdOutLines := DrushCommandOut // strings.Split(string(DrushCommandOut), "\n")
					for _, StdOutLine := range StdOutLines {
						//log.Print(StdOutLine)
						FinalOutput = append(FinalOutput, StdOutLine)
					}
				} else {
					log.Infof("%v, %v, successful.", DrushCommand.GetAlias(), DrushCommand.GetCommand())
					StdOutLines := DrushCommandOut // strings.Split(string(DrushCommandOut), "\n")
					for _, StdOutLine := range StdOutLines {
						//log.Print(StdOutLine)
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
	}
}
