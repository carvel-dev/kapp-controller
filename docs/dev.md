## Development & Deploy

Install ytt, kbld, kapp beforehand (https://k14s.io).

```
./hack/build.sh # to build locally

# add `-v image_repo=docker.io/username/kapp-controller` with your registry to ytt invocation inside
./hack/deploy.sh # to deploy

export KAPPCTRL_E2E_NAMESPACE=kappctrl-test
./hack/test-all.sh
```
