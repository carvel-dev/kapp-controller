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

	if newPath != path && !strings.HasPrefix(newPath, path+string(filepath.Separator)) {
		return "", fmt.Errorf("Invalid path: %s", subPath)
	}

	return newPath, nil
}
