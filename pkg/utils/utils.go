package utils

import (
	"fmt"
	"bytes"
	"os/exec"
)

// CmdExists - make sure a command exists before executing it
func CmdExists(cmdlet string) error {
	cmd := exec.Command("which", cmdlet)	// we use `which` for the probing
	var out bytes.Buffer
	var stderr bytes.Buffer
	// redirect the stdout and stderr to the created bytebuffer.
	// this ensures the user doesn't see anything from the command probing.
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	// run the command and check for errors.
	err := cmd.Run()
	if err != nil {
		fmt.Println("")
		return fmt.Errorf(fmt.Sprint(err) + ": " + stderr.String())
	}
	return nil
}

// InRange - make sure that the next element is in range before assigning it.
func InRange(min int, max int) bool {
	if min < max {
		return true
	}
	return false
}
