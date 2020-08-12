 // Copyright 2020 VMware, Inc.
 // SPDX-License-Identifier: Apache-2.0

package memdir

import (
	"fmt"
	"path/filepath"
	"strings"
)

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

	return newPath, nil
}
