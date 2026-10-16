// Copyright 2026 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_truncateOutput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxBytes int
		check    func(t *testing.T, result string)
	}{
		{
			name:     "under limit is unchanged",
			input:    "small output",
			maxBytes: 1024,
			check: func(t *testing.T, result string) {
				assert.Equal(t, "small output", result)
			},
		},
		{
			name:     "exactly at limit is unchanged",
			input:    strings.Repeat("x", 100),
			maxBytes: 100,
			check: func(t *testing.T, result string) {
				assert.Equal(t, strings.Repeat("x", 100), result)
			},
		},
		{
			name:     "empty string is unchanged",
			input:    "",
			maxBytes: 100,
			check: func(t *testing.T, result string) {
				assert.Equal(t, "", result)
			},
		},
		{
			name:     "over limit has marker prefix and tail suffix",
			input:    strings.Repeat("a", 200) + strings.Repeat("z", 50),
			maxBytes: 50,
			check: func(t *testing.T, result string) {
				assert.True(t, strings.HasPrefix(result, TruncationMarker))
				assert.True(t, strings.HasSuffix(result, strings.Repeat("z", 50)))
				// Marker does not count toward maxBytes.
				assert.Equal(t, len(TruncationMarker)+50, len(result))
			},
		},
		{
			name: "tail content preserved for realistic kapp output",
			// Many per-resource lines followed by an actionable summary.  maxBytes
			// is set to hold only the last ~100 bytes so truncation is guaranteed,
			// and the summary at the tail must survive.
			input: strings.Repeat("op add configmap/cm-001 (v1) namespace: default\n", 100) +
				"Op: 100 add, 0 delete\nWait to: 100 reconcile\n",
			maxBytes: 100,
			check: func(t *testing.T, result string) {
				assert.True(t, strings.HasPrefix(result, TruncationMarker),
					"large output must start with the truncation marker")
				assert.Contains(t, result, "100 add",
					"the actionable summary at the tail must be preserved after truncation")
			},
		},
		{
			name: "UTF-8 boundary: continuation byte is skipped",
			// s = 98 'a's + "éé"  (102 bytes)
			// maxBytes = 3 → start = 99 → s[99] = 0xA9 (continuation of first 'é')
			// → advance to 100 → s[100] = 0xC3 (leading byte of second 'é') → stop
			// kept = s[100:] = "é"
			input:    strings.Repeat("a", 98) + "éé",
			maxBytes: 3,
			check: func(t *testing.T, result string) {
				require.Equal(t, 102, len(strings.Repeat("a", 98)+"éé"))
				assert.True(t, utf8.ValidString(result), "result must be valid UTF-8")
				assert.Equal(t, TruncationMarker+"é", result)
			},
		},
		{
			name:     "pure multibyte runes produce valid UTF-8",
			input:    strings.Repeat("あ", 100), // 300 bytes; each rune is 3 bytes
			maxBytes: 100,
			check: func(t *testing.T, result string) {
				// 100 / 3 = 33 complete runes (99 bytes); byte 100 is a continuation
				// byte and is skipped.
				assert.True(t, utf8.ValidString(result))
				assert.Equal(t, TruncationMarker+strings.Repeat("あ", 33), result)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.check(t, truncateOutput(tc.input, tc.maxBytes))
		})
	}
}

func Test_WithTruncatedStrings(t *testing.T) {
	const limit = 20
	r := CmdRunResult{
		Stdout:   strings.Repeat("o", 100),
		Stderr:   strings.Repeat("e", 100),
		ExitCode: 1,
		Finished: true,
	}
	got := r.WithTruncatedStrings(limit)

	assert.True(t, strings.HasPrefix(got.Stdout, TruncationMarker))
	assert.True(t, strings.HasPrefix(got.Stderr, TruncationMarker))
	assert.Equal(t, r.ExitCode, got.ExitCode, "ExitCode must be preserved")
	assert.Equal(t, r.Finished, got.Finished, "Finished must be preserved")
	assert.Equal(t, r.Error, got.Error, "Error must be preserved")
}

func TestRedactSecrets(t *testing.T) {
	tests := []struct {
		input       string
		expected    string
		description string
	}{
		{"https://user:token@github.com/repo", "https://[REDACTED]@github.com/repo", "HTTPS URL with user:password"},
		{"http://admin:password@registry:5000/path", "http://[REDACTED]@registry:5000/path", "HTTP registry URL with credentials"},
		{"ssh://git@github.com/org/repo.git", "ssh://[REDACTED]@github.com/org/repo.git", "SSH URL with user"},
		{"git@github.com:org/repo.git", "git@github.com:org/repo.git", "Git SSH shorthand (no scheme)"},
		{"plain error message", "plain error message", "Plain text unchanged"},
		{"https://[REDACTED]@github.com/repo", "https://[REDACTED]@github.com/repo", "Already redacted"},
		{"https://user:token@host1 and https://user2:pass2@host2", "https://[REDACTED]@host1 and https://[REDACTED]@host2", "Multiple URLs"},
		{"GOPROXY=https://token:x@proxy.example.com/path", "GOPROXY=https://[REDACTED]@proxy.example.com/path", "GOPROXY variable"},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			result := RedactSecrets(tt.input)
			if result != tt.expected {
				t.Errorf("RedactSecrets(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}

	t.Run("idempotence", func(t *testing.T) {
		input := "https://user:password@github.com/repo"
		once := RedactSecrets(input)
		twice := RedactSecrets(once)
		if once != twice {
			t.Errorf("RedactSecrets is not idempotent: once=%q, twice=%q", once, twice)
		}
	})
}
