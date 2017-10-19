package make

import (
	"github.com/fubarhouse/golang-drush/vhost"
	"strings"
	"log"
)

// VhostPathSet sets a virtual host path
func (Site *Site) VhostPathSet(value string) {
	Site.Vhostpath = value
}

// VhostInstall install a virtual host
func (Site *Site) VhostInstall() {
	var vhostPath string
	if Site.Composer {
		vhostPath = strings.Replace(Site.Path+Site.TimeStampGet(), Site.TimeStampGet(),  ".latest/docroot", -1)
	} else {
		vhostPath = strings.Replace(Site.Path+Site.TimeStampGet(), Site.TimeStampGet(), ".latest", -1)
	}
	vhostFile := vhost.NewVirtualHost(Site.Name, vhostPath, Site.Webserver, Site.Domain, Site.Vhostpath)
	vhostFile.Install(Site.Template)
}

// VhostUninstall un-installs a virtual host
func (Site *Site) VhostUninstall() {
	var vhostPath string
	if Site.Composer {
		vhostPath = strings.Replace(Site.Path+Site.TimeStampGet(), Site.TimeStampGet(), ".latest" + "/docroot", -1)
	} else {
		vhostPath = strings.Replace(Site.Path+Site.TimeStampGet(), Site.TimeStampGet(), ".latest", -1)
	}
	vhostFile := vhost.NewVirtualHost(Site.Name, vhostPath, Site.Webserver, Site.Domain, Site.Vhostpath)
	vhostFile.Uninstall()
}

// WebserverSet sets the webserver field for a site struct
func (Site *Site) WebserverSet(value string) {
	Site.Webserver = value
}
