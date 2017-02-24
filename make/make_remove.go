package make

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func (Site *Site) ActionDestroyDatabases() {
	for _, database := range Site.DatabasesGet() {
		sqlQuery := fmt.Sprintf("DROP DATABASE %v;", database)
		sqlUser := fmt.Sprintf("--user=%v", Site.database.getUser())
		sqlPass := fmt.Sprintf("--password=%v", Site.database.getPass())
		_, err := exec.Command("mysql", sqlUser, sqlPass, "-e", sqlQuery).Output()
		if err == nil {
			log.Printf("OK: Successfully dropped database %v.\n", database)
		} else {
			log.Printf("WARN: Could not drop database %v: %v\n", database, err)
		}
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
			log.Printf("WARN: Could not remove file system for %v at %v\n", Site.Name, Site.Path)
		} else {
			log.Printf("OK: Successfully removed file system for %v at %v\n", Site.Name, Site.Path)
		}
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
