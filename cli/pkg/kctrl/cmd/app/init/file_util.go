// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package init

import (
	"fmt"
	"os"
)

// IsFileExists checks if file exists
func IsFileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, fmt.Errorf("Checking file: %s", err.Error())
	}
}

// WriteFile writes binary content to file
func WriteFile(filePath string, data []byte) error {
	// Create creates or truncates the named file. If the file already exists, it is truncated.
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}
