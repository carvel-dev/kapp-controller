package template

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	goexec "os/exec"
	"path/filepath"
	"strings"

	"github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/k14s/kapp-controller/pkg/exec"
)

type Sops struct {
	opts        v1alpha1.AppTemplateSops
	genericOpts GenericOpts
}

var _ Template = &Sops{}

func walkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func NewSops(opts v1alpha1.AppTemplateSops, genericOpts GenericOpts) *Sops {

	return &Sops{opts, genericOpts}
}

func (t *Sops) TemplateDir(dirPath string) exec.CmdRunResult {

	files, err := walkMatch(dirPath, "*")
	result := exec.CmdRunResult{}

	if err != nil {

		result.AttachErrorf("Templating dir: %s", err)
		return result
	}

	for _, file := range files {

		templateResult := t.template(file)
		content, err := ioutil.ReadFile(file) // just pass the file name
		if err != nil {
			result.AttachErrorf("Templating dir: %s", err)
			return result
		}

		result.Stdout = fmt.Sprintf("%s\n---\n%s", result.Stdout, content)
		result.Stderr = fmt.Sprintf("%s \n file: %s \n %s", result.Stderr, file, templateResult.Stderr)

	}

	return result
}

func (t *Sops) TemplateStream(_ io.Reader) exec.CmdRunResult {
	return exec.NewCmdRunResultWithErr(fmt.Errorf("Templating data is not supported"))
}

func (t *Sops) template(input string) exec.CmdRunResult {

	var stdoutBs, stderrBs bytes.Buffer

	args := t.addArgs(input, t.opts)
	cmd := goexec.Command("sops", args...)
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs

	err := cmd.Run()

	errStr := stderrBs.String()
	if strings.Contains(errStr, "sops metadata not found") {
		errStr = ""
	}

	result := exec.CmdRunResult{
		Stdout: stdoutBs.String(),
		Stderr: errStr,
	}
	result.AttachErrorf("Templating dir: %s", err)

	return result
}

func (t *Sops) addArgs(inputFile string, inputArgs v1alpha1.AppTemplateSops) []string {
	var args = []string{}

	if inputArgs.IgnoreMac {
		args = append(args, "--ignore-mac")
	}
	if len(inputArgs.PGP) > 0 {
		args = append(args, "--pgp", strings.Join(inputArgs.PGP, ","))
	}
	if len(inputArgs.AWSProifle) > 0 {
		args = append(args, "--aws-profile", inputArgs.AWSProifle)
	}
	if len(inputArgs.KMSKeys) > 0 {
		args = append(args, "--kms", strings.Join(inputArgs.KMSKeys, ","))
	}
	if len(inputArgs.GCPKms) > 0 {
		args = append(args, "--gcp-kms", strings.Join(inputArgs.GCPKms, ","))
	}
	if len(inputArgs.AzureKV) > 0 {
		args = append(args, "--azure-kv", strings.Join(inputArgs.AzureKV, ","))
	}
	args = append(args, "--input-type=yaml", "--output-type=yaml", "-d", "-i", inputFile)
	return args
}
