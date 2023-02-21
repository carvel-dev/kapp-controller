---
name: Release Checklist
about: Checklist for release kapp-controller
title: ''
labels: carvel, release
assignees: ''

---

## Releasing a new minor / major:
- [ ] OSS Release
    - [ ] [Releasing via workflow](https://github.com/carvel-dev/kapp-controller/blob/develop/docs/dev.md#release).
    - [ ] Close any GitHub issues that have been delivered.
    - [ ] Add a link to the release on the issue.
    - [ ] Communicate to the kctrl maintainers, so they can update the release notes.
    - [ ] Press the Publish Release button
- [ ] Update the packaging repositories [kctrl]
    - [ ] [Check if Update Homebrew](https://hackmd.io/uVpvITUuR4Cbwzkzb7MEpQ?view#Update-Homebrew)
    - [ ] [Check if Update Website Installation Script](https://hackmd.io/uVpvITUuR4Cbwzkzb7MEpQ?view#Update-Website-Installation-Script)
    - [ ] [Check Github Action](https://hackmd.io/uVpvITUuR4Cbwzkzb7MEpQ?view#Update-Github-Action)
- [ ] Update Documentation by [generating a new docs version](https://hackmd.io/uVpvITUuR4Cbwzkzb7MEpQ?view#Generate-new-docs-version)
- [ ] [Push any artifacts to a registry](https://hackmd.io/uVpvITUuR4Cbwzkzb7MEpQ?view#Push-OCI-Images-to-Registry)
- [ ] [Communicate in Slack](https://hackmd.io/uVpvITUuR4Cbwzkzb7MEpQ?view#Communicate-in-Slack)
- [ ] [Add to "Announcements" in Next Community Meeting Agenda](https://hackmd.io/uVpvITUuR4Cbwzkzb7MEpQ?view#Announce-in-community-meeting)

## Releasing a patch version and backporting a CVE:
- [ ] Validate which branch lines to backport the CVE to. Based on our [private confluence page](https://confluence.eng.vmware.com/x/FyIuSQ).
- [ ] For each line, e.g `v0.30.x`, `v0.38.x`, do the following:
    - [ ] Validate that the branch contains the latest patches, that no newer code was forgotten to be merged back in.
    - [ ] `git checkout v0.38.x`.
    - [ ] `git checkout -b v0.38.<next-patch-version>`.
    - [ ] Make the necessary fixes / cherry-picks.
    - [ ] `git push origin v0.38.<next-patch-version>`.
    - [ ] Make a PR.
    - [ ] Once approved, merge the changes back to the `v0.38.x` branch and `git push` the branch and delete your temporary branch used in the PR.
    - [ ] To Release: follow the instructions FROM THE BRANCH YOU ARE UPDATING at `docs/dev.md#release` in the repository. These will contain the relevant steps at each point of time in the project's history, e.g when updating `v0.25.x` the url will look like: https://github.com/carvel-dev/kapp-controller/blob/v0.25.x/docs/dev.md#release
