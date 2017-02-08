package make

import (
	"fmt"
	"github.com/fubarhouse/golang-drush/alias"
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

func RestartWebServer(webserver string) {
	_, stdErr := exec.Command("sudo", "service", webserver, "restart").Output()
	if stdErr != nil {
		log.Fatalln(stdErr)
	} else {
		log.Printf("Webserver %v successfully restarted.\n", webserver)
	}
}

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

type Site struct {
	Timestamp string
	Path      string
	Make      string
	Name      string
	Alias     string
	database  *makeDB
}

func NewSite(make, name, path, alias string) *Site {
	Site := &Site{}
	Site.TimeStampReset()
	Site.Make = make
	Site.Name = name
	Site.Alias = alias
	Site.Path = path
	return Site
}

func (Site *Site) AliasExists(filter string) bool {
	y := aliases.NewAliasList()
	y.Generate(filter)
	for _, z := range y.GetNames() {
		if strings.Contains(z, Site.Alias) {
			return true
		}
	}
	return false
}

func (Site *Site) AliasInstall() {
	siteAlias := alias.NewAlias(Site.Name, Site.Path+"_latest", Site.Alias)
	siteAlias.Install()
}

func (Site *Site) AliasUninstall() {
	siteAlias := alias.NewAlias(Site.Name, Site.Path+"_latest", Site.Alias)
	siteAlias.Uninstall()
}

func (Site *Site) ActionBuild() {
	// TODO: Define purpose with the existence of ProcessMake()
	if Site.AliasExists(Site.Name) == true {
		Site.Path = fmt.Sprintf("%v%v", Site.Path, Site.TimeStampGet())
		//Site.ProcessMakes([]string{"core.make", "libraries.make", "contrib.make", "custom.make"})
		Site.ActionInstall()
	}
}

func (Site *Site) ActionDestroy() {
	// Destroy will remove all traces of said site.
	for _, database := range Site.DatabasesGet() {
		sqlQuery := fmt.Sprintf("DROP DATABASE %v;", database)
		sqlUser := fmt.Sprintf("--user=%v", Site.database.getUser())
		sqlPass := fmt.Sprintf("--password=%v", Site.database.getPass())
		_, err := exec.Command("mysql", sqlUser, sqlPass, "-e", sqlQuery).Output()
		if err == nil {
			log.Printf("Database %v was dropped.\n")
		} else {
			log.Printf("Database %v was notdropped: %v\n", err)
		}
	}
	// TODO: Delete folders
	// TODO: Delete drush aliases
	// TODO: Delete virtual hosts
	// TODO: Delete symlink
}

func (Site *Site) ActionInstall() {
	sqlQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %v_%v;", Site.Name, Site.Timestamp)
	sqlUser := fmt.Sprintf("--user=%v", Site.database.getUser())
	sqlPass := fmt.Sprintf("--password=%v", Site.database.getPass())
	_, err := exec.Command("mysql", sqlUser, sqlPass, "-e", sqlQuery).Output()
	if err != nil {
		log.Println("MySQL Error:", err)
	}
	output, _ := exec.Command("mysql", sqlUser, sqlPass, "-e", "show databases;").Output()
	if strings.Contains(string(output), Site.Name+"_"+Site.Timestamp) == false {
		log.Printf("Database %v_%v could not be created.\n", Site.Name, Site.Timestamp)
	} else {
		log.Printf("Database %v_%v was successfully created.\n", Site.Name, Site.Timestamp)
	}
	thisCmd := fmt.Sprintf("-y site-install standard --sites-subdir=%v --db-url=mysql://%v:%v@%v:%v/%v_%v install_configure_form.update_status_module='array(FALSE,FALSE)'", Site.Name, Site.database.getUser(), Site.database.getPass(), Site.database.getHost(), Site.database.getPort(), Site.Name, Site.Timestamp)
	output, err = exec.Command("sh", "-c", "cd "+Site.Path+"_"+Site.Timestamp+" && drush "+thisCmd).Output()
	_, cpErr := exec.Command("cp", "-f", Site.Path+"_"+Site.Timestamp+"/sites/"+Site.Name+"/settings.php", Site.Path+"_"+Site.Timestamp+"/sites/default/settings.php").Output()
	if cpErr != nil {
		panic("copy failed")
	}
	if err != nil {
		_, statErr := os.Stat(Site.Path + "_" + Site.Timestamp + "sites/" + Site.Name + "/settings.php")
		if statErr == nil {
			log.Println("Drush error:", err)
			log.Println(string(output))
			log.Println("cd " + Site.Path + "_" + Site.Timestamp + " && drush " + thisCmd)
		} else {
			log.Println("Drush install succeeded")
		}
	} else {
		log.Println("Drush install succeeded")
	}
}

func (Site *Site) ActionKill() {
	// Kill will delete a single site instance.
	// What to do with the default...
	if Site.AliasExists(Site.Name) == true {
		Site.Path = fmt.Sprintf("%v", Site.Path)
		_, err := os.Stat(Site.Path)
		if err == nil {
			os.Remove(Site.Path)
		}
	}
}

func (Site *Site) ActionRebuild() {
	// TODO: Define purpose with the existence of ProcessMake()
	if Site.AliasExists(Site.Name) == true {
		Site.TimeStampReset()
		Site.Path = fmt.Sprintf("%v%v", Site.Path, Site.TimeStampGet())
		//Site.ProcessMake()
		//Site.ActionInstall()
	}
}

func (Site *Site) ActionDatabaseDumpLocal(path string) {
	srcAlias := strings.Replace(Site.Alias, "@", "", -1)
	x := command.NewDrushCommand()
	x.Set(srcAlias, fmt.Sprintf("sql-dump %v", path), true)
	_, err := x.Output()
	if err == nil {
		log.Println("Database dump complete.")
		log.Println("Dump can be found at", path)
	} else {
		log.Println("Database dump could not complete.")
	}
}

func (Site *Site) ActionDatabaseDumpRemote(alias, path string) {
	srcAlias := strings.Replace(alias, "@", "", -1)
	x := command.NewDrushCommand()
	x.Set(srcAlias, fmt.Sprintf("sql-dump %v", path), true)
	_, err := x.Output()
	if err == nil {
		log.Println("Database dump complete.")
		log.Println("Dump can be found at", path)
	} else {
		log.Println("Database dump could not complete.")
	}
}

func (Site *Site) ActionDatabaseSyncLocal(alias string) {
	x := command.NewDrushCommand()
	srcAlias := strings.Replace(alias, "@", "", -1)
	destAlias := strings.Replace(Site.Alias, "@", "", -1)
	x.Set("", fmt.Sprintf("sql-sync @%v @%v -y", srcAlias, destAlias), true)
	_, err := x.Output()
	if err == nil {
		log.Println("Database syncronise complete.")
	} else {
		log.Println("Database syncronise could not complete.")
	}
}

func (Site *Site) ActionDatabaseSyncRemote(alias string) {
	x := command.NewDrushCommand()
	srcAlias := strings.Replace(Site.Alias, "@", "", -1)
	destAlias := strings.Replace(alias, "@", "", -1)
	x.Set("", fmt.Sprintf("sql-sync @%v @%v -y", srcAlias, destAlias), true)
	_, err := x.Output()
	if err == nil {
		log.Println("Database syncronise complete.")
	} else {
		log.Println("Database syncronise could not complete.")
	}
}

func (Site *Site) ActionFilesSyncLocal(alias string) {
	x := command.NewDrushCommand()
	srcAlias := strings.Replace(alias, "@", "", -1)
	destAlias := strings.Replace(Site.Alias, "@", "", -1)
	x.Set("", fmt.Sprintf("rsync -y --exclude-other-sites --exclude-conf @%v:%%files @%v:%%files", srcAlias, destAlias), true)
	_, err := x.Output()
	if err == nil {
		log.Println("Public file system has been synced.")
	} else {
		log.Println("Public file system has not been synced.")
	}
	x.Set("", fmt.Sprintf("rsync -y --exclude-other-sites --exclude-conf @%v:%%private @%v:%%private", srcAlias, destAlias), true)
	_, err = x.Output()
	if err == nil {
		log.Println("Private file system has been synced.")
	} else {
		log.Println("Private file system has not been synced.")
	}
}

func (Site *Site) ActionFilesSyncRemote(alias string) {
	x := command.NewDrushCommand()
	srcAlias := strings.Replace(Site.Alias, "@", "", -1)
	destAlias := strings.Replace(alias, "@", "", -1)
	x.Set("", fmt.Sprintf("rsync -y --exclude-other-sites --exclude-conf @%v:%%files @%v:%%files", srcAlias, destAlias), true)
	_, err := x.Output()
	if err == nil {
		log.Println("Public file system has been synced.")
	} else {
		log.Println("Public file system has not been synced.")
	}
	x.Set("", fmt.Sprintf("rsync -y --exclude-other-sites --exclude-conf @%v:%%private @%v:%%private", srcAlias, destAlias), true)
	_, err = x.Output()
	if err == nil {
		log.Println("Private file system has been synced.")
	} else {
		log.Println("Private file system has not been synced.")
	}
}

func (Site *Site) DatabaseSet(database *makeDB) {
	Site.database = database
}

func (Site *Site) DatabasesGet() []string {
	values, _ := exec.Command("mysql", "--user="+Site.database.dbUser, "--password="+Site.database.dbPass, "-e", "show databases").Output()
	databases := strings.Split(string(values), "\n")
	siteDbs := []string{}
	for _, database := range databases {
		if strings.Contains(database, Site.Name) == true {
			siteDbs = append(siteDbs, database)
		}
	}
	return siteDbs
}

func (Site *Site) SymInstall(timestamp string) {
	Symlink := Site.Path + "_latest"
	err := os.Symlink(Site.Path+"_"+Site.TimeStampGet(), Symlink)
	if err == nil {
		log.Println("Symlink has been created.")
	} else {
		log.Println("Symlink has not been created:", err)
	}
}

func (Site *Site) SymUninstall(timestamp string) {
	Symlink := Site.Path + "_latest"
	_, statErr := os.Stat(Symlink)
	if statErr == nil {
		os.Remove(Symlink)
		log.Println("Symlink has been removed.")
	}
}

func (Site *Site) SymReinstall(timestamp string) {
	Site.SymUninstall(timestamp)
	Site.SymInstall(timestamp)
}

func (Site *Site) TimeStampGet() string {
	return Site.Timestamp
}

func (Site *Site) TimeStampSet(value string) {
	Site.Timestamp = fmt.Sprintf("%v", value)
}

func (Site *Site) TimeStampReset() {
	now := time.Now()
	Site.Timestamp = fmt.Sprintf("%v", now.Format("20060102150405"))
}

func (Site *Site) ProcessMake(makeFile string) {

	fullPath := makeFile
	_, err := os.Stat(fullPath)
	if err != nil {
		log.Println("Error! File not found:", err)
		os.Exit(1)
	}

	_, err = os.Stat(Site.Path)
	if err == nil {
		log.Println("Creating directory for site at", Site.Path+"_"+Site.Timestamp)
		os.MkdirAll(Site.Path, 0755)
	}

	drushCommand := ""
	// @TODO: Figure out a way to run make without core, but optionally based on makefile.
	if strings.Contains(makeFile, "core") == true {
		drushCommand = fmt.Sprintf("make -y --working-copy %v %v_%v", fullPath, Site.Path, Site.Timestamp)
	} else {
		drushCommand = fmt.Sprintf("make -y --no-core --working-copy %v %v_%v", fullPath, Site.Path, Site.Timestamp)
	}
	log.Println("Building from", makeFile)
	drushMake := command.NewDrushCommand()
	drushMake.Set("", drushCommand, true)
	cmd, err := drushMake.Output()
	if err != nil {
		if string(err.Error()) == "exit status 1" {
			log.Println("Processed make file was completed with errors. :", drushCommand)

		}
	} else {
		fmt.Sprintln(cmd)
	}
}

func (Site *Site) VhostInstall(webserver, path string) {
	vhostPath := strings.Replace(Site.Path+"_"+Site.TimeStampGet(), "_"+Site.TimeStampGet(), "_latest", -1)
	vhostFile := vhost.NewVirtualHost(Site.Name, vhostPath, webserver, path)
	vhostFile.Install()
}

func (Site *Site) VhostUninstall(webserver, path string) {
	vhostPath := strings.Replace(Site.Path+"_"+Site.TimeStampGet(), "_"+Site.TimeStampGet(), "_latest", -1)
	vhostFile := vhost.NewVirtualHost(Site.Name, vhostPath, webserver, path)
	vhostFile.Uninstall()
}
