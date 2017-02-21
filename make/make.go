package make

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/fubarhouse/golang-drush/command"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
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

func (Site *Site) RestartWebServer() {
	_, stdErr := exec.Command("sudo", "service", Site.Webserver, "restart").Output()
	if stdErr != nil {
		log.Fatalln(stdErr)
	} else {
		log.Printf("Webserver %v successfully restarted.\n", Site.Webserver)
	}
}

type Site struct {
	Timestamp string
	Path      string
	Make      string
	Name      string
	Alias     string
	database  *makeDB
	Webserver string
	Vhostpath string
}

func NewSite(make, name, path, alias, webserver, vhostpath string) *Site {
	Site := &Site{}
	Site.TimeStampReset()
	Site.Make = make
	Site.Name = name
	Site.Alias = alias
	Site.Path = path
	Site.Webserver = webserver
	Site.Vhostpath = vhostpath
	return Site
}

func (Site *Site) ActionBuild() {
	// TODO: Define purpose with the existence of ProcessMake()
	if Site.AliasExists(Site.Name) == true {
		Site.Path = fmt.Sprintf("%v%v", Site.Path, Site.TimeStampGet())
		//Site.ProcessMakes([]string{"core.make", "libraries.make", "contrib.make", "custom.make"})
		Site.ActionInstall()
	}
}

func (Site *Site) ActionInstall() {
	// Obtain a string value from the Port value in db config.
	stringPort := fmt.Sprintf("%v", Site.database.getPort())
	// Open a mysql connection
	db, dbErr := sql.Open("mysql", Site.database.getUser()+":"+Site.database.getPass()+"@tcp("+Site.database.dbHost+":"+stringPort+")/")
	// Defer the connection
	defer db.Close()
	// Report any connection errors
	if dbErr != nil {
		log.Println(dbErr)
	}
	// Create database
	dbName := strings.Replace(Site.Name+Site.Timestamp, ".", "_", -1)
	_, dbErr = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if dbErr != nil {
		panic(dbErr)
	}
	// Drush site-install
	thisCmd := fmt.Sprintf("-y site-install standard --sites-subdir=%v --db-url=mysql://%v:%v@%v:%v/%v install_configure_form.update_status_module='array(FALSE,FALSE)'", Site.Name, Site.database.getUser(), Site.database.getPass(), Site.database.getHost(), Site.database.getPort(), dbName)
	output, _ := exec.Command("sh", "-c", "cd "+Site.Path+Site.Timestamp+" && drush "+thisCmd).Output()
	drushLog := strings.Split(string(output), "\n")
	for _, logEntry := range drushLog {
		// Print output in a fairly standardized format.
		log.Println(logEntry)
	}
	// Reset file system permissions...
	err := os.Chmod(Site.Path+Site.Timestamp+"/sites/"+Site.Name, os.FileMode(0755))
	if err != nil {
		log.Println("Unable to reset file system permissions for", Site.Path+Site.Timestamp+"/sites/"+Site.Name)
	} else {
		log.Println("Successfully reset file system permissions for", Site.Path+Site.Timestamp+"/sites/"+Site.Name)
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

func (Site *Site) ActionRebuildCodebase(Makefiles []string) {
	// This function exists for the sole purpose of
	// rebuilding a specific Drupal codebase in a specific
	// directory for Release management type work.
	// TODO Add a way to specify the branch for cloning in an independent way
	log.Println("Generating temporary make file...")
	newMakeFilePath := "/tmp/tmp.make"
	file, crErr := os.Create(newMakeFilePath)
	if crErr != nil {
		log.Println("Error creating "+newMakeFilePath+":", crErr)
	}
	writer := bufio.NewWriter(file)
	defer file.Close()

	fmt.Fprintln(writer, "core = 7.x")
	fmt.Fprintln(writer, "api = 2")

	for _, Makefile := range Makefiles {
		cmdOut, _ := exec.Command("cat", Makefile).Output()
		output := strings.Split(string(cmdOut), "\n")
		for _, line := range output {
			if strings.HasPrefix(line, "core") == false && strings.HasPrefix(line, "api") == false {
				if strings.HasPrefix(line, "projects") == true || strings.HasPrefix(line, "libraries") == true {
					fmt.Fprintln(writer, line)
				}
			}
		}
	}

	writer.Flush()
	//replaceTextInFile(newMakeFilePath, "rewriteme", "master")
	Site.ProcessMake(newMakeFilePath)
	log.Println("Removing temporary make file", newMakeFilePath)
	os.Remove(newMakeFilePath)
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
	// TODO make this relative to the lowest possible level
	Symlink := Site.Path + ".latest"
	err := os.Symlink(Site.Path+Site.TimeStampGet(), Symlink)
	if err == nil {
		log.Println("Symlink has been created.")
	} else {
		log.Println("Symlink has not been created:", err)
	}
}

func (Site *Site) SymUninstall(timestamp string) {
	Symlink := Site.Path + ".latest"
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
	Site.Timestamp = fmt.Sprintf(".%v", value)
}

func (Site *Site) TimeStampReset() {
	now := time.Now()
	Site.Timestamp = fmt.Sprintf(".%v", now.Format("20060102150405"))
}

func (Site *Site) ProcessMake(makeFile string) {

	fullPath := makeFile
	_, err := os.Stat(fullPath)
	if err != nil {
		log.Println("Error! File not found:", err)
		os.Exit(1)
	}

	drushCommand := fmt.Sprintf("make -y --overwrite --working-copy %v %v%v", fullPath, Site.Path, Site.Timestamp)
	log.Println("Building from", makeFile)
	drushMake := command.NewDrushCommand()
	drushMake.Set("", drushCommand, true)
	cmd, err := drushMake.Output()
	if err != nil {
		if string(err.Error()) == "exit status 1" {
			log.Println("Processed make file was completed with errors. :", err.Error())
			log.Printf("Could not execute `drush %v`", drushCommand)
		}
	}
	log.Printf("Creating directory for site at %v", Site.Path+Site.Timestamp)
	drushLog := cmd
	for _, logEntry := range drushLog {
		// Print output in a fairly standardized format.
		log.Println(logEntry)
	}
}

func (Site *Site) RebuildRegistry() {
	drushCommand := command.NewDrushCommand()
	drushCommand.Set(Site.Alias, "rr", false)
	_, err := drushCommand.Output()
	if err != nil {
		log.Println("Could not rebuild registry...")
	} else {
		log.Println("Rebuilt registry...")
	}
}

func (Site *Site) InstallSiteRef() {
	log.Println("Adding", Site.Path+Site.Timestamp+"/sites/sites.php")
	data := map[string]string{
		"Name":  Site.Name,
		"Alias": Site.Alias,
	}
	filename := Site.Path + Site.Timestamp + "/sites/sites.php"
	tpl, err := template.ParseFiles("templates/sites-template.gotpl")
	if err != nil {
		log.Fatalln(err)
	}

	nf, err := os.Create(filename)
	if err != nil {
		log.Fatalln("error creating file", err)
	}
	defer nf.Close()

	err = tpl.Execute(nf, data)
	if err != nil {
		log.Fatalln(err)
	}
}
