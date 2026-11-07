// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package release

import (
	"fmt"
	"testing"

	cmdcore "carvel.dev/kapp-controller/cli/pkg/kctrl/cmd/core"
	"carvel.dev/kapp-controller/cli/pkg/kctrl/cmd/package/repository/release/build"
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	namePrompt     = "Enter the package repository name"
	registryPrompt = "Enter the registry url"
)

// fakeAuthoringUI answers prompts by label and records the defaults that were offered
type fakeAuthoringUI struct {
	cmdcore.AuthoringUI
	answers map[string]string
	asked   map[string]string
}

func (*fakeAuthoringUI) PrintHeaderText(string)        {}
func (*fakeAuthoringUI) PrintInformationalText(string) {}

func (u *fakeAuthoringUI) AskForText(opts ui.TextOpts) (string, error) {
	u.asked[opts.Label] = opts.Default
	answer, found := u.answers[opts.Label]
	if !found {
		return "", fmt.Errorf("unexpected prompt: %s", opts.Label)
	}
	return answer, nil
}

func newPkgRepoBuild(name, image string) *build.PackageRepoBuild {
	return &build.PackageRepoBuild{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec: build.PackageRepoBuildSpec{
			Export: &build.PackageRepoBuildExport{
				ImgpkgBundle: &build.PackageRepoBuildExportImgpkgBundle{Image: image},
			},
		},
	}
}

func TestConfigureBasicInformationUsesFlagsWithoutPrompting(t *testing.T) {
	fakeUI := &fakeAuthoringUI{asked: map[string]string{}}
	o := &ReleaseOptions{ui: fakeUI, pkgRepoName: "flag.carvel.dev", registryURL: " registry.example.com/flag "}
	pkgRepoBuild := newPkgRepoBuild("file.carvel.dev", "registry.example.com/file")

	require.NoError(t, o.configureBasicInformation(pkgRepoBuild))

	require.Empty(t, fakeUI.asked)
	require.Equal(t, "flag.carvel.dev", pkgRepoBuild.Name)
	require.Equal(t, "registry.example.com/flag", pkgRepoBuild.Spec.Export.ImgpkgBundle.Image)
}

func TestConfigureBasicInformationPromptsForValuesWithoutFlags(t *testing.T) {
	fakeUI := &fakeAuthoringUI{
		asked:   map[string]string{},
		answers: map[string]string{namePrompt: "answer.carvel.dev", registryPrompt: "registry.example.com/answer"},
	}
	o := &ReleaseOptions{ui: fakeUI}
	pkgRepoBuild := newPkgRepoBuild("file.carvel.dev", "registry.example.com/file")

	require.NoError(t, o.configureBasicInformation(pkgRepoBuild))

	require.Equal(t, map[string]string{namePrompt: "file.carvel.dev", registryPrompt: "registry.example.com/file"}, fakeUI.asked)
	require.Equal(t, "answer.carvel.dev", pkgRepoBuild.Name)
	require.Equal(t, "registry.example.com/answer", pkgRepoBuild.Spec.Export.ImgpkgBundle.Image)
}

func TestConfigureBasicInformationPromptsOnlyForMissingFlag(t *testing.T) {
	fakeUI := &fakeAuthoringUI{asked: map[string]string{}, answers: map[string]string{namePrompt: DefaultPkgRepoName}}
	o := &ReleaseOptions{ui: fakeUI, registryURL: "registry.example.com/flag"}
	pkgRepoBuild := newPkgRepoBuild("", "")

	require.NoError(t, o.configureBasicInformation(pkgRepoBuild))

	require.Equal(t, map[string]string{namePrompt: DefaultPkgRepoName}, fakeUI.asked)
	require.Equal(t, DefaultPkgRepoName, pkgRepoBuild.Name)
	require.Equal(t, "registry.example.com/flag", pkgRepoBuild.Spec.Export.ImgpkgBundle.Image)
}

func TestConfigureBasicInformationRequiresRegistryURL(t *testing.T) {
	fakeUI := &fakeAuthoringUI{
		asked:   map[string]string{},
		answers: map[string]string{registryPrompt: "  "},
	}
	o := &ReleaseOptions{ui: fakeUI, pkgRepoName: "flag.carvel.dev"}

	err := o.configureBasicInformation(newPkgRepoBuild("", ""))

	require.ErrorContains(t, err, "expected registry URL to be specified")
}
