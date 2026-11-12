package sources

import (
	"bytes"
	"fmt"
	goexec "os/exec"

	cmdcore "carvel.dev/kapp-controller/cli/pkg/kctrl/cmd/core"
)

type VendirRunner struct {
	ui cmdcore.AuthoringUI
}

func NewVendirRunner(ui cmdcore.AuthoringUI) VendirRunner {
	return VendirRunner{ui: ui}
}

func (r VendirRunner) Sync(fetchMode string) error {
	if fetchMode == LocalDirectory {
		return nil
	}

	r.ui.PrintInformationalText("We will use vendir to fetch the data from the source to the local directory.\n" +
		"Vendir allows us to declaratively state what should be in a directory and sync data sources into it.\n" +
		"All the information entered above has been persisted into a vendir.yml file.")
	err := r.sync()
	if err != nil {
		return err
	}
	return nil
}

// TODO vendir sync failure. Reproduce: In case of 429 from github, we dont show errors today.
func (r VendirRunner) sync() error {
	r.ui.PrintInformationalText("\nNext step is to run 'vendir sync' to fetch the data from the source to the local directory.\n" +
		"Vendir will sync the data into the upstream folder.")
	r.ui.PrintActionableText("Running vendir sync")
	r.ui.PrintCmdExecutionText("vendir sync -f vendir.yml\n")

	var stderrBs bytes.Buffer
	cmd := goexec.Command("vendir", []string{"sync", "-f", vendirFileName}...)
	cmd.Stdin = nil
	cmd.Stderr = &stderrBs
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("running vendir sync failed; check the source configuration and authentication: %w", err)
	}
	return nil
}

func (r VendirRunner) listFiles(dir string) error {
	var stdoutBs, stderrBs bytes.Buffer

	cmd := goexec.Command("ls", []string{"-lR", dir}...)
	cmd.Stdin = nil
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Listing files.\n %s", stderrBs.String())
	}
	r.ui.PrintCmdExecutionOutput(stdoutBs.String())
	return nil
}
