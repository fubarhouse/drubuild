package make

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
)

func (Site *Site) ActionDestroyDatabases() {
	var dbDeleteCount int
	for _, database := range Site.DatabasesGet() {
		sqlQuery := fmt.Sprintf("DROP DATABASE %v;", database)
		sqlUser := fmt.Sprintf("--user=%v", Site.database.getUser())
		sqlPass := fmt.Sprintf("--password=%v", Site.database.getPass())
		_, err := exec.Command("mysql", sqlUser, sqlPass, "-e", sqlQuery).Output()
		if err == nil {
			log.Infoln("Dropped database", database)
			dbDeleteCount++
		} else {
			log.Warnln("Could not drop database", database, err)
		}
	}
	if dbDeleteCount == 0 {
		log.Warnln("No database was found")
	} else {
		log.Infof("%v databases were removed", dbDeleteCount)
	}
}

func (Site *Site) ActionDestroyAlias() {
	Site.AliasUninstall()
}

func (Site *Site) ActionDestroyVhost() {
	Site.VhostUninstall()
}

func (Site *Site) ActionDestroySym() {
	Site.SymUninstall(Site.Timestamp)
}

func (Site *Site) ActionDestroyFiles() {
	_, statErr := os.Stat(Site.Path)
	if statErr == nil {
		err := os.RemoveAll(Site.Path)
		if err != nil {
			log.Warnf("Could not remove file system for %v at %v\n", Site.Name, Site.Path)
		} else {
			log.Infof("Removed file system for %v at %v\n", Site.Name, Site.Path)
		}
	} else {
		log.Warnln("Site directory was not found.")
	}
}

func (Site *Site) ActionDestroy() {
	// Destroy will remove all traces of said site.
	Site.ActionDestroyDatabases()
	Site.ActionDestroyAlias()
	Site.ActionDestroyVhost()
	Site.ActionDestroySym()
	Site.ActionDestroyFiles()
}
