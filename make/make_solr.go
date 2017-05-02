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

// SolrCore is a struct for solr core
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

// LogCLI prints the status of the solr CLI tool
func (SolrCore *SolrCore) LogCLI() bool {
	if SolrCore.VerifyCLI() {
		log.Infoln("Found Solr command-line tools")
	} else {
		log.Warnln("Could not find Solr command-line tools")
	}
	return SolrCore.VerifyCLI()
}

// VerifyCLI returns the status of the solr CLI tool
func (SolrCore *SolrCore) VerifyCLI() bool {
	returnValue := false
	if SolrCore.Binary != "" {
		_, err := os.Stat(SolrCore.Binary)
		if err == nil {
			returnValue = true
		} else {
		}
	}
	return returnValue
}

// SolrInstallLog prints the status of the solr installation
func (SolrCore *SolrCore) SolrInstallLog() bool {
	if SolrCore.VerifyInstall() {
		log.Infoln("Found Solr installation")
	} else {
		log.Errorln("Could not find Solr installation")
	}
	return SolrCore.VerifyInstall()
}

// SolrResourcesLog returns the status of the solr installation
func (SolrCore *SolrCore) SolrResourcesLog(Template string) bool {
	if SolrCore.VerifyResources() {
		log.Infoln("Found configuration folder")
	} else {
		log.Errorln("Could not find configuration folder")
	}
	return SolrCore.VerifyResources()
}

// LogCore prints the status of a solr core
func (SolrCore *SolrCore) LogCore() bool {
	if SolrCore.VerifyCore() {
		log.Infoln("Solr core is installed.")
	} else {
		log.Warnln("Solr core is not installed.")
	}
	return SolrCore.VerifyCore()
}

// VerifyInstall returns the status of a solr core
func (SolrCore *SolrCore) VerifyInstall() bool {
	_, err := os.Stat(SolrCore.Path)
	if err == nil {
		return true
	}
	return false
}

// VerifyResources returns the availability of input Resources
func (SolrCore *SolrCore) VerifyResources() bool {
	_, err := os.Stat(SolrCore.Template)
	if err == nil {
		return true
	}
	return false
}

// VerifyCore will attempt to verify the install status of a solr sore
func (SolrCore *SolrCore) VerifyCore() bool {

	returnValue := false
	response, err := http.Get(SolrCore.Address + "/solr/admin/cores?action=STATUS")
	content, err := ioutil.ReadAll(response.Body)
	response.Body.Close()

	if err == nil {
		if strings.Contains(string(content), `<int name="status">500</int>`) {
			log.Fatalln("Verification process has returned error 500.")
		}
		if strings.Contains(string(content), `<str name="name">`+SolrCore.Name+`</str>`) {
			returnValue = true
		} else {
		}
	} else {
		log.Errorln("Solr did not respond:", err.Error(), strings.Replace(string(content), "\n", "", -1))
	}
	return returnValue
}

// _createDataDir will create the data directory for a solr core
func (SolrCore *SolrCore) _createDataDir() bool {
	returnValue := false
	directoryArgs := []string{SolrCore.Path, SolrCore.SubDir, SolrCore.Name, SolrCore.DataDir}
	actualPath := strings.Join(directoryArgs, "/")
	err := os.MkdirAll(actualPath, SolrCore.SolrUserMode)
	if err == nil {
		log.Infof("Data directory %v has been created.", actualPath)
		returnValue = true
	} else {
		log.Errorf("Data directory %v has not been created. %v", actualPath, err.Error())
	}
	return returnValue
}

// _createConfigDir will create the config directory for a solr core
func (SolrCore *SolrCore) _createConfigDir() bool {
	returnValue := false
	directoryArgs := []string{SolrCore.Path, SolrCore.SubDir, SolrCore.Name, "conf"}
	actualPath := strings.Join(directoryArgs, "/")
	err := os.MkdirAll(actualPath, SolrCore.SolrUserMode)
	if err == nil {
		log.Infof("Data directory %v has been created.", actualPath)
		returnValue = true
	} else {
		log.Errorf("Data directory %v has not been created. %v", actualPath, err.Error())
	}
	return returnValue
}

// _copyTemplateData will copy resources for a solr core
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

// _resetModeTemplateData resets mode on a solr core
func (SolrCore *SolrCore) _resetModeTemplateData() bool {
	returnValue := false
	directoryArgs := []string{SolrCore.Path, SolrCore.SubDir, SolrCore.Name}
	actualPath := strings.Join(directoryArgs, "/")
	FileMode := fmt.Sprintf("%04o", SolrCore.SolrUserMode)
	_, err := exec.Command("chmod", "-R", FileMode, actualPath).Output()
	if err == nil {
		log.Infof("Changed mode of %v to %v", actualPath, SolrCore.SolrUserMode)
		returnValue = true
	} else {
		log.Warnf("Could not change mode of %v to %v: %v", actualPath, SolrCore.SolrUserMode, err)
	}
	return returnValue
}

