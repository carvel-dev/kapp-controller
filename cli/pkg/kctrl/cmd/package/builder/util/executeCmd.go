package util

import (
	"bytes"
	goexec "os/exec"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
)

func Execute(cmd string, args []string) exec.CmdRunResult {
	var stdoutBs, stderrBs bytes.Buffer
	command := goexec.Command(cmd, args...)
	command.Stdout = &stdoutBs
	command.Stderr = &stderrBs
	err := command.Run()

	result := exec.CmdRunResult{
		Stdout: stdoutBs.String(),
		Stderr: stderrBs.String(),
	}
	result.AttachErrorf("Running Command %s", err)
	return result
}
