## Development & Deploy

Install ytt, kbld, kapp beforehand (https://k14s.io).

```
./hack/build.sh # to build locally

# deploys secretgen-controller with kapp-controller
# and also runs tests where kapp-controller integrates 
# with secretgen-controller
export KAPPCTRL_E2E_SECRETGEN_CONTROLLER=true

# add `-v image_repo=docker.io/username/kapp-controller` with your registry to ytt invocation inside
./hack/deploy.sh # to deploy

# deploys test assets in addition to kapp-controller for e2e tests
./hack/deploy-test.sh

export KAPPCTRL_E2E_NAMESPACE=kappctrl-test
./hack/test-all.sh
```

### Troubleshooting tips

1. If testing against a `minikube` cluster, run `eval $(minikube docker-env)` before development.

   This prevents the following error, which is a result of the docker daemon being unable to pull the `kapp-controller` dev image.

```
11:01:16AM:     ^ Pending: ImagePullBackOff (message: Back-off pulling image "kbld:kapp-controller-sha256-1bb8a9169c8265defc094a0220fa51d8c69a621d778813e4c4567d8cabde0e45")
11:01:05AM:     ^ Pending: ErrImagePull (message: rpc error: code = Unknown desc = Error response from daemon: pull access denied for kbld, repository does not exist or may require 'docker login': denied: requested access to the resource is denied)
```

### Release

Release versions are scraped from git tags in the same style as the goreleaser
tool.

Tag the release - it's necessary to do this first because the release process uses the latest tag to record the version.
```
git tag "v1.2.3"
```

Authenticate to the image registry where the image will be pushed (`i.e. index.docker.io/k14s`).

Build and push the kapp-controller image and generate the release YAML.
```
./hack/build-release.sh
```

The release YAML will be available as `./tmp/release.yml`.

Verify the release deploys successfully to a Kubernetes cluster.
```
kapp deploy -a kc -f ./tmp/release.yml
```

After verifying, push the tag to GitHub.
```
git push --tags
```

After pushing up the tag, you can `Draft a new release` through the GitHub UI and 
add release notes in the format shown [here](https://github.com/vmware-tanzu/carvel-kapp-controller/releases/tag/v0.20.0). 
Make sure to always thank external contributors for their additions to kapp-controller 
in the release notes.

As part of drafting the release through the GitHub UI, include the generated `release.yml` 
file and make sure to document the file checksum. This checksum is generated as part of 
the `./hack/build-release.sh` but can be rerun as `shasum -a 256 ./tmp/release*.yml`.

### Packaging Development

Due to the fact the one of our resources is named package, which is a golang
keyword, we were not able to use the code-generation binaries. To get around
this, we generated the code using the name pkg, and then manually edited those
files to enable us to use the name package. To avoid breaking this code, we are
commenting out the gen script on the packaging branch for extra safety. We will
have to come up with a long term solution to enable us to use the
code-generation binaries again.
