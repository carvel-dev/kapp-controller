package build

import (
	"fmt"
	"os"
	"strings"
)

func validatePathExists(allPaths string) (bool, string, error) {
	if len(allPaths) == 0 {
		return false, "Path cannot be empty", nil
	}
	paths := strings.Split(allPaths, ",")
	for _, path := range paths {
		path = strings.TrimSpace(path)
		_, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				return false, fmt.Sprintf("Location %s doesn't exist", path), nil
			} else {
				return false, fmt.Sprintf("Invalid Location. Failed while checking for existence of location %s.Error is: %s", path, err.Error()), err
			}
		}
	}

	return true, "", nil
}
