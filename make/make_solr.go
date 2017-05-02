package make

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// struct for solr core
type SolrCore struct {
	Address       string      // HTTP address of solr.
	Binary        string      // Path to solr binary.
	Name          string      // Name of core as subject.
	Template      string      // Directory which contains data.
	Path          string      // Path of Solr.
	SubDir        string      // Subdirectory in the solr directory (Path).
	DataDir       string      // Path to core data inside of solr data directory.
	ConfigFile    string      // Name of config file.
	SchemaFile    string      // Name of schema file.
	SolrUserName  string      // The name of the solr user.
	SolrUserGroup string      // The name of the solr user.
	SolrUserMode  os.FileMode // The mode of the solr core.
}

// Prints the status of the solr CLI tool
func (SolrCore *SolrCore) LogCLI() bool {
	if SolrCore.VerifyCLI() {
		log.Infoln("Found Solr command-line tools")
		return true
	} else {
		log.Warnln("Could not find Solr command-line tools")
		return false
	}
}

// Returns the status of the solr CLI tool
func (SolrCore *SolrCore) VerifyCLI() bool {
	if SolrCore.Binary != "" {
		_, err := os.Stat(SolrCore.Binary)
		if err == nil {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

// Prints the status of the solr installation
func (SolrCore *SolrCore) SolrInstallLog() bool {
	if SolrCore.VerifyInstall() {
		log.Infoln("Found Solr installation")
		return true
	} else {
		log.Errorln("Could not find Solr installation")
		return false
	}
}

// Returns the status of the solr installation
func (SolrCore *SolrCore) SolrResourcesLog(Template string) bool {
	if SolrCore.VerifyResources() {
		log.Infoln("Found configuration folder")
		return true
	} else {
		log.Errorln("Could not find configuration folder")
		return false
	}
}

// Prints the status of a solr core
func (SolrCore *SolrCore) LogCore() bool {
	if SolrCore.VerifyCore() {
		log.Infoln("Solr core is installed.")
		return true
	} else {
		log.Warnln("Solr core is not installed.")
		return false
	}
}

// Returns the status of a solr core
func (SolrCore *SolrCore) VerifyInstall() bool {
	_, err := os.Stat(SolrCore.Path)
	if err == nil {
		return true
	} else {
		return false
	}
}

// Returns the availability of input Resources
func (SolrCore *SolrCore) VerifyResources() bool {
	_, err := os.Stat(SolrCore.Template)
	if err == nil {
		return true
	} else {
		return false
	}
}

// Attempt to verify the install status of a solr sore
func (SolrCore *SolrCore) VerifyCore() bool {

	response, err := http.Get(SolrCore.Address + "/solr/admin/cores?action=STATUS")
	content, err := ioutil.ReadAll(response.Body)
	response.Body.Close()

	if err == nil {
		if strings.Contains(string(content), `<int name="status">500</int>`) {
			log.Fatalln("Verification process has returned error 500.")
			return false
		}
		if strings.Contains(string(content), `<str name="name">`+SolrCore.Name+`</str>`) {
			return true
		} else {
			return false
		}
	} else {
		log.Errorln("Solr did not respond:", err.Error(), strings.Replace(string(content), "\n", "", -1))
		return false
	}
}

// Create the data directory for a solr core
func (SolrCore *SolrCore) _createDataDir() bool {
	directoryArgs := []string{SolrCore.Path, SolrCore.SubDir, SolrCore.Name, SolrCore.DataDir}
	actualPath := strings.Join(directoryArgs, "/")
	err := os.MkdirAll(actualPath, SolrCore.SolrUserMode)
	if err == nil {
		log.Infof("Data directory %v has been created.", actualPath)
		return true
	} else {
		log.Errorf("Data directory %v has not been created. %v", actualPath, err.Error())
		return false
	}
}

// Create the config directory for a solr core
func (SolrCore *SolrCore) _createConfigDir() bool {
	directoryArgs := []string{SolrCore.Path, SolrCore.SubDir, SolrCore.Name, "conf"}
	actualPath := strings.Join(directoryArgs, "/")
	err := os.MkdirAll(actualPath, SolrCore.SolrUserMode)
	if err == nil {
		log.Infof("Data directory %v has been created.", actualPath)
		return true
	} else {
		log.Errorf("Data directory %v has not been created. %v", actualPath, err.Error())
		return false
	}
}

// Copy resources for a solr core
func (SolrCore *SolrCore) _copyTemplateData() {
	directoryArgs := []string{SolrCore.Path, SolrCore.SubDir, SolrCore.Name, "conf"}
	actualPath := strings.Join(directoryArgs, "/")
	err := new(error)
	_ = filepath.Walk(SolrCore.Template, func(path string, FileInfo os.FileInfo, _ error) error {
		realpath := strings.Split(string(path), "\n")
		for _, name := range realpath {
			if !FileInfo.IsDir() {
				data, err := ioutil.ReadFile(name)
				if err != nil {
					log.Errorf("Could not read from %v: %v", name, err.Error())
				}
				err = ioutil.WriteFile(actualPath+"/"+FileInfo.Name(), data, SolrCore.SolrUserMode)
				if err == nil {
					log.Infof("Copied %v to %v", name, actualPath+"/"+FileInfo.Name())
				} else {
					log.Errorf("Could not copy %v to %v: %v", name, actualPath+"/"+FileInfo.Name(), err.Error())
				}
			}
		}
		return *err
	})
}

// Reset mode on a solr core
func (SolrCore *SolrCore) _resetModeTemplateData() bool {
	directoryArgs := []string{SolrCore.Path, SolrCore.SubDir, SolrCore.Name}
	actualPath := strings.Join(directoryArgs, "/")
	FileMode := fmt.Sprintf("%04o", SolrCore.SolrUserMode)
	_, err := exec.Command("chmod", "-R", FileMode, actualPath).Output()
	if err == nil {
		log.Infof("Changed mode of %v to %v", actualPath, SolrCore.SolrUserMode)
		return true
	} else {
		log.Warnf("Could not change mode of %v to %v: %v", actualPath, SolrCore.SolrUserMode, err)
		return false
	}
}

// Reset owner on a solr core
func (SolrCore *SolrCore) _resetOwnerTemplateData() bool {
	directoryArgs := []string{SolrCore.Path, SolrCore.SubDir, SolrCore.Name}
	actualPath := strings.Join(directoryArgs, "/")
	_, err := exec.Command("chown", "-R", SolrCore.SolrUserName, actualPath).Output()
	if err == nil {
		log.Infof("Changed ownership of %v to %v", actualPath, SolrCore.SolrUserName)
		return true
	} else {
		log.Warnf("Could not change ownership of %v to %v: %v", actualPath, SolrCore.SolrUserName, err)
		return false
	}
}

// Reset group on a solr core
func (SolrCore *SolrCore) _resetGroupTemplateData() bool {
	directoryArgs := []string{SolrCore.Path, SolrCore.SubDir, SolrCore.Name}
	actualPath := strings.Join(directoryArgs, "/")
	_, err := exec.Command("chgrp", "-R", SolrCore.SolrUserGroup, actualPath).Output()
	if err == nil {
		log.Infof("Changed group of %v to %v", actualPath, SolrCore.SolrUserGroup)
		return true
	} else {
		log.Warnf("Could not change group of %v to %v: %v", actualPath, SolrCore.SolrUserGroup, err)
		return false
	}
}

// Remove resources for a given core
func (SolrCore *SolrCore) _deleteTemplateData() bool {
	dataDir := strings.Join([]string{SolrCore.Path, SolrCore.SubDir, SolrCore.Name}, "/")
	_, err := os.Stat(dataDir)
	if err == nil {
		err := os.RemoveAll(dataDir)
		if err == nil {
			log.Infoln("Core resources have been removed.")
			return true
		} else {
			log.Errorln("Core resources could not be removed:", err)
			return false
		}
	} else {
		return true
	}
}

// Run the CLI tool to create a solr core
func (SolrCore *SolrCore) _createCore() bool {
	if !SolrCore.VerifyCore() {
		if SolrCore.VerifyCLI() {
			log.Infoln("Installing with specified Solr binary")
			out, err := exec.Command(SolrCore.Binary, "create", "-c", SolrCore.Name).Output()
			if err == nil {
				log.Infoln("Creation command completed successfully.")
			} else {
				log.Errorln("Creation command could not complete:", err, string(out))
			}
			if SolrCore.VerifyCore() {
				log.Infoln("Solr core installation succeeded.")
				return true
			} else {
				log.Errorln("Solr core installation failed.")
				return false
			}
		} else {
			response, err := http.Get(SolrCore.Address + "/solr/admin/cores?action=CREATE&name=" + SolrCore.Name + "&instanceDir=" + SolrCore.Name + "&dataDir=" + SolrCore.DataDir + "&config=" + SolrCore.ConfigFile + "&schema=" + SolrCore.SchemaFile)
			_, err = ioutil.ReadAll(response.Body)
			response.Body.Close()
			if err == nil {
				fmt.Sprint("Creation command completed successfully.")
			} else {
				fmt.Sprint("Creation command could not complete:", err)
			}
			if SolrCore.VerifyCore() {
				log.Infoln("Solr core installation succeeded.")
				return true
			} else {
				log.Errorln("Solr core installation failed.")
				return false
			}
		}
	} else {
		log.Infof("Solr core %v is already created.", SolrCore.Name)
	}
	return SolrCore.VerifyCore()
}

// Run the CLI tool for a solr core removal
func (SolrCore *SolrCore) _deleteCore() bool {
	if SolrCore.VerifyCore() {
		if SolrCore.VerifyCLI() {
			log.Infoln("Uninstalling using specified Solr binary")
			_, err := exec.Command(SolrCore.Binary, "delete", "-c", SolrCore.Name).Output()
			if err == nil {
				log.Infoln("Core has been successfully uninstalled.")
				return true
			} else {
				log.Errorln("Core could not be uninstalled:", err)
				return false
			}
		} else {
			response, err := http.Get(SolrCore.Address + "/solr/admin/cores?action=UNLOAD&core=" + SolrCore.Name + "&deleteIndex=true&deleteDataDir=true&deleteInstanceDir=true")
			_, err = ioutil.ReadAll(response.Body)
			response.Body.Close()
			if err == nil {
				log.Infoln("Core has been successfully uninstalled.")
				return true
			} else {
				log.Errorln("Core could not be uninstalled:", err)
				return false
			}
		}
	} else {
		log.Infof("Solr core %v does not exist.", SolrCore.Name)
		return true
	}
}

// API call to install a solr core.
func (SolrCore *SolrCore) Install() {
	if SolrCore.VerifyInstall() && SolrCore.VerifyResources() {
		log.Infoln("All checks have passed.")

		SolrCore._deleteTemplateData()
		SolrCore._createDataDir()
		SolrCore._createConfigDir()
		SolrCore._copyTemplateData()
		SolrCore._resetModeTemplateData()
		SolrCore._resetOwnerTemplateData()
		SolrCore._resetGroupTemplateData()
		SolrCore._deleteCore()
		SolrCore._createCore()

	} else {
		if !SolrCore.VerifyInstall() {
			log.Errorln("An error was found trying to verify the installation of Solr.")
		}
		if !SolrCore.VerifyResources() {
			log.Errorln("An error was found trying to find the specified templated Resources.")
		}
	}
}

// API call to un-install a solr core.
func (SolrCore *SolrCore) Uninstall() {
	SolrCore._deleteCore()
	SolrCore._deleteTemplateData()
}

// Instantiate a new Solr core
func NewCore(Address, Binary, Name, Template, Path, SubDir, DataDir, ConfigFile, SchemaFile, Subpath, Datapath, SolrUserName, SolrUserGroup string, SolrUserMode os.FileMode) SolrCore {
	return SolrCore{Address, Binary, Name, Template, Path, SubDir, DataDir, ConfigFile, SchemaFile, SolrUserName, SolrUserGroup, SolrUserMode}
}
