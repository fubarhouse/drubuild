package make

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
)

// ActionDestroyDatabases will destroy all databases associated to the site struct.
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
		log.Infof("Database(s) removed: %v", dbDeleteCount)
	}
}

// ActionDestroyAlias is an API call for alias un-installation.
func (Site *Site) ActionDestroyAlias() {
	Site.AliasUninstall()
}

// ActionDestroyVhost is an API call for virtual-host un-installation.
func (Site *Site) ActionDestroyVhost() {
	Site.VhostUninstall()
}

// ActionDestroyPermissions is an API call for site file system un-installation.
func (Site *Site) ActionDestroyPermissions() {
	privateFilesPath := Site.Path
	_, statErr := os.Stat(privateFilesPath)
	if statErr == nil {
		files, _ := ioutil.ReadDir(privateFilesPath)
		for _, file := range files {
			privateFilesPathTarget := privateFilesPath + "/" + file.Name() + "/sites/" + Site.Name
			chmodErr := os.Chmod(privateFilesPathTarget, 0777)
			if chmodErr != nil {
				log.Warnf("Could not set permissions of %v to %v: %v", privateFilesPathTarget, "0777", chmodErr)
			} else {
				log.Infof("Set permissions of %v to %v", privateFilesPathTarget, "0777")
			}
		}
	} else {
		log.Warnln("Could not find target folders", privateFilesPath)
	}
}

// ActionDestroySym is an API call for symlink un-installation.
func (Site *Site) ActionDestroySym() {
	Site.SymUninstall()
}

// ActionDestroyFiles is an API call for file system removal.
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
		log.Warnln("Site directory was not found: ", Site.Path)
	}
}

// ActionDestroy is an API call for site removal.
func (Site *Site) ActionDestroy() {
	// Destroy will remove all traces of said site.
	Site.ActionDestroyDatabases()
	Site.ActionDestroyAlias()
	Site.ActionDestroyVhost()
	Site.ActionDestroyPermissions()
	Site.ActionDestroyFiles()
	Site.ActionDestroySym()
}
