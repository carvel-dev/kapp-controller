# Setup a local Carvel Kapp Controller in KIND

The purpose of this doc is to walk through the necessary steps to setup Kapp-Controller inside of [KIND.](https://kind.sigs.k8s.io/docs/user/quick-start/)

This doc also provides examples on creating:

- `packageRepository`
- `package`
- `packageInstall`

## Prerequisite's

### KIND

Install [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)

mac: `brew install kind`

### Docker

Install [docker](https://docs.docker.com/get-docker/)

Turn it on.

### jq

Install [jq](https://stedolan.github.io/jq/download/)

mac: `brew install jq`


## KIND + docker registry

Setup KIND cluster + docker registry.

Use the setup script:

`./scripts/kind-with-registry.sh`

Ours is slightly updated, but original script can be found [here.](https://raw.githubusercontent.com/kubernetes-sigs/kind/main/site/static/examples/kind-with-registry.sh)

## KAPP Controller

Install the latest release of `kapp-controller`:

```sh
kubectl apply -f https://github.com/vmware-tanzu/carvel-kapp-controller/releases/latest/download/release.yml
```

### install RBAC for admin SA account

```
kubectl apply -f https://raw.githubusercontent.com/vmware-tanzu/carvel-kapp-controller/develop/examples/rbac/cluster-admin.yml
```

## Hack the gibson

Update `/etc/hosts` on our local machine to make accessing the local registry easier for all tools.

**NOTE**

We use the registry name `kind-registry.local` as it forces `imgpkg` to use http instead of https. [slack thread](https://kubernetes.slack.com/archives/CH8KCCKA5/p1654541811762389)


```sh
sudo vim /etc/hosts

# ADD
127.0.0.1 kind-registry.local
```

** REMEMBER TO DELETE THIS AFTER **

## Build test bundle + push to local docker registry

We use Istio's basic httpbin [example](https://raw.githubusercontent.com/istio/istio/master/samples/httpbin/httpbin.yaml) for our test service.

```sh
kbld -f dist/service --imgpkg-lock-output dist/.imgpkg/images.yml
```

setup latest *and* `0.1.0`

```sh
for tag in latest 0.1.0 
do
  imgpkg push -b kind-registry.local:5000/http-bin:$tag -f dist
done

```

### check container and tags

List all images in local registry

```sh
curl -s -X GET kind-registry.local:5000/v2/_catalog | jq .

{
  "repositories": [
    "http-bin"
  ]
}
```

check tags for `http-bin` in local registry

```sh
curl -s -X GET kind-registry.local:5000/v2/http-bin/tags/list | jq .

{
  "name": "http-bin",
  "tags": [
    "0.1.0",
    "sha256-463b022efdddc68b4bee71012d14be0d0ff22cf0ac53754fbb7400990106b3a1.imgpkg",
    "latest"
  ]
}
```

## Build package repository and push to docker registry

Build package repository bundle

```sh
kbld -f ./package-repository/packages --imgpkg-lock-output "./package-repository/.imgpkg/images.yml"
```

Push package repository to registry

```sh
for tag in latest 0.1.0
do
  imgpkg push -b kind-registry.local:5000/package-repository:$tag -f package-repository
done
```

## Install package repository

```sh
kubectl apply -f package_repository_install.yml
```

## Install the package

```sh
kubectl apply -f package_install.yml
```

## Verify the install(s)

Package Repository

```sh
❯ kubectl get packagerepository -n kapp-controller-packaging-global
NAME                            AGE    DESCRIPTION
package-repository.example.com   125m   Reconcile succeeded
```

Packages

```sh
❯ kubectl get packages -n default
NAME                         PACKAGEMETADATA NAME   VERSION   AGE
example.com.http-bin.0.1.0   example.com.http-bin   0.1.0     1m6s
```

PackageInstalls

```sh
❯ kubectl get packageinstall -n default
NAME       PACKAGE NAME           PACKAGE VERSION   DESCRIPTION           AGE
http-bin   example.com.http-bin   0.1.0             Reconcile succeeded   83s
```

Apps

```sh
❯ kubectl get apps -n default
NAME       DESCRIPTION           SINCE-DEPLOY   AGE
http-bin   Reconcile succeeded   113s           114s
```

Pods

```sh
❯ kubectl get pods -n default
NAME                       READY   STATUS    RESTARTS   AGE
httpbin-76778749f4-srfhq   1/1     Running   0          2m9s
```

