// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package local

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	vendirconf "carvel.dev/vendir/pkg/vendir/config"
)

type localVendirConf struct {
	// Indexed by numeric fetch index
	localPaths map[int]string
}

func newLocalVendirConf(resAnnotations map[string]string) (localVendirConf, error) {
	cwdPath, err := os.Getwd()
	if err != nil {
		return localVendirConf{}, err
	}

	const (
		prefix = "kctrl.carvel.dev/local-fetch-"
	)

	localPaths := map[int]string{}

	for key, val := range resAnnotations {
		if strings.HasPrefix(key, prefix) {
			fetchIdx, err := strconv.Atoi(strings.TrimPrefix(key, prefix))
			if err != nil {
				return localVendirConf{}, err
			}
			localPaths[fetchIdx] = filepath.Join(cwdPath, val)
		}
	}

	return localVendirConf{localPaths}, nil
}

func (c localVendirConf) Adjust(conf vendirconf.Config) vendirconf.Config {
	for fetchIdx, localPath := range c.localPaths {
		if fetchIdx >= len(conf.Directories) {
			// Ignore invalid indexes
			continue
		}
		conf.Directories[fetchIdx].Contents[0] = vendirconf.DirectoryContents{
			Path: conf.Directories[fetchIdx].Contents[0].Path,
			Directory: &vendirconf.DirectoryContentsDirectory{
				Path: localPath,
			},
		}
	}
	return conf
}
