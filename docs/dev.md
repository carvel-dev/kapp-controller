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

### Release

Release versions are scraped from git tags in the same style as the goreleaser
tool.

```
# create and push the release tag (see `git tag --list` for examples)
./hack/build-release.sh
```

### Packaging Development

Due to the fact the one of our resources is named package, which is a golang
keyword, we were not able to use the code-generation binaries. To get around
this, we generated the code using the name pkg, and then manually edited those
files to enable us to use the name package. To avoid breaking this code, we are
commenting out the gen script on the packaging branch for extra safety. We will
have to come up with a long term solution to enable us to use the
code-generation binaries again.
