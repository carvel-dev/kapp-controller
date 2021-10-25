// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

// import (
// 	"fmt"
// 	"strings"

// 	"github.com/pkg/errors"
// 	"github.com/spf13/cobra"
// 	apierrors "k8s.io/apimachinery/pkg/api/errors"

// 	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/component"
// 	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/kappclient"
// 	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/log"
// 	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/tkgpackageclient"
// )

// var packageAvailableGetCmd = &cobra.Command{
// 	Use:   "get PACKAGE_NAME or PACKAGE_NAME/VERSION",
// 	Short: "Get details for an available package or the openAPI schema of a package with a specific version",
// 	Args:  cobra.ExactArgs(1),
// 	Example: `
//     # Get package details for a package without specifying the version
//     tanzu package available get contour.tanzu.vmware.com --namespace test-ns

//     # Get package details for a package with specified version
//     tanzu package available get contour.tanzu.vmware.com/1.15.1-tkg.1-vmware1 --namespace test-ns

//     # Get openAPI schema of a package with specified version
//     tanzu package available get contour.tanzu.vmware.com/1.15.1-tkg.1-vmware1 --namespace test-ns --values-schema`,
// 	RunE:         packageAvailableGet,
// 	PreRunE:      validatePackage,
// 	SilenceUsage: true,
// }

// func init() {
// 	packageAvailableGetCmd.Flags().BoolVarP(&packageAvailableOp.ValuesSchema, "values-schema", "", false, "Values schema of the package, optional")
// 	packageAvailableCmd.AddCommand(packageAvailableGetCmd)
// }

// var pkgName string
// var pkgVersion string

// func validatePackage(cmd *cobra.Command, args []string) error {
// 	pkgNameVersion := strings.Split(args[0], "/")
// 	if len(pkgNameVersion) == 2 {
// 		pkgName = pkgNameVersion[0]
// 		pkgVersion = pkgNameVersion[1]
// 	} else if len(pkgNameVersion) == 1 {
// 		pkgName = pkgNameVersion[0]
// 	} else {
// 		return fmt.Errorf("package should be of the format name or name/version")
// 	}
// 	return nil
// }

// func packageAvailableGet(cmd *cobra.Command, args []string) error {
// 	kc, kcErr := kappclient.NewKappClient(packageAvailableOp.KubeConfig)
// 	if kcErr != nil {
// 		return kcErr
// 	}
// 	if packageAvailableOp.AllNamespaces {
// 		packageAvailableOp.Namespace = ""
// 	}

// 	if packageAvailableOp.ValuesSchema {
// 		if err := getValuesSchema(cmd, args, kc); err != nil {
// 			return err
// 		}
// 		return nil
// 	}

// 	t, err := component.NewOutputWriterWithSpinner(cmd.OutOrStdout(), getOutputFormat(),
// 		fmt.Sprintf("Retrieving package details for %s...", args[0]), true)
// 	if err != nil {
// 		return err
// 	}

// 	pkgMetadata, err := kc.GetPackageMetadataByName(pkgName, packageAvailableOp.Namespace)
// 	if err != nil {
// 		t.StopSpinner()
// 		if apierrors.IsNotFound(err) {
// 			log.Warningf("package '%s' does not exist in the '%s' namespace", pkgName, packageAvailableOp.Namespace)
// 			return nil
// 		}
// 		return err
// 	}

// 	if pkgVersion != "" {
// 		pkg, err := kc.GetPackage(fmt.Sprintf("%s.%s", pkgName, pkgVersion), packageAvailableOp.Namespace)
// 		if err != nil {
// 			if apierrors.IsNotFound(err) {
// 				return errors.Errorf("package '%s/%s' does not exist in the '%s' namespace", pkgName, pkgVersion, packageAvailableOp.Namespace)
// 			}
// 			return err
// 		}
// 		t.SetKeys("name", "version", "released-at", "display-name", "short-description", "package-provider", "minimum-capacity-requirements",
// 			"long-description", "maintainers", "release-notes", "license", "support", "category")
// 		t.AddRow(pkg.Spec.RefName, pkg.Spec.Version, pkg.Spec.ReleasedAt, pkgMetadata.Spec.DisplayName, pkgMetadata.Spec.ShortDescription,
// 			pkgMetadata.Spec.ProviderName, pkg.Spec.CapactiyRequirementsDescription, pkgMetadata.Spec.LongDescription, pkgMetadata.Spec.Maintainers,
// 			pkg.Spec.ReleaseNotes, pkg.Spec.Licenses, pkgMetadata.Spec.SupportDescription, pkgMetadata.Spec.Categories)

// 		t.RenderWithSpinner()
// 	} else {
// 		t.SetKeys("name", "display-name", "short-description", "package-provider", "long-description", "maintainers", "support", "category")
// 		t.AddRow(pkgMetadata.Name, pkgMetadata.Spec.DisplayName, pkgMetadata.Spec.ShortDescription,
// 			pkgMetadata.Spec.ProviderName, pkgMetadata.Spec.LongDescription, pkgMetadata.Spec.Maintainers, pkgMetadata.Spec.SupportDescription, pkgMetadata.Spec.Categories)

// 		t.RenderWithSpinner()
// 	}
// 	return nil
// }

// func getValuesSchema(cmd *cobra.Command, args []string, kc kappclient.Client) error {
// 	if pkgVersion == "" {
// 		return errors.New("version is required when --values-schema flag is declared. Please specify <PACKAGE-NAME>/<VERSION>")
// 	}
// 	pkg, pkgGetErr := kc.GetPackage(fmt.Sprintf("%s.%s", pkgName, pkgVersion), packageAvailableOp.Namespace)
// 	if pkgGetErr != nil {
// 		if apierrors.IsNotFound(pkgGetErr) {
// 			return errors.Errorf("package '%s/%s' does not exist in the '%s' namespace", pkgName, pkgVersion, packageAvailableOp.Namespace)
// 		}
// 		return pkgGetErr
// 	}

// 	t, err := component.NewOutputWriterWithSpinner(cmd.OutOrStdout(), outputFormat,
// 		fmt.Sprintf("Retrieving package details for %s...", args[0]), true)
// 	if err != nil {
// 		return err
// 	}

// 	var parseErr error
// 	dataValuesSchemaParser, parseErr := tkgpackageclient.NewValuesSchemaParser(pkg.Spec.ValuesSchema)
// 	if parseErr != nil {
// 		return parseErr
// 	}
// 	parsedProperties, parseErr := dataValuesSchemaParser.ParseProperties()
// 	if parseErr != nil {
// 		return parseErr
// 	}

// 	t.SetKeys("KEY", "DEFAULT", "TYPE", "DESCRIPTION")
// 	for _, v := range parsedProperties {
// 		t.AddRow(v.Key, v.Default, v.Type, v.Description)
// 	}
// 	t.RenderWithSpinner()

// 	return nil
// }
