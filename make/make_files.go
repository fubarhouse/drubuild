package make

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"strings"
)

// InstallFileSystem installs a basic private file system for any given input.
func (Site *Site) InstallFileSystem(DirectoryPath string) {
	// Test the file system, create it if it doesn't exist!
	dirPath := fmt.Sprintf(strings.Join([]string{Site.Path, DirectoryPath}, "/"))
	_, err := os.Stat(dirPath + "/" + dirPath)
	if err != nil {
		dirErr := os.MkdirAll(dirPath, 0755)
		if dirErr != nil {
			log.Errorln("Couldn't create file system at", dirPath, dirErr)
		} else {
			log.Infoln("Created file system at", dirPath)
		}
	}
}
