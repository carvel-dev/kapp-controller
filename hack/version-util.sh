# source this file:
# use this via `source ./hack/version-util.sh` because `./hack/version-me.sh` won't make the functions available to other scripts

# We extract version information from git tags
# the implicit contract is that our git tags will be in ~semver (three-part) format and prefaced with the letter 'v'.
# this contract is required by the goreleaser tool and used throughout Carvel suite.

# git tag version extraction graciously provided by https://github.com/vmware-tanzu/carvel-imgpkg/blob/develop/hack/build-binaries.sh
function get_latest_git_tag {
  git describe --tags | grep -Eo 'v[0-9]+\.[0-9]+\.[0-9]+(-alpha\.[0-9]+)?'
}

function get_kappctrl_ver {
  echo "${1:-`get_latest_git_tag`}"
}

