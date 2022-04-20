package util

import (
	"bytes"
	"fmt"
	"os/exec"
)

func Execute(cmd string, args []string) (string, error) {
	cmdFound := isExecutableInstalled(cmd)
	if !cmdFound {
		fmt.Printf("Executable \"%s\" not installed", cmd)
	}
	command := exec.Command(cmd, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &out
	command.Stderr = &stderr

	err := command.Run()
	if err != nil {
		return "", fmt.Errorf("%s: %s", fmt.Sprint(err), stderr.String())
	}
	return out.String(), nil
}

func isExecutableInstalled(cmd string) bool {
	_, err := exec.LookPath(cmd)
	if err != nil {
		//TODO should log an error here instead of sending it up
		return false
	}
	return true

}
