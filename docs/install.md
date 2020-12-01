## Install

Grab the latest copy of YAML from the [Releases page](https://github.com/vmware-tanzu/carvel-kapp-controller/releases) and use your favorite deployment tool (such as [kapp](https://get-kapp.io) or kubectl) to install it.

Example:

```bash
$ kapp deploy -a kc -f https://github.com/vmware-tanzu/carvel-kapp-controller/releases/latest/download/release.yml
or
$ kubectl apply -f https://github.com/vmware-tanzu/carvel-kapp-controller/releases/latest/download/release.yml
```

### Advanced

`release.yml` is produced with [ytt](https://get-ytt.io) and [kbld](https://get-kbld.io) at the time of the release. You can use these tools yourself and customize the kapp controller configuration if the defaults do not not fit your needs.

Example:

```
$ git clone ...
$ kapp deploy -a kc -f <(ytt -f config/ | kbld -f-)
```

Next: [Walkthrough](walkthrough.md)
