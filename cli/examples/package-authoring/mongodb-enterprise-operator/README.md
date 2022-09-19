# mongodb-enterprise-operator

This example demonstrates how to package [`mongodb-enterprise-operator`](https://github.com/mongodb/helm-charts/tree/main/charts/enterprise-operator) helm chart.

This example uses the version `1.16.0` of the `mongodb` helm chart.

## Change helm chart version

In case you want to use some other version, rerun `pkg init`. During rerun, `kctrl` will read the values from `package-build.yml`, `vendir.yml` and present them as default values. When asked for the `Helm Chart version`, enter the desired version.
```shell
$ cd cli/examples/mongodb-enterprise-operator
$ kctrl pkg init
```

## Run kctrl pkg release to create new package.

To create a new package, rerun `kctrl pkg release`. You can provide a version to the package while running `kctrl pkg release` via flag.

During rerun, it will ask for a registry URL to push the imgpkg bundle. Ensure to provide a URL where you have push access.
```shell
$ cd cli/examples/mongodb-enterprise-operator
$ kctrl pkg release -v 2.0.0
```

Newly created `package.yml` and `metadata.yml` will be present in the `carvel-artifacts/packages/mongodb-enterprise-operator.carvel.dev/` directory.