// _resetOwnerTemplateData resets owner on a solr core
func (SolrCore *SolrCore) _resetOwnerTemplateData() bool {
	returnValue := false
	directoryArgs := []string{SolrCore.Path, SolrCore.SubDir, SolrCore.Name}
	actualPath := strings.Join(directoryArgs, "/")
	_, err := exec.Command("chown", "-R", SolrCore.SolrUserName, actualPath).Output()
	if err == nil {
		log.Infof("Changed ownership of %v to %v", actualPath, SolrCore.SolrUserName)
		returnValue = true
	} else {
		log.Warnf("Could not change ownership of %v to %v: %v", actualPath, SolrCore.SolrUserName, err)
	}
	return returnValue
}

// _resetGroupTemplateData resets group on a solr core
func (SolrCore *SolrCore) _resetGroupTemplateData() bool {
	returnValue := false
	directoryArgs := []string{SolrCore.Path, SolrCore.SubDir, SolrCore.Name}
	actualPath := strings.Join(directoryArgs, "/")
	_, err := exec.Command("chgrp", "-R", SolrCore.SolrUserGroup, actualPath).Output()
	if err == nil {
		log.Infof("Changed group of %v to %v", actualPath, SolrCore.SolrUserGroup)
		returnValue = true
	} else {
		log.Warnf("Could not change group of %v to %v: %v", actualPath, SolrCore.SolrUserGroup, err)
	}
	return returnValue
}

// _deleteTemplateData removes resources for a given core
func (SolrCore *SolrCore) _deleteTemplateData() bool {
	returnValue := false
	dataDir := strings.Join([]string{SolrCore.Path, SolrCore.SubDir, SolrCore.Name}, "/")
	_, err := os.Stat(dataDir)
	if err == nil {
		err := os.RemoveAll(dataDir)
		if err == nil {
			log.Infoln("Core resources have been removed.")
			returnValue = true
		} else {
			log.Errorln("Core resources could not be removed:", err)
		}
	} else {
		returnValue = true
	}
	return returnValue
}

// _createCore runs the CLI tool to create a solr core
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
			} else {
				log.Errorln("Solr core installation failed.")
			}
		} else {
			response, err := http.Get(SolrCore.Address + "/solr/admin/cores?action=CREATE&name=" + SolrCore.Name + "&instanceDir=" + SolrCore.Name + "&dataDir=" + SolrCore.DataDir + "&config=" + SolrCore.ConfigFile + "&schema=" + SolrCore.SchemaFile)
			_, err = ioutil.ReadAll(response.Body)
			response.Body.Close()
			if err == nil {
				log.Infoln("Creation command completed successfully.")
			} else {
				log.Errorln("Creation command could not complete:", err)
			}
			if SolrCore.VerifyCore() {
				log.Infoln("Solr core installation succeeded.")
			} else {
				log.Errorln("Solr core installation failed.")
			}
		}
	} else {
		log.Infof("Solr core %v is already created.", SolrCore.Name)
	}
	return SolrCore.VerifyCore()
}

// _deleteCore runs the CLI tool for a solr core removal
func (SolrCore *SolrCore) _deleteCore() bool {
	returnValue := false
	if SolrCore.VerifyCore() {
		if SolrCore.VerifyCLI() {
			log.Infoln("Uninstalling using specified Solr binary")
			_, err := exec.Command(SolrCore.Binary, "delete", "-c", SolrCore.Name).Output()
			if err == nil {
				log.Infoln("Core has been successfully uninstalled.")
				returnValue = true
			} else {
				log.Errorln("Core could not be uninstalled:", err)
			}
		} else {
			response, err := http.Get(SolrCore.Address + "/solr/admin/cores?action=UNLOAD&core=" + SolrCore.Name + "&deleteIndex=true&deleteDataDir=true&deleteInstanceDir=true")
			_, err = ioutil.ReadAll(response.Body)
			response.Body.Close()
			if err == nil {
				log.Infoln("Core has been successfully uninstalled.")
				returnValue = true
			} else {
				log.Errorln("Core could not be uninstalled:", err)
			}
		}
	} else {
		log.Infof("Solr core %v does not exist.", SolrCore.Name)
		returnValue = true
	}
	return returnValue
}

// Install is an API call to install a solr core.
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

// NewCore is an API call to un-install a solr core.
func (SolrCore *SolrCore) NewCore() {
	SolrCore._deleteCore()
	SolrCore._deleteTemplateData()
}

// NewCore instantiates a new Solr core
func NewCore(Address, Binary, Name, Template, Path, SubDir, DataDir, ConfigFile, SchemaFile, Subpath, Datapath, SolrUserName, SolrUserGroup string, SolrUserMode os.FileMode) SolrCore {
	return SolrCore{Address, Binary, Name, Template, Path, SubDir, DataDir, ConfigFile, SchemaFile, SolrUserName, SolrUserGroup, SolrUserMode}
}
