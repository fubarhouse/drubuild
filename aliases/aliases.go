package aliases

import (
	"fmt"
	"github.com/fubarhouse/golang-drush/alias"
	"github.com/fubarhouse/golang-drush/command"
	"strings"
)

// AliasList struct to contain a slice of Alias structs
type AliasList struct {
	value []*alias.Alias
}

// NewAliasList will instantiate an AliasList struct
func NewAliasList() *AliasList {
	return &AliasList{}
}

// Add an alias to an AliasList
func (list *AliasList) Add(item *alias.Alias) {
	list.value = append(list.value, item)
}

// Generate an AliasList from a given key from all available aliases
func (list *AliasList) Generate(key string) {
	sites := command.NewDrushCommand()
	sites.Set("", "sa", false)
	values, _ := sites.Output()
	values = strings.Split(fmt.Sprintf("%v", values), "\n")
	for _, currAlias := range values {
		if strings.Contains(currAlias, key) == true {
			thisAlias := alias.NewAlias(currAlias, "", "")
			list.Add(thisAlias)
		}
	}
}

// Filter an AliasList by a given key.
func (list *AliasList) Filter(key string) {
	values := list.GetNames()
	newList := NewAliasList()
	for _, currAlias := range values {
		if strings.Contains(currAlias, key) == true {
			thisAlias := alias.NewAlias(currAlias, "", "")
			newList.Add(thisAlias)
		} else {
			fmt.Sprintln("Filtered out", currAlias)
		}
	}
	*list = *newList
}

// Count will return how many aliases are in the AliasList
func (list *AliasList) Count() int {
	count := 0
	for _, thisAlias := range list.value {
		fmt.Sprintln(thisAlias)
		count++
	}
	return count
}

// GetNames gets a list of alias names from the AliasList items
func (list *AliasList) GetNames() []string {
	returnVals := []string{}
	for _, val := range list.value {
		returnVals = append(returnVals, val.GetName())
	}
	return returnVals
}

// GetAliasNames gets alias uri fields from AliasList items
func (list *AliasList) GetAliasNames() []string {
	returnVals := []string{}
	for _, val := range list.value {
		returnVals = append(returnVals, val.GetUri())
	}
	return returnVals
}

// GetAliases gets value field from AliasList items
func (list *AliasList) GetAliases() *AliasList {
	returnVals := NewAliasList()
	for _, val := range list.value {
		returnVals.value = append(returnVals.value, val)
	}
	return returnVals
}
