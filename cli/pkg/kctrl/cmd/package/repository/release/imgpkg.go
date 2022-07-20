package release

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	goexec "os/exec"
	"path/filepath"
	"strings"

	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
)

type ImgpkgRunner struct {
	BundlePath        string
	Paths             []string
	UseKbldImagesLock bool
	ImgLockFilepath   string
	UI                cmdcore.AuthoringUI
}

func (r ImgpkgRunner) Run() (string, error) {
	dir, err := os.MkdirTemp(".", fmt.Sprintf("bundle-%s-*", strings.Replace(r.BundlePath, "/", "-", -1)))
	if err != nil {
		return "", err
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for _, path := range r.Paths {
		var stderrBuf bytes.Buffer
		cmd := goexec.Command("cp", "-r", filepath.Join(workingDirectory, path), dir)
		cmd.Stderr = &stderrBuf
		err := cmd.Run()
		if err != nil {
			return "", fmt.Errorf("%s", stderrBuf.String())
		}
	}
	if r.UseKbldImagesLock {
		err = goexec.Command("mkdir", filepath.Join(dir, LockOutputFolder)).Run()
		if err != nil {
			return "", err
		}
		err = goexec.Command("cp", r.ImgLockFilepath, filepath.Join(dir, LockOutputFolder, LockOutputFile)).Run()
		if err != nil {
			return "", err
		}
	}
	defer os.RemoveAll(dir)

	r.UI.PrintInformationalText("\nAn imgpkg bundle consists of all required YAML configuration bundled into an OCI image " +
		"that can be pushed to an image registry and consumed by the package.\n")
	r.UI.PrintHeaderText("Pushing imgpkg bundle (Step 3/3)")

	imgpkgCmdRunner := exec.NewPlainCmdRunner()
	cmd := goexec.Command("imgpkg", "push", "-b", r.BundlePath, "-f", dir, "--tty=true")
	r.UI.PrintCmdExecutionOutput(fmt.Sprintf("$ %s", strings.Join(cmd.Args, " ")))

	var stdoutBuf, stderrBuf bytes.Buffer
	inMemoryStdoutWriter := bufio.NewWriter(&stdoutBuf)
	cmd.Stdout = io.MultiWriter(inMemoryStdoutWriter)
	cmd.Stderr = &stderrBuf
	err = imgpkgCmdRunner.Run(cmd)
	if err != nil {
		return stdoutBuf.String(), fmt.Errorf("%s", stderrBuf.String())
	}
	err = inMemoryStdoutWriter.Flush()
	if err != nil {
		return "", err
	}
	r.UI.PrintCmdExecutionOutput(stdoutBuf.String())

	return r.imgpkgBundleURLFromStdout(stdoutBuf.String())
}

func (r *ImgpkgRunner) imgpkgBundleURLFromStdout(imgpkgStdout string) (string, error) {
	lines := strings.Split(imgpkgStdout, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Pushed") {
			line = strings.TrimPrefix(line, "Pushed")
			line = strings.Replace(line, "'", "", -1)
			line = strings.Replace(line, " ", "", -1)
			return line, nil
		}
	}
	return "", fmt.Errorf("Could not get imgpkg bundle location")
}
