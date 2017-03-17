package vhost

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

type VirtualHost struct {
	serverName            string
	serverRoot            string
	serverDomain          string
	webserver             string
	installationDirectory string
}

func NewVirtualHost(name, root, webserver, domain, installationDirectory string) *VirtualHost {
	return &VirtualHost{name, root, domain, webserver, installationDirectory}
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

func (VirtualHost *VirtualHost) SetDomain(value string) {
	VirtualHost.serverDomain = value
}

func (VirtualHost *VirtualHost) GetDomain() string {
	return VirtualHost.serverDomain
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

func (VirtualHost *VirtualHost) Install(Templates ...string) {
	data := map[string]string{
		"Domain":     VirtualHost.GetDomain(),
		"ServerRoot": VirtualHost.GetUri(),
	}
	filename := VirtualHost.GetInstallationDirectory() + "/" + VirtualHost.GetDomain() + ".conf"
	buffer := []byte{}
	for _, val := range Templates {
		moreBytes, moreBytesErr := ioutil.ReadFile(val)
		if moreBytesErr != nil {
			log.Warnf("Could not read from %v: %v", val, moreBytesErr)
		}
		for _, value := range moreBytes {
			buffer = append(buffer, value)
		}
	}
	if len(buffer) > 0 {
		if VirtualHost.webserver == "nginx" {
			buffer = []byte{115, 101, 114, 118, 101, 114, 32, 123, 10, 32, 32, 32, 32, 108, 105, 115, 116, 101, 110, 32, 56, 48, 59, 10, 10, 32, 32, 32, 32, 115, 101, 114, 118, 101, 114, 95, 110, 97, 109, 101, 32, 68, 111, 109, 97, 105, 110, 10, 32, 32, 32, 32, 101, 114, 114, 111, 114, 95, 108, 111, 103, 32, 47, 118, 97, 114, 47, 108, 111, 103, 47, 110, 103, 105, 110, 120, 47, 101, 114, 114, 111, 114, 46, 108, 111, 103, 32, 105, 110, 102, 111, 59, 10, 32, 32, 32, 32, 114, 111, 111, 116, 32, 83, 101, 114, 118, 101, 114, 82, 111, 111, 116, 59, 10, 32, 32, 32, 32, 105, 110, 100, 101, 120, 32, 105, 110, 100, 101, 120, 46, 112, 104, 112, 32, 105, 110, 100, 101, 120, 46, 104, 116, 109, 108, 32, 105, 110, 100, 101, 120, 46, 104, 116, 109, 59, 10, 10, 32, 32, 32, 32, 108, 111, 99, 97, 116, 105, 111, 110, 32, 47, 32, 123, 10, 32, 32, 32, 32, 32, 32, 32, 32, 35, 32, 68, 111, 110, 39, 116, 32, 116, 111, 117, 99, 104, 32, 80, 72, 80, 32, 102, 111, 114, 32, 115, 116, 97, 116, 105, 99, 32, 99, 111, 110, 116, 101, 110, 116, 46, 10, 32, 32, 32, 32, 32, 32, 32, 32, 116, 114, 121, 95, 102, 105, 108, 101, 115, 32, 36, 117, 114, 105, 32, 64, 114, 101, 119, 114, 105, 116, 101, 59, 10, 32, 32, 32, 32, 125, 10, 10, 32, 32, 32, 32, 35, 32, 68, 111, 110, 39, 116, 32, 97, 108, 108, 111, 119, 32, 100, 105, 114, 101, 99, 116, 32, 97, 99, 99, 101, 115, 115, 32, 116, 111, 32, 80, 72, 80, 32, 102, 105, 108, 101, 115, 32, 105, 110, 32, 116, 104, 101, 32, 118, 101, 110, 100, 111, 114, 32, 100, 105, 114, 101, 99, 116, 111, 114, 121, 46, 10, 32, 32, 32, 32, 108, 111, 99, 97, 116, 105, 111, 110, 32, 126, 32, 47, 118, 101, 110, 100, 111, 114, 47, 46, 42, 92, 46, 112, 104, 112, 36, 32, 123, 10, 32, 32, 32, 32, 32, 32, 32, 32, 100, 101, 110, 121, 32, 97, 108, 108, 59, 10, 32, 32, 32, 32, 32, 32, 32, 32, 114, 101, 116, 117, 114, 110, 32, 52, 48, 52, 59, 10, 32, 32, 32, 32, 125, 10, 10, 32, 32, 32, 32, 35, 32, 85, 115, 101, 32, 102, 97, 115, 116, 99, 103, 105, 32, 102, 111, 114, 32, 97, 108, 108, 32, 112, 104, 112, 32, 102, 105, 108, 101, 115, 46, 10, 32, 32, 32, 32, 108, 111, 99, 97, 116, 105, 111, 110, 32, 126, 32, 92, 46, 112, 104, 112, 36, 32, 123, 10, 32, 32, 32, 32, 32, 32, 32, 32, 35, 32, 83, 101, 99, 117, 114, 101, 32, 42, 46, 112, 104, 112, 32, 102, 105, 108, 101, 115, 46, 10, 32, 32, 32, 32, 32, 32, 32, 32, 116, 114, 121, 95, 102, 105, 108, 101, 115, 32, 36, 117, 114, 105, 32, 61, 32, 52, 48, 52, 59, 10, 32, 32, 32, 32, 32, 32, 32, 32, 105, 110, 99, 108, 117, 100, 101, 32, 47, 101, 116, 99, 47, 110, 103, 105, 110, 120, 47, 102, 97, 115, 116, 99, 103, 105, 95, 112, 97, 114, 97, 109, 115, 59, 10, 32, 32, 32, 32, 32, 32, 32, 32, 102, 97, 115, 116, 99, 103, 105, 95, 115, 112, 108, 105, 116, 95, 112, 97, 116, 104, 95, 105, 110, 102, 111, 32, 94, 40, 46, 43, 92, 46, 112, 104, 112, 41, 40, 47, 46, 43, 41, 36, 59, 10, 32, 32, 32, 32, 32, 32, 32, 32, 35, 32, 102, 97, 115, 116, 99, 103, 105, 95, 112, 97, 115, 115, 32, 32, 49, 50, 55, 46, 48, 46, 48, 46, 49, 58, 57, 48, 48, 48, 59, 10, 32, 32, 32, 32, 32, 32, 32, 32, 102, 97, 115, 116, 99, 103, 105, 95, 105, 110, 100, 101, 120, 32, 105, 110, 100, 101, 120, 46, 112, 104, 112, 59, 10, 32, 32, 32, 32, 32, 32, 32, 32, 102, 97, 115, 116, 99, 103, 105, 95, 112, 97, 115, 115, 32, 117, 110, 105, 120, 58, 47, 118, 97, 114, 47, 114, 117, 110, 47, 112, 104, 112, 47, 112, 104, 112, 53, 46, 54, 45, 102, 112, 109, 46, 115, 111, 99, 107, 59, 10, 32, 32, 32, 32, 32, 32, 32, 32, 102, 97, 115, 116, 99, 103, 105, 95, 112, 97, 114, 97, 109, 32, 83, 67, 82, 73, 80, 84, 95, 70, 73, 76, 69, 78, 65, 77, 69, 32, 36, 100, 111, 99, 117, 109, 101, 110, 116, 95, 114, 111, 111, 116, 36, 102, 97, 115, 116, 99, 103, 105, 95, 115, 99, 114, 105, 112, 116, 95, 110, 97, 109, 101, 59, 10, 32, 32, 32, 32, 32, 32, 32, 32, 102, 97, 115, 116, 99, 103, 105, 95, 114, 101, 97, 100, 95, 116, 105, 109, 101, 111, 117, 116, 32, 49, 50, 48, 59, 10, 32, 32, 32, 32, 125, 10, 10, 32, 32, 32, 32, 108, 111, 99, 97, 116, 105, 111, 110, 32, 64, 114, 101, 119, 114, 105, 116, 101, 32, 123, 10, 32, 32, 32, 32, 32, 32, 32, 32, 35, 32, 70, 111, 114, 32, 68, 55, 32, 97, 110, 100, 32, 97, 98, 111, 118, 101, 58, 10, 32, 32, 32, 32, 32, 32, 32, 32, 114, 101, 119, 114, 105, 116, 101, 32, 94, 32, 47, 105, 110, 100, 101, 120, 46, 112, 104, 112, 59, 10, 10, 32, 32, 32, 32, 32, 32, 32, 32, 35, 32, 70, 111, 114, 32, 68, 114, 117, 112, 97, 108, 32, 54, 32, 97, 110, 100, 32, 98, 101, 108, 111, 119, 58, 10, 32, 32, 32, 32, 32, 32, 32, 32, 35, 114, 101, 119, 114, 105, 116, 101, 32, 94, 47, 40, 46, 42, 41, 36, 32, 47, 105, 110, 100, 101, 120, 46, 112, 104, 112, 63, 113, 61, 36, 49, 59, 10, 32, 32, 32, 32, 125, 10, 10, 32, 32, 32, 32, 108, 111, 99, 97, 116, 105, 111, 110, 32, 126, 32, 94, 47, 115, 105, 116, 101, 115, 47, 46, 42, 47, 102, 105, 108, 101, 115, 47, 115, 116, 121, 108, 101, 115, 47, 32, 123, 10, 32, 32, 32, 32, 32, 32, 32, 32, 116, 114, 121, 95, 102, 105, 108, 101, 115, 32, 36, 117, 114, 105, 32, 64, 114, 101, 119, 114, 105, 116, 101, 59, 10, 32, 32, 32, 32, 125, 10, 10, 32, 32, 32, 32, 108, 111, 99, 97, 116, 105, 111, 110, 32, 61, 32, 47, 102, 97, 118, 105, 99, 111, 110, 46, 105, 99, 111, 32, 123, 10, 32, 32, 32, 32, 32, 32, 32, 32, 108, 111, 103, 95, 110, 111, 116, 95, 102, 111, 117, 110, 100, 32, 111, 102, 102, 59, 10, 32, 32, 32, 32, 32, 32, 32, 32, 97, 99, 99, 101, 115, 115, 95, 108, 111, 103, 32, 111, 102, 102, 59, 10, 32, 32, 32, 32, 125, 10, 10, 32, 32, 32, 32, 108, 111, 99, 97, 116, 105, 111, 110, 32, 61, 32, 47, 114, 111, 98, 111, 116, 115, 46, 116, 120, 116, 32, 123, 10, 32, 32, 32, 32, 32, 32, 32, 32, 97, 108, 108, 111, 119, 32, 97, 108, 108, 59, 10, 32, 32, 32, 32, 32, 32, 32, 32, 108, 111, 103, 95, 110, 111, 116, 95, 102, 111, 117, 110, 100, 32, 111, 102, 102, 59, 10, 32, 32, 32, 32, 32, 32, 32, 32, 97, 99, 99, 101, 115, 115, 95, 108, 111, 103, 32, 111, 102, 102, 59, 10, 32, 32, 32, 32, 125, 10, 10, 32, 32, 32, 32, 108, 111, 99, 97, 116, 105, 111, 110, 32, 126, 32, 40, 94, 124, 47, 41, 92, 46, 32, 123, 10, 32, 32, 32, 32, 32, 32, 32, 32, 114, 101, 116, 117, 114, 110, 32, 52, 48, 51, 59, 10, 32, 32, 32, 32, 125, 10, 10, 32, 32, 32, 32, 108, 111, 99, 97, 116, 105, 111, 110, 32, 126, 42, 32, 92, 46, 40, 106, 115, 124, 99, 115, 115, 124, 112, 110, 103, 124, 106, 112, 103, 124, 106, 112, 101, 103, 124, 103, 105, 102, 124, 105, 99, 111, 41, 36, 32, 123, 10, 32, 32, 32, 32, 32, 32, 32, 32, 101, 120, 112, 105, 114, 101, 115, 32, 109, 97, 120, 59, 10, 32, 32, 32, 32, 32, 32, 32, 32, 108, 111, 103, 95, 110, 111, 116, 95, 102, 111, 117, 110, 100, 32, 111, 102, 102, 59, 10, 32, 32, 32, 32, 125, 10, 10, 32, 32, 32, 32, 103, 122, 105, 112, 32, 111, 110, 59, 10, 32, 32, 32, 32, 103, 122, 105, 112, 95, 112, 114, 111, 120, 105, 101, 100, 32, 97, 110, 121, 59, 10, 32, 32, 32, 32, 103, 122, 105, 112, 95, 115, 116, 97, 116, 105, 99, 32, 111, 110, 59, 10, 32, 32, 32, 32, 103, 122, 105, 112, 95, 104, 116, 116, 112, 95, 118, 101, 114, 115, 105, 111, 110, 32, 49, 46, 48, 59, 10, 32, 32, 32, 32, 103, 122, 105, 112, 95, 100, 105, 115, 97, 98, 108, 101, 32, 34, 77, 83, 73, 69, 32, 91, 49, 45, 54, 93, 92, 46, 34, 59, 10, 32, 32, 32, 32, 103, 122, 105, 112, 95, 118, 97, 114, 121, 32, 111, 110, 59, 10, 32, 32, 32, 32, 103, 122, 105, 112, 95, 99, 111, 109, 112, 95, 108, 101, 118, 101, 108, 32, 54, 59, 10, 32, 32, 32, 32, 103, 122, 105, 112, 95, 116, 121, 112, 101, 115, 10, 32, 32, 32, 32, 32, 32, 32, 32, 116, 101, 120, 116, 47, 112, 108, 97, 105, 110, 10, 32, 32, 32, 32, 32, 32, 32, 32, 116, 101, 120, 116, 47, 99, 115, 115, 10, 32, 32, 32, 32, 32, 32, 32, 32, 116, 101, 120, 116, 47, 120, 109, 108, 10, 32, 32, 32, 32, 32, 32, 32, 32, 116, 101, 120, 116, 47, 106, 97, 118, 97, 115, 99, 114, 105, 112, 116, 10, 32, 32, 32, 32, 32, 32, 32, 32, 97, 112, 112, 108, 105, 99, 97, 116, 105, 111, 110, 47, 106, 97, 118, 97, 115, 99, 114, 105, 112, 116, 10, 32, 32, 32, 32, 32, 32, 32, 32, 97, 112, 112, 108, 105, 99, 97, 116, 105, 111, 110, 47, 120, 45, 106, 97, 118, 97, 115, 99, 114, 105, 112, 116, 10, 32, 32, 32, 32, 32, 32, 32, 32, 97, 112, 112, 108, 105, 99, 97, 116, 105, 111, 110, 47, 106, 115, 111, 110, 10, 32, 32, 32, 32, 32, 32, 32, 32, 97, 112, 112, 108, 105, 99, 97, 116, 105, 111, 110, 47, 120, 109, 108, 10, 32, 32, 32, 32, 32, 32, 32, 32, 97, 112, 112, 108, 105, 99, 97, 116, 105, 111, 110, 47, 120, 109, 108, 43, 114, 115, 115, 10, 32, 32, 32, 32, 32, 32, 32, 32, 97, 112, 112, 108, 105, 99, 97, 116, 105, 111, 110, 47, 120, 104, 116, 109, 108, 43, 120, 109, 108, 10, 32, 32, 32, 32, 32, 32, 32, 32, 97, 112, 112, 108, 105, 99, 97, 116, 105, 111, 110, 47, 120, 45, 102, 111, 110, 116, 45, 116, 116, 102, 10, 32, 32, 32, 32, 32, 32, 32, 32, 97, 112, 112, 108, 105, 99, 97, 116, 105, 111, 110, 47, 120, 45, 102, 111, 110, 116, 45, 111, 112, 101, 110, 116, 121, 112, 101, 10, 32, 32, 32, 32, 32, 32, 32, 32, 105, 109, 97, 103, 101, 47, 115, 118, 103, 43, 120, 109, 108, 10, 32, 32, 32, 32, 32, 32, 32, 32, 105, 109, 97, 103, 101, 47, 120, 45, 105, 99, 111, 110, 59, 10, 32, 32, 32, 32, 103, 122, 105, 112, 95, 98, 117, 102, 102, 101, 114, 115, 32, 49, 54, 32, 56, 107, 59, 10, 32, 32, 32, 32, 103, 122, 105, 112, 95, 109, 105, 110, 95, 108, 101, 110, 103, 116, 104, 32, 53, 49, 50, 59, 10, 125, 10}
		} else if VirtualHost.webserver == "apache" || VirtualHost.webserver == "apache2" || VirtualHost.webserver == "httpd" {
			buffer = []byte{68, 105, 114, 101, 99, 116, 111, 114, 121, 73, 110, 100, 101, 120, 32, 105, 110, 100, 101, 120, 46, 112, 104, 112, 10, 60, 86, 105, 114, 116, 117, 97, 108, 72, 111, 115, 116, 32, 42, 58, 56, 48, 62, 10, 32, 32, 32, 32, 83, 101, 114, 118, 101, 114, 78, 97, 109, 101, 32, 68, 111, 109, 97, 105, 110, 10, 32, 32, 32, 32, 68, 111, 99, 117, 109, 101, 110, 116, 82, 111, 111, 116, 32, 83, 101, 114, 118, 101, 114, 82, 111, 111, 116, 10, 32, 32, 32, 32, 60, 68, 105, 114, 101, 99, 116, 111, 114, 121, 32, 34, 83, 101, 114, 118, 101, 114, 82, 111, 111, 116, 34, 62, 10, 32, 32, 32, 32, 32, 32, 79, 112, 116, 105, 111, 110, 115, 32, 73, 110, 100, 101, 120, 101, 115, 32, 70, 111, 108, 108, 111, 119, 83, 121, 109, 76, 105, 110, 107, 115, 32, 77, 117, 108, 116, 105, 86, 105, 101, 119, 115, 10, 32, 32, 32, 32, 32, 32, 65, 108, 108, 111, 119, 79, 118, 101, 114, 114, 105, 100, 101, 32, 65, 108, 108, 10, 32, 32, 32, 32, 32, 32, 79, 112, 116, 105, 111, 110, 115, 32, 45, 73, 110, 100, 101, 120, 101, 115, 32, 43, 70, 111, 108, 108, 111, 119, 83, 121, 109, 76, 105, 110, 107, 115, 10, 32, 32, 32, 32, 32, 32, 82, 101, 113, 117, 105, 114, 101, 32, 97, 108, 108, 32, 103, 114, 97, 110, 116, 101, 100, 10, 32, 32, 32, 32, 60, 47, 68, 105, 114, 101, 99, 116, 111, 114, 121, 62, 10, 60, 47, 86, 105, 114, 116, 117, 97, 108, 72, 111, 115, 116, 62}
		}
	}
	tpl := fmt.Sprintf("%v", string(buffer[:]))
	tpl = strings.Replace(tpl, "Domain", data["Domain"], -1)
	tpl = strings.Replace(tpl, "ServerRoot", data["ServerRoot"], -1)
	tpl = strings.Replace(tpl, ".latest", "/"+VirtualHost.GetName()+".latest", -1)

	_, statErr := os.Stat(filename)
	if statErr != nil {
		nf, err := os.Create(filename)
		if err != nil {
			log.Fatalln("Error creating file", err)
		}
		_, err = nf.WriteString(tpl)
		if err != nil {
			log.Warnln("Could not add vhost", VirtualHost.GetDomain())
		} else {
			log.Infoln("Added vhost", VirtualHost.GetDomain())
		}
		defer nf.Close()
	} else {
		log.Warnln("Vhost already created")
	}
}

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
