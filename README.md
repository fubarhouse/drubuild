<img style="float:left" alight="left" height="128px" width="100px" src="https://github.com/fubarhouse/ansible-role-golang/raw/master/gopher.png">

# Drubuild

[![Go Report Card](https://goreportcard.com/badge/github.com/fubarhouse/drubuild)](https://goreportcard.com/report/github.com/fubarhouse/drubuild)

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
drubuild build --name mysite --domain mysite.dev --destination /sites/mysite.dev --composer /composer.json
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
drubuild sync --source mysite.prod --destination mysite.local  --database --files
```

### User

The user command will act much like the runner, which can create, block, unblock, delete, reset password, add roles, and general maintenance accross a set of sites using alias pattern matching.

This one has proven critical when managing literally hundreds of Drupal accounts.
```
drubuild user --name TestUser --email test@user.com --password MyPassword --aliases dev,test,preprod,prod --pattern mysite.%v --create
```

## Author Information

This product was originally created in 2016 by [Karl Hepworth](https://twitter.com/fubarhouse).

Image of Go's mascot was created by [Takuya Ueda](https://twitter.com/tenntenn). Licenced under the Creative Commons 3.0 Attributions license. This image has been resized for purpose, but is otherwise unchanged.

## License

MIT - Free to use and manipulate with no guaranteed support by the creator.

Obviously we want to make the best product available so PR's, bug reports and feature requests are welcome! 