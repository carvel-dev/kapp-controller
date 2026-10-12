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

func Test_RedactSecrets(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no credentials, unchanged",
			input:    "This is normal output",
			expected: "This is normal output",
		},
		{
			name:     "empty string unchanged",
			input:    "",
			expected: "",
		},
		{
			name:     "simple https URL with basic auth",
			input:    "https://user:password@example.com/path",
			expected: "https://[REDACTED]@example.com/path",
		},
		{
			name:     "simple http URL with basic auth",
			input:    "http://user:token@registry.io",
			expected: "http://[REDACTED]@registry.io",
		},
		{
			name:     "GOPROXY with auth",
			input:    "GOPROXY=https://user:token@proxy.example.com",
			expected: "GOPROXY=https://[REDACTED]@proxy.example.com",
		},
		{
			name:     "git URL with credentials",
			input:    "git clone https://myuser:mypass@github.com/repo.git",
			expected: "git clone https://[REDACTED]@github.com/repo.git",
		},
		{
			name:     "password containing @ symbol",
			input:    "https://user:pass@word@example.com/path",
			expected: "https://[REDACTED]@example.com/path",
		},
		{
			name:     "multiple URLs in output, all redacted",
			input:    "fetching https://user:pass1@host1.com and https://user2:pass2@host2.com",
			expected: "fetching https://[REDACTED]@host1.com and https://[REDACTED]@host2.com",
		},
		{
			name:     "URL with only username (no password)",
			input:    "https://user@example.com/path",
			expected: "https://[REDACTED]@example.com/path",
		},
		{
			name:     "multiline output with credentials",
			input:    "Error fetching\nhttps://user:password@example.com\nConnection failed",
			expected: "Error fetching\nhttps://[REDACTED]@example.com\nConnection failed",
		},
		{
			name:     "URL with special characters in password",
			input:    "https://user:p%40ss!word#$@example.com",
			expected: "https://[REDACTED]@example.com",
		},
		{
			name:     "ftp scheme also supported",
			input:    "ftp://admin:secret@ftp.example.com/files",
			expected: "ftp://[REDACTED]@ftp.example.com/files",
		},
		{
			name:     "custom scheme with +/. characters",
			input:    "custom+scheme://user:pass@host.com",
			expected: "custom+scheme://[REDACTED]@host.com",
		},
		{
			name:     "URL-like text without scheme is not redacted",
			input:    "user:password@example.com",
			expected: "user:password@example.com",
		},
		{
			name:     "text containing @ but not a URL",
			input:    "email: test@example.com not redacted",
			expected: "email: test@example.com not redacted",
		},
		{
			name:     "whitespace stops URL matching",
			input:    "https://user:pass@host.com and some other text",
			expected: "https://[REDACTED]@host.com and some other text",
		},
		{
			name:     "ssh URL with user",
			input:    "ssh://user@example.com/path",
			expected: "ssh://[REDACTED]@example.com/path",
		},
		{
			name:     "ssh git URL (common git over SSH)",
			input:    "ssh://git@github.com/carvel-dev/kapp-controller.git",
			expected: "ssh://[REDACTED]@github.com/carvel-dev/kapp-controller.git",
		},
		{
			name:     "ssh URL with port",
			input:    "Cloning ssh://user:token@git.example.com:2222/repo.git",
			expected: "Cloning ssh://[REDACTED]@git.example.com:2222/repo.git",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := RedactSecrets(tc.input)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func Test_RedactSecrets_Idempotence(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "simple https URL",
			input: "https://user:password@example.com/path",
		},
		{
			name:  "multiple URLs",
			input: "fetching https://user:pass1@host1.com and https://user2:pass2@host2.com",
		},
		{
			name:  "ssh git URL",
			input: "Cloning ssh://git@github.com/repo.git",
		},
		{
			name:  "GOPROXY with auth",
			input: "GOPROXY=https://user:token@proxy.example.com",
		},
		{
			name:  "password with @ symbol",
			input: "https://user:pass@word@example.com/path",
		},
		{
			name:  "no credentials",
			input: "This is normal output",
		},
		{
			name:  "empty string",
			input: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// First redaction
			redacted := RedactSecrets(tc.input)
			// Second redaction should be identical (idempotent)
			redactedTwice := RedactSecrets(redacted)
			assert.Equal(t, redacted, redactedTwice,
				"RedactSecrets should be idempotent: calling it twice should produce the same result")
		})
	}
}
