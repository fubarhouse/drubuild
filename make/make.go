package make

import (
	"bufio"
	"database/sql"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/golang-drush/command"
	_ "github.com/go-sql-driver/mysql" // mysql is assumed under this system (for now).
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ReplaceTextInFile is a utility function to replace all instances of a string in a file.
func ReplaceTextInFile(fullPath string, oldString string, newString string) {
	read, err := ioutil.ReadFile(fullPath)
	if err != nil {
		log.Panicln(err)
	}
	newContents := strings.Replace(string(read), oldString, newString, -1)
	err = ioutil.WriteFile(fullPath, []byte(newContents), 0)
	if err != nil {
		log.Panicln(err)
	}
}

// RestartWebServer is a function to run a command to restart the given web service.
func (Site *Site) RestartWebServer() {
	_, stdErr := exec.Command("sudo", "service", Site.Webserver, "restart").Output()
	if stdErr != nil {
		log.Errorf("Could not restart webserver %v. %v\n", Site.Webserver, stdErr)
	} else {
		log.Infof("Restarted webserver %v.\n", Site.Webserver)
	}
}

// StartWebServer is a function to run a command to start the given web service.
func (Site *Site) StartWebServer() {
	_, stdErr := exec.Command("sudo", "service", Site.Webserver, "start").Output()
	if stdErr != nil {
		log.Errorf("Could not start webserver %v. %v\n", Site.Webserver, stdErr)
	} else {
		log.Infof("Started webserver %v.\n", Site.Webserver)
	}
}

// StopWebServer is a function to run a command to stop the given web service.
func (Site *Site) StopWebServer() {
	_, stdErr := exec.Command("sudo", "service", Site.Webserver, "stop").Output()
	if stdErr != nil {
		log.Errorf("Could not stop webserver %v. %v\n", Site.Webserver, stdErr)
	} else {
		log.Infof("Stopped webserver %v.\n", Site.Webserver)
	}
}

// DrupalProject struct which represents a Drupal project on drupal.org
type DrupalProject struct {
	Type   string
	Name   string
	Subdir string
	Status bool
}

// Site struct which represents a build website being used.
type Site struct {
	Timestamp                  string
	Path                       string
	Make                       string
	Name                       string
	Alias                      string
	Domain                     string
	database                   *makeDB
	Webserver                  string
	Vhostpath                  string
	Template                   string
	MakeFileRewriteSource      string
	MakeFileRewriteDestination string
	FilePathPrivate            string
	FilePathPublic             string
	FilePathTemp               string
	WorkingCopy                bool
}

// NewSite instantiates an instance of the struct Site
func NewSite(make, name, path, alias, webserver, domain, vhostpath, template string) *Site {
	Site := &Site{}
	Site.TimeStampReset()
	Site.Make = make
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
	Site.MakeFileRewriteSource = ""
	Site.MakeFileRewriteDestination = ""
	Site.WorkingCopy = false
	return Site
}

// ActionBuild is a superseded build action, requires action, documentation or removal.
func (Site *Site) ActionBuild() {
	// TODO: Define purpose with the existence of ProcessMake()
	if Site.AliasExists(Site.Name) == true {
		Site.Path = fmt.Sprintf("%v%v", Site.Path, Site.TimeStampGet())
		//Site.ProcessMakes([]string{"core.make", "libraries.make", "contrib.make", "custom.make"})
		Site.ActionInstall()
	}
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
	thisCmd := fmt.Sprintf("-y site-install standard --sites-subdir=%v --db-url=mysql://%v:%v@%v:%v/%v install_configure_form.update_status_module='array(FALSE,FALSE)'", Site.Name, Site.database.getUser(), Site.database.getPass(), Site.database.getHost(), Site.database.getPort(), dbName)
	_, installErr := exec.Command("sh", "-c", "cd "+Site.Path+"/"+Site.Name+Site.Timestamp+" && drush "+thisCmd).Output()
	if installErr != nil {
		log.Warnln("Unable to install Drupal.")
		log.Debugln("drush", thisCmd)
	} else {
		log.Infof("Installed Drupal.")
	}
}

// ActionKill will delete a single site instance.
func (Site *Site) ActionKill() {
	// What to do with the default...
	if Site.AliasExists(Site.Name) == true {
		Site.Path = fmt.Sprintf("%v", Site.Path)
		_, err := os.Stat(Site.Path)
		if err == nil {
			os.Remove(Site.Path)
		}
	}
}

// ActionRebuild rebuild site structs, needs action, documentation or purging.
func (Site *Site) ActionRebuild() {
	// TODO: Define purpose with the existence of ProcessMake()
	if Site.AliasExists(Site.Name) == true {
		Site.TimeStampReset()
		Site.Path = fmt.Sprintf("%v%v", Site.Path, Site.TimeStampGet())
		//Site.ProcessMake()
		//Site.ActionInstall()
	}
}

// ActionRebuildProject purges a specific project from a specified path, and re-download it
// Re-downloading will use drush dl, or git clone depending on availability.
func (Site *Site) ActionRebuildProject(Makefiles []string, Project string, GitPath, Branch string, RemoveGit bool) {
	log.Infoln("Searching for module/theme...")
	moduleFound := false
	var moduleType string
	var moduleCat string
	err := new(error)
	_ = filepath.Walk(Site.Path, func(path string, _ os.FileInfo, _ error) error {
		realpath := strings.Split(string(path), "\n")
		for _, name := range realpath {
			if strings.Contains(name, "/contrib/"+Project+"/") || strings.Contains(name, "/custom/"+Project+"/") {
				if strings.Contains(name, "/contrib/"+Project+"/") {
					moduleType = "contrib"
				} else {
					moduleType = "custom"
				}
				if strings.Contains(name, "/modules/"+moduleType+"/"+Project+"/") {
					moduleCat = "modules"
				} else if strings.Contains(name, "/themes/"+moduleType+"/"+Project+"/") {
					moduleCat = "themes"
				}
				moduleFound = true
			}
		}
		return nil
	})

	if moduleFound {
		log.Infoln("Found module at", Site.Path+"/sites/all/"+moduleCat+"/"+moduleType+"/"+Project+"/")
	}

	if moduleType != "" && moduleCat != "" {
		ProjectDir := Site.Path + "/sites/all/" + moduleCat + "/" + moduleType + "/" + Project + "/"
		_, errMod := os.Stat(ProjectDir)
		if errMod == nil {
			*err = os.RemoveAll(ProjectDir)
			if *err == nil {
				log.Infoln("Removed", ProjectDir)
			} else {
				log.Warn("Could not remove ", ProjectDir)
			}
		}
	}

	if moduleFound == false {
		log.Infof("Could not find project %v in %v", Project, Site.Path)
	}
	if moduleCat == "" || moduleType == "" {
		// By this point, we should fall back to the input make file.
		for _, val := range Makefiles {
			unprocessedMakes, unprocessedMakeErr := ioutil.ReadFile(val)
			if unprocessedMakeErr != nil {
				log.Warnf("Could not read from %v: %v", val, unprocessedMakeErr)
			}
			Projects := strings.Split(string(unprocessedMakes), "\n")
			for _, ThisProject := range Projects {
				if strings.Contains(ThisProject, "projects["+Project+"][subdir] = ") {
					moduleType = strings.Replace(ThisProject, "projects["+Project+"][subdir] = ", "", -1)
					moduleType = strings.Replace(moduleType, "\"", "", -1)
					moduleType = strings.Replace(moduleType, " ", "", -1)
				}
				if strings.Contains(ThisProject, "projects["+Project+"][type] = ") {
					moduleCat = strings.Replace(ThisProject, "projects["+Project+"][type] = ", "", -1)
					moduleCat = strings.Replace(moduleCat, "\"", "", -1)
					moduleCat = strings.Replace(moduleCat, " ", "", -1)
				}
			}
		}
		if moduleCat == "" {
			log.Warnln("Project category could not be detected.")
		} else {
			log.Infoln("Project category was found to be", moduleCat)
		}
		if moduleType == "" {
			log.Warnln("Project type could not be detected.")
		} else {
			log.Infoln("Project type was found to be", moduleCat)
		}
	}
	path := Site.Path + "/" + "/sites/all/" + moduleCat + "/" + moduleType + "/"
	if moduleType == "contrib" {
		command.DrushDownloadToPath(path, Project)
	} else {
		gitCmd := exec.Command("git", "clone", "-b", Branch, GitPath, path+"/"+Project)
		_, *err = gitCmd.Output()
		if *err == nil {
			log.Infof("Downloaded package %v from %v to %v", Project, GitPath, path+"/"+Project)
			if RemoveGit {
				*err = os.RemoveAll(path + "/" + Project + "/.git")
				if *err == nil {
					log.Infoln("Removed .git folder from file system.")
				} else {
					log.Warnln("Unable to remove .git folder from file system.")
				}
			}
		} else {
			log.Errorf("Could not clone %v from %v: %v\n", Project, GitPath, *err)
		}
	}
}

// ActionRebuildCodebase re-runs drush make on a specified path.
func (Site *Site) ActionRebuildCodebase(Makefiles []string) {
	// This function exists for the sole purpose of
	// rebuilding a specific Drupal codebase in a specific
	// directory for Release management type work.
	var newMakeFilePath string
	if Site.Timestamp == "." {
		Site.Timestamp = ""
		newMakeFilePath = "/tmp/drupal-" + Site.Name + Site.TimeStampGenerate() + ".make"
	} else {
		newMakeFilePath = "/tmp/drupal-" + Site.Name + Site.TimeStampGet() + ".make"
	}
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

	chmodErr := os.Chmod(Site.Path, 0777)
	if chmodErr != nil {
		log.Warnln("Could not change permissions on codebase directory")
	} else {
		log.Infoln("Changed docroot permissions to 0777 for file removal.")
	}

	_ = filepath.Walk(Site.Path, func(path string, Info os.FileInfo, _ error) error {
		realpath := strings.Split(Site.Path, "\n")
		err := new(error)
		for _, name := range realpath {
			fmt.Sprintln(name)
			if !strings.Contains(path, "/sites") || strings.Contains(path, "/sites/all") {
				if Info.IsDir() && !strings.HasSuffix(path, Site.Path) {
					os.Chmod(path, 0777)
					delErr := os.RemoveAll(path)
					if delErr != nil {
						log.Warnln("Could not remove", path)
					}
				} else if !Info.IsDir() {
					delErr := os.Remove(path)
					if delErr != nil {
						log.Warnln("Could not remove", path)
					}
				}
			}
		}
		return *err
	})

	Site.ProcessMake(newMakeFilePath)
	err := os.Remove(newMakeFilePath)
	if err != nil {
		log.Warnln("Could not remove temporary make file", newMakeFilePath)
	} else {
		log.Infoln("Removed temporary make file", newMakeFilePath)
	}
}

// ActionDatabaseDumpLocal run drush sql-dump to a specified path on a site struct.
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

// ActionDatabaseDumpRemote run drush sql-dump to a specified path on a site alias.
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

// TimeStampSet sets the timestamp field for the site struct to a given value
func (Site *Site) TimeStampSet(value string) {
	Site.Timestamp = fmt.Sprintf(".%v", value)
}

// TimeStampReset sets the timestamp field for the site struct to a new value
func (Site *Site) TimeStampReset() {
	now := time.Now()
	Site.Timestamp = fmt.Sprintf(".%v", now.Format("20060102150405"))
}

// TimeStampGenerate generates a new timestamp and returns it, does not latch to site struct
func (Site *Site) TimeStampGenerate() string {
	return fmt.Sprintf(".%v", time.Now().Format("20060102150405"))
}

// VerifyProcessedMake requires documentation, @TODO for revisitation.
func (Site *Site) VerifyProcessedMake(makeFile string) []DrupalProject {
	unprocessedMakes, unprocessedMakeErr := ioutil.ReadFile(makeFile)
	Projects := make([]DrupalProject, 50)
	if unprocessedMakeErr != nil {
		log.Infoln("Could not read from", unprocessedMakeErr)
	}
	for _, Line := range strings.Split(string(unprocessedMakes), "\n") {
		var Type string
		if strings.Contains(Line, "subdir") || strings.Contains(Line, "directory_name") {
			currentType := strings.SplitAfter(Line, "=")
			Type = strings.Replace(currentType[1], "\"", "", -1)
			Type = strings.Replace(Type, " ", "", -1)
		}
		if Type != "" {
			if strings.HasPrefix(Line, "projects") {
				Project := strings.SplitAfter(Line, "[")
				Project[1] = strings.Replace(Project[1], "[", "", -1)
				Project[1] = strings.Replace(Project[1], "]", "", -1)
				thisProject := DrupalProject{"modules", Project[1], Type, false}
				Projects = append(Projects, thisProject)
			}
			if strings.HasPrefix(Line, "libraries") {
				Library := strings.SplitAfter(Line, "[")
				Library[1] = strings.Replace(Library[1], "[", "", -1)
				Library[1] = strings.Replace(Library[1], "]", "", -1)
				thisProject := DrupalProject{"libraries", Library[1], Type, false}
				Projects = append(Projects, thisProject)
			}
		}
	}
	var foundModules int
	for index, Project := range Projects {
		if Project.Name != "" {
			//log.Printf("Package %v is of type %v, belonging to subdir %v", Project.Name, Project.Type, Project.Subdir)
			err := new(error)
			_ = filepath.Walk(Site.Path, func(path string, _ os.FileInfo, _ error) error {
				realpath := strings.Split(Site.Path, "\n")
				for _, name := range realpath {
					if strings.Contains(path, "custom/"+Project.Name+"/") || strings.Contains(path, "contrib/"+Project.Name+"/") || strings.Contains(path, "libraries/"+Project.Subdir+"/") {
						fmt.Sprintln(name)
						foundModules++
						Projects[index].Status = true
						break
					}
				}
				return *err
			})
		}
	}
	return Projects
}

// ProcessMake processes a make file at a particular path.
func (Site *Site) ProcessMake(makeFile string) bool {

	// Test the make file exists
	fullPath := makeFile
	_, err := os.Stat(fullPath)
	if err != nil {
		log.Fatalln("File not found:", err)
		os.Exit(1)
	}
	if Site.MakeFileRewriteSource != "" && Site.MakeFileRewriteDestination != "" {
		log.Printf("Applying specified rewrite string on temporary makefile: %v -> %v", Site.MakeFileRewriteSource, Site.MakeFileRewriteDestination)
		ReplaceTextInFile(makeFile, Site.MakeFileRewriteSource, Site.MakeFileRewriteDestination)
	} else {
		log.Println("No rewrite string was configured, continuing without additional parsing.")
	}

	log.Infof("Building from %v...", makeFile)
	drushMake := command.NewDrushCommand()
	drushCommand := ""
	if Site.WorkingCopy {
		drushCommand = fmt.Sprintf("make --yes --working-copy %v", makeFile)
	} else {
		drushCommand = fmt.Sprintf("make --yes %v", makeFile)
	}
	drushMake.Set("", drushCommand, false)
	if Site.Timestamp == "" {
		drushMake.SetWorkingDir(Site.Path + "/")
	} else {
		drushMake.SetWorkingDir(Site.Path + "/" + Site.Name + Site.Timestamp)
	}
	mkdirErr := os.MkdirAll(drushMake.GetWorkingDir(), 0755)
	if mkdirErr != nil {
		log.Warnln("Could not create directory", drushMake.GetWorkingDir())
	} else {
		log.Infoln("Created directory", drushMake.GetWorkingDir())
	}
	_ = drushMake.LiveOutput()
	return true
}

// InstallSiteRef installs the Drupal multisite sites.php file for the site struct.
func (Site *Site) InstallSiteRef() {

	data := map[string]string{
		"Name":   Site.Name,
		"Domain": Site.Domain,
	}
	dirPath := Site.Path + "/" + Site.Name + Site.Timestamp + "/sites/"
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
	buffer := []byte{60, 63, 112, 104, 112, 10, 10, 47, 42, 42, 10, 32, 42, 32, 64, 102, 105, 108, 101, 10, 32, 42, 32, 67, 111, 110, 102, 105, 103, 117, 114, 97, 116, 105, 111, 110, 32, 102, 105, 108, 101, 32, 102, 111, 114, 32, 68, 114, 117, 112, 97, 108, 39, 115, 32, 109, 117, 108, 116, 105, 45, 115, 105, 116, 101, 32, 100, 105, 114, 101, 99, 116, 111, 114, 121, 32, 97, 108, 105, 97, 115, 105, 110, 103, 32, 102, 101, 97, 116, 117, 114, 101, 46, 10, 32, 42, 47, 10, 10, 32, 32, 32, 36, 115, 105, 116, 101, 115, 91, 39, 68, 111, 109, 97, 105, 110, 39, 93, 32, 61, 32, 39, 78, 97, 109, 101, 39, 59, 10, 10, 63, 62, 10}
	tpl := fmt.Sprintf("%v", string(buffer[:]))
	tpl = strings.Replace(tpl, "Name", data["Name"], -1)
	tpl = strings.Replace(tpl, "Domain", data["Domain"], -1)

	nf, err := os.Create(filename)
	nf.Chmod(0755)
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
