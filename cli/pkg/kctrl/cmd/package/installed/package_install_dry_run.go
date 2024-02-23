package installed

import (
	"fmt"

	versions "carvel.dev/vendir/pkg/vendir/versions/v1alpha1"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	kcpkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

type PackageInstalledDryRun struct {
	*CreateOrUpdateOptions
}

func (d PackageInstalledDryRun) PrintResources() error {
	const yamlSeperator = "---\n"
	packageInstall := &kcpkgv1alpha1.PackageInstall{
		TypeMeta:   metav1.TypeMeta{Kind: "PackageInstall", APIVersion: "packaging.carvel.dev/v1alpha1"},
		ObjectMeta: metav1.ObjectMeta{Name: d.Name, Namespace: d.NamespaceFlags.Name},
		Spec: kcpkgv1alpha1.PackageInstallSpec{
			PackageRef: &kcpkgv1alpha1.PackageRef{
				RefName: d.CreateOrUpdateOptions.packageName,
				VersionSelection: &versions.VersionSelectionSemver{
					Constraints: d.version,
					Prereleases: &versions.VersionSelectionSemverPrereleases{},
				},
			},
		},
	}

	rbacResourcesYAML := ""
	createRBAC := d.serviceAccountName == ""
	if createRBAC {
		packageInstall.Spec.ServiceAccountName = d.createdAnnotations.ServiceAccountAnnValue()

		serviceAccount := &corev1.ServiceAccount{
			TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "ServiceAccount"},
			ObjectMeta: metav1.ObjectMeta{
				Name:      d.createdAnnotations.ServiceAccountAnnValue(),
				Namespace: d.NamespaceFlags.Name,
				Annotations: map[string]string{
					KctrlPkgAnnotation: d.createdAnnotations.PackageAnnValue(),
					TanzuPkgAnnotation: d.createdAnnotations.PackageAnnValue(),
				},
			},
		}
		serviceAccountYAML, err := yaml.Marshal(serviceAccount)
		if err != nil {
			return fmt.Errorf("Marshaling ServiceAccount YAML: %s", err)
		}

		clusterRole := &rbacv1.ClusterRole{
			TypeMeta: metav1.TypeMeta{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "ClusterRole"},
			ObjectMeta: metav1.ObjectMeta{
				Name: d.createdAnnotations.ClusterRoleAnnValue(),
				Annotations: map[string]string{
					KctrlPkgAnnotation: d.createdAnnotations.PackageAnnValue(),
					TanzuPkgAnnotation: d.createdAnnotations.PackageAnnValue(),
				},
			},
			Rules: []rbacv1.PolicyRule{
				{APIGroups: []string{"*"}, Verbs: []string{"*"}, Resources: []string{"*"}},
			},
		}
		clusterRoleYAML, err := yaml.Marshal(clusterRole)
		if err != nil {
			return fmt.Errorf("Marshaling ClusterRole YAML: %s", err)
		}

		clusterRoleBinding := &rbacv1.ClusterRoleBinding{
			TypeMeta: metav1.TypeMeta{APIVersion: "rbac.authorization.k8s.io/v1", Kind: "ClusterRoleBinding"},
			ObjectMeta: metav1.ObjectMeta{
				Name: d.createdAnnotations.ClusterRoleBindingAnnValue(),
				Annotations: map[string]string{
					KctrlPkgAnnotation: d.createdAnnotations.PackageAnnValue(),
					TanzuPkgAnnotation: d.createdAnnotations.PackageAnnValue(),
				},
			},
			Subjects: []rbacv1.Subject{{Kind: KindServiceAccount.AsString(), Name: d.createdAnnotations.ServiceAccountAnnValue(), Namespace: d.NamespaceFlags.Name}},
			RoleRef: rbacv1.RoleRef{
				APIGroup: rbacv1.SchemeGroupVersion.Group,
				Kind:     KindClusterRole.AsString(),
				Name:     d.createdAnnotations.ClusterRoleAnnValue(),
			},
		}
		clusterRoleBindingYAML, err := yaml.Marshal(clusterRoleBinding)
		if err != nil {
			return fmt.Errorf("Marshaling ClusterRoleBinding YAML: %s", err)
		}

		rbacResourcesYAML = yamlSeperator + string(serviceAccountYAML) + yamlSeperator + string(clusterRoleYAML) + yamlSeperator + string(clusterRoleBindingYAML)
	} else {
		packageInstall.Spec.ServiceAccountName = d.serviceAccountName
	}

	secretResourcesYAML := ""
	createSecret := d.values && d.valuesFile != ""
	if createSecret {
		var err error
		dataValues := make(map[string]string)

		datavaluesBytes, err := cmdcore.NewInputFile(d.valuesFile).Bytes()
		if err != nil {
			return fmt.Errorf("Reading data values file '%s': %s", d.valuesFile, err.Error())
		}
		dataValues[valuesFileKey] = string(datavaluesBytes)

		secret := &corev1.Secret{
			TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Secret"},
			ObjectMeta: metav1.ObjectMeta{
				Name:      d.createdAnnotations.SecretAnnValue(),
				Namespace: d.NamespaceFlags.Name,
				Annotations: map[string]string{
					KctrlPkgAnnotation: d.createdAnnotations.PackageAnnValue(),
					TanzuPkgAnnotation: d.createdAnnotations.PackageAnnValue(),
				},
			},
			StringData: dataValues,
		}
		secretYAML, err := yaml.Marshal(secret)
		if err != nil {
			return fmt.Errorf("Marshaling Secret YAML: %s", err)
		}
		secretResourcesYAML = yamlSeperator + string(secretYAML)

		packageInstall.Spec.Values = []kcpkgv1alpha1.PackageInstallValues{
			{
				SecretRef: &kcpkgv1alpha1.PackageInstallValuesSecretRef{
					Name: d.createdAnnotations.SecretAnnValue(),
				},
			},
		}
	}
	d.addCreatedResourceAnnotations(&packageInstall.ObjectMeta, createRBAC, createSecret, false)

	packageInstallYAML, err := yaml.Marshal(packageInstall)
	if err != nil {
		return fmt.Errorf("Marshaling PackageInstall YAML: %s", err)
	}

	d.ui.PrintBlock([]byte(rbacResourcesYAML + secretResourcesYAML + yamlSeperator + string(packageInstallYAML)))
	return nil
}
