// Copyright © 2017 Karl Hepworth karl.hepworth@gmail.com
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"log"
	"os"

	c "github.com/fubarhouse/drubuild/util/composer"

		"github.com/spf13/cobra"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Install or remove a project.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if remove {
			c.Run([]string{"remove", name})
		}
		if add {
			var r string
			if preferSource {
				r = "require --prefer-source"
			} else {
				r = "require"
			}
			if version != "" {
				name += ":" + version
			}
			c.Run([]string{r, name})
		}

		//if !workingCopy {
		//	x, e := c.GetPath(destination, name)
		//	if e != nil {
		//		err := errors.New("could not find path associated to " + name)
		//		err.Error()
		//	}
		//	if x != "" {
		//		removeGitFromPath(x)
		//	}
		//}

		if !add && !remove {
			fmt.Println("No action selected, add --add or --remove for effect.")
		}

	},
}

func init() {
	RootCmd.AddCommand(projectCmd)

	// Get the current working directory.
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	projectCmd.Flags().StringVarP(&name, "name", "n", "", "The human-readable name for this site")
	projectCmd.Flags().StringVarP(&destination, "path", "p", dir, "The path to where the site(s) exist.")
	projectCmd.Flags().StringVarP(&version, "version", "v", "", "Version of the package.")
	projectCmd.Flags().BoolVarP(&add, "add", "a", false, "Flag to trigger add action.")
	projectCmd.Flags().BoolVarP(&remove, "remove", "r", false, "Flag to trigger remove action.")
	projectCmd.Flags().BoolVarP(&workingCopy, "working-copy", "w", false, "Mark as a working-copy.")
	projectCmd.Flags().BoolVarP(&preferSource, "prefer-source", "s", false, "Build with preference to source packages.")

	projectCmd.MarkFlagRequired("name")
	projectCmd.MarkFlagRequired("path")

}
