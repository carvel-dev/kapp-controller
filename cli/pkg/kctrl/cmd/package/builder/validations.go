package builder

import (
	"fmt"
	"strings"

	"github.com/vmware-tanzu/carvel-vendir/pkg/vendir/versions"
	"k8s.io/apimachinery/pkg/util/validation"
)

//TODO should we use the same validation used in kapp controller. But that accepts other parameter. ValidatePackageMetadataName in validations.go file
func validateFQName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("Fully Qualified Name of a package cannot be empty")
	}
	if errs := validation.IsDNS1123Subdomain(name); len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, ","))
	}
	if len(strings.Split(name, ".")) < 3 {
		return fmt.Errorf("should be a fully qualified name with at least three segments separated by dots")
	}
	return nil
}

func validatePackageSpecVersion(version string) error {
	if version == "" {
		return fmt.Errorf("Version cannot be empty")
	}
	if _, err := versions.NewSemver(version); err != nil {
		return fmt.Errorf("must be valid semver: %v", err)
	}
	return nil
}
