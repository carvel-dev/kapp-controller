package init

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	goexec "os/exec"

	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
)

type VendirRunner struct {
	ui cmdcore.AuthoringUI
}

func NewVendirRunner(ui cmdcore.AuthoringUI) VendirRunner {
	return VendirRunner{ui: ui}
}

func (r VendirRunner) RunSync() error {
	_, err := os.Stat(vendirFileName)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if os.IsNotExist(err) {
		return nil
	}

	r.ui.PrintInformationalText("We will use vendir to fetch the data from the source to the local directory. Vendir allows us to declaratively state what should be in a directory and sync data sources into it. All the information entered above has been persisted into a vendir.yml file.")
	err = r.printVendirFile()
	if err != nil {
		return err
	}
	err = r.runVendirSync()
	if err != nil {
		return err
	}
	return nil
}

func (r VendirRunner) printVendirFile() error {
	r.ui.PrintActionableText(fmt.Sprintf("Printing %s \n", vendirFileName))
	err := r.printFile(vendirFileName)
	if err != nil {
		return err
	}
	return nil
}

func (r VendirRunner) printFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Printing file: %w", err)
	}
	defer func() {
		file.Close()
	}()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		r.ui.PrintCmdExecutionOutput(scanner.Text())
	}
	return nil
}

// TODO vendir sync failure. Reproduce: In case of 429 from github, we dont show errors today.
func (r VendirRunner) runVendirSync() error {
	r.ui.PrintInformationalText("\nNext step is to run `vendir sync` to fetch the data from the source to the local directory. Vendir will sync the data into the upstream folder.")
	r.ui.PrintActionableText("Running vendir sync")
	r.ui.PrintCmdExecutionText("vendir sync -f vendir.yml\n")
	var stdoutBs, stderrBs bytes.Buffer

	localCmdRunner := exec.NewPlainCmdRunner()
	cmd := goexec.Command("vendir", []string{"sync", "-f", vendirFileName}...)
	cmd.Stdin = nil
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs
	err := localCmdRunner.Run(cmd)
	result := exec.CmdRunResult{
		Stdout: stdoutBs.String(),
		Stderr: stderrBs.String(),
	}
	result.AttachErrorf("Fetching resources: %s", err)
	if result.Error != nil {
		return fmt.Errorf("Running vendir sync: %s", result.Stderr)
	}
	return nil
}

func (r VendirRunner) listFiles(dir string) error {
	var stdoutBs, stderrBs bytes.Buffer

	localCmdRunner := exec.NewPlainCmdRunner()
	cmd := goexec.Command("ls", []string{"-lR", dir}...)
	cmd.Stdin = nil
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs
	localCmdRunner.Run(cmd)

	result := exec.CmdRunResult{
		Stdout: stdoutBs.String(),
		Stderr: stderrBs.String(),
	}
	if result.Error != nil {
		return fmt.Errorf("Listing files.\n %s", result.Stderr)
	}
	r.ui.PrintCmdExecutionOutput(result.Stdout)
	return nil
}
