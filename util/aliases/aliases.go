package aliases

import (
	"fmt"
	"github.com/fubarhouse/drubuild/util/alias"
		"strings"
	"github.com/fubarhouse/drubuild/util/drush"
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
	aliases, _ := drush.Run([]string{"sa"})
	values := strings.Split(fmt.Sprintf("%v", aliases), "\n")
	for _, currAlias := range values {
		if strings.Contains(currAlias, key) == true {
			thisAlias := alias.NewAlias(currAlias, "", "")
			list.Add(thisAlias)
		}
	}
}

// GetNames gets a list of alias names from the AliasList items
func (list *AliasList) GetNames() []string {
	returnVals := []string{}
	for _, val := range list.value {
		returnVals = append(returnVals, val.GetName())
	}
	return returnVals
}