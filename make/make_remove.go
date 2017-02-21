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
			log.Printf("Database %v was dropped.\n", database)
		} else {
			log.Printf("Database %v was not dropped: %v\n", database, err)
		}
	}
}

func (Site *Site) ActionDestroyAlias() {
	log.Printf("Removing file system for %v at %v\n", Site.Name, Site.Path)
	Site.AliasUninstall()
}

func (Site *Site) ActionDestroyVhost() {
	log.Printf("Removing Virtual Host file for site %v\n", Site.Name)
	Site.VhostUninstall()
}

func (Site *Site) ActionDestroySym() {
	log.Printf("Removing symlink for site %v\n", Site.Name)
	Site.SymUninstall(Site.Timestamp)
}

func (Site *Site) ActionDestroyFiles() {
	log.Printf("Removing file system for %v at %v\n", Site.Name, Site.Path)
	os.RemoveAll(Site.Path + "/../")
}

func (Site *Site) ActionDestroy() {
	// Destroy will remove all traces of said site.
	Site.ActionDestroyDatabases()
	Site.ActionDestroyAlias()
	Site.ActionDestroyVhost()
	Site.ActionDestroySym()
	Site.ActionDestroyFiles()
}
