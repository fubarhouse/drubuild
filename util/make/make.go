package make

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql" // mysql is assumed under this system (for now).
	"github.com/fubarhouse/drubuild/util/drush"
)

// Site struct which represents a build website being used.
type Site struct {
	Timestamp string
	Path      string
	Name          string
	Alias         string
	Domain        string
	Docroot       string
	database      *makeDB
	Webserver     string
	Vhostpath     string
	Template      string
	AliasTemplate string
	FilePathPrivate            string
	FilePathPublic             string
	FilePathTemp               string
	WorkingCopy                bool
	Composer                   bool
}

// NewSite instantiates an instance of the struct Site
func NewSite(make, name, path, alias, webserver, domain, vhostpath, template string) *Site {
	Site := &Site{}
	Site.TimeStampReset()
	Site.Name = name
	Site.Path = path
	Site.Webserver = webserver
	Site.Alias = alias
	Site.Domain = domain
	Site.Vhostpath = vhostpath
	Site.Template = template
	Site.FilePathPrivate = "files/private"
	Site.FilePathPublic = "" // For later implementation
	Site.FilePathTemp = "files/private/temp"
	Site.WorkingCopy = false
	return Site
}

// ActionInstall runs drush site-install on a Site struct
func (Site *Site) ActionInstall() {
	// Obtain a string value from the Port value in db config.
	stringPort := fmt.Sprintf("%v", Site.database.getPort())
	// Open a mysql connection
	db, dbErr := sql.Open("mysql", Site.database.getUser()+":"+Site.database.getPass()+"@tcp("+Site.database.dbHost+":"+stringPort+")/")
	// Defer the connection
	defer db.Close()
	// Report any connection errors
	if dbErr != nil {
		log.Warnf("WARN:", dbErr)
	}
	// Create database
	dbName := strings.Replace(Site.Name+Site.Timestamp, ".", "_", -1)
	_, dbErr = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if dbErr != nil {
		panic(dbErr)
	}
	// Drush site-install
	var sitePath string
	sitePath = Site.Path + string(os.PathSeparator) + Site.Name + Site.Timestamp + string(os.PathSeparator) + Site.Docroot
	drush.Run([]string{"site-install", "--root="+sitePath, "--yes", "--sites-subdir="+Site.Name, fmt.Sprintf("--db-url=mysql://%v:%v@%v:%v/%v", Site.database.getUser(), Site.database.getPass(), Site.database.getHost(), Site.database.getPort(), dbName)})
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

// DatabaseSet sets the database field to an inputted *makeDB struct.
func (Site *Site) DatabaseSet(database *makeDB) {
	Site.database = database
}

// DatabasesGet returns a list of databases associated to local builds from the site struct
func (Site *Site) DatabasesGet() []string {
	values, _ := exec.Command("mysql", "--user="+Site.database.dbUser, "--password="+Site.database.dbPass, "-e", "show databases").Output()
	databases := strings.Split(string(values), "\n")
	siteDbs := []string{}
	for _, database := range databases {
		if strings.HasPrefix(database, Site.Name+"_2") {
			siteDbs = append(siteDbs, database)
		}
	}
	return siteDbs
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

// InstallSiteRef installs the Drupal multisite sites.php file for the site struct.
func (Site *Site) InstallSiteRef(Template string) {

	if Template == "" {
		log.Warnln("no template specified for sites.php")
		return
	}

	data := map[string]string{
		"Name":  Site.Name,
		"Alias": Site.Alias,
	}
	var dirPath string
	dirPath = Site.Path + "/" + Site.Name + Site.Timestamp + "/" + Site.Docroot + "/sites/"
	dirErr := os.MkdirAll(dirPath+Site.Name, 0755)
	if dirErr != nil {
		log.Errorln("Unable to create directory", dirPath+Site.Name, dirErr)
	} else {
		log.Infoln("Created directory", dirPath+Site.Name)
	}

	dirErr = os.Chmod(dirPath+Site.Name, 0775)
	if dirErr != nil {
		log.Errorln("Could not set permissions 0755 on", dirPath+Site.Name, dirErr)
	} else {
		log.Infoln("Permissions set to 0755 on", dirPath+Site.Name)
	}

	filename := dirPath + "/sites.php"

	t := template.New("sites.php")
	defaultData, _ := ioutil.ReadFile(Template)
	t.Parse(string(defaultData))
	file, _ := os.Create(filename)
	tplErr := t.Execute(file, data)

	if tplErr == nil {
		log.Infof("Successfully templated multisite config to file %v", filename)
	} else {
		log.Warnf("Error templating multisite config to file %v", filename)
	}
}

// ReplaceTextInFile reinstalls and verifies the ctools cache folder for the site struct.
func (Site *Site) ReplaceTextInFile() {
	// We need to remove and re-add the ctools cache directory as 0777.
	cToolsDir := fmt.Sprintf("%v/%v%v/sites/%v/files/ctools", Site.Path, Site.Name, Site.Timestamp, Site.Name)
	// Remove the directory!
	cToolsErr := os.RemoveAll(cToolsDir)
	if cToolsErr != nil {
		log.Errorln("Couldn't remove", cToolsDir)
	} else {
		log.Infoln("Created", cToolsDir)
	}
	// Add the directory!
	cToolsErr = os.Mkdir(cToolsDir, 0777)
	if cToolsErr != nil {
		log.Errorln("Couldn't remove", cToolsDir)
	} else {
		log.Infoln("Created", cToolsDir)
	}
}
