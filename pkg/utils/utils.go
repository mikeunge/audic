package utils

import (
	"bytes"
	"fmt"
	"os/exec"
)

// CmdExists - make sure a command exists before executing it
func CmdExists(cmdlet string) error {
	cmd := exec.Command("which", cmdlet) // we use `which` for the probing
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// capture old stdout & stderr
	oldStdout := cmd.Stdout
	oldStderr := cmd.Stderr

	// redirect the stdout and stderr to the created bytebuffer.
	// this ensures the user doesn't see anything from the command probing.
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// run the command and check for errors.
	err := cmd.Run()

	// Switch Stdout & Stderr back
	cmd.Stdout = oldStdout
	cmd.Stderr = oldStderr
	if err != nil {
		return fmt.Errorf(fmt.Sprint(err) + ": " + stderr.String())
	}
	return nil
}
