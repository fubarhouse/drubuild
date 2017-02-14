package make

import (
	"bufio"
	"fmt"
	"github.com/fubarhouse/golang-drush/command"
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
	sqlQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %v_%v;", Site.Name, Site.Timestamp)
	sqlUser := fmt.Sprintf("--user=%v", Site.database.getUser())
	sqlPass := fmt.Sprintf("--password=%v", Site.database.getPass())
	_, err := exec.Command("mysql", sqlUser, sqlPass, "-e", sqlQuery).Output()
	if err != nil {
		log.Println("MySQL Error:", err)
	}
	output, _ := exec.Command("mysql", sqlUser, sqlPass, "-e", "show databases;").Output()
	if strings.Contains(string(output), Site.Name+Site.Timestamp) == false {
		log.Printf("Database %v_%v could not be created.\n", Site.Name, Site.Timestamp)
	} else {
		log.Printf("Database %v_%v was successfully created.\n", Site.Name, Site.Timestamp)
	}
	thisCmd := fmt.Sprintf("-y site-install standard --sites-subdir=%v --db-url=mysql://%v:%v@%v:%v/%v_%v install_configure_form.update_status_module='array(FALSE,FALSE)'", Site.Name, Site.database.getUser(), Site.database.getPass(), Site.database.getHost(), Site.database.getPort(), Site.Name, Site.Timestamp)
	output, err = exec.Command("sh", "-c", "cd "+Site.Path+Site.Timestamp+" && drush "+thisCmd).Output()
	_, cpErr := exec.Command("cp", "-f", Site.Path+Site.Timestamp+"/sites/"+Site.Name+"/settings.php", Site.Path+Site.Timestamp+"/sites/default/settings.php").Output()
	if cpErr != nil {
		panic("copy failed")
	}
	if err != nil {
		_, statErr := os.Stat(Site.Path + Site.Timestamp + "sites/" + Site.Name + "/settings.php")
		if statErr == nil {
			log.Println("Drush error:", err)
			log.Println(string(output))
			log.Println("cd " + Site.Path + Site.Timestamp + " && drush " + thisCmd)
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
				if strings.HasPrefix(line, "projects") == true || strings.HasPrefix(line, "libraries") == true || strings.HasPrefix(line, "defaults") == true {
					fmt.Fprintln(writer, line)
				}
			}
		}
	}

	writer.Flush()
	Site.ProcessMake(newMakeFilePath)
	//os.Remove(newMakeFilePath)
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
	err := os.Symlink(Site.Path+Site.TimeStampGet(), Symlink)
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

	_, err = os.Stat(Site.Path)
	if err != nil {
		log.Println("Creating directory for site at", Site.Path+Site.Timestamp)
		os.MkdirAll(Site.Path, 0755)
	}

	drushCommand := ""
	// @TODO: Figure out a way to run make without core, but optionally based on makefile.
	if strings.Contains(makeFile, "core") == true {
		drushCommand = fmt.Sprintf("make -y --overwrite --working-copy %v %v%v", fullPath, Site.Path, Site.Timestamp)
	} else {
		drushCommand = fmt.Sprintf("make -y --overwrite --no-core --working-copy %v %v%v", fullPath, Site.Path, Site.Timestamp)
	}
	log.Println("Building from", makeFile)
	drushMake := command.NewDrushCommand()
	drushMake.Set("", drushCommand, true)
	cmd, err := drushMake.Output()
	if err != nil {
		if string(err.Error()) == "exit status 1" {
			log.Println("Processed make file was completed with errors. :", err.Error())

		}
	} else {
		fmt.Sprintln(cmd)
	}
}
