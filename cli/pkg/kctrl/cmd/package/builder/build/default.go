package build

import (
	"sigs.k8s.io/yaml"
)

const defaultPackageBuildYAML = `
--- 
apiVersion: kctrl.carvel.dev/v1alpha1
kind: Config
spec: 
  imgpkg: 
    registryUrl: ~
    repoName: ~
    tag: v1.0.0
  package: 
    apiVersion: data.packaging.carvel.dev/v1alpha1
    kind: Package
    metadata: 
      name: samplepackage.corp.com.1.0.0
      namespace: default
    spec: 
      licenses: 
        - "Apache 2.0"
        - MIT
      refName: samplepackage.corp.com
      releaseNotes: "Initial release of the simple app package\n"
      template: 
        spec: 
          deploy: 
            - 
              kapp: {}
          fetch: 
            - 
              imgpkgBundle: 
                image: ~
          template: 
            - 
              ytt: 
                paths: 
                  - config/
            - 
              kbld: 
                paths: 
                  - .imgpkg/images.yml
                  - "-"
      version: "1.0.0"
  packageMetadata: 
    apiVersion: data.packaging.carvel.dev/v1alpha1
    kind: PackageMetadata
    metadata:
      name: samplepackage.corp.com
      namespace: default
    spec: 
      categories: 
        - demo
      displayName: "Simple App"
      longDescription: "Simple app consisting of a k8s deployment and service"
      shortDescription: "Simple app for demoing"
  vendir: 
    apiVersion: vendir.k14s.io/v1alpha1
    kind: Config
    minimumRequiredVersion: "0.12.0"
`

func NewDefaultPackageBuild() (PackageBuild, error) {
	var packageBuild PackageBuild
	err := yaml.Unmarshal([]byte(defaultPackageBuildYAML), &packageBuild)
	if err != nil {
		return PackageBuild{}, err
	}
	return packageBuild, nil
}
