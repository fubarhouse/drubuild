# Drush

A wide array of tools designed for system and drush specific use during scripting and dev-ops.

**Package rundown**:

* aliase
    ````
    Provides the ability to create an alias object and provides a range of alias specific methods including installation and uninstallation.
    ````
* aliases
    ````
    Execute 'drush sa' and grab all aliases, and provides convenient ways to filter and store those results.
    ````
* command
    ````
    Execute a Drush command via go
    ````
* commandlist
    ````
    An API to line up an infinite amount of Drush commands in various objects methods for the API including the ability to execute them.
    ````
* make
    ````
    Provides a way to communicate drush make and site-install commands to the command package.
    ````
* makeupdater
    ````
    Very basic functions which you can pass an absolute path, and for any contibuted modules the version number will be updated. It will also return a []string with those values, should you need to use those.
    ````
* sites
    ````
    An API to create a list of sites available to Drush, and ways to add and remove individual sites, or to apply sets of sites via filters.
    ````
* vhost
    ````
    An unfinished tool which can create and remove virtual host files.
    ````

## Purpose

This repository features a broad array of Drush-related packages designed to be portable, useful and platform independent (mostly). It's intended so that I have a multitude of tools available to me during the course of my career when relating to Drupal. Dev-ops has become a critical part of my job, and this tool-set provides many rewritten tools based on previous experience.

The intent here is strictly to have the tool I want available to me in the public domain. Absolutely none of the work in this repository can be claimed as an asset of anybody.

The packages have been designed for ease-of scripting, with everything configurable.

This repository does not include the packages associated to my Drush Version Management tool (also hosted on GitHub).

## Install

```console
$ go get github.com/fubarhouse/golang-drush
```

## License

MIT