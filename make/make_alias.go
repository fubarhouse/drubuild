package make

import (
	"strings"

	"github.com/fubarhouse/golang-drush/alias"
	"github.com/fubarhouse/golang-drush/aliases"
)

// AliasExists returns a boolean for the status of a given alias in a given list.
func (Site *Site) AliasExists(filter string) bool {
	y := aliases.NewAliasList()
	y.Generate(filter)
	for _, z := range y.GetNames() {
		if strings.Contains(z, Site.Alias) {
			return true
		}
	}
	return false
}

// AliasInstall installs an alias for a given site struct
func (Site *Site) AliasInstall(docroot string) {
	var siteAlias alias.Alias
	siteAlias.Docroot = docroot
	siteAlias.SetName(Site.Name)
	siteAlias.SetPath(Site.Path)
	siteAlias.SetURI(Site.Alias)
	siteAlias.SetTemplate(Site.AliasTemplate)
	siteAlias.Install()
}

// AliasUninstall un-installs an alias for a given site struct
func (Site *Site) AliasUninstall() {
	var siteAlias alias.Alias
	siteAlias.SetName(Site.Name)
	siteAlias.SetPath(Site.Path)
	siteAlias.SetURI(Site.Alias)
	siteAlias.Uninstall()
}
