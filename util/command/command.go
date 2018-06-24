package command

import (
	"os"
	"os/exec"

	log "github.com/Sirupsen/logrus"
	"bytes"
	"io"
)


// Run a composer command with the input arguments.
func Run(name string, args []string) (string, error) {

	// Generate the command, based on input.
	cmd := exec.Cmd{}
	cmd.Path = name
	cmd.Args = []string{cmd.Path}

	// Add our arguments to the command.
	for _, arg := range args {
		cmd.Args = append(cmd.Args, arg)
	}

	// Create a buffer for the output.
	var out bytes.Buffer
	multi := io.MultiWriter(&out, os.Stdout)

	// Assign the output to the writer.
	cmd.Stdout = multi

	// Check the errors, return as needed.
	if err := cmd.Run(); err != nil {
		log.Errorln(err)
		return out.String(), err
	}
	cmd.Wait()

	// Return out output as a string.
	return out.String(), nil

}
