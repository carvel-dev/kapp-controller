# kapp-controller

- Slack: [#k14s in Kubernetes slack](https://slack.kubernetes.io)
- [Docs](docs/README.md) with topics about installation, config, etc.
- Install: see [Install instructions](docs/install.md)

kapp controller provides a way to specify which applications should run on your K8s cluster. It will install, and continiously apply updates.

Features:

- supports fetching Helm charts (via `helm fetch`), git repos (via `git`), Docker images (via [imgpkg](https://github.com/k14s/imgpkg)), inline content within resource
- supports templating of Helm charts, [ytt](https://get-ytt.io) configuration (let us know what else we should support...)
- installs and syncs resources with [kapp](https://get-kapp.io)

More details in [docs](docs/README.md).
