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
		"log"
	"os"

	"github.com/spf13/cobra"
		"github.com/fubarhouse/drubuild/util/drush"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Take a archive-dump snapshot of a local site",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		drush.Run([]string{source, "archive-dump", "--destination=" + destination})
	},
}

func init() {

	RootCmd.AddCommand(backupCmd)

	// Get the current working directory.
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Flags
	backupCmd.Flags().StringVarP(&source, "source", "s", "", "Drush alias to use for operation")
	backupCmd.Flags().StringVarP(&destination, "destination", "d", dir, "Path to Drush archive-dump destination")
	// Mark flags as required
	backupCmd.MarkFlagRequired("source")
	backupCmd.MarkFlagRequired("destination")
}
