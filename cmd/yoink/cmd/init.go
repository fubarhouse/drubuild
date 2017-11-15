// Copyright © 2017 Karl Hepworth <Karl.Hepworth@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"
	"os"

	"fmt"

	"io"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var (

	// Home is the package variable to store the location of the config files.
	Home string

	// sites_php_template_date is the data for sites.php file.
	// it was taken from the templates folder and serves as a backup
	// when that file isn't available (not a go get install).
	sites_php_template_date = `
<?php

/**
 * @file
 * Configuration file for Drupal's multi-site directory aliasing feature.
 */

 $sites['default'] = '{{ .Name }}';
 $sites['{{ .Alias }}'] = '{{ .Name }}';

?>
`

	// drush_alias_template is the data for drush alias file.
	// it was taken from the templates folder and serves as a backup
	// when that file isn't available (not a go get install).
	drush_alias_template = `
<?php
	$aliases['{{ .Alias }}'] = array(
		'root' => '{{ .Root }}',
		'uri' => '{{ .Domain }}',
		'path-aliases' => array(
			'%files' => 'sites/{{ .Name }}/files',
			'%private' => 'sites/{{ .Name }}/private',
		),
	);
?>
`
	// vhost_template_data is the data for an nginx virtualhost config.
	// nginx is the default configured webserver, so other web server
	// templates will need to be added/changed as required.
	// it was taken from the templates folder and serves as a backup
	// when that file isn't available (not a go get install).
	vhost_template_data = `
server {
    listen 80;

    server_name {{ .Domain }}
    error_log /var/log/nginx/error.log info;
    root {{ .Root }};
    index index.php index.html index.htm;

    location / {
        # Don't touch PHP for static content.
        try_files $uri @rewrite;
    }

    # Don't allow direct access to PHP files in the vendor directory.
    location ~ /vendor/.*\.php$ {
        deny all;
        return 404;
    }

    # Use fastcgi for all php files.
    location ~ \.php$ {
        # Secure *.php files.
        try_files $uri = 404;
        include /etc/nginx/fastcgi_params;
        fastcgi_split_path_info ^(.+\.php)(/.+)$;
        fastcgi_pass  127.0.0.1:9000;
        fastcgi_index index.php;
        # fastcgi_pass unix:/var/run/php/php5.6-fpm.sock;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        fastcgi_read_timeout 120;
    }

    location @rewrite {
        # For D7 and above:
        rewrite ^ /index.php;

        # For Drupal 6 and below:
        #rewrite ^/(.*)$ /index.php?q=$1;
    }

    location ~ ^/sites/.*/files/styles/ {
        try_files $uri @rewrite;
    }

    location = /favicon.ico {
        log_not_found off;
        access_log off;
    }

    location = /robots.txt {
        allow all;
        log_not_found off;
        access_log off;
    }

    location ~ (^|/)\. {
        return 403;
    }

    location ~* \.(js|css|png|jpg|jpeg|gif|ico)$ {
        expires max;
        log_not_found off;
    }

    gzip on;
    gzip_proxied any;
    gzip_static on;
    gzip_http_version 1.0;
    gzip_disable "MSIE [1-6]\.";
    gzip_vary on;
    gzip_comp_level 6;
    gzip_types
        text/plain
        text/css
        text/xml
        text/javascript
        application/javascript
        application/x-javascript
        application/json
        application/xml
        application/xml+rss
        application/xhtml+xml
        application/x-font-ttf
        application/x-font-opentype
        image/svg+xml
        image/x-icon;
    gzip_buffers 16 8k;
    gzip_min_length 512;
}
`

	// config_yml_template_data is the data for sites.php file.
	// it was taken from the templates folder and serves as a backup
	// when that file isn't available (not a go get install).
	config_yml_template_data = `
---

db_user: drupal
db_pass: drupal
db_host: localhost
db_port: 3306

webserver: nginx
alias_template: $HOME/alias.tmpl
sites_php_template: $HOME/sites.php.tmpl
virtualhost_template: $HOME/vhost.tmpl

virtualhost_path: /etc/nginx/sites-enabled/
`
)

func WriteStringToFile(filepath, s string) error {
	fo, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer fo.Close()

	_, err = io.Copy(fo, strings.NewReader(s))
	if err != nil {
		return err
	}

	return nil
}

// syncCmd represents the backup command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise a set of templates in the provided destination path",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			Home = home
		}

		r := strings.Join([]string{Home, "yoink", ""}, string(os.PathSeparator))

		if _, s := os.Stat(r); s != nil {
			e := os.Mkdir(r, 0755)
			if e != nil {
				log.Fatalf("error in creating directory %v\n%v", r, e)
			}
		}

		config_yml_template_data = strings.Replace(config_yml_template_data, "$HOME", r, -1)
		config_yml_template_data = strings.Replace(config_yml_template_data, "//", string(os.PathSeparator), -1)

		log.Printf("Creating/replacing %vconfig.yml", r)
		WriteStringToFile(r+"config.yml", config_yml_template_data)
		log.Printf("Creating/replacing %vsites.php.tmpl", r)
		WriteStringToFile(r+"sites.php.tmpl", sites_php_template_date)
		log.Printf("Creating/replacing %valias.tmpl", r)
		WriteStringToFile(r+"alias.tmpl", drush_alias_template)
		log.Printf("Creating/replacing %vvhost.tmpl", r)
		WriteStringToFile(r+"vhost.tmpl", vhost_template_data)
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
