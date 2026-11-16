#!/usr/bin/env python3

"""Resolve the nightly source revisions once, before any image builds start."""

import json
import re
import subprocess
import sys


def remote_refs(repository, *patterns):
    result = subprocess.run(
        ["git", "ls-remote", "--refs", f"https://github.com/{repository}.git", *patterns],
        check=True, capture_output=True, text=True,
    )
    refs = {}
    for line in result.stdout.splitlines():
        sha, ref = line.split()
        if not re.fullmatch(r"[0-9a-f]{40}", sha):
            raise ValueError(f"Invalid SHA for {repository}: {ref}")
        refs[ref] = sha
    if not refs:
        raise ValueError(f"Failed to resolve refs for {repository}")
    return refs


def kubernetes_matrix(refs):
    releases = {}
    for ref in refs:
        match = re.fullmatch(r"refs/heads/release-1\.(\d+)", ref)
        if match:
            releases[int(match[1])] = ref
    if not releases:
        raise ValueError("Failed to resolve Kubernetes release branches")

    # Exercise every minor from N-3 through N, plus the development branch.
    # Require a complete window rather than silently selecting older branches.
    latest = max(releases)
    selected = []
    for minor in range(latest - 3, latest + 1):
        if minor not in releases:
            raise ValueError(f"Missing Kubernetes release-1.{minor} in support window")
        selected.append(releases[minor])
    if "refs/heads/master" not in refs:
        raise ValueError("Failed to resolve Kubernetes master")
    selected.append("refs/heads/master")
    return {"include": [
        {"k8s_branch": ref.removeprefix("refs/heads/"), "k8s_sha": refs[ref]}
        for ref in selected
    ]}


def controller_refs(refs):
    # Ignore prereleases and choose the newest patch of each of the last two minors.
    releases = {}
    for ref in refs:
        match = re.fullmatch(r"refs/tags/v(\d+)\.(\d+)\.(\d+)", ref)
        if match:
            major, minor, patch = map(int, match.groups())
            releases[(major, minor)] = max(patch, releases.get((major, minor), -1))
    if len(releases) < 2:
        raise ValueError("Failed to resolve two kapp-controller release minors")
    return ["develop"] + [
        f"v{major}.{minor}.{releases[(major, minor)]}"
        for major, minor in sorted(releases)[-2:]
    ]


def main():
    matrix = kubernetes_matrix(remote_refs(
        "kubernetes/kubernetes", "refs/heads/release-1.*", "refs/heads/master",
    ))
    kind_refs = remote_refs("kubernetes-sigs/kind", "refs/heads/main")
    kind_sha = kind_refs["refs/heads/main"]
    targets = controller_refs(remote_refs("carvel-dev/kapp-controller", "refs/tags/v*"))
    # Lines can be appended directly to GITHUB_OUTPUT; all values are single-line.
    print(f"matrix={json.dumps(matrix, separators=(',', ':'))}")
    print(f"kind_sha={kind_sha}")
    print(f"controller_refs={json.dumps(targets, separators=(',', ':'))}")


if __name__ == "__main__":
    try:
        main()
    except (ValueError, KeyError, subprocess.CalledProcessError) as error:
        print(f"::error::Nightly version resolution failed: {error}", file=sys.stderr)
        sys.exit(1)
