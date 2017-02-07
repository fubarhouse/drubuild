package make

import (
	"fmt"
	"github.com/fubarhouse/golang-drush/aliases"
	"github.com/fubarhouse/golang-drush/command"
	"github.com/fubarhouse/golang-drush/vhost"
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

func NewmakeDB(dbHost, dbUser, dbPass string, dbPort int) *makeDB {
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
	timestamp string
	path      string
	make      string
	name      string
	alias     string
	database  *makeDB
}

func NewSite(make, name, path, alias string) *Site {
	Site := &Site{}
	Site.SetTimeStamp()
	Site.make = make
	Site.name = name
	Site.alias = alias
	Site.path = path
	return Site
}

func (Site *Site) AliasExists(filter string) bool {
	y := aliases.NewAliasList()
	y.Generate(filter)
	for _, z := range y.GetNames() {
		if strings.Contains(z, Site.alias) {
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
	_, err := os.Stat(Site.path)
	if err == nil {
		log.Println("Creating new directory for site")
		os.MkdirAll(Site.path, 0755)
	}
	for _, makefile := range []string{"core.make"} {
		fullPath := Site.make
		_, err := os.Stat(fullPath)
		if err != nil {
			log.Println("Error! File not found:", err)
			os.Exit(1)
		}
		log.Println("Building from", makefile)
		// TODO: Consider copying from codebase here...
		drushCommand := fmt.Sprintf("make %v %v_%v -y --working-copy", fullPath, Site.path, Site.timestamp)
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
		fullPath := fmt.Sprintf("%v/%v/%v", Site.path, makefile)
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
			// TODO: Consider copying from codebase instead of drush making...
			// TODO: Rewrite files separately elsewhere...
			drushCommand := fmt.Sprintf("-y --no-core --working-copy make %v %v", fullPath, Site.path)
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
		}
	}
}

func (Site *Site) Install() {
	sqlQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %v_%v;", Site.name, Site.timestamp)
	sqlUser := fmt.Sprintf("--user=%v", Site.database.getUser())
	sqlPass := fmt.Sprintf("--password=%v", Site.database.getPass())
	_, err := exec.Command("mysql", sqlUser, sqlPass, "-e", sqlQuery).Output()
	if err != nil {
		log.Println("MySQL Error:", err)
	}
	output, _ := exec.Command("mysql", sqlUser, sqlPass, "-e", "show databases;").Output()
	if strings.Contains(string(output), Site.name+"_"+Site.timestamp) == false {
		log.Printf("Database %v_%v could not be created.\n", Site.name, Site.timestamp)
	} else {
		log.Printf("Database %v_%v was successfully created.\n", Site.name, Site.timestamp)
	}
	thisCmd := fmt.Sprintf("-y site-install standard --sites-subdir=%v --db-url=mysql://%v:%v@%v:%v/%v_%v", Site.name, Site.database.getUser(), Site.database.getPass(), Site.database.getHost(), Site.database.getPort(), Site.name, Site.timestamp)
	output, err = exec.Command("sh", "-c", "cd "+Site.path+"_"+Site.timestamp+" && drush "+thisCmd).Output()
	if err != nil {
		_, statErr := os.Stat(Site.path + "_" + Site.timestamp + "sites/aussiejobs/settings.php")
		if statErr == nil {
			log.Println("Drush error:", err)
			log.Println(string(output))
			log.Println("cd " + Site.path + "_" + Site.timestamp + " && drush " + thisCmd)
		} else {
			log.Println("Drush install succeeded")
		}
	} else {
		log.Println("Drush install succeeded")
	}
	vhost := vhost.NewVirtualHost(Site.name, Site.path+"_"+Site.timestamp, "nginx", "/etc/nginx/sites-enabled")
	vhost.Install()
}

func (Site *Site) SetDatabase(database *makeDB) {
	Site.database = database
}

func (Site *Site) Build() {
	if Site.AliasExists(Site.name) == true {
		Site.path = fmt.Sprintf("%v%v", Site.path, Site.GetTimeStamp())
		//Site.ProcessMakes([]string{"core.make", "libraries.make", "contrib.make", "custom.make"})
		Site.Install()
	}
}

func (Site *Site) Rebuild() {
	if Site.AliasExists(Site.name) == true {
		Site.SetTimeStamp()
		Site.path = fmt.Sprintf("%v%v", Site.path, Site.GetTimeStamp())
		Site.ProcessMakes([]string{"core.make", "libraries.make", "contrib.make", "custom.make"})
		Site.Install()
	}
}

func (Site *Site) Destroy() {
	if Site.AliasExists(Site.name) == true {
		Site.path = fmt.Sprintf("%v", Site.path)
		_, err := os.Stat(Site.path)
		if err == nil {
			os.Remove(Site.path)
		}
	}
}
