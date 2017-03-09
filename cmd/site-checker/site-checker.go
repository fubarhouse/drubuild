package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/golang-drush/command"
	"strings"
)

func main() {
	var strAliases = flag.String("aliases", "", "alias1,alias2,alias3")
	flag.Parse()

	if *strAliases != "" {
		for _, Alias := range strings.Split(*strAliases, "\n") {

			cmdQ := command.NewDrushCommand()
			cmdQ.SetAlias(Alias)
			cmdQ.SetCommand("rr")
			_, errQ := cmdQ.Output()
			if errQ != nil {
				log.Warnf("%v, %v, unsuccessful.\n", cmdQ.GetAlias(), cmdQ.GetCommand())
			} else {
				log.Infof("%v, %v, unsuccessful.\n", cmdQ.GetAlias(), cmdQ.GetCommand())
			}

			cmdR := command.NewDrushCommand()
			cmdR.SetAlias(Alias)
			cmdR.SetCommand("updb --yes")
			_, errR := cmdR.Output()
			if errR != nil {
				log.Warnf("%v, %v, unsuccessful.\n", cmdR.GetAlias(), cmdR.GetCommand())
			} else {
				log.Infof("%v, %v, unsuccessful.\n", cmdR.GetAlias(), cmdR.GetCommand())
			}

			cmdS := command.NewDrushCommand()
			cmdS.SetAlias(Alias)
			cmdS.SetCommand("cc all")
			_, errS := cmdS.Output()
			if errS != nil {
				log.Warnf("%v, %v, unsuccessful.\n", cmdS.GetAlias(), cmdS.GetCommand())
			} else {
				log.Infof("%v, %v, unsuccessful.\n", cmdS.GetAlias(), cmdS.GetCommand())
			}

		}
	} else {
		flag.Usage()
	}
}
