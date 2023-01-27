## Development & Deploy

### Prerequisites

You will need the following tools to build and deploy kapp-controller: 
* ytt
* kbld
* kapp

For linux/mac users, all the tools below can be installed by running `./hack/install-deps.sh`.

For windows users, please download the binaries from the respective GitHub repositories:
* https://github.com/carvel-dev/ytt
* https://github.com/carvel-dev/kbld
* https://github.com/carvel-dev/kapp

### Build

To build the kapp-controller project locally, run the following:
```
./hack/build.sh
```

### Deploy

The kapp-controller source can be built and deployed to a Kubernetes cluster using one of the options below.

#### minikube

```
eval $(minikube docker-env)
./hack/deploy.sh
```

*Note:* for rapid iteration while developing, you can run the script
`./hack/dev-deploy` which is much faster (0.5 minutes vs. 4+ minutes), but it
requires that you have previously run the full deploy at least once to set
things up.

#### Non-minikube environment

1. Change the [push_images property](https://github.com/carvel-dev/kapp-controller/blob/develop/config/values.yml#L10) to true
2. Change the [image_repo property](https://github.com/carvel-dev/kapp-controller/blob/develop/config/values.yml#L12) to the location to push the kapp-controller image
3. Run `./hack/deploy.sh`

*Note:* As above, while iterating you may prefer to run
`./hack/dev-deploy` which is much faster (0.5 minutes vs. 4+ minutes), but it
requires that you have previously run the full deploy at least once to set
things up. Additionally you may need to make changes directly to the overlay in
config-dev-deploy/build.yml.

#### secretgen-controller for private auth workflows

See more on kapp-controller's integration with secretgen-controller [here](https://carvel.dev/kapp-controller/docs/latest/private-registry-auth/).

```
# deploys secretgen-controller with kapp-controller
export KAPPCTRL_E2E_SECRETGEN_CONTROLLER=true

# use one of the methods above for where/how to deploy kapp-controller
./hack/deploy.sh
```

### Testing

kapp-controller has unit tests and e2e tests that can be run as documented below.

#### Unit Testing

```
./hack/test.sh
```

#### e2e Testing

```
# deploy kapp-controller to cluster with test assets
./hack/deploy-test.sh

# namespace where tests will be run on cluster
export KAPPCTRL_E2E_NAMESPACE=kappctrl-test

# run e2e test suite
./hack/test-all.sh
```

The `hack/test-e2e.sh` script (also run by `test-all`) will tee its output
to both your stdout and the (gitignored) file `tmp/e2eoutput.txt` so that you
can more easily grep/search/parse the output.

#### Benchmark Testing

Benchmark tests are run via github action. They can also be run locally via go
toolchain `go test ./test/bench/... -bench=.`

Benchmarks run on develop branch are
graphed in [github
pages](https://carvel-dev.github.io/kapp-controller/dev/bench/index.html).

### Profiling
1.) Enable profiling by editing config/values.yaml and setting `dangerous_enable_pprof`
to true
2.) deploy (see above)
3.) install graphviz: `brew install graphviz`
4.) If you're using minikube you can then get the url for pprof via `minikube service --url
pprof -n kapp-controller` - then append `/debug/pprof/` as there is no redirect
for `/`.
5.) consume data from the pprof server with your local toolchain. For instance
the below will show you the memory usage profile:
```
> export PROFURL=`minikube service --url pprof -n kapp-controller`
> go tool pprof -png $PROFURL/debug/pprof/heap > heap.png
> open heap.png
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

Push the tag to GitHub.
```
git push --tags
```

After pushing the tag to GitHub, the release process will be carried out by a GitHub workflow. 
This workflow will:
* Build and push the kapp-controller container image to ghcr.io
* Generate the release YAML (i.e. `release.yml`)
* Create a GitHub release draft 
* Upload `release.yml` to the GitHub release

