package make

import (
	"bufio"
	"database/sql"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/golang-drush/command"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
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

func (Site *Site) RestartWebServer() {
	_, stdErr := exec.Command("sudo", "service", Site.Webserver, "restart").Output()
	if stdErr != nil {
		log.Errorf("Could not restart webserver %v. %v\n", Site.Webserver, stdErr)
	} else {
		log.Infof("Restarted webserver %v.\n", Site.Webserver)
	}
}

type Site struct {
	Timestamp string
	Path      string
	Make      string
	Name      string
	Alias     string
	Domain    string
	database  *makeDB
	Webserver string
	Vhostpath string
}

func NewSite(make, name, path, alias, webserver, domain, vhostpath string) *Site {
	Site := &Site{}
	Site.TimeStampReset()
	Site.Make = make
	Site.Name = name
	Site.Path = path
	Site.Webserver = webserver
	Site.Alias = alias
	Site.Domain = domain
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
		log.Warnf("WARN:", dbErr)
	}
	// Create database
	dbName := strings.Replace(Site.Name+Site.Timestamp, ".", "_", -1)
	_, dbErr = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if dbErr != nil {
		panic(dbErr)
	}
	// Drush site-install
	thisCmd := fmt.Sprintf("-y site-install standard --sites-subdir=%v --db-url=mysql://%v:%v@%v:%v/%v install_configure_form.update_status_module='array(FALSE,FALSE)'", Site.Name, Site.database.getUser(), Site.database.getPass(), Site.database.getHost(), Site.database.getPort(), dbName)
	_, installErr := exec.Command("sh", "-c", "cd "+Site.Path+"/"+Site.Name+Site.Timestamp+" && drush "+thisCmd).Output()
	if installErr != nil {
		log.Warnln("Unable to install Drupal.")
		log.Debugln("drush", thisCmd)
	} else {
		log.Infof("Installed Drupal.")
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
	if Site.Timestamp == "." {
		Site.Timestamp = ""
	}
	newMakeFilePath := "/tmp/drupal-" + Site.Name + Site.TimeStampGet() + ".make"
	file, crErr := os.Create(newMakeFilePath)
	if crErr == nil {
		log.Infoln("Generated temporary make file...")
	} else {
		log.Errorln("Error creating "+newMakeFilePath+":", crErr)
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
	err := os.Remove(newMakeFilePath)
	if err != nil {
		log.Warnln("Could not remove temporary make file", newMakeFilePath)
	} else {
		log.Infoln("Removed temporary make file", newMakeFilePath)
	}
}

func (Site *Site) ActionDatabaseDumpLocal(path string) {
	srcAlias := strings.Replace(Site.Alias, "@", "", -1)
	x := command.NewDrushCommand()
	x.Set(srcAlias, fmt.Sprintf("sql-dump %v", path), true)
	_, err := x.Output()
	if err == nil {
		log.Println("Dump complete. Dump can be found at", path)
	} else {
		log.Println("Could not dump database.", err)
	}
}

func (Site *Site) ActionDatabaseDumpRemote(alias, path string) {
	srcAlias := strings.Replace(alias, "@", "", -1)
	x := command.NewDrushCommand()
	x.Set(srcAlias, fmt.Sprintf("sql-dump %v", path), true)
	_, err := x.Output()
	if err == nil {
		log.Infoln("Dump complete. Dump can be found at", path)
	} else {
		log.Errorln("Could not dump database.", err)
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
	Symlink := Site.Path + "/" + Site.Name + ".latest"
	err := os.Symlink(Site.Path+"/"+Site.Name+Site.TimeStampGet(), Symlink)
	if err == nil {
		log.Infoln("Created symlink")
	} else {
		log.Warnln("Could not create symlink:", err)
	}
}

func (Site *Site) SymUninstall(timestamp string) {
	Symlink := Site.Path + "/" + Site.Name + ".latest"
	_, statErr := os.Stat(Site.Path + "/" + Symlink)
	if statErr == nil {
		err := os.Remove(Symlink)
		if err != nil {
			log.Errorln("Could not remove symlink.")
		} else {
			log.Infoln("Removed symlink.")
		}
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

	// Test the make file exists
	fullPath := makeFile
	_, err := os.Stat(fullPath)
	if err != nil {
		log.Fatalln("File not found:", err)
		os.Exit(1)
	}

	log.Infof("Building from %v...", makeFile)
	drushMake := command.NewDrushCommand()
	drushCommand := ""
	drushCommand = fmt.Sprintf("make -y --no-core --overwrite --working-copy %v %v/%v%v", fullPath, Site.Path, Site.Name, Site.Timestamp)
	drushMake.Set("", drushCommand, true)
	cmd, err := drushMake.Output()
	if err != nil {
		log.Warnln("Could not execute Drush make without errors.", err.Error())
		log.Warnln("drush", drushCommand)
		drushLog := cmd
		for _, logEntry := range drushLog {
			// Print output in a fairly standardized format.
			logEntryLines := strings.Split(logEntry, "\n")
			for _, logEntryLine := range logEntryLines {
				log.Infoln(logEntryLine)
			}
		}
	} else {
		log.Infoln("Finished building new codebase without errors")
	}
}

func (Site *Site) InstallSiteRef() {

	data := map[string]string{
		"Name":   Site.Name,
		"Domain": Site.Domain,
	}
	filename := Site.Path + "/" + Site.Name + Site.Timestamp + "/sites/sites.php"
	buffer := []byte{60, 63, 112, 104, 112, 10, 10, 47, 42, 42, 10, 32, 42, 32, 64, 102, 105, 108, 101, 10, 32, 42, 32, 67, 111, 110, 102, 105, 103, 117, 114, 97, 116, 105, 111, 110, 32, 102, 105, 108, 101, 32, 102, 111, 114, 32, 68, 114, 117, 112, 97, 108, 39, 115, 32, 109, 117, 108, 116, 105, 45, 115, 105, 116, 101, 32, 100, 105, 114, 101, 99, 116, 111, 114, 121, 32, 97, 108, 105, 97, 115, 105, 110, 103, 32, 102, 101, 97, 116, 117, 114, 101, 46, 10, 32, 42, 47, 10, 10, 32, 32, 32, 36, 115, 105, 116, 101, 115, 91, 39, 68, 111, 109, 97, 105, 110, 39, 93, 32, 61, 32, 39, 78, 97, 109, 101, 39, 59, 10, 10, 63, 62, 10}
	tpl := fmt.Sprintf("%v", string(buffer[:]))
	tpl = strings.Replace(tpl, "Name", data["Name"], -1)
	tpl = strings.Replace(tpl, "Domain", data["Domain"], -1)

	nf, err := os.Create(filename)
	if err != nil {
		log.Fatalln("Could not create", err)
	}
	_, err = nf.WriteString(tpl)
	if err != nil {
		log.Errorln("Could not add", filename)
	} else {
		log.Infoln("Added", filename)
	}
	defer nf.Close()
}

func (Site *Site) InstallPrivateFileSystem() {
	// Test the file system, create it if it doesn't exist!
	filename := "sites/" + Site.Name + "/private"
	dirPath := fmt.Sprintf("%v/%v%v", Site.Path, Site.Name, Site.Timestamp)
	_, err := os.Stat(dirPath + "/" + filename)
	if err != nil {
		dirErr := os.MkdirAll(dirPath, 0755)
		if dirErr != nil {
			log.Errorln("Could not create private file system", filename)
		} else {
			log.Infoln("Created file system", filename)
		}
	}
}
