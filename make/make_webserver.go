package make

import (
	"strings"

	"github.com/fubarhouse/drubuild/vhost"
)

// VhostInstall install a virtual host
func (Site *Site) VhostInstall() {
	var vhostPath string
	vhostPath = strings.Replace(Site.Path+Site.TimeStampGet(), Site.TimeStampGet(), "/" + Site.Domain + ".latest/"+Site.Docroot, -1)
	vhostFile := vhost.NewVirtualHost(Site.Name, vhostPath, Site.Webserver, Site.Domain, Site.Vhostpath)

	vhostFile.Install(Site.Template)
}

// VhostUninstall un-installs a virtual host
func (Site *Site) VhostUninstall() {
	var vhostPath string
	vhostPath = strings.Replace(Site.Path+Site.TimeStampGet(), Site.TimeStampGet(), Site.Name + "/" + Site.Domain+".latest/"+Site.Docroot, -1)
	vhostFile := vhost.NewVirtualHost(Site.Name, vhostPath, Site.Webserver, Site.Domain, Site.Vhostpath)
	vhostFile.Uninstall()
}
