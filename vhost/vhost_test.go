package vhost

import (
	"fmt"
	"os"
	"os/user"
	"testing"
	"time"
)

func TestInstansiation(t *testing.T) {
	obj := NewVirtualHost("testing", "~/tmp-vhost-package", "apache", "testing.local", "~/tmp-vhost-package")
	if obj.serverName != "testing" {
		t.Error("Expected value for field 'name' was not as expected.")
	}
}

func TestFieldNameGet(t *testing.T) {
	obj := NewVirtualHost("testing", "~/tmp-vhost-package", "apache", "testing.local", "~/tmp-vhost-package")
	if obj.serverName != "testing" {
		t.Error("Expected value for field 'serverName' was not as expected.")
	}
}
func TestFieldNameSet(t *testing.T) {
	obj := NewVirtualHost("testing", "~/tmp-vhost-package", "apache", "testing.local", "~/tmp-vhost-package")
	obj.SetName("changed value")
	if obj.serverName == "testing" {
		t.Error("Expected value for field 'serverName' was not as expected.")
	}
}
func TestFieldServerRootGet(t *testing.T) {
	obj := NewVirtualHost("testing", "~/tmp-vhost-package", "apache", "testing.local", "~/tmp-vhost-package")
	if obj.serverRoot != "~/tmp-vhost-package" {
		t.Error("Expected value for field 'serverRoot' was not as expected.")
	}
}

func TestFieldServerRootSet(t *testing.T) {
	obj := NewVirtualHost("testing", "~/tmp-vhost-package", "apache", "testing.local", "~/tmp-vhost-package")
	obj.SetURI("changed value")
	if obj.serverRoot == "~/tmp-vhost-package" {
		t.Error("Expected value for field 'serverRoot' was not as expected.")
	}
}
func TestFieldWebserverGet(t *testing.T) {
	obj := NewVirtualHost("testing", "~/tmp-vhost-package", "apache", "testing.local", "~/tmp-vhost-package")
	if obj.webserver != "apache" {
		t.Error("Expected value for field 'webserver' was not as expected.")
	}
}
func TestFieldWebserverSet(t *testing.T) {
	obj := NewVirtualHost("testing", "~/tmp-vhost-package", "apache", "testing.local", "~/tmp-vhost-package")
	obj.SetWebServer("changed value")
	if obj.webserver == "apache" {
		t.Error("Expected value for field 'webserver' was not as expected.")
	}
}
func TestFieldDomainGet(t *testing.T) {
	obj := NewVirtualHost("testing", "~/tmp-vhost-package", "apache", "testing.local", "~/tmp-vhost-package")
	if obj.serverDomain != "testing.local" {
		t.Error("Expected value for field 'serverDomain' was not as expected.")
	}
}
func TestFieldDomainSet(t *testing.T) {
	obj := NewVirtualHost("testing", "~/tmp-vhost-package", "apache", "testing.local", "~/tmp-vhost-package")
	obj.SetDomain("changed value")
	if obj.serverDomain == "apache" {
		t.Error("Expected value for field 'serverDomain' was not as expected.")
	}
}
func TestFieldInstallDirGet(t *testing.T) {
	obj := NewVirtualHost("testing", "~/tmp-vhost-package", "apache", "testing.local", "~/tmp-vhost-package")
	if obj.installationDirectory != "~/tmp-vhost-package" {
		t.Error("Expected value for field 'installationDirectory' was not as expected.")
	}
}
func TestFieldInstallDirSet(t *testing.T) {
	obj := NewVirtualHost("testing", "~/tmp-vhost-package", "apache", "testing.local", "~/tmp-vhost-package")
	obj.SetInstallationDirectory("changed value")
	if obj.installationDirectory == "~/tmp-vhost-package" {
		t.Error("Expected value for field 'installationDirectory' was not as expected.")
	}
}

func TestVhostInstall(t *testing.T) {
	thisUser, thisErr := user.Current()
	Home := thisUser.HomeDir
	if thisErr != nil {
		t.Error("Could not determine current user.")
	}
	obj := NewVirtualHost("testing", Home, "apache", "testing.local", Home)
	obj.Install()
	_, FilCreateErr := os.Stat(Home + "/testing.local.conf")
	if FilCreateErr != nil {
		t.Error("File could not be created.")
	}
	obj.Uninstall()
	_, FilDeleteErr := os.Stat(Home + "/testing.local.conf")
	if FilDeleteErr == nil {
		t.Error("File could not be deleted.")
	}
}

func TestVhostStatus(t *testing.T) {
	thisUser, thisErr := user.Current()
	Home := thisUser.HomeDir
	if thisErr != nil {
		t.Error("Could not determine current user.")
	}
	obj := NewVirtualHost("testing", Home, "apache", "testing.local", Home)
	obj.Install()
	_, FilCreateErr := os.Stat(Home + "/testing.local.conf")
	if FilCreateErr != nil {
		t.Error("File could not be created.")
	}
	time.Sleep(1000)
	if !obj.GetStatus() {
		obj.PrintStatus()
		fmt.Println(obj.GetURI())
		t.Error("Status returned false, expected true.")
	}
	obj.Uninstall()
	_, FilDeleteErr := os.Stat(Home + "/testing.local.conf")
	if FilDeleteErr == nil {
		t.Error("File could not be deleted.")
	}
	time.Sleep(1000)
	if obj.GetStatus() {
		t.Error("Status returned true, expected false.")
	}
}
