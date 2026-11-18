// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package memdir_test

import (
	"os"
	"path/filepath"
	"testing"

	"carvel.dev/kapp-controller/pkg/memdir"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ScopedPath(t *testing.T) {
	baseDir := t.TempDir()
	outsideDir := t.TempDir()

	require.NoError(t, os.MkdirAll(filepath.Join(baseDir, "subdir"), 0700))
	require.NoError(t, os.WriteFile(filepath.Join(baseDir, "subdir", "file.yml"), []byte("a: 1"), 0600))
	require.NoError(t, os.WriteFile(filepath.Join(outsideDir, "file.yml"), []byte("a: 1"), 0600))
	// symlinks within baseDir, relative and absolute
	require.NoError(t, os.Symlink("subdir", filepath.Join(baseDir, "link-to-subdir")))
	require.NoError(t, os.Symlink(filepath.Join(baseDir, "subdir", "file.yml"), filepath.Join(baseDir, "link-to-file.yml")))
	// symlinks pointing outside of baseDir
	require.NoError(t, os.Symlink(outsideDir, filepath.Join(baseDir, "link-to-outside")))
	require.NoError(t, os.Symlink(filepath.Join("..", filepath.Base(outsideDir)), filepath.Join(baseDir, "relative-link-to-outside")))
	require.NoError(t, os.Symlink(filepath.Join(outsideDir, "file.yml"), filepath.Join(baseDir, "link-to-outside-file.yml")))

	testCases := []struct {
		subPath     string
		expectedErr string
	}{
		{subPath: "."},
		{subPath: "subdir"},
		{subPath: "subdir/file.yml"},
		{subPath: "subdir/does-not-exist-yet.yml"},
		{subPath: "new-dir/new-file.yml"},
		{subPath: "link-to-subdir/file.yml"},
		{subPath: "link-to-subdir/does-not-exist-yet.yml"},
		{subPath: "link-to-file.yml"},
		{subPath: "../file.yml", expectedErr: "Invalid path: ../file.yml"},
		{subPath: "link-to-outside", expectedErr: "resolves outside of the base directory through a symlink"},
		{subPath: "link-to-outside/file.yml", expectedErr: "resolves outside of the base directory through a symlink"},
		{subPath: "link-to-outside/does-not-exist-yet.yml", expectedErr: "resolves outside of the base directory through a symlink"},
		{subPath: "relative-link-to-outside/file.yml", expectedErr: "resolves outside of the base directory through a symlink"},
		{subPath: "link-to-outside-file.yml", expectedErr: "resolves outside of the base directory through a symlink"},
	}

	for _, tc := range testCases {
		t.Run(tc.subPath, func(t *testing.T) {
			result, err := memdir.ScopedPath(baseDir, tc.subPath)
			if tc.expectedErr != "" {
				assert.ErrorContains(t, err, tc.expectedErr)
				return
			}
			require.NoError(t, err)
			// returned path is the lexical path, symlinks are not resolved
			assert.Equal(t, filepath.Join(baseDir, tc.subPath), result)
		})
	}
}

func Test_ScopedPath_BaseDirIsSymlink(t *testing.T) {
	rootDir := t.TempDir()
	actualBaseDir := filepath.Join(rootDir, "actual-base")
	linkedBaseDir := filepath.Join(rootDir, "linked-base")
	require.NoError(t, os.Mkdir(actualBaseDir, 0700))
	require.NoError(t, os.Symlink(actualBaseDir, linkedBaseDir))
	require.NoError(t, os.WriteFile(filepath.Join(actualBaseDir, "file.yml"), []byte("a: 1"), 0600))

	for _, subPath := range []string{"file.yml", "does-not-exist-yet.yml", "new-dir/new-file.yml"} {
		result, err := memdir.ScopedPath(linkedBaseDir, subPath)
		require.NoError(t, err, "subPath %s", subPath)
		assert.Equal(t, filepath.Join(linkedBaseDir, subPath), result)
	}
}
