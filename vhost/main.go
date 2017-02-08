package vhost

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

type VirtualHost struct {
	serverName            string
	serverRoot            string
	webserver             string
	installationDirectory string
}

func NewVirtualHost(name, root, webserver, installationDirectory string) *VirtualHost {
	return &VirtualHost{name, root, webserver, installationDirectory}
}

func (VirtualHost *VirtualHost) SetName(value string) {
	VirtualHost.serverName = value
}

func (VirtualHost *VirtualHost) GetName() string {
	return VirtualHost.serverName
}

func (VirtualHost *VirtualHost) SetUri(value string) {
	VirtualHost.serverRoot = value
}

func (VirtualHost *VirtualHost) GetUri() string {
	return VirtualHost.serverRoot
}

func (VirtualHost *VirtualHost) SetWebServer(value string) {
	VirtualHost.webserver = value
}

func (VirtualHost *VirtualHost) GetWebServer() string {
	return VirtualHost.webserver
}

func (VirtualHost *VirtualHost) SetInstallationDirectory(value string) {
	VirtualHost.installationDirectory = value
}

func (VirtualHost *VirtualHost) GetInstallationDirectory() string {
	return VirtualHost.installationDirectory
}

func (VirtualHost *VirtualHost) Install() {
	log.Println("Adding vhost", VirtualHost.GetName())
	tpl, err := template.ParseGlob("templates/" + VirtualHost.webserver + "*")
	file, err := os.Create(VirtualHost.GetInstallationDirectory() + "/" + VirtualHost.GetName() + ".conf")
	fmt.Println(VirtualHost.GetInstallationDirectory() + "/" + VirtualHost.GetName() + ".conf")
	tpl.Execute(file, nil)
	//		struct {
	//	Name                  string
	//	Root                  string
	//	Webserver             string
	//	InstallationDirectory string
	//}{VirtualHost.GetName(), VirtualHost.GetUri(), VirtualHost.GetWebServer(), VirtualHost.GetInstallationDirectory()})

	if err != nil {
		log.Println("Error reading files:", err)
	}
	defer file.Close()
}

func (VirtualHost *VirtualHost) Uninstall() {
	log.Println("Removing vhost", VirtualHost.GetName())
	os.Remove(VirtualHost.installationDirectory + "/" + VirtualHost.GetUri() + ".conf")

}

func (VirtualHost *VirtualHost) Reinstall() {
	VirtualHost.Uninstall()
	VirtualHost.Install()

}

func (VirtualHost *VirtualHost) GetStatus() bool {
	_, err := os.Stat(VirtualHost.installationDirectory + "/" + VirtualHost.GetUri() + ".conf")
	if err != nil {
		return false
	} else {
		return true
	}
}

func (VirtualHost *VirtualHost) PrintStatus() {
	_, err := os.Stat(VirtualHost.installationDirectory + "/" + VirtualHost.GetUri() + ".conf")
	if err != nil {
		fmt.Println("false")
	} else {
		fmt.Println("true")
	}
}
