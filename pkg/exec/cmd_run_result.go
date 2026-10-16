// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	trailingSpace = regexp.MustCompile("\\s+\n")
)

// TruncationMarker is prepended to any output field that was clipped.
// Operators can detect truncation by checking for this prefix.
const TruncationMarker = "[output truncated]\n"

type CmdRunResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Error    error
	Finished bool
}

func NewCmdRunResultWithErr(err error) CmdRunResult {
	var result CmdRunResult
	result.AttachErrorf("%s", err)
	return result
}

func (r *CmdRunResult) AttachErrorf(msg string, err error) {
	r.Finished = true
	if err != nil {
		r.ExitCode = -1
		if exitError, ok := err.(interface{ ExitCode() int }); ok {
			r.ExitCode = exitError.ExitCode()
		}

		if err.Error() == "exit status 1" {
			r.Error = fmt.Errorf(msg, "Error (see .status.usefulErrorMessage for details)")
		} else {
			r.Error = fmt.Errorf(msg, err)
		}
	}
}

func (r CmdRunResult) ErrorStr() string {
	if r.Error != nil {
		return r.Error.Error()
	}
	return ""
}

func (r CmdRunResult) WithFriendlyYAMLStrings() CmdRunResult {
	// YAML can format muliline strings nicely
	// if they do not have trailing spaces right before newlines
	return CmdRunResult{
		Stdout:   trailingSpace.ReplaceAllString(strings.TrimSpace(r.Stdout), "\n"),
		Stderr:   trailingSpace.ReplaceAllString(strings.TrimSpace(r.Stderr), "\n"),
		ExitCode: r.ExitCode,
		Error:    r.Error,
		Finished: r.Finished,
	}
}

func (r CmdRunResult) IsEmpty() bool {
	return r == (CmdRunResult{})
}

// WithRedactedSecrets returns a copy of r with Stdout, Stderr, and Error fields
// redacted to mask credential-bearing substrings before truncation or output.
// Redaction must happen before truncation to avoid losing the URL scheme prefix,
// which would prevent pattern matching (e.g., https://user:pass@host becomes
// pass@host after tail-truncation, and the regex won't match without the scheme).
func (r CmdRunResult) WithRedactedSecrets() CmdRunResult {
	return CmdRunResult{
		Stdout:   RedactSecrets(r.Stdout),
		Stderr:   RedactSecrets(r.Stderr),
		ExitCode: r.ExitCode,
		Error:    r.Error,
		Finished: r.Finished,
	}
}

// WithTruncatedStrings returns a copy of r with Stdout and Stderr each capped
// to maxBytes by keeping the tail (the most actionable kapp output always
// appears last).  When a field is under the limit it is left unchanged.
func (r CmdRunResult) WithTruncatedStrings(maxBytes int) CmdRunResult {
	return CmdRunResult{
		Stdout:   truncateOutput(r.Stdout, maxBytes),
		Stderr:   truncateOutput(r.Stderr, maxBytes),
		ExitCode: r.ExitCode,
		Error:    r.Error,
		Finished: r.Finished,
	}
}

// truncateOutput returns s unchanged when len(s) <= maxBytes.  Otherwise it
// returns TruncationMarker followed by the last maxBytes bytes of s, adjusted
// forward to a valid UTF-8 rune boundary.
func truncateOutput(s string, maxBytes int) string {
	if len(s) <= maxBytes {
		return s
	}
	start := len(s) - maxBytes
	// Advance past UTF-8 continuation bytes to avoid splitting a multi-byte rune.
	for start < len(s) && !utf8.RuneStart(s[start]) {
		start++
	}
	return TruncationMarker + s[start:]
}

var (
	// urlUserinfo matches user:password@ and user@ patterns in URLs (http, https, ssh, git schemes)
	urlUserinfo = regexp.MustCompile(`([a-zA-Z][a-zA-Z0-9+.-]*://)[^\s/@]+@|ssh://[^\s/@]+@`)
)

// RedactSecrets masks credential-bearing substrings in subprocess output.
// It redacts URL userinfo (user:password@host or user@host) to prevent credential
// leakage in resource status fields, logs, and debug output.
// For example: https://token:x@registry → https://[REDACTED]@registry
//
// RedactSecrets is idempotent: calling it multiple times produces the same result.
func RedactSecrets(s string) string {
	return urlUserinfo.ReplaceAllString(s, "${1}[REDACTED]@")
}
