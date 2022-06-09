// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
)

func TestValues(t *testing.T) {
	subject := Values{
		ValuesFrom:  nil,
		genericOpts: GenericOpts{},
		coreClient:  nil,
	}

	t.Run("Downward API values", func(t *testing.T) {
		subject := subject
		subject.genericOpts.Metadata = &PartialObjectMetadata{
			ObjectMeta: ObjectMeta{
				Name:        "some-name",
				Namespace:   "some-namespace",
				UID:         "some-uid",
				Labels:      map[string]string{"a_label": "a_label_val"},
				Annotations: map[string]string{"a_ann": "a_ann_val"},
			},
		}

		subject.ValuesFrom = []v1alpha1.AppTemplateValuesSource{{DownwardAPI: &v1alpha1.AppTemplateValuesDownwardAPI{
			Items: []v1alpha1.DownwardAPIAppTemplateValues{
				{Name: "a_some-annotation-key", FieldPath: "metadata.annotations"},
				{Name: "b_some-name-key", FieldPath: "metadata.name"},
				{Name: "c_some-namespace-key", FieldPath: "metadata.namespace"},
				{Name: "d_some-uid-key", FieldPath: "metadata.uid"},
				{Name: "e_some-label-key", FieldPath: "metadata.labels"},
			}},
		}}

		paths, cleanup, err := subject.AsPaths(os.TempDir())
		require.NoError(t, err)
		t.Cleanup(cleanup)

		require.Len(t, paths, 5)
		expectedValues := []string{
			"a_some-annotation-key: \n  a_ann: \"a_ann_val\"",
			"b_some-name-key: \"some-name\"",
			"c_some-namespace-key: \"some-namespace\"",
			"d_some-uid-key: \"some-uid\"",
			"e_some-label-key: \n  a_label: \"a_label_val\"",
		}

		for i, p := range paths {
			assertFileContents(t, p, expectedValues[i])
		}

		t.Run("map field paths should allow subpaths", func(t *testing.T) {
			subject := subject
			subject.ValuesFrom = []v1alpha1.AppTemplateValuesSource{{DownwardAPI: &v1alpha1.AppTemplateValuesDownwardAPI{
				Items: []v1alpha1.DownwardAPIAppTemplateValues{
					{Name: "a_some-annotation-key", FieldPath: "metadata.annotations['a_ann']"},
					{Name: "b_some-label-key", FieldPath: "metadata.labels['a_label']"},
				}},
			}}

			paths, cleanup, err := subject.AsPaths(os.TempDir())
			require.NoError(t, err)
			require.Len(t, paths, 2)
			t.Cleanup(cleanup)

			expectedValues := []string{
				"a_some-annotation-key: \"a_ann_val\"",
				"b_some-label-key: \"a_label_val\"",
			}

			for i, p := range paths {
				assertFileContents(t, p, expectedValues[i])
			}

		})

		t.Run("file names have encoding suffix due to templating engines such as cue requirement", func(t *testing.T) {
			subject := subject
			paths, cleanup, err := subject.AsPaths(os.TempDir())
			require.NoError(t, err)
			require.Len(t, paths, 5)
			t.Cleanup(cleanup)

			for _, path := range paths {
				assert.True(t, strings.HasSuffix(path, ".yaml"))
			}
		})

		t.Run("map field paths should return in stable order", func(t *testing.T) {
			subject := subject
			subject.genericOpts.Metadata.Annotations = map[string]string{
				"z_ann": "z_ann_val",
				"s_ann": "s_ann_val",
				"a_ann": "a_ann_val",
			}
			subject.ValuesFrom = []v1alpha1.AppTemplateValuesSource{{DownwardAPI: &v1alpha1.AppTemplateValuesDownwardAPI{
				Items: []v1alpha1.DownwardAPIAppTemplateValues{
					{Name: "some-annotation-key", FieldPath: "metadata.annotations"},
				}},
			}}

			paths, cleanup, err := subject.AsPaths(os.TempDir())
			require.NoError(t, err)
			require.Len(t, paths, 1)
			t.Cleanup(cleanup)

			assertFileContents(t, paths[0], "some-annotation-key: \n  a_ann: \"a_ann_val\"\n  s_ann: \"s_ann_val\"\n  z_ann: \"z_ann_val\"")
		})

		t.Run("items with same key name, latter keys specified should clobber the earlier ones", func(t *testing.T) {
			subject := subject
			subject.ValuesFrom = []v1alpha1.AppTemplateValuesSource{{DownwardAPI: &v1alpha1.AppTemplateValuesDownwardAPI{
				Items: []v1alpha1.DownwardAPIAppTemplateValues{
					{Name: "a_name", FieldPath: "metadata.name"},
					{Name: "a_name", FieldPath: "metadata.namespace"},
				}},
			}}

			paths, cleanup, err := subject.AsPaths(os.TempDir())
			require.NoError(t, err)
			require.Len(t, paths, 2)
			t.Cleanup(cleanup)

			assertFileContents(t, paths[1], "a_name: \"some-namespace\"")
		})

		t.Run("return helpful error if subpath is not found", func(t *testing.T) {
			subject := subject
			subject.ValuesFrom = []v1alpha1.AppTemplateValuesSource{{DownwardAPI: &v1alpha1.AppTemplateValuesDownwardAPI{
				Items: []v1alpha1.DownwardAPIAppTemplateValues{
					{Name: "a_some-annotation-key", FieldPath: "metadata.annotations['INVALID_SUBPATH']"},
				}},
			}}

			_, _, err := subject.AsPaths(os.TempDir())
			require.Error(t, err)
			assert.ErrorContains(t, err, "Writing paths: INVALID_SUBPATH is not found")
		})

		t.Run("return helpful error if an unsupported downward api field spec is provided", func(t *testing.T) {
			subject := subject
			subject.ValuesFrom = []v1alpha1.AppTemplateValuesSource{{DownwardAPI: &v1alpha1.AppTemplateValuesDownwardAPI{
				Items: []v1alpha1.DownwardAPIAppTemplateValues{
					{Name: "unsupportedFieldPath", FieldPath: "metadata.creationTimestamp"},
				}},
			}}

			_, _, err := subject.AsPaths(os.TempDir())
			require.Error(t, err)
			assert.ErrorContains(t, err, "Writing paths: creationTimestamp is not found")
		})

		t.Run("return helpful error if multiple field spec is provided", func(t *testing.T) {
			subject := subject
			subject.ValuesFrom = []v1alpha1.AppTemplateValuesSource{{DownwardAPI: &v1alpha1.AppTemplateValuesDownwardAPI{
				Items: []v1alpha1.DownwardAPIAppTemplateValues{
					{Name: "a_some-annotation-key", FieldPath: "metadata['name', 'uid']"},
				}},
			}}

			_, _, err := subject.AsPaths(os.TempDir())
			require.Error(t, err)
			assert.ErrorContains(t, err, "Writing paths: invalid field spec provided to DownwardAPI. Only single supported fields are allowed")
		})

		//TODO: ensure only metadata.name namespace, uid etc are allowed
		//TODO: test relaxed fieldspec paths
	})
}

func assertFileContents(t *testing.T, path string, expectedVal string) {
	valueContents, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, expectedVal, string(valueContents))
}