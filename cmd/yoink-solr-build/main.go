package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/golang-drush/solr"
	"os"
)

func main() {

	// Required fields.
	var Name = flag.String("name", "", "Name of core to create.")
	var Template = flag.String("resources", "", "Path to Solr resources for new cores.")
	var Binary = flag.String("binary", "", "Path to the solr binary, if available.")
	var Path = flag.String("path", "", "Path to Solr folder containing instance directories.")

	// Optional fields, which should be changed where appropriate.
	var Address = flag.String("address", "http://127.0.0.1:8983", "http address of solr installation where solr version >= 5.")
	var SubDir = flag.String("sub-dir", "", "If the core belongs in a subdirectory of the given path value, enter it here.")
	var DataDir = flag.String("data-dir", "data", "The data directory name inside the subject core.")
	var ConfigFile = flag.String("configfile", "solrconfig.xml", "The data directory name inside the subject core.")
	var SchemaFile = flag.String("schemafile", "schema.xml", "The data directory name inside the subject core.")
	var SolrUserName = flag.String("user-name", "solr", "The user name of which the core needs to belong to.")
	var SolrUserGroup = flag.String("user-group", "solr", "The user group of which the core needs to belong to.")
	var SolrUserMode = flag.Int("user-mode", 0777, "The mode of the core affected.")

	flag.Parse()

	if *Name == "" || *Template == "" || *Path == "" {
		flag.Usage()
		log.Fatalln("A value for 'path', 'name' and 'resources' has not beed specified, exiting...")
	}
	if *Binary == "" {
		log.Warnln("A binary path has not been specified, results may vary.")
	}

	// Convert file mode input uint to FileMode
	SolrUserFileMode := os.FileMode(*SolrUserMode)

	SolrCore := solr.SolrCore{*Address, *Binary, *Name, *Template, *Path, *SubDir, *DataDir, *ConfigFile, *SchemaFile, *SolrUserName, *SolrUserGroup, SolrUserFileMode}
	log.Infoln("Starting Solr core uninstallation task.")
	SolrCore.Uninstall()
	log.Infoln("Starting Solr core installation task.")
	SolrCore.Install()
}
