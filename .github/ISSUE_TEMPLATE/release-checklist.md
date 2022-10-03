---
name: Release Checklist
about: Checklist for release kapp-controller
title: ''
labels: carvel, release
assignees: ''

---

- [ ] OSS Release
    - [ ] [Releasing via workflow](https://github.com/vmware-tanzu/carvel-kapp-controller/blob/develop/docs/dev.md#release)
    - [ ] Close any GitHub issues that have been delivered.
    - [ ] Add a link to the release on the issue.
- [ ] Update the packaging repositories [kctrl]
    - [ ] [Check if Update Homebrew](https://hackmd.io/uVpvITUuR4Cbwzkzb7MEpQ?view#Update-Homebrew)
    - [ ] [Check if Update Website Installation Script](https://hackmd.io/uVpvITUuR4Cbwzkzb7MEpQ?view#Update-Website-Installation-Script)
    - [ ] [Check Github Action](https://hackmd.io/uVpvITUuR4Cbwzkzb7MEpQ?view#Update-Github-Action)
- [ ] Update Documentation by [generating a new docs version](#Generate-new-docs-version)
- [ ] [Push any artifacts to a registry](https://hackmd.io/uVpvITUuR4Cbwzkzb7MEpQ?view#Push-OCI-Images-to-Registry)
- [ ] [Communicate in Slack](https://hackmd.io/uVpvITUuR4Cbwzkzb7MEpQ?view#Communicate-in-Slack)
- [ ] Communicate to the Carvel India team so they can update the release notes.
- [ ] [Add to "Announcements" in Next Community Meeting Agenda](https://hackmd.io/uVpvITUuR4Cbwzkzb7MEpQ?view#Announce-in-community-meeting)

Post Release:
- [ ] Create a Pull Request for [Tanzu Community Edition](https://github.com/vmware-tanzu/community-edition)
    - [ ] Create a new folder, and copy all the contents from the latest previous version e.g `cp -r addons/packages/kapp-controller/0.30.1 addons/packages/kapp-controller/0.30.2`.
    - [ ] Update `bundle/vendir.yaml` to the newly released tag version.
    - [ ] Update `addons/packages/kapp-controller/<your-version>/package.yaml` to the newly released tag version.
    - [ ] Run `make vendir-sync-package PACKAGE=kapp-controller VERSION=<your-version>`
    - [ ] Run `make lock-package-images PACKAGE=kapp-controller VERSION=<your-version>`. This will update the contents of `addons/packages/kapp-controller/0.30.1/bundle/.imgpkg/images.yml`.
    - [ ] Ensure there is only one item in the `images` array for `image.yml` above and that it is the correct sha for the released version. If there are multiple images in the `image.yml` please rerun the previous step.
    - [ ] Update `spec.template.spec.initContainers[0].image` to the release image sha if the value is not the same as above sha.
    - [ ] Run `make push-package PACKAGE=kapp-controller VERSION=<your-version> TAG=<your-version>`.
         - > NOTE: Ensure you are logged into the registry `docker login projects.registry.vmware.com`. Contact the [#tanzu-community-edition](https://kubernetes.slack.com/archives/C02GY94A8KT) slack on the kubernetes workspace if you do not have access to push.
    - [ ] Once you push the image, copy the above SHA generated and replace the image sha in `addons/packages/kapp-controller/<your-version>/package.yaml`.
    - [ ] Verify the generated package looks correct by running `ytt --ignore-unknown-comments -f addons/packages/kapp-controller/<your-version>/bundle/config > test.yaml`. (Don't include this test.yaml file in the PR)
