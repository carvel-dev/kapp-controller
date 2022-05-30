package build

import (
	"fmt"
	"os"
)

func validatePathExists(path string) (bool, string, error) {
	if len(path) == 0 {
		return false, "Path cannot be empty", nil
	}
	_, err := os.Stat(path)
	if err == nil {
		return true, "", nil
	} else if os.IsNotExist(err) {
		return false, fmt.Sprintf("Location %s doesn't exist", path), nil
	} else {
		return false, fmt.Sprintf("Invalid Location. Failed while checking for existence of location %s.Error is: %s", path, err.Error()), err
	}
	return true, "", nil
}
