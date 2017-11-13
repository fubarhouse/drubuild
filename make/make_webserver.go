package make

import (
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/fubarhouse/golang-drush/vhost"
)

// VhostPathSet sets a virtual host path
func (Site *Site) VhostPathSet(value string) {
	Site.Vhostpath = value
}

// VhostInstall install a virtual host
func (Site *Site) VhostInstall() {
	var vhostPath string
	vhostPath = strings.Replace(Site.Path+Site.TimeStampGet(), Site.TimeStampGet(), ".latest/"+Site.Docroot, -1)
	vhostFile := vhost.NewVirtualHost(Site.Name, vhostPath, Site.Webserver, Site.Domain, Site.Vhostpath)

	if Site.Template == "" {
		Site.Template = fmt.Sprintf("%v/src/github.com/fubarhouse/golang-drush/cmd/yoink/templates/vhost-%v.gotpl", os.Getenv("GOPATH"), Site.Webserver)
		log.Printf("No input vhost file, using %v", Site.Template)
	}

	vhostFile.Install(Site.Template)
}

// VhostUninstall un-installs a virtual host
func (Site *Site) VhostUninstall() {
	var vhostPath string
	vhostPath = strings.Replace(Site.Path+Site.TimeStampGet(), Site.TimeStampGet(), Site.Name + "/" + Site.Domain+".latest/"+Site.Docroot, -1)
	vhostFile := vhost.NewVirtualHost(Site.Name, vhostPath, Site.Webserver, Site.Domain, Site.Vhostpath)
	vhostFile.Uninstall()
}

// WebserverSet sets the webserver field for a site struct
func (Site *Site) WebserverSet(value string) {
	Site.Webserver = value
}
