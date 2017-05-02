package sites

import (
	"fmt"
	"github.com/fubarhouse/golang-drush/alias"
	"github.com/fubarhouse/golang-drush/aliases"
	"strings"
)

type SiteList struct {
	// A SiteList struct with an attached key for attaching methods.
	value []string
	key   string
}

func NewSiteList() *SiteList {
	// Create an empty SiteList object.
	return &SiteList{}
}

func (list *SiteList) SetKey(key string) {
	// Set the key to a string value.
	// The key is used as a filter based off the output of 'drush sa'
	list.key = key
}

func (list *SiteList) Rewrite(oldString string, newString string) {
	// Rewrite a set of values from a SiteList
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

func (list *SiteList) SetList() {
	// Adds a set of aliases to a SiteList.
	aliases := aliases.NewAliasList()
	aliases.Generate(list.key)
	for _, alias := range aliases.GetNames() {
		list.value = append(list.value, alias)
	}
}

func (list *SiteList) GetList() []string {
	// Return the dataset in the SiteList object.
	return list.value
}

func (list *SiteList) Count() int {
	// Return the quantity of items in the SiteList object.
	return len(list.value)
}
