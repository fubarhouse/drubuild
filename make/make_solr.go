package make

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

type SolrCore struct {
	Address  string
	Name     string
	Template string
	Path     string
}

func logSolrInstall() bool {
	if verifySolrInstall() {
		log.Infoln("Found Solr installation")
		return true
	} else {
		log.Errorln("Could not find Solr installation")
		return false
	}
}
func logSolrCLI() bool {
	if verifySolrCLI() {
		log.Infoln("Found Solr command-line tools")
		return true
	} else {
		log.Warnln("Could not find Solr command-line tools")
		return false
	}
}
func logResources(Template string) bool {
	if verifyResources(Template) {
		log.Infoln("Found configuration folder")
		return true
	} else {
		log.Errorln("Could not find configuration folder")
		return false
	}
}
func logSolrCore(SolrCore *SolrCore) bool {
	if verifySolrCore(SolrCore) {
		log.Infoln("Solr core is installed.")
		return true
	} else {
		log.Warnln("Solr core is not installed.")
		return false
	}
}
func verifySolrInstall() bool {
	_, err := os.Stat("/opt/solr")
	if err == nil {
		return true
	} else {
		return false
	}
}
func verifySolrCLI() bool {
	_, err := os.Stat("/opt/solr/bin/solr")
	if err == nil {
		return true
	} else {
		return false
	}
}
func verifyResources(Template string) bool {
	_, err := os.Stat(Template)
	if err == nil {
		return true
	} else {
		return false
	}
}
func verifySolrCore(SolrCore *SolrCore) bool {
	curlResponse, err := exec.Command("curl", "http://"+SolrCore.Address+"/solr/admin/cores?action=STATUS").Output()
	if err == nil {
		if strings.Contains(string(curlResponse), `<str name="name">`+SolrCore.Name+`</str>`) == true {
			return true
		} else {
			return false
		}
	} else {
		log.Errorln("Solr could not be accessed using CURL:", err.Error())
	}
	return false
}

func NewCore(Address, Name, Template, Path string) SolrCore {
	return SolrCore{Address, Name, Template, Path}
}

func (SolrCore *SolrCore) Install() {
	if logSolrInstall() && logResources(SolrCore.Template) {
		log.Infoln("All checks have passed.")
		if logSolrCLI() == true {

			// Solr version >= 5

			// Create data directories

			err := os.MkdirAll(SolrCore.Path+"/"+SolrCore.Name+"/data/index", 0777)
			if err == nil {
				log.Infoln("Directory has been created.")
			} else {
				log.Errorln("Directory has not been created:", err.Error())
			}

			// Sync

			_, err = exec.Command("rsync", "-a", SolrCore.Template+"/", SolrCore.Path+"/"+SolrCore.Name+"/conf/").Output()
			if err == nil {
				log.Infoln("Configuration has been synced with boilerplate resources.")
			} else {
				log.Errorln("Configuration could not be synced with boilerplate resources:", err.Error())
			}

			// Install core via CLI

			cliOut, err := exec.Command("/opt/solr/bin/solr", "create", "-c", SolrCore.Name).Output()
			if err == nil && strings.Contains(string(cliOut), "Unable to create core") == false {
				if verifySolrCore(SolrCore) {
					log.Infoln("Core has been installed and verified successfully.")
				} else {
					log.Errorln("Core could not be installed, check logs for details.")
				}
			} else {
				log.Errorln("Core could not be installed:", err.Error(), string(cliOut))
			}

		} else {

			// Solr version <= 4

			// Create data directories
			err := os.MkdirAll(SolrCore.Path+"/"+SolrCore.Name+"/", 0777)
			if err == nil {
				log.Infoln("Directory has been created.")
			} else {
				log.Errorln("Directory has not been created:", err.Error())
			}

			// Sync

			_, err = exec.Command("rsync", "-a", SolrCore.Template+"/", SolrCore.Path+"/"+SolrCore.Name+"/conf/").Output()
			if err == nil {
				log.Infoln("Configuration has been synced with boilerplate resources.")
			} else {
				log.Errorln("Configuration could not be synced with boilerplate resources:", err.Error())
			}

			// Install core via CURL

			_, err = exec.Command("curl", "http://"+SolrCore.Address+"/solr/admin/cores?action=CREATE&name="+SolrCore.Name+"&instanceDir="+SolrCore.Name+"&dataDir=data&config=solrconfig.xml&schema=schema.xml").Output()
			if err == nil {
				log.Infoln("Core has been successfully installed.")
			} else {
				log.Errorln("Core could not be installed:", err)
			}
		}
	}
	verifySolrCore(SolrCore)
}

func (SolrCore *SolrCore) Uninstall() {
	if logSolrInstall() && logSolrCLI() && logResources(SolrCore.Template) && logSolrCore(SolrCore) || !logSolrCore(SolrCore) {
		log.Infoln("All checks have passed.")
		// Uninstall core
		_, err := exec.Command("sh", "-c", "/opt/solr/bin/solr", "delete", "-c", SolrCore.Name).Output()
		if err == nil {
			log.Infoln("Core has been successfully uninstalled.")
		} else {
			log.Errorln("Core could not be uninstalled:", err)
		}
		// Delete resources.
		err = os.RemoveAll(SolrCore.Path + "/" + SolrCore.Name)
		if err == nil {
			log.Infoln("Core resources have been removed.")
		} else {
			log.Errorln("Core resources could not be removed:", err)
		}
	}
}
