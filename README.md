<img style="float:left" alight="left" height="128px" width="100px" src="https://github.com/fubarhouse/ansible-role-golang/raw/master/gopher.png">

# Golang Drush

[![Go Report Card](https://goreportcard.com/badge/github.com/fubarhouse/golang-drush)](https://goreportcard.com/report/github.com/fubarhouse/golang-drush)

## Purpose

Golang-Drush is a suite of tools designed for use behind the terminal via scripts or to tie into CI configurations which are designed to build and maintain Drupal websites and support the needs of development teams in building Drupal sites with a vast array of Drush integration tasks to remove the complexities of having developers behind a command-line.

History proves that not all developers are comfortable or fo not completely embrace command-line tools and these non-interactive tools allow use in many different configurations and platforms. 

## Package rundown

* alias:
  Provides types and functions associated to managing a single Drush alias.
* aliases
  Provides types and functions associated to grouping a collection of Drush aliases.
* command
  Provides types and functions associated to execution of a Drush command, in a range of ways including live pipelines.
* commandlist
  Provides types and functions for grouping a list/group of Drush commands for execution.
* make
  Provides many types and functions associated to the creation and removal of many aspects of Drupal sites via Drush.
* makeupdater
  Provides types and functions associated to updating make files, including make file creation, recreation and generation.
* sites
  Provides types and functions for a collection of sites, such as finding an available site via `drush sa`
* user
  Provides types and functions for user management, including creating, blocking and verification.
* vhost
  Provides types and functions for creating and removing virtual hosts.

## Included binaries

All of the following binaries should be used with the -h flag to invoke the usage. Usage will not be supplied here explicitly as there are a lot of binaries and a lot of potential parameters for them.

* module-auditor  
  Uses Drush to run run a report against make files and Drush aliases.
* rewrite-make  
  Completely rewrites a make file - supports contributed make files only.
* site-checker  
  Binary name needs to be rewritten, but it will run a series of commands on a series of aliases matching a specified pattern and report results and verbose output if desired.
* update-make  
  Runs through a make file with pm-info and updates version numbers to the latest available recommended version for each project.
* user-block  
  Performs a set of actions to block a given user on aliases matching a specific pattern.
* user-create  
  Performs a set of actions to create a given user on aliases matching a specific pattern.
* user-unblock  
  Performs a set of actions to unblock a given user on aliases matching a specific pattern.
* user-verify  
  Performs a set of actions to validate and change the information on a given user on aliases matching a specific pattern.
* yoink-backup-site  
  Performs a Drush archive-dump command on an alias to a given destination
* yoink-build-site  
  Builds a Drupal website based upon MySQL with a drush alias, virtualhosts from specified make files, and supports infinite amount of builds per site.
* yoink-destroy-site  
  Removes everything put in place by the build program. 
* yoink-rebuild-site  
  Rebuilds a site without virtual hosts, aliases or anything, it will build a site at a specific location with given make files.
* yoink-solr-build  
  Installs a solr core with provided resource files.
* yoink-solr-destroy  
  Removes a solr core installed by the solr build program.
* yoink-sync-remote-site  
  Syncs files and/or database between a source and destination alias.
* yoink-sync-site  
  Syncs files and/or database between a source and destination alias, which performs basic checks on the destination upon completion.
* yoink-validate  
  Runs some basic system tests to ensure funtionality will execute.

## Install

### Install the entire package
```console
$ go get github.com/fubarhouse/golang-drush/...
```

### Installing individual binaries
```console
$ go get github.com/fubarhouse/golang-drush/cmd/module-auditor
$ go get github.com/fubarhouse/golang-drush/cmd/rewrite-make
$ go get github.com/fubarhouse/golang-drush/cmd/site-checker
$ go get github.com/fubarhouse/golang-drush/cmd/update-make
$ go get github.com/fubarhouse/golang-drush/cmd/user-block
$ go get github.com/fubarhouse/golang-drush/cmd/user-create
$ go get github.com/fubarhouse/golang-drush/cmd/user-unblock
$ go get github.com/fubarhouse/golang-drush/cmd/user-verify
$ go get github.com/fubarhouse/golang-drush/cmd/yoink-backup-site
$ go get github.com/fubarhouse/golang-drush/cmd/yoink-build-site
$ go get github.com/fubarhouse/golang-drush/cmd/yoink-destroy-site
$ go get github.com/fubarhouse/golang-drush/cmd/yoink-rebuild-site
$ go get github.com/fubarhouse/golang-drush/cmd/yoink-solr-build
$ go get github.com/fubarhouse/golang-drush/cmd/yoink-solr-destroy
$ go get github.com/fubarhouse/golang-drush/cmd/yoink-sync-remote-site
$ go get github.com/fubarhouse/golang-drush/cmd/yoink-sync-site
$ go get github.com/fubarhouse/golang-drush/cmd/yoink-validate
```

## Author Information

This product was originally created in 2016 by [Karl Hepworth](https://twitter.com/fubarhouse).

Image of Go's mascot was created by [Takuya Ueda](https://twitter.com/tenntenn). Licenced under the Creative Commons 3.0 Attributions license. This image has been resized for purpose, but is otherwise unchanged.

## License

MIT - Free to use and manipulate with no guaranteed support by the creator.

Obviously we want to make the best product available so PR's, bug reports and feature requests are welcome! 