package aliases

import (
	"fmt"
	"github.com/fubarhouse/golang-drush/alias"
	"github.com/fubarhouse/golang-drush/command"
	"strings"
)

type AliasList struct {
	// A simple Alias List for attaching methods.
	value []*alias.Alias
}

func NewAliasList() *AliasList {
	// Create a new but empty Alias List
	return &AliasList{}
}

func (list *AliasList) Add(item *alias.Alias) {
	// Add an alias to the alias list.
	list.value = append(list.value, item)
}

func (list *AliasList) Generate(key string) {
	// Add a set of aliases to an Alias List based of a string value.
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

func (list *AliasList) Filter(key string) {
	// Filter an existing list with a key string
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

func (list *AliasList) Count() int {
	count := 0
	for _, thisAlias := range list.value {
		fmt.Sprintln(thisAlias)
		count++
	}
	return count
}

func (list *AliasList) GetNames() []string {
	// Return values from the Alias List object
	returnVals := []string{}
	for _, val := range list.value {
		returnVals = append(returnVals, val.GetName())
	}
	return returnVals
}

func (list *AliasList) GetAliasNames() []string {
	// Return values from the Alias List object
	returnVals := []string{}
	for _, val := range list.value {
		returnVals = append(returnVals, val.GetUri())
	}
	return returnVals
}

func (list *AliasList) GetAliases() *AliasList {
	// Return values from the Alias List object
	returnVals := NewAliasList()
	for _, val := range list.value {
		returnVals.value = append(returnVals.value, val)
	}
	return returnVals
}
