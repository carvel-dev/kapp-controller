# Return a GOPATH to a temp directory. Works around the out-of-GOPATH issues
# for k8s client gen mixed with go mod.
# Intended to be used like:
#   export GOPATH=$(go_mod_gopath_hack)
function go_mod_gopath_hack() {
  local tmp_dir=$(mktemp -d)
  local module="$(go list -m)"

  local tmp_repo="${tmp_dir}/src/${module}"
  mkdir -p "$(dirname ${tmp_repo})"
  ln -s "$PWD" "${tmp_repo}"

  echo "${tmp_dir}"
}
