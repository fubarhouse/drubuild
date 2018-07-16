package make

import (
	"strings"
	"github.com/fubarhouse/drubuild/util/aliases"
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
