package make

import (
	"fmt"
	"github.com/fubarhouse/golang-drush/aliases"
	"github.com/fubarhouse/golang-drush/command"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type makeDB struct {
	dbHost string
	dbUser string
	dbPass string
	dbPort int
}

func newmakeDB(dbHost, dbUser, dbPass string, dbPort int) *makeDB {
	newDB := &makeDB{}
	newDB.setHost(dbHost)
	newDB.setUser(dbUser)
	newDB.setPass(dbPass)
	newDB.setPort(dbPort)
	return newDB
}

func (db *makeDB) setHost(dbHost string) {
	db.dbHost = dbHost
}

func (db *makeDB) getHost() string {
	return db.dbHost
}

func (db *makeDB) setUser(dbUser string) {
	db.dbUser = dbUser
}

func (db *makeDB) getUser() string {
	return db.dbUser
}

func (db *makeDB) setPass(dbPass string) {
	db.dbPass = dbPass
}

func (db *makeDB) getPass() string {
	return db.dbPass
}

func (db *makeDB) setPort(dbPort int) {
	db.dbPort = dbPort
}

func (db *makeDB) getPort() int {
	return db.dbPort
}

func (db *makeDB) getInfo() *makeDB {
	return db
}

func replaceTextInFile(fullPath string, oldString string, newString string) {
	read, err := ioutil.ReadFile(fullPath)
	if err != nil {
		panic(err)
	}
	newContents := strings.Replace(string(read), oldString, newString, -1)
	err = ioutil.WriteFile(fullPath, []byte(newContents), 0)
	if err != nil {
		panic(err)
	}
}

type Site struct {
	timestamp     string
	site_name     string
	site_path     string
	make_path     string
	make_dir      string
	codebase_name string
	codebase_path string
	branch        string
	alias         string
	drushVersion  int
}

func newSite(codebase_name, codebase_path, make_dir, make_path, branch, site_name, site_path, alias string, drushVersion int) *Site {
	Site := &Site{}
	Site.SetTimeStamp()
	Site.codebase_name = codebase_name
	Site.codebase_path = codebase_path
	Site.make_dir = make_dir
	Site.make_path = make_path
	Site.branch = branch
	Site.alias = alias
	Site.drushVersion = drushVersion
	Site.site_name = site_name
	Site.site_path = site_path
	return Site
}

func (Site *Site) AliasExists(filter string) bool {
	y := aliases.NewAliasList()
	y.Generate(filter)
	for _, z := range y.GetNames() {
		if strings.Contains(z, Site.site_name) {
			return true
		}
	}
	return false
}

func (Site *Site) SetTimeStamp() {
	now := time.Now()
	Site.timestamp = fmt.Sprintf("%v", now.Format("20060102150405"))
}

func (Site *Site) GetTimeStamp() string {
	return Site.timestamp
}

func (Site *Site) ProcessCoreMake() {
	// Function to build core.make
	_, err := os.Stat(Site.site_path)
	if err == nil {
		log.Println("Creating new directory for site")
		os.MkdirAll(Site.site_path, 0755)
	}
	for _, makefile := range []string{"core.make"} {
		fullPath := fmt.Sprintf("%v/%v/%v", Site.make_path, Site.make_dir, makefile)
		_, err := os.Stat(fullPath)
		if err != nil {
			log.Println("Error! File not found:", err)
			os.Exit(1)
		}
		log.Println("Building from", makefile)
		// TODO: Consider copying from codebase here...
		drushCommand := fmt.Sprintf("make %v %v -y --working-copy", fullPath, Site.site_path)
		drushMake := command.NewDrushCommand()
		drushMake.Set("", drushCommand, true)
		cmd, err := drushMake.Output()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Sprintln(cmd)
		}
	}
}

func (Site *Site) ProcessMakes(makeFiles []string) {

	// TODO: Consider doing away with makes and/or copying the data from the codebase folder.

	for _, makefile := range makeFiles {
		fullPath := fmt.Sprintf("%v/%v/%v", Site.make_path, Site.make_dir, makefile)
		_, err := os.Stat(fullPath)
		if err != nil {
			log.Println("Error! File not found:", err)
			os.Exit(1)
		}

		if strings.Contains(makefile, "core") == true {
			Site.ProcessCoreMake()
		} else {
			log.Println("Building from", makefile)
			// So this works - even on the host.
			replaceTextInFile(fullPath, "rewriteme", Site.branch)
			replaceTextInFile(fullPath, "master", Site.branch)
			replaceTextInFile(fullPath, "production", Site.branch)
			// TODO: Consider copying from codebase instead of drush making...
			drushCommand := fmt.Sprintf("-y --no-core --working-copy make %v %v", fullPath, Site.site_path)
			drushMake := command.NewDrushCommand()
			drushMake.Set("", drushCommand, true)
			cmd, err := drushMake.Output()
			if err != nil {
				if string(err.Error()) == "exit status 1" {
					fmt.Println("Completed with errors.")
				} else {
					fmt.Println(err)
				}
			} else {
				fmt.Sprintln(cmd)
			}
			replaceTextInFile(fullPath, "rewriteme", "rewriteme")
			replaceTextInFile(fullPath, "master", "rewriteme")
			replaceTextInFile(fullPath, "production", "rewriteme")
		}
	}
}

func (Site *Site) Install(database *makeDB) {
	sqlQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %v_%v;", Site.site_name, Site.timestamp)
	sqlArgs := fmt.Sprintf("--user=`%v` --password=`%v` -e", database.getUser(), database.getPass())
	_, err := exec.Command("sh", "-c", "mysql", sqlArgs, sqlQuery).Output()
	if err != nil {
		log.Println("MySQL Error:", err)
		os.Exit(1)
	}
	x := command.NewDrushCommand()
	x.SetAlias(Site.alias)
	thisCmd := fmt.Sprintf("-y site-install standard --sites-subdir=`%v` --db-url=`mysql://%v:%v@%v:%v/%v_%v`", Site.site_name, database.getUser(), database.getPass(), database.getHost(), database.getPort(), Site.site_name, Site.timestamp)
	x.SetCommand(thisCmd)
	_, err = x.Output()
	if err != nil {
		log.Println(err)
	}
}

func (Site *Site) Build(database *makeDB) {
	if Site.AliasExists(Site.site_name) == true {
		Site.site_path = fmt.Sprintf("%v%v", Site.site_path, Site.GetTimeStamp())
		//Site.ProcessMakes([]string{"core.make", "libraries.make", "contrib.make", "custom.make"})
		Site.Install(database)
	}
}

func (Site *Site) Rebuild(database *makeDB) {
	if Site.AliasExists(Site.site_name) == true {
		Site.SetTimeStamp()
		Site.site_path = fmt.Sprintf("%v%v", Site.site_path, Site.GetTimeStamp())
		Site.ProcessMakes([]string{"core.make", "libraries.make", "contrib.make", "custom.make"})
		Site.Install(database)
	}
}

func (Site *Site) Destroy(database *makeDB) {
	if Site.AliasExists(Site.site_name) == true {
		Site.site_path = fmt.Sprintf("%v", Site.site_path)
		_, err := os.Stat(Site.site_path)
		if err == nil {
			os.Remove(Site.site_path)
		}
	}
}
