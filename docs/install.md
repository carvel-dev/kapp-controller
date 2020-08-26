## Install

Grab the latest copy of YAML from the [Releases page](https://github.com/k14s/kapp-controller/releases) and use your favorite deployment tool (such as [kapp](https://get-kapp.io) or kubectl) to install it.

Example:

```bash
$ kapp deploy -a kc -f https://github.com/k14s/kapp-controller/releases/latest/download/release.yml
or
$ kubectl apply -f https://github.com/k14s/kapp-controller/releases/latest/download/release.yml
```

**Note**: As of v0.6.0+, kapp-controller requires each App CR to specify a dedicated service account. This enables kapp-controller to be used _securely_ by users with different levels of privelege (namespace admin vs cluster admin) within the same cluster. If you want to configure kapp-controller to allow cluster admin level access for any user of App CR (not recommended!) you can _temporarily_ use `release-dangerous-allow-shared-sa.yml`. We are planning to remove this configuration in next few releases.

### Advanced

`release.yml` is produced with [ytt](https://get-ytt.io) and [kbld](https://get-kbld.io) at the time of the release. You can use these tools yourself and customize the kapp controller configuration if the defaults do not not fit your needs.

Example:

```
$ git clone ...
$ kapp deploy -a kc -f <(ytt -f config/ | kbld -f-)
```

Next: [Walkthrough](walkthrough.md)
