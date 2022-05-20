// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/sidecarexec"
)

func main() {
	localCmdRunner := exec.NewPlainCmdRunner()
	server := sidecarexec.NewServer(localCmdRunner, sidecarexec.ServerOpts{
		AllowedCmdNames: []string{
			// Fetch (calls impgkg and others internally)
			"vendir",
			// Template
			"ytt", "kbld", "sops", "helm", "cue",
		},
	})

	err := server.Serve()
	if err != nil {
		log.Fatal("Serving:", err)
	}
}
