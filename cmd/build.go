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
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
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

		if _, err := os.Stat(destination); err != nil {
			log.Infoln("Creating directory", destination)
			dirErr := os.MkdirAll(destination, 0755)
			if dirErr != nil {
				log.Errorln("Unable to create directory", destination, dirErr)
			} else {
				log.Infoln("Created directory", destination)
			}
		}

		timestamp, _ = strconv.ParseInt(time.Now().Format("20060102150405"), 0, 0)
		c.Copy(composer, destination)

		cargs := []string{"install", "--working-dir", destination}
		if preferSource {
			cargs = append(cargs, "--prefer-source")
		}
		if workingCopy {
			cargs = append(cargs, "--prefer-dist")
		}

		log.Infoln("Running composer...")
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
