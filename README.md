<img style="float:left" alight="left" height="128px" width="100px" src="https://github.com/fubarhouse/ansible-role-golang/raw/master/gopher.png">

# Drubuild 0.3.x

[![stability-experimental](https://img.shields.io/badge/stability-experimental-orange.svg?style=for-the-badge)](https://github.com/orangemug/stability-badges)
[![Go Report Card](https://goreportcard.com/badge/github.com/fubarhouse/drubuild?style=for-the-badge)](https://goreportcard.com/report/github.com/fubarhouse/drubuild)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg?style=for-the-badge)](https://raw.githubusercontent.com/fubarhouse/brand/master/LICENSE.txt)

```
version 0.3.x is experimental, please avoid use of this version.
```

## Purpose

Drubuild is a command-line application which builds and manages sites via some common composer and drush commands.

This application was born of the personal desire of the creator to be a useful tool for CI and automation to be used at work, however the stability and reliability was a catalyst for a lot of work here. 

This application is fully-intended to work out of the box with [DrupalVM](https://www.drupalvm.com/), and a [fork of DrupalVM](https://github.com/fubarhouse/drupal-vm) with the configuration of the maintainer's Ansible role [fubarhouse.golang](https://github.com/fubarhouse/ansible-role-golang) is available for quick opportunities to use this in an isolated environment.

The application has previously been used for CI tooling for the management of more than 40 websites simultaneously, supporting developers during their local build and development cycle with Jenkins. 

## Install

It is *highly* recommended to install this using Go, and currently no other options are available. 

```console
$ go get -u github.com/fubarhouse/drubuild
```

After installing via `go get`, or with [the configured fork of DrupalVM](https://github.com/fubarhouse/drupal-vm), you should focus your attention on getting Drush Aliases and SSH working.

## Configuration

There are a couple of steps to configuring Drubuild.

  1. Run `drubuild init --global` to establish a new set of templates and config files with default settings.
  2. Modify `config.yml` and all templates as required.
  3. If you are executing Drubuild from teh destination folder without the destination flag, you *can* run `drubuild init` to copy the global config and templates or to setup brand new config or templates (from default values)
  4. When you first start the build process, please ensure the docroot flag is specified *correctly*.

## Usage

Before you can begin, you will need to run `drubuild init` to establish a set of templates and configuration settings at `$HOME/drubuild/`.

After having made the templates and config available, you can run `drubuild` like any other console application.

```
Using config file: /Users/karl/drubuild/config.yml
A Drupal build system.

Usage:
  drubuild [command]

Available Commands:
  backup      Take a archive-dump snapshot of a local site
  build       The build process for Drubuild
  destroy     Remove all traces of an installed site.
  help        Help about any command
  init        Initialise a set of templates in the provided destination path
  project     Install or remove a project.
  runner      Runs a series of drush commands consecutively.
  sync        Execute drush sql-sync or rsync between two drush aliases
  user        User management tasks

Flags:
  -h, --help   help for drubuild

Use "drubuild [command] --help" for more information about a command.
```

### Backup

A very simple terminal interface which allows for `drush archive-dump`.
```
drubuild backup --source @mysite.dev --destination /mysite.dev.tar.gz
```

### Build

The actual build process, which uses a composer.json file, and some other basic information and build you a site, including installation, multi-site setup, drush alias and virtualhost - all provided by the templates used by the system.

There are many options here to control the process for site builds, folder hierarchy (specify a docroot in the folder) etc.
```
drubuild build --name mysite --domain mysite.dev --destination /sites/mysite.dev --composer /composer.json --vhost --install
```

### Destroy

The reverse of building, to remove all traces of a previously built site, including file system, databases, aliases, virtual hosts etc.
```
drubuild destroy --name mysite --domain mysite.dev --destination /sites/mysite.dev
```

### Help

Basically provide the help text listed above.
```
drubuild --help
```

### Init

Initialise or replace configurables and templates with defaults.
```
drubuild init
```

### Project

Add or remove a project to a built site using information available in the composer.json file.
```
drubuild project --name drupal/core --path /sites/mysite.dev --remove --add --version ^8
```

### Runner

The drush command runner is a scalable task runner which will execute a comma-separated list of drush commands on a comma-separated list of drush aliases with pattern matching supported for scaling.
```
drubuild runner --aliases dev,test,preprod,prod --pattern mysite.%v --commands "updb -y,cache-rebuild"
```

### Sync

The syncer will attempt to sync the databases and files between two Drush aliases. Note that this can be very dangerous if used incorrectly, and Drush doesn't support remote to remote syncs. The exact settings for this can be found in the Drush alias template.
```
drubuild sync --source mysite.prod --destination mysite.local  --database --files --yes
```

### User

The user command will act much like the runner, which can create, block, unblock, delete, reset password, add roles, and general maintenance accross a set of sites using alias pattern matching.

This one has proven critical when managing literally hundreds of Drupal accounts.
```
drubuild user --name TestUser --email test@user.com --password MyPassword --aliases dev,test,preprod,prod --pattern mysite.%v --create
```

## Example

Using our [fork](https://github.com/fubarhouse/drupal-vm) of DrupalVM, here's a full example of how a site could be stood up from the default setup.

<span style="color:red">Do not forget to apply the --docroot flag!</span>

The docroot using the example composer.json file supplied with DrupalVM is "web".

It should also be noted that there is a bug with using dashes in the name flag, please avoid it. 

````
vagrant@drupalvm2:~$ mkdir -p /vagrant/sites/mysiteone
vagrant@drupalvm2:~$ cp /vagrant/example.drupal.composer.json /vagrant/sites/mysiteone/composer.json
vagrant@drupalvm2:~$ cd /vagrant/sites/mysiteone/
vagrant@drupalvm2:/vagrant/sites/mysiteone$ drubuild init
Using config file: /home/vagrant/drubuild/config.yml
2017/11/24 08:11:20 Templating /vagrant/sites/mysiteone/drubuild/config.yml from defaults.
2017/11/24 08:11:20 Replacing /vagrant/sites/mysiteone/drubuild/sites.php.tmpl with /home/vagrant/drubuild/sites.php.tmpl.
2017/11/24 08:11:20 Replacing /vagrant/sites/mysiteone/drubuild/alias.tmpl with /home/vagrant/drubuild/alias.tmpl.
2017/11/24 08:11:20 Replacing /vagrant/sites/mysiteone/drubuild/vhost.tmpl with /home/vagrant/drubuild/vhost.tmpl.
vagrant@drupalvm2:/vagrant/sites/mysiteone$ drubuild build --name mysiteone --domain mysiteone.test --docroot web --vhost --install
Using config file: /vagrant/sites/mysiteone/drubuild/config.yml
INFO[0000] Timestamp not specified, using 20171124081132
INFO[0000] composer.json not found, copying from /vagrant/sites/mysiteone/composer.json
INFO[0000] Copied /vagrant/sites/mysiteone/composer.json to /vagrant/sites/mysiteone/mysiteone.20171124081132_7047/composer.json

    1/1:	https://packages.drupal.org/8/drupal/provider-2017-3$6850a9263c4bed0fd003e1e84e14c71f4aa9b8304a453f7754d6c1edbd25c638.json
    Finished: success: 1, skipped: 0, failure: 0, total: 1
    1/3:	http://packagist.org/p/provider-latest$40d2eeb0a214664f06b7959191e50bfd93e2b6b3adb023ebc5f18e8f95fb7134.json
    2/3:	http://packagist.org/p/provider-2017-10$3706041cdaaf6c9ce37bc29a6f8460c5b6e26806ed76eed4f4e182baece607a1.json
    3/3:	http://packagist.org/p/provider-2017-07$ffb45f4bf7108849406afd662d60ef80cffec40a957318122b362f8b1b99c9d1.json
    Finished: success: 3, skipped: 0, failure: 0, total: 3
Loading composer repositories with package information
Updating dependencies (including require-dev)
Package operations: 44 installs, 0 updates, 0 removals
  - Installing composer/installers (v1.4.0): Loading from cache
  - Installing drupal-composer/drupal-scaffold (2.3.0): Loading from cache
  - Installing zendframework/zend-stdlib (3.1.0): Loading from cache
  - Installing zendframework/zend-escaper (2.5.2): Loading from cache
  - Installing zendframework/zend-feed (2.8.0): Loading from cache
  - Installing psr/http-message (1.0.1): Loading from cache
  - Installing zendframework/zend-diactoros (1.6.1): Loading from cache
  - Installing twig/twig (v1.35.0): Loading from cache
  - Installing symfony/yaml (v3.2.14): Loading from cache
  - Installing symfony/polyfill-mbstring (v1.6.0): Loading from cache
  - Installing symfony/translation (v3.2.14): Loading from cache
  - Installing symfony/validator (v3.2.14): Loading from cache
  - Installing symfony/serializer (v3.2.14): Loading from cache
  - Installing symfony/routing (v3.2.14): Loading from cache
  - Installing paragonie/random_compat (v2.0.11): Loading from cache
  - Installing symfony/http-foundation (v3.2.14): Loading from cache
  - Installing symfony/psr-http-message-bridge (v1.0.0): Loading from cache
  - Installing symfony/process (v3.2.14): Loading from cache
  - Installing symfony/polyfill-iconv (v1.6.0): Loading from cache
  - Installing symfony/event-dispatcher (v3.2.14): Loading from cache
  - Installing psr/log (1.0.2): Loading from cache
  - Installing symfony/debug (v3.3.13): Loading from cache
  - Installing symfony/http-kernel (v3.2.14): Loading from cache
  - Installing symfony/dependency-injection (v3.2.14): Loading from cache
  - Installing symfony/console (v3.2.14): Loading from cache
  - Installing symfony/class-loader (v3.2.14): Loading from cache
  - Installing symfony-cmf/routing (1.4.1): Loading from cache
  - Installing stack/builder (v1.0.5): Loading from cache
  - Installing masterminds/html5 (2.3.0): Loading from cache
  - Installing guzzlehttp/psr7 (1.4.2): Loading from cache
  - Installing guzzlehttp/promises (v1.3.1): Loading from cache
  - Installing guzzlehttp/guzzle (6.3.0): Loading from cache
  - Installing doctrine/lexer (v1.0.1): Loading from cache
  - Installing egulias/email-validator (1.2.14): Loading from cache
  - Installing easyrdf/easyrdf (0.9.1): Loading from cache
  - Installing doctrine/inflector (v1.2.0): Loading from cache
  - Installing doctrine/collections (v1.5.0): Loading from cache
  - Installing doctrine/cache (v1.7.1): Loading from cache
  - Installing doctrine/annotations (v1.5.0): Loading from cache
  - Installing doctrine/common (v2.8.1): Loading from cache
  - Installing composer/semver (1.4.2): Loading from cache
  - Installing asm89/stack-cors (1.1.0): Loading from cache
  - Installing drupal/core (8.4.2): Loading from cache
  - Installing drupal/devel (dev-1.x be072e7): Cloning be072e747c from cache
zendframework/zend-feed suggests installing zendframework/zend-cache (Zend\Cache component, for optionally caching feeds between requests)
zendframework/zend-feed suggests installing zendframework/zend-db (Zend\Db component, for use with PubSubHubbub)
zendframework/zend-feed suggests installing zendframework/zend-http (Zend\Http for PubSubHubbub, and optionally for use with Zend\Feed\Reader)
zendframework/zend-feed suggests installing zendframework/zend-servicemanager (Zend\ServiceManager component, for easily extending ExtensionManager implementations)
zendframework/zend-feed suggests installing zendframework/zend-validator (Zend\Validator component, for validating email addresses used in Atom feeds and entries ehen using the Writer subcomponent)
symfony/translation suggests installing symfony/config ()
symfony/validator suggests installing psr/cache-implementation (For using the metadata cache.)
symfony/validator suggests installing symfony/config ()
symfony/validator suggests installing symfony/expression-language (For using the Expression validator)
symfony/validator suggests installing symfony/intl ()
symfony/serializer suggests installing psr/cache-implementation (For using the metadata cache.)
symfony/serializer suggests installing symfony/config (For using the XML mapping loader.)
symfony/serializer suggests installing symfony/property-access (For using the ObjectNormalizer.)
symfony/serializer suggests installing symfony/property-info (To deserialize relations.)
symfony/routing suggests installing symfony/config (For using the all-in-one router or any loader)
symfony/routing suggests installing symfony/expression-language (For using expression matching)
paragonie/random_compat suggests installing ext-libsodium (Provides a modern crypto API that can be used to generate random bytes.)
symfony/http-kernel suggests installing symfony/browser-kit ()
symfony/http-kernel suggests installing symfony/config ()
symfony/http-kernel suggests installing symfony/finder ()
symfony/http-kernel suggests installing symfony/var-dumper ()
symfony/dependency-injection suggests installing symfony/config ()
symfony/dependency-injection suggests installing symfony/expression-language (For using expressions in service container configuration)
symfony/dependency-injection suggests installing symfony/proxy-manager-bridge (Generate service proxies to lazy load them)
symfony/console suggests installing symfony/filesystem ()
symfony/class-loader suggests installing symfony/polyfill-apcu (For using ApcClassLoader on HHVM)
easyrdf/easyrdf suggests installing ml/json-ld (~1.0)
doctrine/cache suggests installing alcaeus/mongo-php-adapter (Required to use legacy MongoDB driver)
drupal/devel suggests installing symfony/var-dumper (Pretty print complex values better with var-dumper available)
Writing lock file
Generating autoload files
INFO[0089] Found template /vagrant/sites/mysiteone/drubuild/sites.php.tmpl for usage
INFO[0089] Created directory /vagrant/sites/mysiteone/mysiteone.20171124081132_7047/web/sites/mysiteone
INFO[0089] Permissions set to 0755 on /vagrant/sites/mysiteone/mysiteone.20171124081132_7047/web/sites/mysiteone
INFO[0089] Successfully templated multisite config to file /vagrant/sites/mysiteone/mysiteone.20171124081132_7047/web/sites//sites.php
INFO[0089] Created symlink




You are about to create a sites/mysiteone/settings.php file and DROP all tables in your 'mysiteone_20171124081132_7047' database. Do you want to continue? (y/n): y






Starting Drupal installation. This takes a while. Consider using the [ok]
--notify global option.



Installation complete.  User name: admin  User password: o2pJvkvfoM  [ok]
Congratulations, you installed Drupal!                               [status]
INFO[0119] Found template /vagrant/sites/mysiteone/drubuild/vhost.tmpl for usage
INFO[0119] Successfully templated /vagrant/sites/mysiteone/drubuild/vhost.tmpl to file /etc/nginx/sites-enabled//mysiteone.test.conf
INFO[0119] Found template /vagrant/sites/mysiteone/drubuild/alias.tmpl for usage
INFO[0119] Successfully templated alias to file /home/vagrant/.drush/mysiteone.test.alias.drushrc.php
INFO[0119] Based upon the output above, you may need to restart the web service.
vagrant@drupalvm2:/vagrant/sites/mysiteone$ sudo service nginx restart
vagrant@drupalvm2:/vagrant/sites/mysiteone$ curl http://mysiteone.test/ | grep drupal
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0<meta name="Generator" content="Drupal 8 (https://www.drupal.org)" />
        <a href="/user/login" data-drupal-link-system-path="user/login">Log in</a>
100  8518    0  8518    0     0   225k      0 --:--:-- --:--:-- --:--:--  231k
        <a href="/" data-drupal-link-system-path="&lt;front&gt;" class="is-active">Home</a>
    <div class="search-block-form block block-search container-inline" data-drupal-selector="search-block-form" id="block-bartik-search" role="search">
        <input title="Enter the terms you wish to search for." data-drupal-selector="edit-keys" type="search" id="edit-keys" name="keys" value="" size="15" maxlength="128" class="form-search" />
<div data-drupal-selector="edit-actions" class="form-actions js-form-wrapper form-wrapper" id="edit-actions"><input class="search-form__submit button js-form-submit form-submit" data-drupal-selector="edit-submit" type="submit" id="edit-submit" value="Search" />
        <a href="/contact" data-drupal-link-system-path="contact">Contact</a>
      <span>Powered by <a href="https://www.drupal.org">Drupal</a></span>
````

From here, you simply need to add a hosts entry for the domain, and a drush alias if accessing the site via drush from the host machine.

## Author Information

This product was originally created in 2016 by [Karl Hepworth](https://twitter.com/fubarhouse).

Image of Go's mascot was created by [Takuya Ueda](https://twitter.com/tenntenn). Licenced under the Creative Commons 3.0 Attributions license. This image has been resized for purpose, but is otherwise unchanged.

## License

MIT - Free to use and manipulate with no guaranteed support by the creator.

Obviously we want to make the best product available so PR's, bug reports and feature requests are welcome! 