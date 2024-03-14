// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// Taken from https://github.com/linkerd/linkerd2/blob/main/cli/cmd/completion.go

func NewCmdCompletion() *cobra.Command {
	example := `  # bash <= 3.2
  source /dev/stdin <<< "$(kctrl completion bash)"

  # bash >= 4.0
  source <(kctrl completion bash)

  # bash <= 3.2 on osx
  brew install bash-completion # ensure you have bash-completion 1.3+
  kctrl completion bash > $(brew --prefix)/etc/bash_completion.d/kctrl

  # bash >= 4.0 on osx
  brew install bash-completion@2
  kctrl completion bash > $(brew --prefix)/etc/bash_completion.d/kctrl

  # zsh
  source <(kctrl completion zsh)

  # zsh on osx / oh-my-zsh
  kctrl completion zsh > "${fpath[1]}/_kctrl"

  # fish:
  kctrl completion fish | source

  # To load completions for each session, execute once:
  kctrl completion fish > ~/.config/fish/completions/kctrl.fish

  # powershell:
  kctrl completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  kctrl completion powershell > kctrl.ps1
  # and source this file from your powershell profile.
`

	cmd := &cobra.Command{
		Use:       "completion [bash|zsh|fish|powershell]",
		Short:     "Output shell completion code for the specified shell (bash, zsh or fish)",
		Long:      "Output shell completion code for the specified shell (bash, zsh or fish).",
		Example:   example,
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		Args:      cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := getCompletion(args[0], cmd.Parent())
			if err != nil {
				return err
			}

			fmt.Print(out)
			return nil
		},
	}

	return cmd
}

// getCompletion will return the auto completion shell script, if supported
func getCompletion(sh string, parent *cobra.Command) (string, error) {
	var err error
	var buf bytes.Buffer

	switch sh {
	case "bash":
		err = parent.GenBashCompletion(&buf)
	case "zsh":
		err = parent.GenZshCompletion(&buf)
	case "fish":
		err = parent.GenFishCompletion(&buf, true)
	case "powershell":
		err = parent.GenPowerShellCompletion(&buf)
	default:
		err = errors.New("unsupported shell type (must be bash, zsh or fish): " + sh)
	}

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
