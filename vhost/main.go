package vhost

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
)

// VirtualHost is a stuct for virtual host types.
type VirtualHost struct {
	serverName            string
	serverRoot            string
	serverDomain          string
	webserver             string
	installationDirectory string
}

// NewVirtualHost instantiates a new VirtualHost struct.
func NewVirtualHost(name, root, webserver, domain, installationDirectory string) *VirtualHost {
	return &VirtualHost{name, root, domain, webserver, installationDirectory}
}

// SetName Sets the name field of the VirtualHost struct
func (VirtualHost *VirtualHost) SetName(value string) {
	VirtualHost.serverName = value
}

// GetName Gets the name field of the VirtualHost struct
func (VirtualHost *VirtualHost) GetName() string {
	return VirtualHost.serverName
}

// SetURI Sets the uri field of the VirtualHost struct
func (VirtualHost *VirtualHost) SetURI(value string) {
	VirtualHost.serverRoot = value
}

// GetURI Gets the uri field of the VirtualHost struct
func (VirtualHost *VirtualHost) GetURI() string {
	return VirtualHost.serverRoot
}

// SetDomain Sets the domain field of the VirtualHost struct
func (VirtualHost *VirtualHost) SetDomain(value string) {
	VirtualHost.serverDomain = value
}

// GetDomain Gets the domain field of the VirtualHost struct
func (VirtualHost *VirtualHost) GetDomain() string {
	return VirtualHost.serverDomain
}

// SetWebServer Sets the webserver field of the VirtualHost struct
func (VirtualHost *VirtualHost) SetWebServer(value string) {
	VirtualHost.webserver = value
}

// GetWebServer Gets the webserver field of the VirtualHost struct
func (VirtualHost *VirtualHost) GetWebServer() string {
	return VirtualHost.webserver
}

// SetInstallationDirectory Sets the installationDirectory field of the VirtualHost struct
func (VirtualHost *VirtualHost) SetInstallationDirectory(value string) {
	VirtualHost.installationDirectory = value
}

// GetInstallationDirectory Gets the installationDirectory field of the VirtualHost struct
func (VirtualHost *VirtualHost) GetInstallationDirectory() string {
	return VirtualHost.installationDirectory
}

// Install installs a virtual host from an optional input template file path.
func (VirtualHost *VirtualHost) Install(Template string) {
	data := map[string]string{
		"Domain": VirtualHost.GetDomain(),
		"Root":   VirtualHost.GetURI(),
	}

	filename := VirtualHost.GetInstallationDirectory() + "/" + VirtualHost.GetDomain() + ".conf"

	t := template.New("vhost")
	defaultData, _ := ioutil.ReadFile(Template)
	t.Parse(string(defaultData))

	os.Remove(filename)
	file, _ := os.Create(filename)
	tplErr := t.Execute(file, data)

	if tplErr == nil {
		log.Infof("Successfully templated vhost to file %v", filename)
	} else {
		log.Warnf("Error templating vhost to file %v: %v", filename, tplErr)
	}
}

// Uninstall un-installs a virtual host.
func (VirtualHost *VirtualHost) Uninstall() {
	filename := VirtualHost.GetInstallationDirectory() + "/" + VirtualHost.GetDomain() + ".conf"
	_, statErr := os.Stat(filename)
	if statErr == nil {
		err := os.Remove(filename)
		if err != nil {
			log.Warnln("Could not remove vhost", VirtualHost.GetDomain())
		} else {
			log.Infoln("Removed vhost", VirtualHost.GetDomain())
		}
	} else {
		log.Warnln("Vhost file was not found.")
	}

}

// Reinstall re-installs a virtual host.
func (VirtualHost *VirtualHost) Reinstall(Template string) {
	VirtualHost.Uninstall()
	VirtualHost.Install(Template)

}

// GetStatus returns the installation status of a virtual host
func (VirtualHost *VirtualHost) GetStatus() bool {
	_, err := os.Stat(VirtualHost.installationDirectory + "/" + VirtualHost.GetURI() + ".conf")
	if err != nil {
		return false
	}
	return true
}

// PrintStatus Prints the installation status of a virtual host.
func (VirtualHost *VirtualHost) PrintStatus() {
	_, err := os.Stat(VirtualHost.installationDirectory + "/" + VirtualHost.GetURI() + ".conf")
	if err != nil {
		fmt.Println("false")
	} else {
		fmt.Println("true")
	}
}
