package make

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"strings"
	"time"
)

// InstallFileSystem installs a basic private file system for any given input.
func (Site *Site) InstallFileSystem(DirectoryPath string) {
	 // Test the file system, create it if it doesn't exist!
	 var dirPath string
	 if Site.Composer {
		 dirPath = fmt.Sprintf(strings.Join([]string{Site.Path, Site.Name + Site.TimeStampGet(), Site.Docroot, "sites", Site.Name, DirectoryPath}, "/"))
	 } else {
		 dirPath = fmt.Sprintf(strings.Join([]string{Site.Path, Site.Name + Site.TimeStampGet(), Site.Docroot, "sites", Site.Name, DirectoryPath}, "/"))
	 }
	_, err := os.Stat(dirPath + "/" + dirPath)
	if err != nil {
		dirErr := os.MkdirAll(dirPath, 0755)
		if dirErr != nil {
			log.Errorln("Couldn't create file system at", dirPath, dirErr)
		} else {
			log.Infoln("Created file system at", dirPath)
			time.Sleep(1 * time.Second)
		}
	}
}
