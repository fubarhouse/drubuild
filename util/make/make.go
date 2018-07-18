package make

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql" // mysql is assumed under this system (for now).
	)

// Site struct which represents a build website being used.
type Site struct {
	Timestamp string
	Path      string
	Name          string
	Alias         string
	Domain        string
	Docroot       string
	Template      string
	AliasTemplate string
	FilePathPrivate            string
	FilePathPublic             string
	FilePathTemp               string
	WorkingCopy                bool
	Composer                   bool
}

// NewSite instantiates an instance of the struct Site
func NewSite(name, path, alias, domain string) *Site {
	Site := &Site{}
	Site.TimeStampReset()
	Site.Name = name
	Site.Path = path
	Site.Alias = alias
	Site.Domain = domain
	Site.FilePathPrivate = "files/private"
	Site.FilePathPublic = "" // For later implementation
	Site.FilePathTemp = "files/private/temp"
	return Site
}

// CleanCodebase will remove all data from the site path other than the /sites folder and contents.
func (Site *Site) CleanCodebase() {
	_ = filepath.Walk(Site.Path, func(path string, Info os.FileInfo, _ error) error {
		realpath := strings.Split(Site.Path, "\n")
		err := new(error)
		for _, name := range realpath {
			if strings.Contains(path, Site.TimeStampGet()) {
				if !strings.Contains(path, "/sites") || strings.Contains(path, "/sites/all") {
					//return nil
					if path != Site.Path {
						if Info.IsDir() && !strings.HasSuffix(path, Site.Path) {
							fmt.Sprintln(name)
							os.Chmod(path, 0777)
							delErr := os.RemoveAll(path)
							if delErr != nil {
								log.Warnln("Could not remove", path)
							}
						} else if !Info.IsDir() {
							log.Infoln("Not removing", path)
						}
					}
				}
			}
		}
		return *err
	})
}

// SymInstall installs a symlink to the site directory of the site struct
func (Site *Site) SymInstall() {
	Target := filepath.Join(Site.Name + Site.TimeStampGet())
	Symlink := filepath.Join(Site.Path, Site.Domain+".latest")
	err := os.Symlink(Target, Symlink)
	if err == nil {
		log.Infoln("Created symlink")
	} else {
		log.Warnln("Could not create symlink:", err)
	}
}

// SymUninstall removes a symlink to the site directory of the site struct
func (Site *Site) SymUninstall() {
	Symlink := Site.Domain + ".latest"
	_, statErr := os.Stat(Site.Path + "/" + Symlink)
	if statErr == nil {
		err := os.Remove(Site.Path + "/" + Symlink)
		if err != nil {
			log.Errorln("Could not remove symlink.", err)
		} else {
			log.Infoln("Removed symlink.")
		}
	}
}

// SymReinstall re-installs a symlink to the site directory of the site struct
func (Site *Site) SymReinstall() {
	Site.SymUninstall()
	Site.SymInstall()
}

// TimeStampGet returns the timestamp variable for the site struct
func (Site *Site) TimeStampGet() string {
	return Site.Timestamp
}

// TimeStampReset sets the timestamp field for the site struct to a new value
func (Site *Site) TimeStampReset() {
	now := time.Now()
	r := rand.Intn(100) * rand.Intn(100)
	Site.Timestamp = fmt.Sprintf(".%v_%v", now.Format("20060102150405"), r)
}
