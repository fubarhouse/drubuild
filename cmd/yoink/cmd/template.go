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

var templateAlias = `
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

var templateVhostApache = `
DirectoryIndex index.php

<VirtualHost *:80>

    ServerName {{ .Domain }}
    DocumentRoot {{ .Root }}

    <Directory "{{ .Root }}">
      Options Indexes FollowSymLinks MultiViews
      AllowOverride All
      Options -Indexes +FollowSymLinks
      Require all granted
    </Directory>

</VirtualHost>
`
var templateVhostHttpd = `
DirectoryIndex index.php
<VirtualHost *:80>
    ServerName {{ .Domain }}
    DocumentRoot {{ .Root }}
    <Directory "{{ .Root }}">
      Options Indexes FollowSymLinks MultiViews
      AllowOverride All
      Options -Indexes +FollowSymLinks
      Require all granted
    </Directory>
</VirtualHost>
`

var templateVhostNginx = `
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
        # fastcgi_pass  127.0.0.1:9000;
        fastcgi_index index.php;
        fastcgi_pass unix:/var/run/php/php5.6-fpm.sock;
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

var templateSitesPhp = `
<?php

/**
 * @file
 * Configuration file for Drupal's multi-site directory aliasing feature.
 */

$sites['{{ .Alias }}'] = '{{ .Name }}';

?>
`