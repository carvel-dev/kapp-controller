![logo](docs/CarvelLogo.png)

# kapp-controller

- Slack: [#carvel in Kubernetes slack](https://slack.kubernetes.io)
- [Docs](docs/README.md) with topics about installation, config, etc.
- Install: see [Install instructions](docs/install.md)

kapp controller provides a way to specify which applications should run on your K8s cluster via one or more App CRs. It will install, and continuously apply updates.

Features:

- supports fetching git repos (via `git`), Helm charts (via `helm fetch`), Docker images (via [imgpkg](https://github.com/k14s/imgpkg)), inline content within resource
- supports templating of Helm charts, [ytt](https://get-ytt.io) configuration (let us know what else we should support...)
- installs and syncs resources with [kapp](https://get-kapp.io)
- [secure multi-tenant usage](docs/security-model.md) via service accounts and RBAC

More details in [docs](docs/README.md).

### Join the Community and Make Carvel Better
Carvel is better because of our contributors and maintainers. It is because of you that we can bring great software to the community.
Please join us during our online community meetings. Details can be found on our [Carvel website](https://carvel.dev/community/).


You can chat with us on Kubernetes Slack in the #carvel channel and follow us on Twitter at @carvel_dev.

Check out which organizations are using and contributing to Carvel: [Adopter's list](https://github.com/vmware-tanzu/carvel/blob/master/ADOPTERS.md)
