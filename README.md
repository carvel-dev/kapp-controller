# kapp-controller

- Slack: [#k14s in Kubernetes slack](https://slack.kubernetes.io)
- [Docs](docs/README.md) with topics about installation, config, etc.
- Install: see [Install instructions](docs/install.md)

kapp controller provides a way to specify which applications should run on your K8s cluster via one or more App CRs. It will install, and continiously apply updates.

Features:

- supports fetching git repos (via `git`), Helm charts (via `helm fetch`), Docker images (via [imgpkg](https://github.com/k14s/imgpkg)), inline content within resource
- supports templating of Helm charts, [ytt](https://get-ytt.io) configuration (let us know what else we should support...)
- installs and syncs resources with [kapp](https://get-kapp.io)
- [secure multi-tenant usage](docs/security-model.md) via service accounts and RBAC

More details in [docs](docs/README.md).
