// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package memdir

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ScopedPath joins path with subpath, returning the result if and only if it is still within the original path.
// (Security measure to prevent '..' shenanigans)
func ScopedPath(path, subPath string) (string, error) {
	newPath, err := filepath.Abs(filepath.Join(path, subPath))
	if err != nil {
		return "", fmt.Errorf("Abs path: %s", err)
	}

	// Check that subPath is contained within path (disallow this scenario):
	//   ScopedPath("/root", "../root-trick/file1")
	//   "/root-trick/file1" == "/root" + "../root-trick/file1"
	if newPath != path && !strings.HasPrefix(newPath, path+string(filepath.Separator)) {
		return "", fmt.Errorf("Invalid path: %s", subPath)
	}

	// The check above is lexical only. Also check that newPath is contained
	// within path once symlinks are resolved, since a symlink component in
	// subPath may point outside of path.
	resolvedPath, err := resolveExistingSymlinks(path)
	if err != nil {
		return "", fmt.Errorf("resolving symlinks: %w", err)
	}

	resolvedNewPath, err := resolveExistingSymlinks(newPath)
	if err != nil {
		return "", fmt.Errorf("resolving symlinks: %w", err)
	}

	if resolvedNewPath != resolvedPath && !strings.HasPrefix(resolvedNewPath, resolvedPath+string(filepath.Separator)) {
		return "", fmt.Errorf("invalid path: %s (resolves outside of the base directory through a symlink)", subPath)
	}

	return newPath, nil
}

// resolveExistingSymlinks resolves symlinks in the deepest existing ancestor
// of path and joins the remaining (not yet existing) components unchanged,
// so that paths that are about to be written can be checked as well.
func resolveExistingSymlinks(path string) (string, error) {
	if _, err := os.Lstat(path); err == nil {
		return filepath.EvalSymlinks(path)
	}

	parent := filepath.Dir(path)
	if parent == path {
		return path, nil
	}

	resolvedParent, err := resolveExistingSymlinks(parent)
	if err != nil {
		return "", err
	}

	return filepath.Join(resolvedParent, filepath.Base(path)), nil
}
