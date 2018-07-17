// Copyright © 2017 Karl Hepworth <Karl.Hepworth@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"path/filepath"
	c "github.com/fubarhouse/drubuild/util/composer"
)

// removeGitFromPath will purge all .git data recursively from the specified path.
func removeGitFromPath(path string) {
	// Generate a list of .git file systems from the input path.
	fileList := []string{}
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".git") && f.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})

	// Loop over the generated list to remove them.
	for _, file := range fileList {
		if err := os.RemoveAll(file); err != nil {
			fmt.Printf("could not delete file system %v: %v\n", file, err)
		} else {
			fmt.Printf("removed %v\n", file)
		}
	}
}

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "The build process for Drubuild",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		timestamp, _ = strconv.ParseInt(time.Now().Format("20060102150405"), 0, 0)
		c.Copy(composer, destination)

		cargs := []string{"install", "-d='"+destination+"'"}
		if preferSource {
			cargs = append(cargs, "--prefer-source")
		}
		if workingCopy {
			cargs = append(cargs, "--working-copy")
		}

		c.Run(cargs)
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)

	// Parameters/Flags
	buildCmd.Flags().StringVarP(&destination, "destination", "p", "", "The destination path")
	buildCmd.Flags().StringVarP(&composer, "composer", "c", "", "Path to composer.json")
	buildCmd.Flags().BoolVarP(&preferSource, "prefer-source", "s", false, "Build with preference to source packages.")
	buildCmd.Flags().BoolVarP(&workingCopy, "working-copy", "w", false, "Build with preference to working-copy")

	// Required flags
	buildCmd.MarkFlagRequired("name")
	buildCmd.MarkFlagRequired("destination")
	buildCmd.MarkFlagRequired("domain")
}
