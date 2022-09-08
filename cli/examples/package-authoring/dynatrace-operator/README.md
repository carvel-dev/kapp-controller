# dynatrace-operator

This example demonstrate how to package [`dynatrace-operator`](https://github.com/Dynatrace/dynatrace-operator) Kubernetes manifest. Dynatrace operator releases `kubernetes.yaml` file which can be packaged into a Carvel package. `kubernetes.yml` can be fetched from the [`github release artifact`](https://github.com/Dynatrace/dynatrace-operator/releases). 

This example uses the `v0.6.0` release.

## Change Release version

In case you want to use some other release version, rerun `pkg init`. During rerun, `kctrl` will read the values from `package-build.yml`, `vendir.yml` and present them as default values. When asked for the `Release tag`, enter the desired release version.
```shell
$ cd cli/examples/dynatrace-operator
$ kctrl pkg init
```

## Run kctrl pkg release to create new package.

To create a new package, rerun `kctrl pkg release`. You can provide a version to the package while running `kctrl pkg release` via flag.

During rerun, it will ask for a registry URL to push the imgpkg bundle. Ensure to provide a URL where you have push access.
```shell
$ cd cli/examples/dynatrace-operator
$ kctrl pkg release -v 2.0.0
```

Newly created `package.yml` and `metadata.yml` will be present in the `carvel-artifacts/packages/dynatrace.carvel.dev/` directory.