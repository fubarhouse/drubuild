package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/golang-drush/alias"
	"github.com/fubarhouse/golang-drush/vhost"
)

func main() {

	var WebserverDir = flag.String("vhost-dir", "/etc/nginx/sites-enabled", "Directory containing virtual host file(s)")
	var Webserver = flag.String("webserver-name", "nginx", "The name of the web service on the server.")

	flag.Parse()

	log.Println("Instanciating Alias")
	Alias := alias.NewAlias("temporaryAlias", "/tmp", "temporaryAlias")
	log.Println("Installing Alias")
	Alias.Install()
	log.Println("Uninstalling Alias")
	Alias.Uninstall()

	log.Println("Instanciating Vhost")
	VirtualHost := vhost.NewVirtualHost("temporaryVhost", "/tmp", *Webserver, "temporary.vhost", *WebserverDir)
	log.Println("Installing Vhost")
	VirtualHost.Install()
	log.Println("Uninstalling Vhost")
	VirtualHost.Uninstall()
}
