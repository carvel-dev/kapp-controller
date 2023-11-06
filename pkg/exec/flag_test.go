// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package exec_test

import (
	"testing"

	"carvel.dev/kapp-controller/pkg/exec"
)

func TestNewFlagFromString(t *testing.T) {
	invalid := []string{
		"",                  // empty
		"--",                // no name
		"-f",                // missing --
		"flag",              // missing --
		"--flag",            // missing val
		"--flag-flag-flag",  // missing val
		"--flAg",            // uppercase
		"--flag flag=val",   // space
		"--flag --flag=val", // space
	}

	valid := map[string]exec.Flag{
		"--f=v":              exec.Flag{Name: "--f", Value: "v"},
		"--flag=val":         exec.Flag{Name: "--flag", Value: "val"},
		"--flag=val val":     exec.Flag{Name: "--flag", Value: "val val"},
		"--flag=VaL":         exec.Flag{Name: "--flag", Value: "VaL"},
		"--flag-flag-flag=v": exec.Flag{Name: "--flag-flag-flag", Value: "v"},
		"--flag='quotes'":    exec.Flag{Name: "--flag", Value: "'quotes'"},
	}

	for _, one := range invalid {
		_, err := exec.NewFlagFromString(one)
		if err == nil {
			t.Fatalf("Expected error for '%s', but was nil", one)
		}
	}

	for str, expectedFlag := range valid {
		flag, err := exec.NewFlagFromString(str)
		if err != nil {
			t.Fatalf("Did not expect error for '%s', but was error", str)
		}
		if flag != expectedFlag {
			t.Fatalf("Expected flag '%s' to equal: %#v vs %#v", str, flag, expectedFlag)
		}
	}
}
