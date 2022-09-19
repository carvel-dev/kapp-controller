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
		ValuesFrom: nil,
		appContext: AppContext{},
		coreClient: nil,
	}

	t.Run("Downward API values", func(t *testing.T) {
		subject := subject
		subject.appContext.Metadata = PartialObjectMetadata{
			ObjectMeta: ObjectMeta{
				Name:        "some-name",
				Namespace:   "some-namespace",
				UID:         "some-uid",
				Labels:      map[string]string{"a_label": "a_label_val"},
				Annotations: map[string]string{"a_ann": "a_ann_val"},
			},
		}

		subject.ValuesFrom = []v1alpha1.AppTemplateValuesSource{{DownwardAPI: &v1alpha1.AppTemplateValuesDownwardAPI{
			Items: []v1alpha1.AppTemplateValuesDownwardAPIItem{
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
			"a_some-annotation-key:\n  a_ann: a_ann_val\n",
			"b_some-name-key: some-name\n",
			"c_some-namespace-key: some-namespace\n",
			"d_some-uid-key: some-uid\n",
			"e_some-label-key:\n  a_label: a_label_val\n",
		}

		for i, p := range paths {
			assertFileContents(t, p, expectedValues[i])
		}

		t.Run("name should allow nested key structure", func(t *testing.T) {
			subject := subject
			subject.ValuesFrom = []v1alpha1.AppTemplateValuesSource{{DownwardAPI: &v1alpha1.AppTemplateValuesDownwardAPI{
				Items: []v1alpha1.AppTemplateValuesDownwardAPIItem{
					{Name: "parent.child", FieldPath: "metadata.name"},
					{Name: "parent.child1.child2", FieldPath: "metadata.namespace"},
					{Name: "parent.childwith\\.dot", FieldPath: "metadata.namespace"},
				}},
			}}

			paths, cleanup, err := subject.AsPaths(os.TempDir())
			require.NoError(t, err)
			require.Len(t, paths, 3)
			t.Cleanup(cleanup)

			assertFileContents(t, paths[0], "parent:\n  child: some-name\n")
			assertFileContents(t, paths[1], "parent:\n  child1:\n    child2: some-namespace\n")
			assertFileContents(t, paths[2], "parent:\n  childwith.dot: some-namespace\n")
		})

		t.Run("map field paths should allow subpaths", func(t *testing.T) {
			subject := subject
			subject.ValuesFrom = []v1alpha1.AppTemplateValuesSource{{DownwardAPI: &v1alpha1.AppTemplateValuesDownwardAPI{
				Items: []v1alpha1.AppTemplateValuesDownwardAPIItem{
					{Name: "a_some-annotation-key", FieldPath: "metadata.annotations['a_ann']"},
					{Name: "b_some-label-key", FieldPath: "metadata.labels['a_label']"},
				}},
			}}

			paths, cleanup, err := subject.AsPaths(os.TempDir())
			require.NoError(t, err)
			require.Len(t, paths, 2)
			t.Cleanup(cleanup)

			expectedValues := []string{
				"a_some-annotation-key: a_ann_val\n",
				"b_some-label-key: a_label_val\n",
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
			subject.appContext.Metadata.Annotations = map[string]string{
				"z_ann": "z_ann_val",
				"s_ann": "s_ann_val",
				"a_ann": "a_ann_val",
			}
			subject.ValuesFrom = []v1alpha1.AppTemplateValuesSource{{DownwardAPI: &v1alpha1.AppTemplateValuesDownwardAPI{
				Items: []v1alpha1.AppTemplateValuesDownwardAPIItem{
					{Name: "some-annotation-key", FieldPath: "metadata.annotations"},
					{Name: "parent.some-annotation-key", FieldPath: "metadata.annotations"},
				}},
			}}

			paths, cleanup, err := subject.AsPaths(os.TempDir())
			require.NoError(t, err)
			require.Len(t, paths, 2)
			t.Cleanup(cleanup)

			assertFileContents(t, paths[0], "some-annotation-key:\n  a_ann: a_ann_val\n  s_ann: s_ann_val\n  z_ann: z_ann_val\n")
			assertFileContents(t, paths[1], "parent:\n  some-annotation-key:\n    a_ann: a_ann_val\n    s_ann: s_ann_val\n    z_ann: z_ann_val\n")
		})

		t.Run("items with same key name, latter keys specified should clobber the earlier ones", func(t *testing.T) {
			subject := subject
			subject.ValuesFrom = []v1alpha1.AppTemplateValuesSource{{DownwardAPI: &v1alpha1.AppTemplateValuesDownwardAPI{
				Items: []v1alpha1.AppTemplateValuesDownwardAPIItem{
					{Name: "a_name", FieldPath: "metadata.name"},
					{Name: "a_name", FieldPath: "metadata.namespace"},
				}},
			}}

			paths, cleanup, err := subject.AsPaths(os.TempDir())
			require.NoError(t, err)
			require.Len(t, paths, 2)
			t.Cleanup(cleanup)

			assertFileContents(t, paths[1], "a_name: some-namespace\n")
		})

		t.Run("return helpful error if subpath is not found", func(t *testing.T) {
			subject := subject
			subject.ValuesFrom = []v1alpha1.AppTemplateValuesSource{{DownwardAPI: &v1alpha1.AppTemplateValuesDownwardAPI{
				Items: []v1alpha1.AppTemplateValuesDownwardAPIItem{
					{Name: "a_some-annotation-key", FieldPath: "metadata.annotations['INVALID_SUBPATH']"},
				}},
			}}

			_, _, err := subject.AsPaths(os.TempDir())
			require.Error(t, err)
			assert.ErrorContains(t, err, "INVALID_SUBPATH is not found")
		})

		t.Run("return helpful error if an unsupported downward api field spec is provided", func(t *testing.T) {
			subject := subject
			subject.ValuesFrom = []v1alpha1.AppTemplateValuesSource{{DownwardAPI: &v1alpha1.AppTemplateValuesDownwardAPI{
				Items: []v1alpha1.AppTemplateValuesDownwardAPIItem{
					{Name: "unsupportedFieldPath", FieldPath: "metadata.creationTimestamp"},
				}},
			}}

			_, _, err := subject.AsPaths(os.TempDir())
			require.Error(t, err)
			assert.ErrorContains(t, err, "creationTimestamp is not found")
		})

		t.Run("return helpful error if invalid nested key structure is provided", func(t *testing.T) {
			subject := subject
			subject.ValuesFrom = []v1alpha1.AppTemplateValuesSource{{DownwardAPI: &v1alpha1.AppTemplateValuesDownwardAPI{
				Items: []v1alpha1.AppTemplateValuesDownwardAPIItem{
					{Name: "parent.", FieldPath: "metadata.name"},
					{Name: ".parent", FieldPath: "metadata.name"},
					{Name: "parent..child", FieldPath: "metadata.name"},
				}},
			}}

			_, _, err := subject.AsPaths(os.TempDir())
			require.Error(t, err)
			assert.ErrorContains(t, err, "Invalid name was provided 'parent.' (hint: separate paths should only use a single '.' character)")
		})

		t.Run("return helpful error if multiple field spec is provided", func(t *testing.T) {
			subject := subject
			subject.ValuesFrom = []v1alpha1.AppTemplateValuesSource{{DownwardAPI: &v1alpha1.AppTemplateValuesDownwardAPI{
				Items: []v1alpha1.AppTemplateValuesDownwardAPIItem{
					{Name: "a_some-annotation-key", FieldPath: "metadata['name', 'uid']"},
				}},
			}}

			_, _, err := subject.AsPaths(os.TempDir())
			require.Error(t, err)
			assert.ErrorContains(t, err, "Invalid field spec provided to DownwardAPI. Only single supported fields are allowed")
		})

		t.Run("return kubernetes cluster version if not supplied", func(t *testing.T) {
			subject := subject
			subject.ValuesFrom = []v1alpha1.AppTemplateValuesSource{{DownwardAPI: &v1alpha1.AppTemplateValuesDownwardAPI{
				Items: []v1alpha1.AppTemplateValuesDownwardAPIItem{
					{Name: "k8s-version", KubernetesVersion: &v1alpha1.Version{}},
				}},
			}}
			subject.AdditionalValues = AdditionalDownwardAPIValues{
				KubernetesVersion: func() (string, error) {
					return "0.20.0", nil
				},
			}

			paths, cleanup, err := subject.AsPaths(os.TempDir())
			require.NoError(t, err)
			t.Cleanup(cleanup)

			require.Len(t, paths, 1)
			assertFileContents(t, paths[0], "k8s-version: 0.20.0\n")
		})

		t.Run("return kapp-controller version", func(t *testing.T) {
			subject := subject
			subject.ValuesFrom = []v1alpha1.AppTemplateValuesSource{{DownwardAPI: &v1alpha1.AppTemplateValuesDownwardAPI{
				Items: []v1alpha1.AppTemplateValuesDownwardAPIItem{
					{Name: "kc-version", KappControllerVersion: &v1alpha1.Version{}},
				}},
			}}
			subject.AdditionalValues = AdditionalDownwardAPIValues{
				KappControllerVersion: func() (string, error) {
					return "0.42.31337", nil
				},
			}

			paths, cleanup, err := subject.AsPaths(os.TempDir())
			require.NoError(t, err)
			t.Cleanup(cleanup)

			require.Len(t, paths, 1)
			assertFileContents(t, paths[0], "kc-version: 0.42.31337\n")
		})
	})
}

func assertFileContents(t *testing.T, path string, expectedVal string) {
	valueContents, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, expectedVal, string(valueContents))
}
