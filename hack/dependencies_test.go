// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"testing"
)

func init() {
	client = http.DefaultClient
}

func TestDependencyDownload(t *testing.T) {
	client = http.DefaultClient
	for _, tc := range []struct {
		name string
		dep  dependency
	}{
		{
			name: "with regular files",
			dep: dependency{
				Name:        "test",
				Repo:        "benmoss/test-resources",
				Version:     "v1.0.0",
				URLTemplate: "https://github.com/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}-{{.Version}}-{{.OS}}-{{.Arch}}",
				Checksums: map[string]map[string]checksum{
					"darwin": {
						"arm64": "d01cccabb10342523dfb4dbacc27a85df54c36873fa9900b194a27d40985005e",
						"amd64": "d51945f0bca8e1b54025a8a18ffebf885edd09a8731a8955100a9b0f03dbd4c0",
					},
					"linux": {
						"arm64": "4968af2083a16b93b2ef5dfb35255fc5591d731260c5477e31e595827aac6bba",
						"amd64": "7ab37ac20f25d25e0522cfb95629e9c0051a71aa4cb489e375e166137a215f3a",
					},
				},
			},
		},
		{
			name: "with tgz files",
			dep: dependency{
				Name:        "test",
				Repo:        "benmoss/test-resources",
				Version:     "v1.0.0",
				URLTemplate: "https://github.com/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}-{{.Version}}-{{.OS}}-{{.Arch}}.tgz",
				Checksums: map[string]map[string]checksum{
					"darwin": {
						"arm64": "06fd0371bdcb7880805d3bbcffce1f5a7f81f15e3aae5253b12264a9b54feb79",
						"amd64": "33eb45047a683824571250728a022769cd505b21473fa86e302f5a79b7acf94d",
					},
					"linux": {
						"arm64": "c40091f96a3c6016d6e7c9e2160df906ca92f928faecd0ef4fdf399dbc5fe8e0",
						"amd64": "794602e675b88a814d5599018b14f0a31c623e58a3f98bf42e0413faec26242b",
					},
				},
				TarballSubpath: stringPtr("{{.Name}}-{{.Version}}-{{.OS}}-{{.Arch}}"),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dep := tc.dep
			for os, arches := range dep.Checksums {
				for arch, expectedChecksum := range arches {
					platform := platform{OS: os, Arch: arch}
					blob, err := dep.download(context.Background(), platform)
					if err != nil {
						t.Fatal(err)
					}
					actualChecksum, err := blob.Checksum()
					if err != nil {
						t.Fatal(err)
					}
					if actualChecksum != expectedChecksum {
						t.Errorf("%s/%s checksum mismatch, expected %q, got %q", os, arch, expectedChecksum, actualChecksum)
					}
					reader, err := blob.Binary()
					if err != nil {
						t.Fatal(err)
					}
					bs, err := io.ReadAll(reader)
					if err != nil {
						t.Fatal(err)
					}
					expectedContent := fmt.Sprintf("%s-%s-%s-%s\n", dep.Name, dep.Version, os, arch)
					actualContent := string(bs)
					if expectedContent != actualContent {
						t.Errorf("%s/%s content mismatch, expected %q, got %q", os, arch, expectedContent, actualContent)
					}
				}
			}
		})
	}
}

func TestDependencyUpdate(t *testing.T) {
	dep := dependency{
		Name:        "test",
		Repo:        "benmoss/test-resources",
		Version:     "v1.0.0",
		URLTemplate: "https://github.com/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}-{{.Version}}-{{.OS}}-{{.Arch}}",
		Checksums: map[string]map[string]checksum{
			"darwin": {
				"arm64": "d01cccabb10342523dfb4dbacc27a85df54c36873fa9900b194a27d40985005e",
				"amd64": "d51945f0bca8e1b54025a8a18ffebf885edd09a8731a8955100a9b0f03dbd4c0",
			},
			"linux": {
				"arm64": "4968af2083a16b93b2ef5dfb35255fc5591d731260c5477e31e595827aac6bba",
				"amd64": "7ab37ac20f25d25e0522cfb95629e9c0051a71aa4cb489e375e166137a215f3a",
			},
		},
	}
	if err := dep.update(context.Background()); err != nil {
		t.Fatal(err)
	}
	expectedChecksums := map[string]map[string]checksum{
		"darwin": {
			"arm64": "1c0ab099d3ffe0986a680faa21df05d38b849dfc705b04017697c3e9a458725b",
			"amd64": "307fe95a5ff2572debd7c6aac7c19023df8815088451282b1a3df78f39d6a16b",
		},
		"linux": {
			"arm64": "ca0a4640f7dd2e8c8044b936bd18ceede178c4f3a716b8984f6fe9082ea598b8",
			"amd64": "5ce8948a28c332f79fad3f6818f10b15d674ff3a7d2fb9772f8bc0195ea55907",
		},
	}
	for os, arches := range dep.Checksums {
		for arch, actualChecksum := range arches {
			expectedChecksum := expectedChecksums[os][arch]
			if expectedChecksum != actualChecksum {
				t.Errorf("%s/%s checksum mismatch, expected %q, got %q", os, arch, expectedChecksum, actualChecksum)
			}
		}
	}
	expectedVersion := "v1.0.1"
	actualVersion := dep.Version
	if expectedVersion != actualVersion {
		t.Errorf("version mismatch, expected %q, got %q", expectedVersion, actualVersion)
	}
}

func TestDownloadAndVerify(t *testing.T) {
	dep := dependency{
		Name:        "test",
		Repo:        "benmoss/test-resources",
		Version:     "v1.0.0",
		URLTemplate: "https://github.com/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}-{{.Version}}-{{.OS}}-{{.Arch}}",
		Checksums: map[string]map[string]checksum{
			"darwin": {
				"amd64": "d51945f0bca8e1b54025a8a18ffebf885edd09a8731a8955100a9b0f03dbd4c0",
			},
		},
	}
	tmpDir, err := os.MkdirTemp("", "hack-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	install := installCommand{
		os:      "darwin",
		arch:    "amd64",
		destDir: tmpDir,
	}
	if err := install.downloadAndVerify(context.Background(), &dep); err != nil {
		t.Fatal(err)
	}
	fileInfo, err := os.Stat(path.Join(tmpDir, dep.Name))
	if err != nil {
		t.Fatal(err)
	}
	if fileInfo.Mode().Perm() != 0777 {
		t.Fatalf("expected file to be executable, got mode %s", fileInfo.Mode().Perm())
	}
}

func stringPtr(str string) *string {
	return &str
}
