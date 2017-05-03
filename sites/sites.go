package sites

import (
	"github.com/fubarhouse/golang-drush/alias"
	"github.com/fubarhouse/golang-drush/aliases"
	"strings"
)

// SiteList struct with an attached key for attaching methods.
type SiteList struct {
	value []string
	key   string
}

// NewSiteList instantiates an empty SiteList struct.
func NewSiteList() *SiteList {
	// Create an empty SiteList object.
	return &SiteList{}
}

// SetKey sets the key field to an input for a SiteList struct.
func (list *SiteList) SetKey(key string) {
	// The key is used as a filter based off the output of 'drush sa'
	list.key = key
}

// Rewrite rewrites a set of values from a SiteList
func (list *SiteList) Rewrite(oldString string, newString string) {
	aliasesList := aliases.NewAliasList()
	aliasesList.Generate(list.key)
	aliasesFiltered := aliases.NewAliasList()
	for _, thisAlias := range list.GetList() {
		//alias = strings.Replace(alias,deleteString, "", -1)
		if strings.Contains(thisAlias, oldString) {
			thisAlias = strings.Replace(thisAlias, oldString, newString, -1)
			newAlias := alias.NewAlias(thisAlias, "", thisAlias)
			aliasesFiltered.Add(newAlias)
		} else {
			newAlias := alias.NewAlias(thisAlias, "", thisAlias)
			aliasesFiltered.Add(newAlias)
		}
	}
	list.value = aliasesFiltered.GetNames()
}

// Remove will remove an entry from a SiteList based on an input
func (list *SiteList) Remove(remove string) {
	// Removes a set of values from a SiteList
	aliasesList := aliases.NewAliasList()
	aliasesList.Generate(list.key)
	aliasesFiltered := aliases.NewAliasList()
	for _, thisAlias := range list.GetList() {
		//alias = strings.Replace(alias,deleteString, "", -1)
		if !strings.Contains(thisAlias, remove) {
			newAlias := alias.NewAlias(thisAlias, "", thisAlias)
			aliasesFiltered.Add(newAlias)
		}
	}
	list.value = aliasesFiltered.GetNames()
}

// FilterBy will filter by an input for an entry from a SiteList
func (list *SiteList) FilterBy(filter string) {
	// Filters a sataset by a set of values from a SiteList
	aliasesList := aliases.NewAliasList()
	aliasesList.Generate(list.key)
	aliasesFiltered := aliases.NewAliasList()
	for _, thisAlias := range list.GetList() {
		//alias = strings.Replace(alias,deleteString, "", -1)
		if strings.Contains(thisAlias, filter) {
			newAlias := alias.NewAlias(thisAlias, "", thisAlias)
			aliasesFiltered.Add(newAlias)
		}
	}
	list.value = aliasesFiltered.GetNames()
}

// SetList adds a set of aliases to a SiteList.
func (list *SiteList) SetList() {
	aliases := aliases.NewAliasList()
	aliases.Generate(list.key)
	for _, alias := range aliases.GetNames() {
		list.value = append(list.value, alias)
	}
}

// GetList returns the dataset in the SiteList object.
func (list *SiteList) GetList() []string {
	return list.value
}

// Count returns the quantity of items in the SiteList object.
func (list *SiteList) Count() int {
	return len(list.value)
}
