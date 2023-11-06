## Examples

Since you need to provide service account for App CRs, we've included two common service account configurations:

- https://github.com/carvel-dev/kapp-controller/blob/master/examples/rbac/default-ns.yml: It creates `default-ns-sa` service account in `default` namespace that allows to change any resource in `default` namespace. (Example usage: `simple-app-http.yml`)

- https://github.com/carvel-dev/kapp-controller/blob/master/examples/rbac/cluster-admin.yml: It creates `cluster-admin-sa` service account within `default` namespace that allows to change _any_ resource in the cluster. (Example usage: `istio-knative.yml`)

```bash
$ kapp deploy -a default-ns-rbac -f https://raw.githubusercontent.com/carvel-dev/kapp-controller/develop/examples/rbac/default-ns.yml
```

Once that's done, deploy any example in this repo.
