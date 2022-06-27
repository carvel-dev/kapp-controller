// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"fmt"
	"os"
)

// Check if file exists
func IsFileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, fmt.Errorf("failed to check for the existence of file. Error is: %s", err.Error())
	}
}

// Write binary content to file
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
