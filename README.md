## kapp-controller

kapp controller provides a way to declaratively specify which applications should run on your K8s cluster.

## Building

```
./hack/build.sh
```

To deploy:

```
ytt -f config/ | kbld -f- | kapp deploy -a kc -f- -c -y
```