After the release process finishes successfully, you can view the newly drafted release via 
the GitHub UI [here](https://github.com/vmware-tanzu/carvel-kapp-controller/releases). 

The newly published release will be available as a Draft (i.e. not available to the public).
Release notes can be autogenerated by GitHub, but make sure to call attention to any points 
that are not immediately clear from the autogenerated notes. Make sure to always thank external 
contributors for their additions to kapp-controller in the release notes.

Once the release notes are ready, clicking the `Publish release` button in the GitHub UI to 
make the release available to users.

#### LTS Releases

We want our releases to be sorted in semver order, but github only sorts on
semver order for releases on the same day (otherwise it sorts preferentially by
date). So we use sneaky post-dated annotated tags for LTS releases to pin them
to the date of the original release. For example:
```
 GIT_COMMITTER_DATE="2022-02-25 2:00" git tag -a -m "v0.30.2" "v0.30.2"
 ```

#### Release development process 

If you are making changes to the release process and want to test the process, it is recommended 
to work on a fork of kapp-controller instead of against the repository in the vmware-tanzu organization.

To do this, you can start by forking this repository. 

Next, head to the `Actions` tab of the fork you are using and enable GitHub Actions to run 
against this fork (i.e. By default, Actions do not run against forked repositories). 

Change the [`config-release/values.yml`](../config-release/values.yml) to point to your forked repository by changing 
`image_repo: ghcr.io/vmware-tanzu/carvel-kapp-controller` to `image_repo: ghcr.io/<YOUR GitHub USERNAME>/carvel-kapp-controller`.

After these steps have been carried out, you can trigger the release process by pushing a tag 
to your forked repository.

Some files to make note of when working on the release process:
* [`./hack/build-release.sh`](../hack/build-release.sh)
  - build-release.sh uses kbld and ytt to build and push the kapp-controller image to ghcr.io 
   and also generates the kapp-controller release.yml
* [`.github/workflows/release-process.yml`](../.github/workflows/release-process.yml)
  - GitHub Action workflow for release process

### Code Generation

kapp-controller relies on Kubernetes-focused code generation tools to generate the following:
* Custom resource definitions
* Clients, deep copy, lister, and informer logic for kapp-controller types
* API Server specifics: protobuf and openapi logic

Code generation should take place when API changes are made to kapp-controller resources 
(e.g. adding new fields). The CI for kapp-controller will always check if code generators 
should have been run and notify users to do so if needed, so do not worry if you are unsure. 
`./hack/verify-no-dirty-files.sh` can also be run locally to make sure any changes are ready 
to be checked in.

#### Prerequisite

Make sure to have [protoc installed](https://grpc.io/docs/protoc-installation/).

#### Run All Generators at Once

To run all generator scripts before checking changes in, use [`./hack/build-and-all-gen.sh`](../hack/build-and-all-gen.sh). 

#### Custom resource generation 

For CRD generation, kapp-controller makes use of [kubebuilder](https://book.kubebuilder.io/reference/generating-crd.html) 
via scripts and configuration.

Running `./hack/build.sh` calls out to [`./hack/gen-crds.sh`](../hack/gen-crds.sh). With every 
build of kapp-controller, the CRDs will be regenerated and output to [config/crds.yml](../config/crds.yml).

The `./hack/gen-crds.sh` script also makes use of a ytt overly in [`./hack/crd-overlay.yml`](../hack/crd-overlay.yml). 
This overlay removes unnecessary properties of the generated YAML.

#### Clients, deep copy, lister, and informer

To regenerate clients, deep copy, lister, and informer code, use [`./hack/gen.sh`](../hack/gen.sh).

#### API Server generation

To regenerate code for API Server updates, use [`./hack/gen-apiserver.sh`](../hack/gen-apiserver.sh).

#### Packaging Development

Due to the fact the one of our resources is named package, which is a golang
keyword, we were not able to use the code-generation binaries. To get around
this, we generated the code using the name pkg, and then manually edited those
files to enable us to use the name package. We will have to come up with a long 
term solution to enable us to use the code-generation binaries again.

### Continuous Integration/Jobs

kapp-controller uses GitHub Actions for all continuous integration for the project. 
You can find these CI processes under the [`.github/workflows`](../.github/workflows) 
folder. 

#### Pull Requests

On each pull request, the following CI processes run:
* [`test-gh`](../.github/workflows/test-gh.yml) - Builds kapp-controller, deploys build to minikube, runs unit tests, runs e2e tests.
* [`golangci-lint`](../.github/workflows/golangci-lint.yml) - Runs project linter. Configuration for linter is in [`.golangci.yml`](../.golangci.yml) file.
* [`test-kctrl-gh`](../.github/workflows/test-kctrl-gh.yml) - Runs build and tests for kapp-controller CLI.
* [`upgrade-testing`](../.github/workflows/upgrade-testing.yml) - This process deploys the latest released version of kapp-controller and then builds and 
  redeploys the changes submitted in the pull request. This helps to assure changes do not break upgrades between kapp-controller versions.
  
#### Daily Jobs

Each day, the following processes run:
* [`Trivy CVE Dependency Scanner`](../.github/workflows/trivy-scan.yml) - This job runs a [`trivy`](https://aquasecurity.github.io/trivy/) scan on 
the kapp-controller code base and latest release to identify CVEs.
* [`Mark issues stale and close stale issues`](../.github/workflows/stale-issues-action.yml) - This job marks any issues without a comment for 40 
days as a stale issue. If no comment is made in the issue, the issue will then be closed in the next 5 days.

#### Jobs Based on Events

The actions below are carried out when a certain event occurs:
* [`Remove label on close`](../.github/workflows/closed-issue.yml) - This job runs whenever an issue is closed. It removes the `carvel-triage` 
label from the closed issue to signal no further attention is needed on the issue.
* [`Closed issue comment labeling`](../.github/workflows/closed-issue-comment.yml) - This job runs whenever a comment is posted to a closed 
issue to signal maintainers should take a look.
* [`kapp-controller release`](../.github/workflows/release-process.yml) - This job carries out the kapp-controller release. More information 
available [here](#release).
