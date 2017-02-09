package make

import (
	"github.com/fubarhouse/golang-drush/alias"
	"github.com/fubarhouse/golang-drush/aliases"
	"strings"
)

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

func (Site *Site) AliasInstall() {
	siteAlias := alias.NewAlias(Site.Name, Site.Path+"_latest", Site.Alias)
	siteAlias.Install()
}

func (Site *Site) AliasUninstall() {
	siteAlias := alias.NewAlias(Site.Name, Site.Path+"_latest", Site.Alias)
	siteAlias.Uninstall()
}
