<img style="float:left" alight="left" height="128px" width="100px" src="https://github.com/fubarhouse/ansible-role-golang/raw/master/gopher.png">

# Drubuild 0.3.x

[![stability-experimental](https://img.shields.io/badge/stability-experimental-orange.svg?style=for-the-badge)](https://github.com/orangemug/stability-badges)
[![Go Report Card](https://goreportcard.com/badge/github.com/fubarhouse/drubuild?style=for-the-badge)](https://goreportcard.com/report/github.com/fubarhouse/drubuild)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg?style=for-the-badge)](https://raw.githubusercontent.com/fubarhouse/brand/master/LICENSE.txt)


> version 0.3.x is experimental, please avoid use of this version.

## Purpose

Drubuild is a command-line application which builds and manages sites via some common composer and drush commands.

This application was born of the personal desire of the creator to be a useful tool for CI and automation to be used at work, however the stability and reliability was a catalyst for a lot of work here. 

This application is fully-intended to work out of the box with [DrupalVM](https://www.drupalvm.com/), and a [fork of DrupalVM](https://github.com/fubarhouse/drupal-vm) with the configuration of the maintainer's Ansible role [fubarhouse.golang](https://github.com/fubarhouse/ansible-role-golang) is available for quick opportunities to use this in an isolated environment.

The application has previously been used for CI tooling for the management of more than 40 websites simultaneously, supporting developers during their local build and development cycle with Jenkins. 

## Install

It is *highly* recommended to install this using Go, and currently no other options are available. 

```sh
$ go get -u github.com/fubarhouse/drubuild
```

## Usage

> Usage documentation for v0.3.x is not available, but will become available closer to the release.

## Responsibilities

A __lot__ has changed since 0.3.x, and due to extended usage of 0.2.x we've been able to drive the project forward in the way of simplicity, conciseness and ultimately value.

Drubuild 0.3.x changes responsibilities from version 0.2.x, the changes include, but are not limited to the following:

* Drubuild no longer handles virtual hosts or webservers to _any_ degree - there is _zero_ support.
* Drubuild no longer interacts directly with databases - databases must be created and dropped by the developer.
* Drubuild no longer is responsible for the features previously provided in the destroy command, the developer is responsible for the database and file system removal. Aliases can be uninstalled with Drubuild, or they can be deleted independently.
* Drubuild no longer creates, changes or removes symlinks.
* Site installation processes have been isolated to their own command, which should improve CI pipelines.
* Drush alias processes have been isolated to their own command, which should improve CI pipelines.
* There's no longer an initialise command, and templates are no longer allowed to stray from default, please submit a PR if you feel this content should deviate.
* API's have been streamlined for drush and composer commands, which generated the interest in a refreshed version. This stemmed initially from the work on [Ansible Role Tester](https://github.com/fubarhouse/ansible-role-tester)

## Author Information

This product was originally created in 2016 by [Karl Hepworth](https://twitter.com/fubarhouse).

Image of Go's mascot was created by [Takuya Ueda](https://twitter.com/tenntenn). Licenced under the Creative Commons 3.0 Attributions license. This image has been resized for purpose, but is otherwise unchanged.

## License

MIT - Free to use and manipulate with no guaranteed support by the creator.

Obviously we want to make the best product available so PR's, bug reports and feature requests are welcome! 