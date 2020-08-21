package template

import (
	"bytes"
	"io"
	"os"
	goexec "os/exec"
	"path/filepath"

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

	files, err := walkMatch(dirPath, t.opts.Match)
	result := exec.CmdRunResult{}

	if err != nil {

		result.AttachErrorf("Templating dir: %s", err)
		return result
	}

	readers := make([]io.Reader, len(files))

	for i, f := range files {

		file, err := os.Open(f) // just pass the file name
		if err != nil {
			result.AttachErrorf("Templating dir: %s", err)
			return result
		}

		readers[i] = file

	}
	if t.opts.MergeFiles == false {
		for _, r := range readers {

			templateResult := t.template(r)
			result.Stdout = result.Stdout + "\n" + templateResult.Stdout
			result.Stderr = result.Stderr + "\n" + templateResult.Stderr

		}
	} else {
		r := io.MultiReader(readers...)

		result = t.template(r)
	}

	return result
}

func (t *Sops) TemplateStream(input io.Reader) exec.CmdRunResult {
	return t.template(input)
}

func (t *Sops) template(input io.Reader) exec.CmdRunResult {

	var stdoutBs, stderrBs bytes.Buffer

	args := t.addArgs(t.opts.Args)
	cmd := goexec.Command("sops", args...)
	cmd.Stdin = input
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs

	err := cmd.Run()

	result := exec.CmdRunResult{
		Stdout: stdoutBs.String(),
		Stderr: stderrBs.String(),
	}
	result.AttachErrorf("Templating dir: %s", err)

	return result
}

func (t *Sops) addArgs(args []string) []string {
	args = append(args, "-d", "/dev/stdin")
	return args
}
