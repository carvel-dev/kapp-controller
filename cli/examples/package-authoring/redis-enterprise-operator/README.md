# redis-enterprise-operator

This example demonstrate how to package [`redis-enterprise-operator`](https://github.com/RedisLabs/redis-enterprise-k8s-docs). It gets the required kubernetes manifest from the [`git`](https://github.com/RedisLabs/redis-enterprise-k8s-docs) repository. 

This example uses the `master` branch to fetch required manifest.

## Change github reference

In case you want to use some other branch, rerun `pkg init`. During rerun, `kctrl` will read the values from `package-build.yml`, `vendir.yml` and present them as default values. When asked for the `Git Reference`, enter the desired git reference.
```shell
$ cd cli/examples/redis-enterprise-operator
$ kctrl pkg init
```

## Run kctrl pkg release to create new package.

To create a new package, rerun `kctrl pkg release`. You can provide a version to the package while running `kctrl pkg release` via flag.

During rerun, it will ask for a registry URL to push the imgpkg bundle. Ensure to provide a URL where you have push access.
```shell
$ cd cli/examples/redis-enterprise-operator
$ kctrl pkg release -v 2.0.0
```

Newly created `package.yml` and `metadata.yml` will be present in the `carvel-artifacts/packages/redis-enterprise-operator.carvel.dev/` directory.