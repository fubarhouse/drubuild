package make

import (
	"github.com/fubarhouse/golang-drush/vhost"
	"strings"
)

// Set a virtual host path
func (Site *Site) VhostPathSet(value string) {
	Site.Vhostpath = value
}

// Install a virtual host
func (Site *Site) VhostInstall() {
	vhostPath := strings.Replace(Site.Path+Site.TimeStampGet(), Site.TimeStampGet(), ".latest", -1)
	vhostFile := vhost.NewVirtualHost(Site.Name, vhostPath, Site.Webserver, Site.Domain, Site.Vhostpath)
	vhostFile.Install(Site.Template)
}

// Un-install a virtual host
func (Site *Site) VhostUninstall() {
	vhostPath := strings.Replace(Site.Path+Site.TimeStampGet(), Site.TimeStampGet(), ".latest", -1)
	vhostFile := vhost.NewVirtualHost(Site.Name, vhostPath, Site.Webserver, Site.Domain, Site.Vhostpath)
	vhostFile.Uninstall()
}

// Set the webserver field for a site struct
func (Site *Site) WebserverSet(value string) {
	Site.Webserver = value
}
