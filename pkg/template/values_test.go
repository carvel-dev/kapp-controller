// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"os"
	"strings"
	"testing"

	"carvel.dev/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
				Name:      "some-name",
				Namespace: "some-namespace",
				UID:       "some-uid",
				Labels: map[string]string{
					"a_label":                 "a_label_val",
					"a_label.a_label2":        "a_label_val2",
					"a_label.a_label2/suffix": "a_label_val2_suffix",
				},
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
			"e_some-label-key:\n  a_label: a_label_val\n  a_label.a_label2: a_label_val2\n  a_label.a_label2/suffix: a_label_val2_suffix\n",
		}

		for i, p := range paths {
			assertFileContents(t, p, expectedValues[i])
		}

		t.Run("name should allow nested key structure", func(t *testing.T) {
			subject := subject
			subject.AdditionalValues = AdditionalDownwardAPIValues{
				KubernetesVersion: func() (string, error) {
					return "test-kubernetes-version", nil
				},
				KubernetesAPIs: func() ([]string, error) {
					return []string{"somegroup.example.com/someversion"}, nil
				},
				KappControllerVersion: func() (string, error) {
					return "test-kapp-controller-version", nil
				},
			}
			subject.ValuesFrom = []v1alpha1.AppTemplateValuesSource{{DownwardAPI: &v1alpha1.AppTemplateValuesDownwardAPI{
				Items: []v1alpha1.AppTemplateValuesDownwardAPIItem{
					{Name: "parent.child", FieldPath: "metadata.name"},
					{Name: "parent.child1.child2", FieldPath: "metadata.namespace"},
					{Name: "parent.childwith\\.dot", FieldPath: "metadata.namespace"},
					{Name: "parent.kubernetes_version", KubernetesVersion: &v1alpha1.Version{}},
					{Name: "parent.kubernetes_version_custom", KubernetesVersion: &v1alpha1.Version{Version: "test-kubernetes-version-custom"}},
					{Name: "parent.kubernetes_apis", KubernetesAPIs: &v1alpha1.KubernetesAPIs{}},
					{Name: "parent.kubernetes_apis_custom", KubernetesAPIs: &v1alpha1.KubernetesAPIs{GroupVersions: []string{"somecustomgroup.example.com/someversion"}}},
					{Name: "parent.kapp_controller_version", KappControllerVersion: &v1alpha1.Version{}},
					{Name: "parent.kapp_controller_version_custom", KubernetesVersion: &v1alpha1.Version{Version: "test-kapp-controller-version-custom"}},
				}},
			}}

			paths, cleanup, err := subject.AsPaths(os.TempDir())
			require.NoError(t, err)
			require.Len(t, paths, 9)
			t.Cleanup(cleanup)

			assertFileContents(t, paths[0], "parent:\n  child: some-name\n")
			assertFileContents(t, paths[1], "parent:\n  child1:\n    child2: some-namespace\n")
			assertFileContents(t, paths[2], "parent:\n  childwith.dot: some-namespace\n")
			assertFileContents(t, paths[3], "parent:\n  kubernetes_version: test-kubernetes-version\n")
			assertFileContents(t, paths[4], "parent:\n  kubernetes_version_custom: test-kubernetes-version-custom\n")
			assertFileContents(t, paths[5], "parent:\n  kubernetes_apis:\n  - somegroup.example.com/someversion\n")
			assertFileContents(t, paths[6], "parent:\n  kubernetes_apis_custom:\n  - somecustomgroup.example.com/someversion\n")
			assertFileContents(t, paths[7], "parent:\n  kapp_controller_version: test-kapp-controller-version\n")
			assertFileContents(t, paths[8], "parent:\n  kapp_controller_version_custom: test-kapp-controller-version-custom\n")
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

		t.Run("deals with complex labels", func(t *testing.T) {
			subject := subject
			subject.ValuesFrom = []v1alpha1.AppTemplateValuesSource{{DownwardAPI: &v1alpha1.AppTemplateValuesDownwardAPI{
				Items: []v1alpha1.AppTemplateValuesDownwardAPIItem{
					{Name: "a_some-label-key", FieldPath: "metadata.labels['a_label\\.a_label2/suffix']"},
				}},
			}}

			paths, cleanup, err := subject.AsPaths(os.TempDir())
			require.NoError(t, err)
			require.Len(t, paths, 1)
			t.Cleanup(cleanup)

			assertFileContents(t, paths[0], "a_some-label-key: a_label_val2_suffix\n")
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
	})
}

func assertFileContents(t *testing.T, path string, expectedVal string) {
	valueContents, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, expectedVal, string(valueContents))
}
