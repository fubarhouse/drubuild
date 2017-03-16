package make

import (
	"github.com/fubarhouse/golang-drush/vhost"
	"strings"
)

func (Site *Site) VhostPathSet(value string) {
	Site.Vhostpath = value
}

func (Site *Site) VhostInstall() {
	vhostPath := strings.Replace(Site.Path+Site.TimeStampGet(), Site.TimeStampGet(), ".latest", -1)
	vhostFile := vhost.NewVirtualHost(Site.Name, vhostPath, Site.Webserver, Site.Domain, Site.Vhostpath)
	vhostFile.Install(Site.Template)
}

func (Site *Site) VhostUninstall() {
	vhostPath := strings.Replace(Site.Path+Site.TimeStampGet(), Site.TimeStampGet(), ".latest", -1)
	vhostFile := vhost.NewVirtualHost(Site.Name, vhostPath, Site.Webserver, Site.Domain, Site.Vhostpath)
	vhostFile.Uninstall()
}

func (Site *Site) WebserverSet(value string) {
	Site.Webserver = value
}
